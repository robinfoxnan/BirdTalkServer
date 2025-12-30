package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// 群组相关的基本操作
func handleGroupOp(msg *pbmodel.Msg, session *Session) {
	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	if groupOpMsg == nil {
		Globals.Logger.Debug("receive wrong group op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group op  is null", nil, session)
		return
	}

	// 都需要验证是否登录与权限
	ok := checkUserLogin(session)
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
		return
	}

	groupOpMsg.MsgId = Globals.snow.GenerateID()

	opCode := groupOpMsg.Operation
	switch opCode {
	case pbmodel.GroupOperationType_GroupCreate: // 创建
		handleGroupCreateOp(msg, session)
		break
	case pbmodel.GroupOperationType_GroupDissolve: // 解散
		handleGroupDissolveOp(msg, session)
		break
	case pbmodel.GroupOperationType_GroupSetInfo: // 设置信息
		handleGroupSetBasicInfo(msg, session)
		break
	case pbmodel.GroupOperationType_GroupKickMember: // 踢人
		handleGroupKickOut(msg, session)
		break
	case pbmodel.GroupOperationType_GroupInviteRequest: // 邀请
		handleInviteSomeone(msg, session)
		break
	case pbmodel.GroupOperationType_GroupInviteAnswer: // 邀请的应答
		handleInviteAnswer(msg, session)
		break // 邀请后处理结果
	case pbmodel.GroupOperationType_GroupJoinRequest: // 加入请求
		handleGroupJoinReq(msg, session)
		break
	case pbmodel.GroupOperationType_GroupJoinAnswer: // 加入请求的处理，同意、拒绝、问题
		handleGroupJoinAnswer(msg, session)
		break
	case pbmodel.GroupOperationType_GroupQuit: // 退出群组
		handleGroupMemberQuit(msg, session)
		break
	case pbmodel.GroupOperationType_GroupAddAdmin: // 增加管理员
		handleGroupSetSomeoneAsAdmin(msg, session)
		break
	case pbmodel.GroupOperationType_GroupDelAdmin: // 删除管理员
		handleGroupRemoveSomeoneFromAdmin(msg, session)
		break
	case pbmodel.GroupOperationType_GroupTransferOwner: // 转让群主
		handleGroupTransferOwner(msg, session)
		break
	case pbmodel.GroupOperationType_GroupSetMemberInfo: // 设置自己在群中的信息
		handleSetMemberInfo(msg, session)
		break

	case pbmodel.GroupOperationType_GroupSearch: // 搜素群组
		handleGroupSearch(msg, session)
		break
	case pbmodel.GroupOperationType_GroupSearchMember: // 人员
		handleGroupSearchMember(msg, session)
	case pbmodel.GroupOperationType_GroupListIn:
		handleListMemberInG(msg, session)
		break
	}

	return
}

// 用户会应答邀请，管理员会应答申请操作，这里需要处理并转发
func handleGroupOpRet(msg *pbmodel.Msg, session *Session) {

	groupOpMsgRet := msg.GetPlainMsg().GetGroupOpRet()
	if groupOpMsgRet == nil {
		Globals.Logger.Debug("receive wrong group op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group op  is null", nil, session)
		return
	}

	// 都需要验证是否登录与权限
	ok := checkUserLogin(session)
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
		return
	}

	opCode := groupOpMsgRet.Operation
	switch opCode {
	case pbmodel.GroupOperationType_GroupInviteAnswer:
		break
	case pbmodel.GroupOperationType_GroupJoinAnswer:
		break
	default:
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTWrongCode), "group op ret has a wrong opcode", nil, session)
	}
	return
}

