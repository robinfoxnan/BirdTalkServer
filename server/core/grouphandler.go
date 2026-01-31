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

	// 为每一条收到的消息重新定义流水号
	groupOpMsg.MsgId = Globals.snow.GenerateID()

	opCode := groupOpMsg.Operation
	switch opCode {
	case pbmodel.GroupOperationType_GroupCreate: // 1创建
		handleGroupCreateOp(msg, session)
		break
	case pbmodel.GroupOperationType_GroupDissolve: // 2解散
		handleGroupDissolveOp(msg, session)
		break
	case pbmodel.GroupOperationType_GroupSetInfo: // 3设置信息
		handleGroupSetBasicInfo(msg, session)
		break
	case pbmodel.GroupOperationType_GroupKickMember: // 4踢人
		handleGroupKickOut(msg, session)
		break
	case pbmodel.GroupOperationType_GroupInviteRequest: // 5邀请
		handleInviteSomeoneDirect(msg, session)
		break
	case pbmodel.GroupOperationType_GroupInviteAnswer: // 6邀请的应答
		handleInviteAnswer(msg, session)
		break // 邀请后处理结果
	case pbmodel.GroupOperationType_GroupJoinRequest: // 7加入请求
		handleGroupJoinReq(msg, session)
		break
	case pbmodel.GroupOperationType_GroupJoinAnswer: // 8加入请求的处理，同意、拒绝、问题
		handleGroupJoinAnswer(msg, session)
		break
	case pbmodel.GroupOperationType_GroupQuit: // 9退出群组
		handleGroupMemberQuit(msg, session)
		break
	case pbmodel.GroupOperationType_GroupAddAdmin: // 10增加管理员
		handleGroupSetSomeoneAsAdmin(msg, session)
		break
	case pbmodel.GroupOperationType_GroupDelAdmin: // 11删除管理员
		handleGroupRemoveSomeoneFromAdmin(msg, session)
		break
	case pbmodel.GroupOperationType_GroupTransferOwner: // 12转让群主
		handleGroupTransferOwner(msg, session)
		break
	case pbmodel.GroupOperationType_GroupSetMemberInfo: // 13设置自己在群中的信息
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

// 1. 创建群
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
			memUser, b, _ := findUser(mem.UserId)
			if memUser == nil || b == false {
				continue
			}
			saveUserJoinGroup(memUser, group)
		}
	}

	// 通知相关用户，建立了新群
	retMsg := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupCreate,
		groupInfo,
		groupOpMsg.GetReqMem(),
		groupOpMsg.GetMembers(),
		groupOpMsg.SendId,
		groupOpMsg.MsgId, // 这里使用刚才生成的号码应答了
		"ok", "", session)

	Globals.Logger.Debug("return create group ret", zap.Any("msg", retMsg))
	//fmt.Println(retMsg)

	notifyGroupMembers(groupInfo.GroupId, retMsg)

	// 保存记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, groupOpMsg.MsgId, groupOpMsg.SendId,
		pbmodel.GroupOperationType_GroupCreate, nil)

	return
}

// 2. 解散群
func handleGroupDissolveOp(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation group info id is wrong", nil, session)
		return
	}
	if !group.IsOwner(session.UserID) {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group dissolve operation, you are not owner", nil, session)
		return
	}

	// 解散群： 在基础信息中设置标记
	param := make(map[string]interface{})
	param["params.visibility"] = "private" // 这个标记是防止再被搜到
	param["params.status"] = "deleted"
	_, err := Globals.mongoCli.UpdateGroupInfoPart(groupInfo.GroupId, param, nil)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save db error", nil, session)
		return
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

	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupDissolve, groupInfo,
		reqMem,
		nil,
		groupOpMsg.SendId,
		groupOpMsg.MsgId,
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

	// 保存记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, groupOpMsg.MsgId, groupOpMsg.SendId, pbmodel.GroupOperationType_GroupDissolve, nil)

	return
}

