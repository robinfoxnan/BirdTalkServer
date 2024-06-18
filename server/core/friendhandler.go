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

// 客户请求的好友的所有的操作
func handleFriendOp(msg *pbmodel.Msg, session *Session) {

	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	if friendOpMsg == nil {
		Globals.Logger.Debug("receive wrong friend op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "friend op  is null", nil, session)
		return
	}

	// 都需要验证是否登录与权限
	ok := checkUserLogin(session)
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
		return
	}

	opCode := friendOpMsg.Operation
	switch opCode {
	case pbmodel.UserOperationType_FindUser:
		handleFriendFind(msg, session)
	case pbmodel.UserOperationType_AddFriend: // 交友模式下需要转发给好友
		handleFriendAdd(msg, session)
	case pbmodel.UserOperationType_ApproveFriend: // 用户应答同意或者拒绝
		handleFriendApprove(msg, session)
	case pbmodel.UserOperationType_RemoveFriend:
		handleFriendRemove(msg, session)
	case pbmodel.UserOperationType_BlockFriend:
		handleFriendBlock(msg, session)
	case pbmodel.UserOperationType_UnBlockFriend:
		handleFriendUnBlock(msg, session)
	case pbmodel.UserOperationType_SetFriendPermission:
		handleFriendPermission(msg, session)
	case pbmodel.UserOperationType_SetFriendMemo:
		handleFriendSetMemo(msg, session)
	case pbmodel.UserOperationType_ListFriends:
		handleFriendList(msg, session)
	default:
		Globals.Logger.Info("receive unknown friend op",
			zap.Int64("sessionId", session.Sid),
			zap.Int64("uid", session.UserID),
			zap.String("opCode", opCode.String()),
		)
	}
}

// 用户回答的好友应答，需要转发
func handleFriendOpRet(msg *pbmodel.Msg, session *Session) {
	session.Ver = int(msg.GetVersion())   // 协议版本号
	ok := checkProtoVersion(msg, session) // 错误会自动应答
	if !ok {
		return
	}

	// 目前先不考虑这里
	keyPrint := msg.GetKeyPrint() // 加密传输的秘钥指纹，这里应该为0
	errCode, errStr := checkKeyPrint(keyPrint)
	if keyPrint != 0 { // 如果设置了秘钥，那么这里需要验证秘钥正确性
		sendBackErrorMsg(errCode, errStr, nil, session)
		return
	}

	UserOpMsg := msg.GetPlainMsg().GetUserOp()
	opCode := UserOpMsg.Operation
	switch opCode {
	case pbmodel.UserOperationType_ApproveFriend: // 这个是处理并转发
	default:
		Globals.Logger.Info("receive unknown friend result op",
			zap.Int64("sessionId", session.Sid),
			zap.Int64("uid", session.UserID),
			zap.String("opCode", opCode.String()),
		)
	}
}

// ////////////////////////////////////////////////////////////////
// 下面是细分的子类型
func handleFriendFind(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	params := friendOpMsg.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"FriendOpMsg  must hava Params field",
			nil,
			session)

		return
	}

	mode, ok1 := params["mode"]
	value, ok2 := params["value"]
	if !ok1 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"search friend mode field is nil",
			map[string]string{
				"field": "mode",
				"value": "nil",
			},
			session)

		return
	}

	if !ok2 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"search friend value field is nil",
			map[string]string{
				"field": "mode",
				"value": "nil",
			},
			session)

		return
	}
	var userList []*pbmodel.UserInfo
	var err error
	switch strings.ToLower(mode) {
	case "id":
		id, _ := strconv.ParseInt(value, 10, 64)
		userList, err = findUserMongoRedis(id)
	case "name":
		userList, err = Globals.mongoCli.FindUserByName(value)
	case "email":
		userList, err = Globals.mongoCli.FindUserByEmail(value)
	case "phone":
		userList, err = Globals.mongoCli.FindUserByPhone(value)
	default:
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"search friend mode is not supported",
			map[string]string{
				"field": "mode",
				"value": mode,
			},
			session)

		return
	}

	if err != nil {
		//sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
		//	"searching in db meet error",
		//	nil,
		//	session)
	}

	// 过滤多余的信息
	filterUserInfo(userList, strings.ToLower(mode))
	// 应答数据
	sendBackFriendOpResult(pbmodel.UserOperationType_FindUser,
		"ok",
		nil,
		userList, nil,
		session, friendOpMsg.SendId, Globals.snow.GenerateID())
}