// 创建群
// .name = "ddd"
// .groupType = "chat" | "channel" | "map"
// .tags = []
// params[visibility"]:  "public" | "private"
// params["brief"]: 简介
// params["icon"]: ""
// params["jointype"] = "direct" | "invite" | "auth" | "question"
func handleGroupCreateOp(msg *pbmodel.Msg, session *Session) {
	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	groupInfo := groupOpMsg.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group create operation group info is null", nil, session)
		return
	}

	if len(groupInfo.GroupType) == 0 {
		groupInfo.GroupType = "chat"
	} else {
		str := strings.ToLower(groupInfo.GroupType)
		if str != "chat" && str != "channel" && str != "map" {
			groupInfo.GroupType = "chat"
		} else {
			groupInfo.GroupType = str
		}
	}

	var err error
	groupInfo.GroupId, err = Globals.redisCli.GetNextGroupId()
	if len(groupInfo.GroupName) == 0 {
		groupInfo.GroupName = "group" + strconv.FormatInt(groupInfo.GroupId, 10)
	}

	params := groupInfo.GetParams()
	if params != nil {
		v, ok := params["visibility"]
		if ok {

			temp := strings.ToLower(v)
			if temp != "private" && temp != "public" {
				params["visibility"] = "public"
			}
		}

		v, ok = params["brief"]
		if !ok {
			params["brief"] = ""
		}

		v, ok = params["icon"]
		if !ok {
			params["icon"] = ""
		}

		v, ok = params["jointype"]
		if ok {
			temp := strings.ToLower(v)
			if temp != "direct" && temp != "invite" && temp != "auth" && temp != "question" {
				params["jointype"] = "direct"
			}
		} else {
			params["jointype"] = "direct"
		}

	} else {
		params = map[string]string{
			"visibility": "public",
			"brief":      "",
			"icon":       "",
			"jointype":   "direct",
		}
		groupInfo.Params = params
	}

	// 保存到数据库
	group, err := saveNewGroup(groupInfo, session)
	if group == nil || err != nil {
		Globals.Logger.Error("saveNewGroup()", zap.Error(err))
		return
	}

	// 如果设置了初始用户，则需要
	groupMems := msg.GetPlainMsg().GetGroupOp().GetMembers()
	if groupMems != nil && len(groupMems) > 0 {
		//members := make([]model.GroupMemberStore, len(groupMems))

		for _, mem := range groupMems {

			memUser, b, _ := findUserInfo(mem.UserId)
			if memUser == nil || b == false {
				continue
			}

			member := model.GroupMemberStore{
				Pk:   db.ComputePk(groupInfo.GroupId),
				Gid:  groupInfo.GroupId,
				Uid:  session.UserID,
				Tm:   utils.GetTimeStamp(),
				Role: model.RoleGroupMember,
				Nick: memUser.NickName,
			}

			item := model.UserInGStore{
				Pk:  db.ComputePk(session.UserID),
				Uid: session.UserID,
				Gid: groupInfo.GroupId,
			}

			// 保存数据库
			Globals.scyllaCli.InsertGroupMember(&member, &item)

			// redis, 用户所在群
			Globals.redisCli.SetUserJoinGroup(session.UserID, groupInfo.GroupId, memUser.NickName)

			// 添加到内存
			group.AddMember(memUser.UserId, memUser.NickName)
		}
	}
	// todo:保存操作记录

	// 通知相关用户，建立了新群
	retMsg := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupCreate,
		groupInfo,
		groupOpMsg.GetReqMem(),
		groupOpMsg.GetMembers(),
		groupOpMsg.SendId,
		groupOpMsg.MsgId,
		"ok", "", session)

	Globals.Logger.Debug("return create group ret", zap.Any("msg", retMsg))
	//fmt.Println(retMsg)

	notifyGroupMembers(groupInfo.GroupId, retMsg)

	return
}

// 解散群
func handleGroupDissolveOp(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation group info id is wrong", nil, session)
		return
	}
	if !group.IsOwner(session.UserID) {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation, you are not owner", nil, session)
		return
	}

	// 解散群： 在基础信息中设置标记
	_, err := Globals.mongoCli.UpdateGroupInfoPart(groupInfo.GroupId, map[string]interface{}{"params.status": "deleted"}, nil)
	if err != nil {

	}
	err = Globals.redisCli.SetGroupInfoPart(groupInfo.GroupId, "Params.status", "deleted")
	group.SetDeleted()

	// 删除群内的所有用户
	err = Globals.scyllaCli.DissolveGroupAllMember(db.ComputePk(groupInfo.GroupId), groupInfo.GroupId)
	err = Globals.redisCli.RemoveAllUserOfGroup(groupInfo.GroupId)

	// 群内还需要通知用户解散
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupDissolve, groupInfo,
		reqMem,
		nil,
		msg.GetPlainMsg().GetGroupOpRet().SendId,
		Globals.snow.GenerateID(),
		"ok",
		"", session)

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	if Globals.Config.Server.ClusterMode {

	}

	// 清理内存
	group.ClearMember()

	// 各个用户所在群，删除一个
	membersId := group.GetMembers()
	for _, mId := range membersId {
		err = Globals.scyllaCli.DeleteUserInG(db.ComputePk(mId), mId, groupInfo.GroupId)
		Globals.redisCli.SetUserLeaveGroup(mId, groupInfo.GroupId)
		user, ok := Globals.uc.GetUser(mId)
		if user != nil && ok {
			user.SetLeaveGroup(groupInfo.GroupId)
		}
	}

	// redis中群组用户分布情况
	Globals.redisCli.RemoveActiveGroupRelated(groupInfo.GroupId)

	return
}

