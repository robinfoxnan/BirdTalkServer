package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"go.uber.org/zap"
	"strings"
)

// 用户的基本操作
func handleUserOp(msg *pbmodel.Msg, session *Session) {
	// 检查数据指针是否为空
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	if userOpMsg == nil {
		Globals.Logger.Debug("receive wrong User op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "user op  is null", nil, session)
		return
	}

	// 二次分发
	opCode := userOpMsg.Operation

	// 除了注册和登录，都需要验证是否登录与权限
	if opCode != pbmodel.UserOperationType_Login && opCode != pbmodel.UserOperationType_RegisterUser {
		ok := checkUserLogin(session) // 内部发送错误通知
		if !ok {
			return
		}
	}

	// 只有管理员才有权限禁用和恢复某个用户
	if opCode == pbmodel.UserOperationType_RecoverUser || opCode == pbmodel.UserOperationType_RecoverUser {
		ok := CheckUserPermission(session)
		if !ok {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "permission error", nil, session)
			return
		}
	}

	switch opCode {
	case pbmodel.UserOperationType_RegisterUser:
		handleUserRegister(msg, session)
	case pbmodel.UserOperationType_UnregisterUser:
		handleUserUnRegister(msg, session)
	case pbmodel.UserOperationType_DisableUser:
		handleUserDisable(msg, session)
	case pbmodel.UserOperationType_RecoverUser:
		handleUserEnable(msg, session)
	case pbmodel.UserOperationType_SetUserInfo:
		handleUserInfo(msg, session)
	case pbmodel.UserOperationType_RealNameVerification: // 更改用户手机邮件时候需要，登录时候也是
		handleUserVerification(msg, session)
	case pbmodel.UserOperationType_Login:
		handleUserLogin(msg, session)
	case pbmodel.UserOperationType_Logout:
		handleUserLogout(msg, session)
	default:
		Globals.Logger.Info("receive unknown user op",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID),
			zap.String("code", opCode.String()),
		)
	}

}

// ////////////////////////////////////////////////////////////////
// 下面是细分的子类型
// 1 注册用户
func handleUserRegister(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"register msg userinfo is nil",
			map[string]string{"field": "userinfo"},
			session)
		return
	}

	regMode := 1 // 1匿名，2邮件，3 手机
	params := userInfo.GetParams()
	if params != nil {
		modeStr, ok := params["regmode"]
		if ok {
			modeStr = strings.ToLower(modeStr)
			switch modeStr {
			case "anonymous":
				regMode = 1
			case "email":
				regMode = 2
			case "phone":
				regMode = 3
			default:
				regMode = 1
			}
		}
	}

	// 邮件是否合法
	ok := utils.IsValidEmail(userInfo.Email)
	if ok {
		regMode = 2
	} else {
		// 通知用户信息不对
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"email is not correct",
			map[string]string{"field": userInfo.Email},
			session)
	}

	if regMode == 2 {
		// 发送验证码
		err := createTempUser(userInfo, session) // 内部记录错误并回执
		if err != nil {
			return
		}
		// 通知用户免登录
		SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
			userInfo,
			true, "waitcode", session)
		session.SetStatus(model.UserStatusOk)

	} else {
		// 匿名注册
		err := createUser(userInfo, session) // 内部记录错误并回执
		if err != nil {
			return
		}
		// 加载用户信息
		err = LoadUserNew(session) // 内部记录错误并回执
		if err != nil {
			// 通知用户重新登录
			SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
				userInfo,
				true, "needlogin", session)
			return
		}

		// 通知用户免登录
		SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
			userInfo,
			true, "loginok", session)

		session.SetStatus(model.UserStatusOk)
	}

}

// 2 注销用户
func handleUserUnRegister(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()

	_, err := Globals.mongoCli.UpdateGroupInfoPart(session.UserID, map[string]interface{}{"params.status": "deleted"}, nil)

	if err != nil {
		Globals.Logger.Fatal("unreg user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"unregister user meet database err",
			nil,
			session)
		return
	}
	Globals.Logger.Info("unreg user ok", zap.Int64("user id", session.UserID))
	SendBackUserOp(pbmodel.UserOperationType_UnregisterUser,
		userInfo,
		true, "unregok", session)

	// todo: 清理资源
}

// 3 禁用用户
func handleUserDisable(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"register msg userinfo is nil",
			map[string]string{"field": "userinfo"},
			session)
		return
	}

	reason := ""
	if userOpMsg.Params != nil {
		reason, _ = userInfo.Params["reason"]
	}

	// 更新数据库
	params := map[string]interface{}{
		"params.status": "deleted",
		"params.reason": reason,
	}
	_, err := Globals.mongoCli.UpdateGroupInfoPart(userInfo.UserId, params, nil)
	if err != nil {
		Globals.Logger.Fatal("disable user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"unregister user meet database err",
			nil,
			session)
		return
	}

	// 更新redis
	err = Globals.redisCli.UpdateUserInfoPart(userInfo.UserId, params, nil)

	// 更新内存
	user, ok := Globals.uc.GetUser(userInfo.UserId)
	if ok && user != nil {
		user.SetBaseExtraKeyValue("status", "disabled")
		user.SetBaseExtraKeyValue("reason", reason)
	}

	Globals.Logger.Info("disable user ok", zap.Int64("user id", userInfo.UserId),
		zap.Int64("admin id", session.UserID))

	SendBackUserOp(pbmodel.UserOperationType_DisableUser,
		userInfo,
		true, "unregok", session)
}

// 4 恢复用户
func handleUserEnable(msg *pbmodel.Msg, session *Session) {

}

// 5 设置用户信息
func handleUserInfo(msg *pbmodel.Msg, session *Session) {

}

// 6 验证码检查
func handleUserVerification(msg *pbmodel.Msg, session *Session) {

}

// 7 登录
func handleUserLogin(msg *pbmodel.Msg, session *Session) {

}

// 8 退出
func handleUserLogout(msg *pbmodel.Msg, session *Session) {

}

// 注册后的回复
func SendBackUserOp(opCode pbmodel.UserOperationType, userInfo *pbmodel.UserInfo,
	ret bool, status string, session *Session) error {

	result := "ok"
	if !ret {
		result = "fail"
	}

	params := map[string]string{
		"status": status,
	}

	msgUserOpRet := pbmodel.UserOpResult{
		Operation: opCode,
		Result:    result,
		Users:     []*pbmodel.UserInfo{userInfo},
		Params:    params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_UserOpRet{
			UserOpRet: &msgUserOpRet,
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
	session.SendMessage(msg)
	return nil
}

// 检查用户是否是管理员，
func CheckUserPermission(session *Session) bool {
	if session.UserID < 10000 {
		return true
	}
	return false
}