// 分为2种模式
func handleFriendAdd(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {

	}

	uid2 := userInfo.UserId
	// 0.0 检测好友是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return
	}

	if isDel := IsUserDeleted(friendInfo); isDel {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user is deleted", nil, session)
		return
	}

	filterUserInfo1(friendInfo, "")
	friendOpMsg.MsgId = Globals.snow.GenerateID()

	// 0.1 如果已经是好友了，应该直接返回；这里主要是为了方式客户端的错误，防止攻击造成计数异常；

	// true为交友模式，否则为社区模式
	if Globals.Config.Server.FriendMode {
		// 如果对方是自己的粉丝
		bFan, _ := checkFriendIsFan(uid2, session.UserID)
		if bFan {
			user := session.GetUser()
			sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
				"ok",
				user.GetUserInfo(),
				[]*pbmodel.UserInfo{friendInfo},
				nil,
				session, friendOpMsg.SendId, friendOpMsg.MsgId)
			return
		}
		onAddFriendStage1(session.UserID, userInfo.UserId, friendInfo, userInfo, friendOpMsg, session)

	} else {
		// 社区模式，如果自己是对方的粉丝
		bFan, _ := checkFriendIsFan(session.UserID, uid2)
		if bFan {
			user := session.GetUser()
			sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
				"ok",
				user.GetUserInfo(),
				[]*pbmodel.UserInfo{friendInfo},
				nil,
				session, friendOpMsg.SendId, friendOpMsg.MsgId)
			return
		}
		onAddFriendOk(session.UserID, userInfo.UserId, friendInfo, userInfo, friendOpMsg, session)
	}

}

// 服务器收到的好友应答，都是session用户同意或者不同意
func handleFriendApprove(msg *pbmodel.Msg, session *Session) {
	friendOpRetMsg := msg.GetPlainMsg().GetFriendOpRet()
	//params := userOpMsg.GetParams()
	reqUserInfo := friendOpRetMsg.GetUser() // 申请者
	result := friendOpRetMsg.Result
	sendId := friendOpRetMsg.SendId
	msgId := friendOpRetMsg.MsgId

	// 保存记录
	if strings.ToLower(result) == "ok" {

		updateFriendOpResult(reqUserInfo.UserId, session.UserID, msgId, true)
	} else {
		updateFriendOpResult(reqUserInfo.UserId, session.UserID, msgId, false)
	}

	user := session.GetUser()
	// 向用户应答好友请求的结果
	msgRet := newFriendOpResultMsg(pbmodel.UserOperationType_AddFriend, strings.ToLower(result),
		reqUserInfo,
		[]*pbmodel.UserInfo{user.GetUserInfo()},
		friendOpRetMsg.GetParams(),
		sendId,
		msgId)

	trySendMsgToUser(reqUserInfo.UserId, msgRet)

}

// 删除好友，如果是交友模式，则双向删除
func handleFriendRemove(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "userinfo is null", nil, session)
		return
	}

	uid2 := userInfo.UserId
	// 0. 检测好友是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return
	}

	// 在数据库中删除
	pk1 := db.ComputePk(session.UserID)
	pk2 := db.ComputePk(userInfo.UserId)
	uid1 := session.UserID
	Globals.scyllaCli.DeleteFollowing(pk1, pk2, uid1, uid2)

	// 更新redis
	Globals.redisCli.RemoveUserFollowing(uid1, []int64{uid2})
	Globals.redisCli.RemoveUserFans(uid2, []int64{uid1})

	// 更新内存
	meUser := session.GetUser()
	meUser.SetFollow(uid2, false)

	friendUser, ok := Globals.uc.GetUser(uid2)
	if ok && friendUser != nil {
		friendUser.SetFan(uid1, false)
	} else {
		// 集群模式
		if Globals.Config.Server.FriendMode {
			// todo: 通知对方服务器更新
		}
	}

	sendBackFriendOpResult(pbmodel.UserOperationType_RemoveFriend,
		"ok",
		userInfo,
		nil,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)

}