// 3. 设置基础信息，这个也是只有管理员才可以设置，不同于微信的
func handleGroupSetBasicInfo(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set info operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
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

	Globals.Logger.Debug("group set info ok", zap.Any("msg", msgRet))

	// 通知所有用户
	notifyGroupMembers(groupInfo.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	if Globals.Config.Server.ClusterMode {

	}

	// 保存记录
	txt := group.GetGroupInfo().String()
	utf8Bytes := []byte(txt)
	saveGroupOpRecord(group.GroupId, session.UserID, 0, msgOp.MsgId, msgOp.SendId,
		pbmodel.GroupOperationType_GroupSetInfo, utf8Bytes)

}

// 4, 踢人
func handleGroupKickOut(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation group info id is wrong", nil, session)
		return
	}

	isAdmin := group.IsAdmin(session.UserID)
	if !isAdmin {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group kick operation, you are not admin", nil, session)
		return
	}

	msgOp := msg.GetPlainMsg().GetGroupOp()
	memList := msgOp.GetMembers()
	if memList == nil || len(memList) < 1 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "must give out the users in members", nil, session)
		return
	}

	// 逐个删除
	for _, mem := range memList {
		uid := mem.UserId
		memUser, _, _ := findUser(uid)
		if memUser == nil {
			continue
		}

		// 同时操作2个表
		Globals.scyllaCli.DeleteGroupMember(db.ComputePk(groupInfo.GroupId), db.ComputePk(uid), groupInfo.GroupId, uid)

		// 从成员表中删除
		// 标记用户的所属群
		Globals.redisCli.SetUserLeaveGroup(uid, groupInfo.GroupId)

		// 从redis分布表中删除
		index := int64(Globals.Config.Server.HostIndex)
		Globals.redisCli.RemoveActiveGroupMembersLua(groupInfo.GroupId, index, []int64{uid})

		// 内存中删除成员
		group.RemoveMember(uid)

		// 内存中，用户退出组，
		memUser.SetLeaveGroup(groupInfo.GroupId)

		// 保存踢人记录
		saveGroupOpRecord(group.GroupId, session.UserID, uid, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupKickMember, []byte("kick ok"))

	}

	// 通知所有的用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    session.GetUser().NickName,
			Icon:    session.GetUser().Icon,
			Role:    "au",
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}

	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupKickMember,
		group.GetGroupInfo(),
		reqMem,
		memList,
		msgOp.SendId,
		msgOp.MsgId,
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

// 5. 邀请某人
// 2026-1-30 这里简化一下，直接拉人进去，和微信一样
func handleInviteSomeoneDirect(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group invite request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group invite request group info id is wrong", nil, session)
		return
	}

	// 私密群只有管理员可以拉人
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

	// 这里是填充邀请者的信息
	user := session.GetUser()
	role := group.GetMemberRoleString(user.UserId)
	reqMem := msgOp.GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    user.NickName,
			Icon:    user.Icon,
			Role:    role,
			GroupId: groupInfo.GroupId,
			Params:  nil,
		}
	}
	// 重新设置，为了转发
	if msgOp.ReqMem == nil {

		msgOp.ReqMem = reqMem
	}

	// 核心：先判断Params是否为nil，为nil则初始化构造map
	if msgOp.Params == nil {
		msgOp.Params = make(map[string]string) // 初始化空map（分配底层内存）
	}

	// 针对申请中的每个人，逐个加入到群中
	for _, mem := range memList {
		// 用户不存在就不能继续操作，否则后续的用户注册进来会造成脏数据
		memUser, _, _ := findUser(mem.UserId)
		if memUser == nil {
			continue
		}

		onJoinGroupOk(memUser, group, msg, session, "invite ok")

		// 保存邀请到个人记录
		//saveGroupOpUserOpRecord(group.GroupId, session.UserID, memUser.UserId, msgOp.MsgId, msgOp.SendId,
		//	pbmodel.GroupOperationType_GroupInviteRequest, nil)

		// 保存到群组记录
		saveGroupOpRecord(group.GroupId, session.UserID, memUser.UserId, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupInviteRequest, []byte("invite ok"))
	}

	// 答复发出邀请的用户
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupInviteRequest,
		group.GetGroupInfo(),
		reqMem,
		memList,
		msgOp.SendId,
		msgOp.MsgId,
		"invite ok",
		"group invitation", session)
	trySendMsgToUser(session.UserID, msgRet)

}