// 设置基础信息，这个也是只有管理员才可以设置，不同于微信的
func handleGroupSetBasicInfo(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set info operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set info operation group info id is wrong", nil, session)
		return
	}

	isAdmin := group.IsAdmin(session.UserID)
	if !isAdmin {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set info operation, you are not admin", nil, session)
		return
	}
	// 先合并到内存，更新到数据库和redis
	group.MergeGroupInfo(groupInfo)
	_, err := Globals.mongoCli.UpdateGroupInfo(group.GetGroupInfo())
	if err != nil {

	}

	err = Globals.redisCli.SetGroupInfo(group.GetGroupInfo())
	if err != nil {

	}
	// 通知所有的用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSetInfo,
		group.GetGroupInfo(),
		reqMem,
		nil,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"", session)

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	if Globals.Config.Server.ClusterMode {

	}

}

// 踢人
func handleGroupKickOut(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation group info id is wrong", nil, session)
		return
	}

	isAdmin := group.IsAdmin(session.UserID)
	if !isAdmin {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation, you are not admin", nil, session)
		return
	}

	opMsg := msg.GetPlainMsg().GetGroupOp()
	uid := int64(0)
	str, ok := opMsg.Params["uid"]
	if !ok || str == "" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation, params.uid is wrong", nil, session)
		return
	}

	uid, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation, params.uid is not an id", nil, session)
		return
	}

	// 同时操作2个表
	Globals.scyllaCli.DeleteGroupMember(db.ComputePk(groupInfo.GroupId), db.ComputePk(uid), groupInfo.GroupId, uid)

	// 从成员表中删除
	// 标记用户的所属群
	Globals.redisCli.SetUserLeaveGroup(uid, groupInfo.GroupId)

	// 从分布表中删除
	index := int64(Globals.Config.Server.HostIndex)
	Globals.redisCli.RemoveActiveGroupMembersLua(groupInfo.GroupId, index, []int64{uid})

	// 内存中删除成员
	group.RemoveMember(uid)

	// 内存中，用户退出组，
	user, ok := Globals.uc.GetUser(uid)
	if user != nil {
		user.SetLeaveGroup(groupInfo.GroupId)
	}

	// 通知所有的用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}

	kickMember := &pbmodel.GroupMember{
		UserId:  uid,
		Nick:    "",
		Icon:    "",
		Role:    "",
		GroupId: groupInfo.GroupId,
		Params:  nil,
	}
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupKickMember,
		group.GetGroupInfo(),
		reqMem,
		[]*pbmodel.GroupMember{
			kickMember,
		},
		msg.GetPlainMsg().GetGroupOpRet().SendId,
		Globals.snow.GenerateID(),
		"ok",
		"", session)

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}

}

// 邀请某人
func handleInviteSomeone(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	isPrivate := group.IsPrivate()
	if isPrivate {
		isAdmin := group.IsAdmin(session.UserID)
		if !isAdmin {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "you are not the admin of the private group", nil, session)
			return
		}
	}

	memList := msgOp.GetMembers()
	if memList == nil || len(memList) < 1 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must give out the invitee in members", nil, session)
		return
	}

	user := session.GetUser()
	reqMem := msgOp.GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    user.NickName,
			Icon:    user.Icon,
			Role:    "",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}

	for _, mem := range memList {
		msgId := Globals.snow.GenerateID()
		// 用户不存在就不能继续操作，否则后续的用户注册进来会造成脏数据
		memUser, _, err := findUserInfo(mem.UserId)
		if memUser == nil {
			continue
		}
		code := utils.GenerateCheckCode(6)
		code32, _ := strconv.Atoi(code)

		record := model.CommonOpStore{
			Pk:   db.ComputePk(session.UserID),
			Uid1: session.UserID,
			Uid2: mem.UserId,
			Gid:  groupInfo.GroupId,
			Id:   msgId,
			Usid: msgOp.SendId,
			Tm:   time.Now().UnixMilli(),
			Tm1:  0,
			Tm2:  0,
			Io:   0,
			St:   0,
			Cmd:  model.CommonGroupOpInviteRequest,
			Ret:  0,
			Mask: int32(code32),
			Ref:  0,
			Draf: nil,
		}
		// todo: save the draft
		pk2 := db.ComputePk(mem.UserId)
		err = Globals.scyllaCli.SaveUserOp(&record, pk2)
		if err != nil {
			continue
		}

		// 通知被邀请人
		msgNotice := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupInviteRequest,
			group.GetGroupInfo(),
			reqMem,
			[]*pbmodel.GroupMember{
				mem,
			},
			msgOp.SendId,
			msgId,
			"notify",
			"group invitation", session)
		trySendMsgToUser(mem.UserId, msgNotice)

		// 恢复发出邀请的用户
		msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupInviteRequest,
			group.GetGroupInfo(),
			reqMem,
			[]*pbmodel.GroupMember{mem},
			msgOp.SendId,
			msgId,
			"wait",
			"group invitation", session)
		trySendMsgToUser(session.UserID, msgRet)
	}

}