// 拉黑某人，好用，这里不区分好友与非好友
func handleFriendBlock(msg *pbmodel.Msg, session *Session) {

	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "userinfo is null", nil, session)
		return
	}

	uid2 := userInfo.UserId
	// 0. 检测对方是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return
	}

	// 更新数据库
	mask := model.PermissionMaskExist
	uid1 := session.UserID
	blockStore := model.BlockStore{
		FriendStore: model.FriendStore{
			Pk:   db.ComputePk(uid1),
			Uid1: uid1,
			Uid2: uid2,
			Nick: userInfo.GetNickName(),
			Tm:   utils.GetTimeStamp(),
		},
		Perm: int32(mask),
	}

	err := Globals.scyllaCli.InsertBlock(&blockStore)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "insert into db error", nil, session)
		return
	}

	// 更新到redis中
	Globals.redisCli.SetUserBlocks(uid1, []model.BlockStore{blockStore})

	// 更新到自己的内存中
	User := session.GetUser()
	User.SetSelfMask(uid2, uint32(mask))

	// 同步到对方的内存中
	friendUser, ok := Globals.uc.GetUser(uid2)
	if ok && friendUser != nil {
		friendUser.SetFriendToMeMask(uid1, uint32(mask))
	} else {

		if Globals.Config.Server.ClusterMode {
			// todo:
		}
	}

	sendBackFriendOpResult(pbmodel.UserOperationType_BlockFriend,
		"ok",
		userInfo,
		nil,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)
}

// 直接删除即可
func handleFriendUnBlock(msg *pbmodel.Msg, session *Session) {

	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "userinfo is null", nil, session)
		return
	}

	uid1 := session.UserID
	uid2 := userInfo.UserId
	// 0. 检测对方是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return
	}

	// 更新数据库
	err := Globals.scyllaCli.DeleteBlock(db.ComputePk(uid1), uid1, uid2)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "insert into db error", nil, session)
		return
	}

	// 更新到redis中
	Globals.redisCli.RemoveUserBlocks(uid1, []int64{uid2})

	// 更新到自己的内存中
	User := session.GetUser()
	// 意思是没有设置，回头再计算
	User.SetSelfMask(uid2, uint32(0))

	// 同步到对方的内存中
	friendUser, ok := Globals.uc.GetUser(uid2)
	if ok && friendUser != nil {
		friendUser.SetFriendToMeMask(uid1, uint32(0))
	} else {

		if Globals.Config.Server.ClusterMode {
			// todo:
		}
	}

	sendBackFriendOpResult(pbmodel.UserOperationType_UnBlockFriend,
		"ok",
		userInfo,
		nil,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)
}

// 给好友设置相应的权限
func handleFriendPermission(msg *pbmodel.Msg, session *Session) {

	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "userinfo is null", nil, session)
		return
	}

	uid1 := session.UserID
	uid2 := userInfo.UserId

	// 0. 检测对方是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return
	}

	params := friendOpMsg.GetParams()
	permStr, ok := params["permission"]
	if !ok || len(permStr) == 0 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "params.permission is null", nil, session)
		return
	}

	permMask, err := strconv.ParseInt(permStr, 10, 16)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "params.permission is too big", nil, session)
		return
	}

	perm := uint32(permMask) | model.PermissionMaskExist

	// 更新数据库
	blockStore := model.BlockStore{
		FriendStore: model.FriendStore{
			Pk:   db.ComputePk(uid1),
			Uid1: uid1,
			Uid2: uid2,
			Nick: userInfo.GetNickName(),
			Tm:   utils.GetTimeStamp(),
		},
		Perm: int32(perm),
	}

	err = Globals.scyllaCli.InsertBlock(&blockStore)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "insert into db error", nil, session)
		return
	}

	// 更新到redis中
	Globals.redisCli.SetUserBlocks(uid1, []model.BlockStore{blockStore})

	// 更新到自己的内存中
	User := session.GetUser()
	User.SetSelfMask(uid2, uint32(perm))

	// 同步到对方的内存中
	friendUser, ok := Globals.uc.GetUser(uid2)
	if ok && friendUser != nil {
		friendUser.SetFriendToMeMask(uid1, uint32(perm))
	} else {

		if Globals.Config.Server.ClusterMode {
			// todo:
		}
	}

	sendBackFriendOpResult(pbmodel.UserOperationType_SetFriendPermission,
		"ok",
		userInfo,
		nil,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)
}