func handleInviteSomeone(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
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
	// 重新设置，为了转发
	if msgOp.ReqMem == nil {
		msgOp.ReqMem = reqMem
	}

	code := utils.GenerateCheckCode(6)
	//code32, _ := strconv.Atoi(code)
	utf8Bytes := []byte(code)

	// 核心：先判断Params是否为nil，为nil则初始化构造map
	if msgOp.Params == nil {
		msgOp.Params = make(map[string]string) // 初始化空map（分配底层内存）
	}
	msgOp.Params["code"] = code

	for _, mem := range memList {
		// 用户不存在就不能继续操作，否则后续的用户注册进来会造成脏数据
		memUser, _, _ := findUserInfo(mem.UserId)
		if memUser == nil {
			continue
		}
		// 生成验证码

		// 保存邀请到个人记录
		saveGroupOpUserOpRecord(group.GroupId, session.UserID, memUser.UserId, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupInviteRequest, utf8Bytes)

		// 保存到群组记录
		saveGroupOpRecord(group.GroupId, session.UserID, memUser.UserId, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupInviteRequest, utf8Bytes)

		// 通知被邀请人，这里使用申请操作
		//msgNotice := createGroupOpMsg(pbmodel.GroupOperationType_GroupInviteRequest,
		//	group.GetGroupInfo(),
		//	reqMem,
		//	[]*pbmodel.GroupMember{
		//		mem,
		//	},
		//	msgOp.SendId,
		//	msgOp.MsgId,
		//	"notify",
		//	"group invitation", session)

		trySendMsgToUser(mem.UserId, msg)

		// 恢复发出邀请的用户
		msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupInviteRequest,
			group.GetGroupInfo(),
			reqMem,
			[]*pbmodel.GroupMember{mem},
			msgOp.SendId,
			msgOp.MsgId,
			"wait",
			"group invitation", session)
		trySendMsgToUser(session.UserID, msgRet)
	}

}

// 6. 邀请的回答，这里的应答，客户端是使用应答发送的，因为这里才有
func handleInviteAnswer(msg *pbmodel.Msg, session *Session) {
	msgOpRet := msg.GetPlainMsg().GetGroupOpRet()

	user := session.GetUser()
	refMsgId := "0"
	if msgOpRet.Params != nil {
		str, ok := msgOpRet.Params["refid"]
		if ok {
			refMsgId = str
		}
	}

	refMsgId64, _ := strconv.ParseInt(refMsgId, 10, 64)

	// 必须要引用了之前的记录才能算数
	record, _ := Globals.scyllaCli.FindUserOpExact(db.ComputePk(session.UserID), session.UserID, refMsgId64)
	if record == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "not find record of invitation", nil, session)
		return
	}

	// record.uid2是邀请人
	// record.gid是群ID

	group, _ := findGroupAndLoad(record.Gid)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group id in record  is wrong", nil, session)
		return
	}

	result := strings.ToLower(msgOpRet.GetResult())
	ok := false
	ret := model.UserOpResultRefuse
	if result == "accept" {
		ret = model.UserOpResultOk
		ok = true
	}

	// 保存邀请到个人记录
	uid2 := record.Uid2
	Globals.scyllaCli.SetUserOpResult(db.ComputePk(uid2), db.ComputePk(session.UserID),
		uid2, session.UserID, refMsgId64, ret)

	// 保存到群组的记录中
	Globals.scyllaCli.SetGroupOpResult(db.ComputePk(record.Gid), record.Gid, refMsgId64, session.UserID, ok)

	// 执行加入的各种操作
	onJoinGroupOk(user, group, msg, session, "invitation")

}

