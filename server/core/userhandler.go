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
		ok := checkUserLogin(session)
		if !ok {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
			return
		}
	}

	// 只有管理员才有权限禁用和恢复某个用户
	if opCode == pbmodel.UserOperationType_RecoverUser || opCode == pbmodel.UserOperationType_RecoverUser {
		ok := checkUserPermission(session)
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

	SendBackUserOp(pbmodel.UserOperationType_UnregisterUser,
		userInfo,
		true, "loginok", session)
}

// 3 禁用用户
func handleUserDisable(msg *pbmodel.Msg, session *Session) {

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