// 这个比较简单，就是给对方设置一个备注
func handleFriendSetMemo(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	params := friendOpMsg.GetParams()

	mode, ok := params["mode"]

	ModeIndex := 1
	if ok && strings.ToLower(mode) == "fans" {
		ModeIndex = 2
	}

	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {

	}
	nickName := userInfo.NickName
	if len(nickName) == 0 {
		return
	}

	pk1 := db.ComputePk(session.UserID)
	uid1 := session.UserID
	uid2 := userInfo.UserId

	// 更新数据库和redis,不用更新内存
	if ModeIndex == 1 {
		Globals.scyllaCli.SetFollowingNick(pk1, uid1, uid2, nickName)
		Globals.redisCli.SetUserFollowingNick(uid1, uid2, nickName)

	} else {
		Globals.scyllaCli.SetFansNick(pk1, uid1, uid2, nickName)
		Globals.redisCli.SetUserFansNick(uid1, uid2, nickName)
	}

	sendBackFriendOpResult(pbmodel.UserOperationType_SetFriendMemo,
		"ok",
		userInfo,
		nil,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)
}

// 查询好友列表
func handleFriendList(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	fromId := int64(0)
	userinfo := friendOpMsg.GetUser()
	if userinfo == nil {
		fromId = 0
	} else {
		fromId = userinfo.UserId
	}

	params := friendOpMsg.GetParams()
	if params == nil {
		return
	}

	mode, ok := params["mode"]
	pk := db.ComputePk(session.UserID)
	var lst []model.FriendStore
	var err error
	if ok && strings.ToLower(mode) == "fans" {
		lst, err = Globals.scyllaCli.FindFans(pk, session.UserID, fromId, db.MaxFriendCacheSize)

	} else {
		lst, err = Globals.scyllaCli.FindFollowing(pk, session.UserID, fromId, db.MaxFriendCacheSize)
	}

	if err != nil {
		Globals.Logger.Error("handleFriendList() query scylladb error", zap.Error(err))
	}
	flst := FriendStore2UserInfo(lst)
	filterUserInfo(flst, "")

	sendBackFriendOpResult(pbmodel.UserOperationType_SetFriendMemo,
		"ok",
		nil,
		flst,
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)

}