// 私有群，或者设置为auth类型的群，需要邀请才能进入
// todo: 应该加入存储个人的好友操作记录表中，否则申请的个人无法查询结果也无法同步到多终端
func handleGroupJoinReq(msg *pbmodel.Msg, session *Session) {

	groupInfo := msg.GetPlainMsg().GetGroupOp().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}
	// 调试
	group.DebugPrint(Globals.Logger)

	_, _, bHas := group.HasMember(session.UserID)
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
	case "auth": // 管理员审核
		onJoinGroupNeedAdmin(group, msg, session)
		break
	default: // ""  | "any"
		onJoinGroupOk(session.GetUser(), group, msg, session, "direct")
	}

}

// 某个管理员对加入申请的应答
func handleGroupJoinAnswer(msg *pbmodel.Msg, session *Session) {
	groupInfo := msg.GetPlainMsg().GetGroupOpRet().GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
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

// 9. 退群申请
func handleGroupMemberQuit(msg *pbmodel.Msg, session *Session) {

	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group quit request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group quit request group info id is wrong", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "you are the owner, transfer owner or dissolve group", nil, session)
		return
	}

	uid := session.UserID
	Globals.Logger.Debug("user quit group", zap.String("user", session.GetUser().NickName), zap.Int64("uid", uid),
		zap.Int64("groupId", groupInfo.GroupId))
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
			Nick:    session.GetUser().NickName,
			Icon:    session.GetUser().Icon,
			Role:    "u",
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

	// 保存记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, msgOp.MsgId, msgOp.SendId,
		pbmodel.GroupOperationType_GroupQuit, nil)
}

// 10 设置某人为群管理员
func handleGroupSetSomeoneAsAdmin(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group add admins request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group add admins request group info id is wrong", nil, session)
		return
	}

	isOwner := group.IsOwner(session.UserID)
	if !isOwner {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group set admin operation, you are not owner", nil, session)
		return
	}

	members := msg.GetPlainMsg().GetGroupOp().GetMembers()
	for _, mem := range members {
		nick, _, ok := group.HasMember(mem.UserId)
		if !ok {
			continue
		}
		// 更新数据库
		err := Globals.scyllaCli.SetGroupMemberRole(db.ComputePk(groupInfo.GroupId), groupInfo.GroupId, mem.UserId,
			model.RoleGroupAdmin|model.RoleGroupMember)
		if err != nil {

		}
		// 更新redis
		data := model.GroupMemberStore{
			Uid:  mem.UserId,
			Role: model.RoleGroupAdmin | model.RoleGroupMember,
			Nick: nick,
		}
		Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{data})

		// 更新内存
		group.AddAdmin(mem.UserId)
		// 保存记录
		saveGroupOpRecord(group.GroupId, session.UserID, mem.UserId, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupAddAdmin, nil)
	}

	// 通知所有用户
	reqMem := msg.GetPlainMsg().GetGroupOp().GetReqMem()
	if reqMem == nil {
		reqMem = &pbmodel.GroupMember{
			UserId:  session.UserID,
			Nick:    session.GetUser().NickName,
			Icon:    session.GetUser().Icon,
			Role:    "ou",
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

// 11. 删除管理员权限
func handleGroupRemoveSomeoneFromAdmin(msg *pbmodel.Msg, session *Session) {
	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request operation group info is null", nil, session)
		return
	}

	group, _ := findGroupAndLoad(groupInfo.GroupId)
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
		nick, _, ok := group.HasMember(mem.UserId)
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

		// 保存记录
		saveGroupOpRecord(group.GroupId, session.UserID, mem.UserId, msgOp.MsgId, msgOp.SendId,
			pbmodel.GroupOperationType_GroupDelAdmin, nil)

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

// 12. 转让群主
func handleGroupTransferOwner(msg *pbmodel.Msg, session *Session) {

	msgOp := msg.GetPlainMsg().GetGroupOp()
	groupInfo := msgOp.GetGroup()
	if groupInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group transfer owner operation, group info is null", nil, session)
		return
	}

	group, err := findGroupAndLoad(groupInfo.GroupId)
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
	nick, _, bHas := group.HasMember(mem.UserId)
	if !bHas {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "new owner is not a member of the group", nil, session)
		return
	}

	err = Globals.scyllaCli.SetGroupMemberRole(db.ComputePk(groupInfo.GroupId), group.GroupId, mem.UserId,
		model.RoleGroupOwner|model.RoleGroupMember|model.RoleGroupAdmin)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "set new owner, db error", nil, session)
		return
	}

	// 更新redis
	memStore := model.GroupMemberStore{
		Pk:   0,
		Role: int16(model.RoleGroupOwner | model.RoleGroupMember | model.RoleGroupAdmin),
		Gid:  groupInfo.GroupId,
		Uid:  mem.UserId,
		Tm:   utils.GetTimeStamp(),
		Nick: nick,
	}
	err = Globals.redisCli.SetGroupMembers(groupInfo.GroupId, []model.GroupMemberStore{memStore})
	if err != nil {
		Globals.Logger.Fatal("SetGroupMembers() redis error", zap.Error(err))
	}

	ownerNew, _, err := findUser(mem.UserId)
	if err != nil {
		Globals.Logger.Fatal("handleGroupTransferOwner()->find user() return nil ", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "set new owner, can't find user error", nil, session)
		return
	}
	// 更新群主
	group.SetOwner(mem.UserId, nick, ownerNew)

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

	// 保存记录
	saveGroupOpRecord(group.GroupId, session.UserID, mem.UserId, msgOp.MsgId, msgOp.SendId,
		pbmodel.GroupOperationType_GroupTransferOwner, nil)
}

