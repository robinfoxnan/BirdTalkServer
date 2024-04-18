package core

import (
	"birdtalk/server/pbmodel"
	"go.uber.org/zap"
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
	case pbmodel.UserOperationType_AddFriend: // 交友模式下需要转发给好友
	case pbmodel.UserOperationType_ApproveFriend: // 用户应答同意或者拒绝
	case pbmodel.UserOperationType_RemoveFriend:
	case pbmodel.UserOperationType_BlockFriend:
	case pbmodel.UserOperationType_UnBlockFriend:
	case pbmodel.UserOperationType_SetFriendPermission:
	case pbmodel.UserOperationType_SetFriendMemo:
	default:
		Globals.Logger.Info("receive unknown friend op",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID),
			zap.String("code", opCode.String()),
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
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID),
			zap.String("code", opCode.String()),
		)
	}
}

//////////////////////////////////////////////////////////////////
// 下面是细分的子类型