// 转发信息，这里主要是用于转发用户好友申请
func sendForwardFriendOpReq(fid int64, params map[string]string, sendId, msgId int64, session *Session) {

	saveAddFriendLog(session.UserID, fid, sendId, msgId, false)
	user := session.GetUser()

	msgFriendOpReq := pbmodel.FriendOpReq{
		Operation: pbmodel.UserOperationType_AddFriend,
		User:      user.GetUserInfo(),
		SendId:    sendId,
		MsgId:     msgId,
		Params:    params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_FriendOp{
			FriendOp: &msgFriendOpReq,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTFriendOp,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	trySendMsgToUser(fid, &msg)
}

func newFriendOpResultMsg(opCode pbmodel.UserOperationType,
	result string,
	user *pbmodel.UserInfo,
	userList []*pbmodel.UserInfo, params map[string]string, sendId, msgId int64) *pbmodel.Msg {

	msgFriendOpRet := pbmodel.FriendOpResult{
		Operation: opCode,
		Result:    result,
		User:      user,
		Users:     userList,
		SendId:    sendId,
		MsgId:     msgId,
		Params:    params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_FriendOpRet{
			FriendOpRet: &msgFriendOpRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTFriendOpRet,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	return &msg

}

// 发送应答信息
func sendBackFriendOpResult(opCode pbmodel.UserOperationType, result string, user *pbmodel.UserInfo,
	userList []*pbmodel.UserInfo, params map[string]string, session *Session, sendId, msgId int64) {

	msg := newFriendOpResultMsg(opCode, result, user, userList, params, sendId, msgId)
	//session.SendMessage(msg)

	// 多终端登录时候，转发到所有的消息
	trySendMsgToUser(session.UserID, msg)
}

// 保存加好友记录
func saveAddFriendLog(uid1, uid2 int64, reqSendId, msgId int64, bFinish bool) error {
	pk1 := db.ComputePk(uid1)

	ret := 0
	if bFinish {
		ret = 1
	}
	OpRecord := &model.CommonOpStore{
		Pk:   pk1,
		Uid1: uid1,
		Uid2: uid2,
		Gid:  0,
		Id:   msgId,
		Usid: reqSendId,
		Tm:   utils.GetTimeStamp(),
		Tm1:  0,
		Tm2:  0,
		Io:   0,
		St:   0,
		Cmd:  int8(model.CommonUserOpAddRequest),
		Ret:  int8(ret),
		Mask: 0,
		Ref:  0,
		Draf: nil,
	}
	pk2 := db.ComputePk(uid2)
	err := Globals.scyllaCli.SaveUserOp(OpRecord, pk2)
	if err != nil {

	}
	return err
}

// 收到应答时候，更新好友操作记录
func updateFriendOpResult(uid1, uid2 int64, msgId int64, bOk bool) {

	pk1 := db.ComputePk(uid1)
	pk2 := db.ComputePk(uid2)
	ret := 1
	if bOk {
		ret = model.UserOpResultOk
	} else {
		ret = model.UserOpResultRefuse
	}
	Globals.scyllaCli.SetUserOpResult(pk1, pk2, uid1, uid2, msgId, ret)
}

// 是否用户已经被删除了
func IsUserDeleted(userInfo *pbmodel.UserInfo) bool {
	if userInfo.Params != nil {
		status, b := userInfo.Params["status"]
		if b {
			if status == "deleted" {
				return true
			}
		}
	}

	return false
}

// 交友模式，阶段1, 权限验证，转发消息
func onAddFriendStage1(uid1, uid2 int64, friendInfo, userInfo *pbmodel.UserInfo, friendOpMsg *pbmodel.FriendOpReq,
	session *Session) {

	// 1) 查看对方设置的权限
	params := friendInfo.GetParams()
	if params == nil { // 如果没有设置附加信息，默认需要同意
		sendForwardFriendOpReq(uid2, friendOpMsg.GetParams(), friendOpMsg.SendId, friendOpMsg.MsgId, session)
		return
	}

	friendAddMode, ok := params["friendaddmode"]
	//"direct" | "require" | "reject" | "question"
	if !ok {
		sendForwardFriendOpReq(uid2, friendOpMsg.GetParams(), friendOpMsg.SendId, friendOpMsg.MsgId, session)
		return
	}

	switch strings.ToLower(friendAddMode) {
	case "direct":
		onAddFriendOkFriendMode(uid1, uid2, friendInfo, userInfo, friendOpMsg, session)
	case "require":
		sendForwardFriendOpReq(uid2, friendOpMsg.GetParams(), friendOpMsg.SendId, friendOpMsg.MsgId, session)
	case "reject":
		// 应答回执

		sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
			"reject",
			nil,
			[]*pbmodel.UserInfo{friendInfo},
			nil,
			session, friendOpMsg.SendId, friendOpMsg.MsgId)

	case "question":
		question, ok := params["friendaddquestion"]
		if !ok {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "friend need a question, but not find question in friend setting", nil, session)
			return
		}
		answer, ok := params["friendaddanswer"]
		if !ok {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "friend need a question, but not find answer in friend setting", nil, session)
			return
		}

		// 如果已经带了答案
		if friendOpMsg.GetParams() != nil {
			answerInReq, ok := friendOpMsg.GetParams()["answer"]
			if ok && answerInReq == answer {
				onAddFriendOkFriendMode(uid1, uid2, friendInfo, userInfo, friendOpMsg, session)
				return
			}
		}

		sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
			"question",
			nil,
			[]*pbmodel.UserInfo{friendInfo},
			map[string]string{
				"question": question,
			},
			session, friendOpMsg.SendId, friendOpMsg.MsgId)
	}

	return
}