// 邀请的回答
func handleInviteAnswer(msg *pbmodel.Msg, session *Session) {
	msgOpRet := msg.GetPlainMsg().GetGroupOpRet()

	user := session.GetUser()
	msgId := msgOpRet.GetMsgId()
	record, _ := Globals.scyllaCli.FindUserOpExact(db.ComputePk(session.UserID), session.UserID, msgId)
	if record == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "not find record of invitation", nil, session)
		return
	}

	group, _ := findGroup(record.Gid)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group id in record  is wrong", nil, session)
		return
	}

	reqMem := msgOpRet.GetReqMem()
	// 这里不能为空，这里是邀请人
	if reqMem == nil || reqMem.UserId != record.Uid2 {

	}

	result := strings.ToLower(msgOpRet.GetResult())
	ret := model.UserOpResultRefuse
	if result == "accept" {
		ret = model.UserOpResultOk
	}

	// 记录到数据库中
	Globals.scyllaCli.SetUserOpResult(db.ComputePk(reqMem.UserId), db.ComputePk(session.UserID),
		reqMem.UserId, session.UserID, msgId, ret)

	// 执行加入的各种操作
	onJoinGroupOk(user, group, msg, session, "invitation")
}

// 公开群加入申请
// todo: 应该加入存储个人的好友操作记录表中，否则申请的个人无法查询结果也无法同步到多终端
func handleGroupJoinReq(msg *pbmodel.Msg, session *Session) {

	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}
	// 调试
	group.DebugPrint(Globals.Logger)

	_, bHas := group.HasMember(session.UserID)
	if bHas {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "already joined the group", nil, session)
		return
	}

	// params["visibility"] = "public" | "private" 默认为public
	visibility := group.GetParamByKey("visibility")
	if strings.ToLower(visibility) == "private" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "group is private, must join the group by invitation", nil, session)
		return
	}

	// 检查这个群的加入类型：params["jointype"] = "any" | "admin"
	joinType := group.GetParamByKey("jointype")
	switch strings.ToLower(joinType) {
	case "invite": // 私群
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "one must join the group by invitation", nil, session)
		return
	case "question":
		onJoinGroupNeedQuestion(group, msg, session)
	case "admin": // 管理员审核
		onJoinGroupNeedAdmin(group, msg, session)
		break
	default: // ""  | "any"
		onJoinGroupOk(session.GetUser(), group, msg, session, "any")
	}

}

// 某个管理员对加入申请的应答
func handleGroupJoinAnswer(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOpRet().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	isAdmin := group.IsAdmin(session.UserID)
	if !isAdmin {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "group join answer should from an admin of group", nil, session)
		return
	}

	msgId := msg.GetPlainMsg().GetGroupOpRet().MsgId
	pk := db.ComputePk(groupInfo.GroupId)
	record, err := Globals.scyllaCli.FindGroupOpExact(pk, groupInfo.GroupId, msgId)
	if record == nil || err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "reference msgId is not correct", nil, session)
		return
	}

	ret := model.UserOpResultOk
	strRet := strings.ToLower(msg.GetPlainMsg().GetGroupOpRet().GetResult())
	if strRet == "refuse" {
		ret = model.UserOpResultRefuse
	} else if strRet == "accept" {
		ret = model.UserOpResultOk
	} else {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "result must be accept or refuse", nil, session)
		return
	}

	if record.Ret != 0 {
		if int(record.Ret) == ret {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNone), "request has been answered", map[string]string{
				"uid": strconv.FormatInt(record.Uid2, 10),
			}, session)
			return
		} else {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTStage), "request has been answered", map[string]string{
				"uid": strconv.FormatInt(record.Uid2, 10),
			}, session)
			return
		}
	}

	// 更新记录
	err = Globals.scyllaCli.SetGroupOpResult(db.ComputePk(group.GroupId), groupInfo.GroupId, msgId, session.UserID, ret == model.UserOpResultOk)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "update scylladb error", nil, session)
		return
	}

	// 这里大概率是直接从缓存加载的，除非申请后下线了，而且过去一周了
	user, _, _ := findUser(record.Uid1)

	if ret == model.UserOpResultOk {
		onJoinGroupOk(user, group, msg, session, "admin")
	} else {
		// 拒绝了用户的请求
		trySendMsgToUser(record.Uid1, msg)
	}
}

