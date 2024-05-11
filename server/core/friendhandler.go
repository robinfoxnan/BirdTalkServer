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
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	params := userOpMsg.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"must hava Params field",
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
	var userList []pbmodel.UserInfo
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
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"searching in db meet error",
			nil,
			session)
	}

	// 应答数据
	sendBackFriendOpResult(pbmodel.UserOperationType_FindUser,
		"ok",
		nil,
		userList, nil,
		session)
}

// 分为2种模式
func handleFriendAdd(msg *pbmodel.Msg, session *Session) {
	friendOpMsg := msg.GetPlainMsg().GetFriendOp()
	//params := userOpMsg.GetParams()
	userInfo := friendOpMsg.GetUser()
	if userInfo == nil {

	}

	// true为交友模式，否则为社区模式
	if Globals.Config.Server.FriendMode {
		onAddFriendStage1(session.UserID, userInfo.UserId, friendOpMsg, userInfo, session)

	} else {
		// 社区模式，直接更新
		onAddFriendCommunity(session.UserID, userInfo.UserId, friendOpMsg, userInfo, session)
	}

}

func handleFriendApprove(msg *pbmodel.Msg, session *Session) {

}

func handleFriendRemove(msg *pbmodel.Msg, session *Session) {

}

func handleFriendBlock(msg *pbmodel.Msg, session *Session) {

}

func handleFriendUnBlock(msg *pbmodel.Msg, session *Session) {

}

func handleFriendPermission(msg *pbmodel.Msg, session *Session) {

}

func handleFriendSetMemo(msg *pbmodel.Msg, session *Session) {

}

// 转发信息，这里主要是用于转发用户好友申请
func sendForwardFriendOpReq() {

}

func newFriendOpResultMsg(opCode pbmodel.UserOperationType,
	result string,
	user *pbmodel.UserInfo,
	userList []pbmodel.UserInfo, params map[string]string) *pbmodel.Msg {

	uList := make([]*pbmodel.UserInfo, len(userList))
	for index, item := range userList {
		uList[index] = &item
	}
	msgFriendOpRet := pbmodel.FriendOpResult{
		Operation: opCode,
		Result:    result,
		User:      user,
		Users:     uList,
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
		MsgType:  pbmodel.ComMsgType_MsgTUserOpRet,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	return &msg

}

// 发送应答信息
func sendBackFriendOpResult(opCode pbmodel.UserOperationType, result string, user *pbmodel.UserInfo,
	userList []pbmodel.UserInfo, params map[string]string, session *Session) {

	msg := newFriendOpResultMsg(opCode, result, user, userList, params)
	session.SendMessage(msg)

}

// 保存加好友记录
func saveAddFriendLog(uid1, uid2 int64, friendOpMsg *pbmodel.FriendOpReq, userInfo *pbmodel.UserInfo, session *Session) error {
	pk1 := db.ComputePk(uid1)
	OpRecord := &model.CommonOpStore{
		Pk:   pk1,
		Uid1: uid1,
		Uid2: uid2,
		Gid:  0,
		Id:   Globals.snow.GenerateID(),
		Usid: friendOpMsg.SendId,
		Tm:   utils.GetTimeStamp(),
		Tm1:  0,
		Tm2:  0,
		Io:   0,
		St:   0,
		Cmd:  int8(pbmodel.UserOperationType_AddFriend),
		Ret:  0,
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
func onAddFriendStage1(uid1, uid2 int64, friendOpMsg *pbmodel.FriendOpReq, userInfo *pbmodel.UserInfo, session *Session) {

}

// 将数据的基础操作在这里组合
// 1.1) 保存日志；
// 1.2）双向保存到db;
// 1.3) 同步到自己的redis, 同步到自己的内存；
// 1.4) 如果redis中有对方的粉丝缓存，保存到对方的redis中；
// 1.5） 如果本机有对方的用户内存块，则保存到内存块；
// 1.6）如果对方在线，则需要通知对方有新粉丝；
func onAddFriendCommunity(uid1, uid2 int64, friendOpMsg *pbmodel.FriendOpReq, userInfo *pbmodel.UserInfo, session *Session) error {
	var err error

	// 0. 检测好友是否存在
	friendInfo, _, _ := findUserInfo(uid2)
	if friendInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user not in db", nil, session)
		return errors.New("friend not exist in db")
	}

	if isDel := IsUserDeleted(friendInfo); isDel {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user is deleted", nil, session)
		return errors.New("friend not exist in db")
	}

	// 1. 日志
	err = saveAddFriendLog(uid1, uid2, friendOpMsg, userInfo, session)
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

	// 2.
	err = Globals.scyllaCli.InsertFollowing(&friend, &fan)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save following and fans  to db fail", nil, session)
		return err
	}

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
	msg := newFriendOpResultMsg(pbmodel.UserOperationType_AddFriend, "notify", user.GetUserInfo(), nil, nil)
	trySendMsgToUser(uid2, msg)

	// 应答回执
	delete(friendInfo.Params, "pwd")
	sendBackFriendOpResult(pbmodel.UserOperationType_AddFriend,
		"ok",
		user.GetUserInfo(),
		[]pbmodel.UserInfo{*friendInfo},
		nil,
		session)

	return err
}