// 将数据的基础操作在这里组合
// 1.0) 检测用户是否被删除了，或者是否存在
// 1.1) 保存日志；
// 1.2）双向保存到db;
// 1.3) 同步到自己的redis, 同步到自己的内存；
// 1.4) 如果redis中有对方的粉丝缓存，保存到对方的redis中；
// 1.5） 如果本机有对方的用户内存块，则保存到内存块；
// 1.6）如果对方在线，则需要通知对方有新粉丝；
func onAddFriendOk(uid1, uid2 int64, friendInfo, userInfo *pbmodel.UserInfo, friendOpMsg *pbmodel.FriendOpReq, session *Session) error {
	var err error

	// 1. 日志
	err = saveAddFriendLog(uid1, uid2, friendOpMsg.SendId, Globals.snow.GenerateID(), true)
	if err != nil {
		return err
	}

	// 2. 保存双向 关注信息，和粉丝信息
	friendName := userInfo.GetNickName()
	if friendName == "" {
		friendName = userInfo.GetUserName()
		if len(friendName) == 0 {
			friendName = friendInfo.NickName
		}
	}
	friend := model.FriendStore{
		Pk:   db.ComputePk(uid1),
		Uid1: uid1,
		Uid2: uid2,
		Tm:   time.Now().UTC().UnixMilli(),
		Nick: friendName,
	}

	// 这里是绝对不应该是空的
	user := session.GetUser()
	if user == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "user not in cache", nil, session)
		return errors.New("user is not in  cache")
	}

	name := "momo"
	name = user.GetNickName()
	if name == "" {
		name = user.GetUserName()
	}

	fan := model.FriendStore{
		Pk:   db.ComputePk(uid2),
		Uid1: uid2,
		Uid2: uid1,
		Tm:   time.Now().UTC().UnixMilli(),
		Nick: name,
	}

	// 2. 保存到scyllaDb中
	err = Globals.scyllaCli.InsertFollowing(&friend, &fan)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save following and fans  to db fail", nil, session)
		return err
	}

	// 应该补充，计算好友的个数和粉丝的个数
	incUserFollowsAndFans(uid1, uid2)

	// 3.关注信息同步到redis
	// 这里登录时候已经加载了，如果一个好友没有，这个地方也是空的，不过这是第一个好友
	err = Globals.redisCli.AddUserFollowing(uid1, []model.FriendStore{friend})
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save following   to redis fail", nil, session)
		return err
	}
	// 3.同步到本地的内存
	user.SetFollow(uid2, true)

	// 4.对方的粉丝信息是否保存需要同步到redis
	var bHasRedisFan = false
	bHasRedisFan, err = Globals.redisCli.ExistFollowing(uid2)
	// 如果存在，这里已经续命了
	if bHasRedisFan {
		err = Globals.redisCli.AddUserFans(uid2, []model.FriendStore{fan})
		if err != nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save peer fans  to redis fail", nil, session)
			return err
		}
	}

	// 5 是否同步到对方的内存
	friendUser, b := Globals.uc.GetUser(uid2)
	if b && friendUser != nil {
		friendUser.SetFan(uid1, true)
	}

	// 6. 如果对方在线，则需要通知对方有新粉丝
	msg := newFriendOpResultMsg(pbmodel.UserOperationType_AddFriend, "notify",
		user.GetUserInfo(), nil, nil, friendOpMsg.SendId, friendOpMsg.MsgId)
	trySendMsgToUser(uid2, msg)

	// 应答回执
	delete(friendInfo.Params, "pwd")
	sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
		"ok",
		user.GetUserInfo(),
		[]*pbmodel.UserInfo{friendInfo},
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)

	return err
}