// 13. 设置自己的在群内的信息
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

	group, err := findGroupAndLoad(groupInfo.GroupId)
	if err != nil || group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "do not find a group with id", nil, session)
		return
	}

	oldNick, _, bHas := group.HasMember(session.UserID)
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

	// 保存记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, msgOp.MsgId, msgOp.SendId,
		pbmodel.GroupOperationType_GroupSetMemberInfo, []byte(nick))

}

// 搜群 params["keyword"]
// 首先转为gid尝试搜，如果关键字不是数字，使用name和tag来搜索，过滤掉 private类型
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

	gid := int64(0)
	gidStr, ok := params["gid"]
	if ok {
		gid, _ = strconv.ParseInt(gidStr, 10, 64)
	}
	/////////////////////////////////////////////////////////////
	// 这一段是按照GID来查找，
	var msgRet *pbmodel.Msg = nil
	gid, err := strconv.ParseInt(keyword, 10, 64)
	if err == nil {
		group, _ := findGroupAndLoad(gid)
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
			// 公开的群，系统用户，或者群成员可以查
			_, _, isMem := group.HasMember(session.UserID)
			if !group.IsPrivate() || session.GetUser().IsSystemUser() || isMem {
				msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
					nil,
					nil,
					nil,
					msgOp.SendId,
					msgOp.MsgId,
					"ok",
					"find it", session)
				msgRet.GetPlainMsg().GetGroupOpRet().Groups = []*pbmodel.GroupInfo{group.GetGroupInfo()}

			} else {
				// 普通用户不让查询私有群,
				msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
					nil,
					nil,
					nil,
					msgOp.SendId,
					msgOp.MsgId,
					"fail",
					"not find id", session)
			}
		}
	} else {
		// 通过关键字搜索
		bFilter := !session.GetUser().IsSystemUser()
		lst, err := Globals.mongoCli.FindGroupByKeyword(keyword, gid, bFilter, int64(Globals.maxPageSize))
		if lst == nil || err != nil {
			msgRet = createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearch,
				nil,
				nil,
				nil,
				msgOp.SendId,
				msgOp.MsgId,
				"fail",
				"not find by keyword", session)
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

	// 返回查询结果
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

	group, _ := findGroupAndLoad(groupInfo.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group join request group info id is wrong", nil, session)
		return
	}

	if _, _, bHas := group.HasMember(session.UserID); !bHas {
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

	members := findGroupMembersFromId(group, fromId, int(Globals.maxPageSize))
	// 某个UID之后就没有了，

	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupSearchMember,
		groupInfo,
		nil,
		members,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"find some", session)

	Globals.Logger.Debug("handleGroupSearchMember", zap.Any("msgRet", msgRet))

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

	// 这里是通用函数，这里没有提供参数，所以在下面填写
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupListIn,
		nil,
		nil,
		nil,
		msgOp.SendId,
		msgOp.MsgId,
		"ok",
		"list groups those user joins in", session)

	ginfoList, _ := LoadUserInGroupList(session.UserID, fromId)

	if ginfoList != nil && len(ginfoList) > 0 {
		msgRet.GetPlainMsg().GetGroupOpRet().Groups = ginfoList
	}

	session.SendMessage(msgRet)

}