// 退群申请
func handleGroupMemberQuit(msg *pbmodel.Msg, session *Session) {

	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "you are the owner, transfer owner or dissolve group", nil, session)
		return
	}

	uid := session.UserID
	// 同时操作2个表
	Globals.scyllaCli.DeleteGroupMember(db.ComputePk(groupInfo.GroupId), db.ComputePk(uid), groupInfo.GroupId, uid)

	// 从成员表中删除
	// 标记用户的所属群
	Globals.redisCli.SetUserLeaveGroup(uid, groupInfo.GroupId)

	// 从分布表中删除
	index := int64(Globals.Config.Server.HostIndex)
	Globals.redisCli.RemoveActiveGroupMembersLua(groupInfo.GroupId, index, []int64{uid})

	// 内存中删除成员
	group.RemoveMember(uid)

	// 内存中，用户退出组，
	user := session.GetUser()
	if user != nil {
		user.SetLeaveGroup(groupInfo.GroupId)
	}

	// 通知所有的用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}

	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupQuit,
		group.GetGroupInfo(),
		reqMem,
		nil,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"notice", session)

	trySendMsgToUser(session.UserID, msgRet)
	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// debug only:
	group.DebugPrint(Globals.Logger)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}
}

// 设置某人为群管理员
func handleGroupSetSomeoneAsAdmin(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if !isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set admin operation, you are not owner", nil, session)
		return
	}

	members := msg.GetPlainMsg().GetGroupOp().GetMembers()
	for _, mem := range members {
		nick, ok := group.HasMember(mem.UserId)
		if !ok {
			continue
		}
		// 更新数据库
		err := Globals.scyllaCli.SetGroupMemberRole(db.ComputePk(groupInfo.GroupId), groupInfo.GroupId, mem.UserId, model.RoleGroupAdmin|model.RoleGroupMember)
		if err != nil {

		}
		// 更新redis
		data := model.GroupMemberStore{
			Uid:  mem.UserId,
			Role: model.RoleGroupAdmin,
			Nick: nick,
		}
		Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{data})

		// 更新内存
		group.AddAdmin(mem.UserId)

	}

	// 通知所有用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: 0,
			Params:  nil,
		}
	}
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupAddAdmin,
		group.GetGroupInfo(),
		reqMem,
		members,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"", session)

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// debug only:
	group.DebugPrint(Globals.Logger)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}
}

// 删除管理员权限
func handleGroupRemoveSomeoneFromAdmin(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if !isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group delete admin operation, you are not owner", nil, session)
		return
	}

	members := msg.GetPlainMsg().GetGroupOp().GetMembers()
	for _, mem := range members {
		nick, ok := group.HasMember(mem.UserId)
		if !ok {
			continue
		}
		// 更新数据库
		err := Globals.scyllaCli.SetGroupMemberRole(db.ComputePk(groupInfo.GroupId), groupInfo.GroupId, mem.UserId, model.RoleGroupMember)
		if err != nil {

		}
		// 更新redis
		data := model.GroupMemberStore{
			Uid:  mem.UserId,
			Role: model.RoleGroupMember,
			Nick: nick,
		}
		Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{data})

		// 更新内存
		group.RemoveAdmin(mem.UserId)

	}

	// 通知所有用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    "",
			Icon:    "",
			Role:    "",
			GroupId: 0,
			Params:  nil,
		}
	}
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupDelAdmin,
		group.GetGroupInfo(),
		reqMem,
		members,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"", session)

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// debug only:
	group.DebugPrint(Globals.Logger)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}

}

// 转让群主
func handleGroupTransferOwner(msg *pbmodel.Msg, session *Session) {

	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group transfer owner operation, group info is null", nil, session)
		return
	}

	group, err := findGroup(groupInfo.GroupId)
	if err != nil || group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "do not find a group with id", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if !isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "you are not the owner", nil, session)
		return
	}

	memList := msgOp.GetMembers()
	if memList == nil || len(memList) < 1 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must give out the new owner in members", nil, session)
		return
	}

	mem := memList[0]
	// 检查member是否是群成员
	nick, bHas := group.HasMember(mem.UserId)
	if !bHas {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "new owner is not a member of the group", nil, session)
		return
	}

	err = Globals.scyllaCli.SetGroupMemberRole(db.ComputePk(groupInfo.GroupId), group.GroupId, mem.UserId, model.RoleGroupOwner|model.RoleGroupMember)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "set new owner, db error", nil, session)
		return
	}

	// 更新redis
	memStore := model.GroupMemberStore{
		Pk:   0,
		Role: int16(model.RoleGroupOwner | model.RoleGroupMember),
		Gid:  groupInfo.GroupId,
		Uid:  mem.UserId,
		Tm:   utils.GetTimeStamp(),
		Nick: nick,
	}
	err = Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{memStore})
	if err != nil {
		Globals.Logger.Fatal("SetGroupMembers() redis error", zap.Error(err))
	}

	// 更新群主
	group.SetOwner(mem.UserId, nick)

	// 通知所有人
	reqMem := msgOp.GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    nick,
			Icon:    "",
			Role:    "",
			GroupId: 0,
			Params:  nil,
		}
	}
	// 生成回复消息
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupTransferOwner,
		groupInfo,
		reqMem,
		memList,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"change owner", session)

	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// debug only:
	group.DebugPrint(Globals.Logger)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}
}