func onAddFriendOkFriendMode(uid1, uid2 int64, friendInfo, userInfo *pbmodel.UserInfo, friendOpMsg *pbmodel.FriendOpReq, session *Session) error {
	var err error

	// 1. 日志
	err = saveAddFriendLog(uid1, uid2, friendOpMsg.SendId, Globals.snow.GenerateID(), true)
	if err != nil {
		return err
	}

	// 2. 保存双向 关注信息，和粉丝信息
	friendName := userInfo.GetNickName()
	if friendName == "" {
		friendName = userInfo.GetUserName()
		if len(friendName) == 0 {
			friendName = friendInfo.NickName
		}
	}
	friend := model.FriendStore{
		Pk:   db.ComputePk(uid1),
		Uid1: uid1,
		Uid2: uid2,
		Tm:   time.Now().UTC().UnixMilli(),
		Nick: friendName,
	}

	// 这里是绝对不应该是空的
	user := session.GetUser()
	if user == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "user not in cache", nil, session)
		return errors.New("user is not in  cache")
	}

	name := "momo"
	name = user.GetNickName()
	if name == "" {
		name = user.GetUserName()
	}

	fan := model.FriendStore{
		Pk:   db.ComputePk(uid2),
		Uid1: uid2,
		Uid2: uid1,
		Tm:   time.Now().UTC().UnixMilli(),
		Nick: name,
	}

	// 2. 保存到scyllaDb中
	err = Globals.scyllaCli.InsertFollowing(&friend, &fan)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save following and fans  to db fail", nil, session)
		return err
	}

	// 应该补充，计算好友的个数和粉丝的个数
	incUserFollowsAndFans(uid1, uid2)

	// 3.关注信息同步到redis
	// 这里登录时候已经加载了，如果一个好友没有，这个地方也是空的，不过这是第一个好友
	err = Globals.redisCli.AddUserFollowing(uid1, []model.FriendStore{friend})
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save following   to redis fail", nil, session)
		return err
	}
	// 3.同步到本地的内存
	user.SetFollow(uid2, true)

	// 4.对方的粉丝信息是否保存需要同步到redis
	var bHasRedisFan = false
	bHasRedisFan, err = Globals.redisCli.ExistFollowing(uid2)
	// 如果存在，这里已经续命了
	if bHasRedisFan {
		err = Globals.redisCli.AddUserFans(uid2, []model.FriendStore{fan})
		if err != nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save peer fans  to redis fail", nil, session)
			return err
		}
	}

	// 5 是否同步到对方的内存
	friendUser, b := Globals.uc.GetUser(uid2)
	if b && friendUser != nil {
		friendUser.SetFan(uid1, true)
	}

	// 6. 如果对方在线，则需要通知对方有新粉丝
	msg := newFriendOpResultMsg(pbmodel.UserOperationType_AddFriend, "notify",
		user.GetUserInfo(), nil, nil, friendOpMsg.SendId, friendOpMsg.MsgId)
	trySendMsgToUser(uid2, msg)

	// 应答回执
	delete(friendInfo.Params, "pwd")
	sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
		"ok",
		user.GetUserInfo(),
		[]*pbmodel.UserInfo{friendInfo},
		nil,
		session, friendOpMsg.SendId, friendOpMsg.MsgId)

	return err
}

// 设置好友关系时候应该，同步更新用户的好友和粉丝个数
// todo: 同步到redis和内存中
func incUserFollowsAndFans(userId, fid int64) {
	Globals.mongoCli.UpdateUserFieldIncNum("params.follows", 1, userId)
	Globals.mongoCli.UpdateUserFieldIncNum("params.fans", 1, fid)
	//
}

func decUserFollowsAndFans(userId, fid int64) {
	Globals.mongoCli.UpdateUserFieldIncNum("params.follows", -1, userId)
	Globals.mongoCli.UpdateUserFieldIncNum("params.fans", -1, fid)
	//
}