// 保存新建立的群的基础信息
func saveNewGroup(groupInfo *pbmodel.GroupInfo, session *Session) (*model.Group, error) {

	// 1) 保存群的基础信息
	// 1.1
	err := Globals.mongoCli.CreateNewGroup(groupInfo)
	if err != nil {
		return nil, err
	}

	/////////////////////////////////////////////////////////////////
	// 2) 保存群成员，2级存储
	user := session.GetUser()
	nick := ""
	if user != nil {
		nick = user.GetNickName()
	}

	mem := model.GroupMemberStore{
		Pk:   db.ComputePk(groupInfo.GroupId),
		Gid:  groupInfo.GroupId,
		Uid:  user.UserId,
		Tm:   utils.GetTimeStamp(),
		Role: model.RoleGroupOwner | model.RoleGroupAdmin | model.RoleGroupMember,
		Nick: nick,
	}

	item := model.UserInGStore{
		Pk:  db.ComputePk(session.UserID),
		Uid: user.UserId,
		Gid: groupInfo.GroupId,
	}
	// 2.1数据库
	err = Globals.scyllaCli.InsertGroupMember(&mem, &item)
	if err != nil {
		return nil, err
	}

	//1.2 群组信息保存到redis
	Globals.redisCli.SetGroupInfo(groupInfo)

	// 2.2redis中
	Globals.redisCli.SetUserJoinGroup(session.UserID, groupInfo.GroupId, model.RoleGroupOwner, user.GetNickName())

	// 1.3 & 2.3 将群信息添加到内存
	g := model.NewGroupFromInfo(groupInfo)
	g.SetOwner(user.UserId, user.GetNickName(), user) // 设置群主
	Globals.grc.InsertGroup(groupInfo.GroupId, g)

	// 3) 设置用户在群中
	// 3.1 内存
	user.SetJoinGroup(groupInfo.GroupId)
	// 3.2 redis
	// 在2.2 函数中已经做了

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

// 加入一个群，成功了，这里是设置了user被加入了群
func onJoinGroupOk(user *model.User, group *model.Group, msg *pbmodel.Msg, session *Session, fromWays string) {
	// 保存成员信息
	msgOp := msg.GetPlainMsg().GetGroupOp()
	nick := ""
	if user != nil {
		nick = user.GetNickName()
	}

	// 2026-01-31 修正
	mem := model.GroupMemberStore{
		Pk:   db.ComputePk(group.GroupId),
		Gid:  group.GroupId,
		Uid:  user.UserId,
		Tm:   utils.GetTimeStamp(),
		Role: model.RoleGroupMember,
		Nick: nick,
	}

	item := model.UserInGStore{
		Pk:  db.ComputePk(user.UserId),
		Uid: user.UserId,
		Gid: group.GroupId,
	}
	// 数据库中加入成员，同时， 保存成员所属的群
	err := Globals.scyllaCli.InsertGroupMember(&mem, &item)
	if err != nil {
		Globals.Logger.Fatal("InsertGroupMember() err", zap.Error(err))
	}

	// redis的群所有成员
	err = Globals.redisCli.SetUserJoinGroup(user.UserId, group.GroupId, model.RoleGroupMember, nick)

	if err != nil {
		Globals.Logger.Fatal("SetUserJoinGroup() err", zap.Error(err))
	}
	index := int64(Globals.Config.Server.HostIndex)
	err = Globals.redisCli.SetActiveGroupMembers(group.GroupId, index, []int64{session.UserID})
	if err != nil {
		Globals.Logger.Fatal("SetActiveGroupMembers() err", zap.Error(err))
	}

	// 更新内存的部分
	group.AddMember(user.UserId, nick, user)
	// 设置用户在群中
	user.SetInGroup([]int64{group.GroupId})

	addedMember := &pbmodel.GroupMember{
		UserId:  user.UserId,
		Nick:    nick,
		Icon:    user.Icon,
		Role:    "u",
		GroupId: group.GroupId,
		Params:  nil,
	}
	// 最后通知所有的成员有新伙伴
	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	msgRet := createGroupOpRetMsg(pbmodel.GroupOperationType_GroupJoinAnswer,
		group.GetGroupInfo(),
		nil,
		[]*pbmodel.GroupMember{
			addedMember,
		},
		groupOpMsg.SendId,
		groupOpMsg.MsgId,
		"ok",
		fromWays, session)

	// 通知所有用户
	notifyGroupMembers(group.GroupId, msgRet)

	// todo:
	// 如果是集群模式，通知其他的服务器同步内存中的信息
	// 通知该用户所在的机器更改user
	if Globals.Config.Server.ClusterMode {

	}
	// 保存用户加入群记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, msgOp.MsgId, msgOp.SendId,
		pbmodel.GroupOperationType_GroupJoinRequest, []byte(fromWays))

}