// 设置自己的在群内的信息
func handleSetMemberInfo(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	params := msgOp.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must have a nick in msg params", nil, session)
		return
	}

	nick, ok := params["nick"]
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must have a nick in msg params", nil, session)
		return
	}

	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, err := findGroup(groupInfo.GroupId)
	if err != nil || group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "do not find a group with id", nil, session)
		return
	}

	oldNick, bHas := group.HasMember(session.UserID)
	if !bHas {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group id in not correct, you are not a member of it", nil, session)
		return
	}

	reqMem := msgOp.GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    nick,
			Icon:    "",
			Role:    "",
			GroupId: 0,
			Params:  nil,
		}
	}

	// 不用改
	if oldNick == nick {
		msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSetInfo,
			groupInfo,
			reqMem,
			nil,
			msgOp.SendId,
			msgOp.MsgId,
			"ok",
			"update nick", session)
		session.SendMessage(msgRet)
		return
	}

	err = Globals.scyllaCli.SetGroupMemberNick(db.ComputePk(groupInfo.GroupId), groupInfo.GroupId, session.UserID, nick)
	if err != nil {
		Globals.Logger.Error("", zap.Error(err))
	}

	role, _ := group.GetMemberRole(session.UserID)
	memStore := model.GroupMemberStore{
		Pk:   0,
		Role: int16(role),
		Gid:  groupInfo.GroupId,
		Uid:  session.UserID,
		Tm:   utils.GetTimeStamp(),
		Nick: nick,
	}
	err = Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{memStore})
	if err != nil {
		Globals.Logger.Error("", zap.Error(err))
	}

	group.SetMemberNick(session.UserID, nick)

	// 生成回复消息
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSetInfo,
		groupInfo,
		reqMem,
		nil,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"update nick", session)

	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// debug only:
	group.DebugPrint(Globals.Logger)

	// 同步到其他的主机
	if Globals.Config.Server.ClusterMode {

	}

}

// 搜群 params["keyword"]
func handleGroupSearch(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	params := msgOp.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must have a keyword in msg params", nil, session)
		return
	}

	keyword, ok := params["keyword"]
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must have a keyword in msg params", nil, session)
		return
	}

	var msgRet *pbmodel.Msg = nil
	id, err := strconv.ParseInt(keyword, 10, 64)
	if err == nil {
		group, _ := findGroup(id)
		if group == nil {
			msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
				nil,
				nil,
				nil,
				msgOp.SendId,
				msgOp.MsgId,
				"fail",
				"not find id", session)

		} else {
			// 普通用户不让查询私有群
			if group.IsPrivate() && !session.GetUser().IsSystemUser() {
				msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
					nil,
					nil,
					nil,
					msgOp.SendId,
					msgOp.MsgId,
					"fail",
					"not find id", session)
			} else {
				msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
					nil,
					nil,
					nil,
					msgOp.SendId,
					msgOp.MsgId,
					"ok",
					"find it", session)
				msgRet.GetPlainMsg().GetGroupOpRet().Groups = []*pbmodel.GroupInfo{group.GetGroupInfo()}
			}
		}
	} else {
		// 通过关键字搜索
		bFilter := !session.GetUser().IsSystemUser()
		lst, err := Globals.mongoCli.FindGroupByKeyword(keyword, bFilter)
		if lst == nil || err != nil {
			msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
				nil,
				nil,
				nil,
				msgOp.SendId,
				msgOp.MsgId,
				"fail",
				"not find id", session)
		} else {

			msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
				nil,
				nil,
				nil,
				msgOp.SendId,
				msgOp.MsgId,
				"ok",
				"find it", session)
			msgRet.GetPlainMsg().GetGroupOpRet().Groups = lst
		}

	}

	// 最后通知所有的成员有新伙伴
	session.SendMessage(msgRet)

}