// 加入一个群时候，需要回答问题，这里直接检查问题的回答是否正确
// params["joinquestion"]
// params["joinanswer"]
func onJoinGroupNeedQuestion(group *model.Group, msg *pbmodel.Msg, session *Session) error {
	question := group.GetParamByKey("joinquestion")
	answer := group.GetParamByKey("joinanswer")

	if answer == "" {
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
// "auth"模型，包括隐私群
func onJoinGroupNeedAdmin(group *model.Group, msg *pbmodel.Msg, session *Session) error {
	// 先把这个消息保存到群的操作记录中

	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	info := ""
	if groupOpMsg.Params != nil {
		// 双值查询+类型断言：一步实现“存在且为字符串则赋值，否则保留空串”
		info = groupOpMsg.Params["info"]
	}

	// 保存到群操作请求记录
	saveGroupOpRecord(group.GroupId, session.UserID, 0, groupOpMsg.MsgId,
		groupOpMsg.SendId, model.CommonGroupOpJoinRequest, []byte(info))

	adminList := group.GetAdminMembers()
	if adminList == nil || len(adminList) == 0 {
		return errors.New("admins and owner is nil")
	}

	// 写入个人请求记录
	for _, admin := range adminList {
		// 保存到个人记录中
		err := saveGroupOpUserOpRecord(group.GroupId, session.UserID, admin, groupOpMsg.MsgId,
			groupOpMsg.SendId, model.CommonGroupOpJoinRequest, []byte(info))

		if err != nil {
			Globals.Logger.Fatal("onJoinGroupNeedAdmin()-> saveGroupOpUserOpRecord() err", zap.Error(err))
			return errors.New("can't record to user op records")
		}
	}

	// 将消息转发给在线的管理员，尝试向左右在线的用户发送消息
	if groupOpMsg.GetReqMem() == nil {
		user := session.GetUser()
		if user == nil {
			Globals.Logger.Fatal("onJoinGroupNeedAdmin() get user from session meet error")
			return errors.New("can't find user in cache")
		}

		groupOpMsg.ReqMem = &pbmodel.GroupMember{
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