// 搜群内的成员
func handleGroupSearchMember(msg *pbmodel.Msg, session *Session) {

	msgOp := msg.GetPlainMsg().GetGroupOp()

	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroup(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	if _, bHas := group.HasMember(session.UserID); !bHas {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "you are not group member", nil, session)
		return
	}

	fromId := int64(0)
	params := msgOp.GetParams()
	if params != nil {
		keyword, ok := params["uid"]
		if !ok {
			uid, err := strconv.ParseInt(keyword, 10, 64)
			if err == nil {
				fromId = uid
			}
		}
	}

	pk := db.ComputePk(groupInfo.GroupId)

	memberStoreList, _ := Globals.scyllaCli.FindGroupMembers(pk, groupInfo.GroupId, fromId, 100)
	if memberStoreList == nil || len(memberStoreList) == 0 {

		msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupJoinAnswer,
			groupInfo,
			nil,
			nil,
			msgOp.SendId,
			msgOp.MsgId,
			"ok",
			"do not has more members", session)

		session.SendMessage(msgRet)
		return
	}

	members := make([]*pbmodel.GroupMember, len(memberStoreList))
	for index, mem := range memberStoreList {
		// 查找资源
		icon := ""
		user, _, _ := findUserInfo(mem.Uid)
		if user != nil {
			icon = user.Icon
		}

		members[index] = &pbmodel.GroupMember{
			UserId:  mem.Uid,
			Nick:    mem.Nick,
			Icon:    icon,
			Role:    mem.GetRoleStr(),
			GroupId: 0,
			Params:  nil,
		}
	}

	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupJoinAnswer,
		groupInfo,
		nil,
		members,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"find some", session)

	session.SendMessage(msgRet)

}

// 查询当前用户所在的各个群列表，同步这个主要为了新登录的终端，同步当前消息，不一定很有用；
// 因为在加入群的时候，正常的客户端都会了解
func handleListMemberInG(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()

	fromId := int64(0)
	params := msgOp.GetParams()
	if params != nil {
		str, ok := params["gid"]
		if ok {
			tempId, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				fromId = tempId
			}
		}
	}

	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupListIn,
		nil,
		nil,
		nil,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"find user in group", session)

	ginfoList, _ := LoadUserInGroupList(session.UserID, fromId)

	if ginfoList != nil && len(ginfoList) > 0 {
		msgRet.GetPlainMsg().GetGroupOpRet().Groups = ginfoList
	}

	session.SendMessage(msgRet)

}

// 保存新建立的群的基础信息
func saveNewGroup(groupInfo *pbmodel.GroupInfo, session *Session) (*model.Group, error) {

	// 保存基础信息
	err := Globals.mongoCli.CreateNewGroup(groupInfo)
	if err != nil {
		return nil, err
	}

	err = Globals.redisCli.SetGroupInfo(groupInfo)
	if err != nil {
		return nil, err
	}

	// 保存群主
	nick := ""
	user := session.GetUser()
	if user != nil {
		nick = user.GetNickName()
	}

	mem := model.GroupMemberStore{
		Pk:   db.ComputePk(groupInfo.GroupId),
		Gid:  groupInfo.GroupId,
		Uid:  session.UserID,
		Tm:   utils.GetTimeStamp(),
		Role: model.RoleGroupOwner,
		Nick: nick,
	}

	item := model.UserInGStore{
		Pk:  db.ComputePk(session.UserID),
		Uid: session.UserID,
		Gid: groupInfo.GroupId,
	}
	// 数据库
	err = Globals.scyllaCli.InsertGroupMember(&mem, &item)

	// redis中
	Globals.redisCli.SetUserJoinGroup(session.UserID, groupInfo.GroupId, "#|"+user.GetNickName())

	// 将群信息添加到内存
	g := model.NewGroupFromInfo(groupInfo)
	Globals.grc.InsertGroup(groupInfo.GroupId, g)
	// 设置群主
	g.SetOwner(user.UserId, user.GetNickName())

	// 设置用户在群中
	user.SetInGroup([]int64{groupInfo.GroupId})

	return g, err
}

func createGroupOpRetMsg(opCode pbmodel.GroupOperationType,
	groupInfo *pbmodel.GroupInfo,
	reqMem *pbmodel.GroupMember,
	members []*pbmodel.GroupMember,
	sendId, msgId int64,
	ret string, detail string, session *Session) *pbmodel.Msg {
	msgGroupOpRet := pbmodel.GroupOpResult{
		ReqMem:    reqMem,
		Operation: opCode,
		SendId:    sendId,
		MsgId:     msgId,
		Result:    ret,
		Detail:    detail,
		Group:     groupInfo,
		Members:   members,
		Params:    nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_GroupOpRet{
			GroupOpRet: &msgGroupOpRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTGroupOpRet,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	return &msg
}

// 加入一个群，成功了，
func onJoinGroupOk(user *model.User, group *model.Group, msg *pbmodel.Msg, session *Session, fromWays string) {
	// 保存成员信息
	nick := ""
	if user != nil {
		nick = user.GetNickName()
	}

	mem := model.GroupMemberStore{
		Pk:   db.ComputePk(group.GroupId),
		Gid:  group.GroupId,
		Uid:  session.UserID,
		Tm:   utils.GetTimeStamp(),
		Role: model.RoleGroupMember,
		Nick: nick,
	}

	item := model.UserInGStore{
		Pk:  db.ComputePk(session.UserID),
		Uid: session.UserID,
		Gid: group.GroupId,
	}
	// 数据库中加入成员，同时， 保存成员所属的群
	err := Globals.scyllaCli.InsertGroupMember(&mem, &item)
	if err != nil {
		Globals.Logger.Fatal("InsertGroupMember() err", zap.Error(err))
	}

	// redis的群所有成员
	err = Globals.redisCli.SetUserJoinGroup(session.UserID, group.GroupId, nick)
	if err != nil {
		Globals.Logger.Fatal("SetUserJoinGroup() err", zap.Error(err))
	}
	index := int64(Globals.Config.Server.HostIndex)
	err = Globals.redisCli.SetActiveGroupMembers(group.GroupId, index, []int64{session.UserID})
	if err != nil {
		Globals.Logger.Fatal("SetActiveGroupMembers() err", zap.Error(err))
	}

	// 更新内存的部分
	group.AddMember(user.UserId, nick)
	// 设置用户在群中
	user.SetInGroup([]int64{group.GroupId})

	addedMember := &pbmodel.GroupMember{
		UserId:  session.UserID,
		Nick:    nick,
		Icon:    user.Icon,
		Role:    "",
		GroupId: group.GroupId,
		Params:  nil,
	}
	// 最后通知所有的成员有新伙伴
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupJoinAnswer,
		group.GetGroupInfo(),
		nil,
		[]*pbmodel.GroupMember{
			addedMember,
		},
		msg.GetPlainMsg().GetGroupOp().SendId,
		Globals.snow.GenerateID(),
		"ok",
		fromWays, session)

	// 通知所有用户
	notifyGroupMembers(group.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}
}

// 加入一个群时候，需要回答问题，这里直接检查问题的回答是否正确
// params["joinquestion"]
// params["joinanswer"]
func onJoinGroupNeedQuestion(group *model.Group, msg *pbmodel.Msg, session *Session) error {
	question := group.GetParamByKey("joinquestion")
	answer := group.GetParamByKey("joinanswer")

	if question == "" || answer == "" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "group question or answer is null", nil, session)
		return errors.New("no answer in msg")
	}

	params := msg.GetPlainMsg().GetGroupOpRet().GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "should answer the question", map[string]string{
			"question": question,
		}, session)
		return errors.New("no answer in msg")
	}
	answerOfUser, ok := params["answer"]
	if !ok || answerOfUser == "" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "should answer the question", map[string]string{
			"question": question,
		}, session)
		return errors.New("no answer in msg")
	}

	if answerOfUser != answer {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "answer is not correct", map[string]string{
			"question": question,
		}, session)
		return errors.New("no answer in msg")
	}

	onJoinGroupOk(session.GetUser(), group, msg, session, "question")

	return nil
}

// 加入一个群，需要管理员审核, 那么需要发消息给所有的的管理员
func onJoinGroupNeedAdmin(group *model.Group, msg *pbmodel.Msg, session *Session) error {
	// 先把这个消息保存到群的操作记录中
	id := Globals.snow.GenerateID()
	pk := db.ComputePk(group.GroupId)
	record := model.CommonOpStore{
		Pk:   pk,
		Gid:  group.GroupId,
		Uid1: session.UserID,
		Uid2: 0,
		Id:   id,
		Usid: 0,
		Tm:   time.Now().UnixMilli(),
		Tm1:  0,
		Tm2:  0,
		Io:   0,
		St:   0,
		Cmd:  model.CommonGroupOpJoinRequest,
		Ret:  0,
		Mask: 0,
		Ref:  0,
		Draf: nil,
	}
	err := Globals.scyllaCli.SaveGroupOp(&record)
	if err != nil {
		return err
	}

	adminList := group.GetAdminMembers()
	if adminList == nil || len(adminList) == 0 {
		return errors.New("admins and owner is nil")
	}

	// 尝试向左右在线的用户发送消息
	if msg.GetPlainMsg().GetGroupOpRet().GetReqMem() == nil {
		user := session.GetUser()
		if user == nil {
			Globals.Logger.Fatal("onJoinGroupNeedAdmin() get user from session meet error")
			return errors.New("can't find user in cache")
		}

		msg.GetPlainMsg().GetGroupOpRet().ReqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    user.GetNickName(),
			Icon:    user.GetIcon(),
			Role:    "",
			GroupId: group.GroupId,
			Params:  nil,
		}
	}
	trySendMsgToUserList(adminList, msg)

	// 通知原用户等待
	// 最后通知所有的成员有新伙伴
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupJoinAnswer,
		group.GetGroupInfo(),
		nil,
		nil,
		msg.GetPlainMsg().GetGroupOpRet().SendId,
		Globals.snow.GenerateID(),
		"wait",
		"", session)

	trySendMsgToUser(session.UserID, msgRet)

	return nil
}
