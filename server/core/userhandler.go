package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"go.uber.org/zap"
	"strconv"
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
	if opCode == pbmodel.UserOperationType_DisableUser || opCode == pbmodel.UserOperationType_RecoverUser {
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
		// 先检查邮箱是否合法
		lst, err := Globals.mongoCli.FindUserByEmail(userInfo.Email)
		if len(lst) > 0 {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"email is already used by user",
				map[string]string{"field": userInfo.Email},
				session)
			return
		}

		// 发送验证码
		err = createTempUser(userInfo, session) // 内部记录错误并回执
		if err != nil {
			return
		}
		// 通知用户
		SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
			userInfo,
			true, "waitcode", session)
		session.SetStatus(model.UserStatusRegister | model.UserStatusValidate)

	} else {
		// 匿名注册
		err := createUser(userInfo, session) // 内部记录错误并回执
		if err != nil {
			return
		}

		// 加载并回执
		onRegisterSuccess(session)
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
			"disable msg userinfo is nil",
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
		"params.status": "disabled",
		"params.reason": reason,
	}
	_, err := Globals.mongoCli.UpdateGroupInfoPart(userInfo.UserId, params, nil)
	if err != nil {
		Globals.Logger.Fatal("disable user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"disable user, meet database err",
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
		true, "disabled", session)
}

// 4 恢复用户
func handleUserEnable(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"recover msg userinfo is nil",
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
		"params.status": "ok",
		"params.reason": reason,
	}
	_, err := Globals.mongoCli.UpdateGroupInfoPart(userInfo.UserId, params, nil)
	if err != nil {
		Globals.Logger.Fatal("recover user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"recover user, meet database err",
			nil,
			session)
		return
	}

	// 更新redis
	err = Globals.redisCli.UpdateUserInfoPart(userInfo.UserId, params, nil)

	// 更新内存
	user, ok := Globals.uc.GetUser(userInfo.UserId)
	if ok && user != nil {
		user.SetBaseExtraKeyValue("status", "ok")
		user.SetBaseExtraKeyValue("reason", reason)
	}

	Globals.Logger.Info("recover user ok", zap.Int64("user id", userInfo.UserId),
		zap.Int64("admin id", session.UserID))

	SendBackUserOp(pbmodel.UserOperationType_DisableUser,
		userInfo,
		true, "ok", session)
}

// 计算需要更新的字段
func mergeUserinfo(user *model.User, userInfo *pbmodel.UserInfo) map[string]interface{} {

	setData := make(map[string]interface{})
	for k, v := range userInfo.Params {
		name := "params." + strings.ToLower(k)
		setData[name] = v
	}

	if user.UserInfo.UserName != userInfo.UserName {
		setData["username"] = userInfo.UserName
	}

	if user.UserInfo.NickName != userInfo.NickName {
		setData["nickname"] = userInfo.NickName
	}

	// 这里需要验证，
	//if user.UserInfo.Email != userInfo.Email{
	//	setData["email"] = userInfo.Email
	//}
	//
	//if user.UserInfo.Phone != userInfo.Phone {
	//	setData["phone"] = userInfo.Phone
	//}

	if user.UserInfo.Gender != userInfo.Gender {
		setData["gender"] = userInfo.Gender
	}

	if user.UserInfo.Age != userInfo.Age {
		setData["age"] = userInfo.Age
	}

	if user.UserInfo.Region != userInfo.Region {
		setData["region"] = userInfo.Region
	}

	if user.UserInfo.Icon != userInfo.Icon {
		setData["icon"] = userInfo.Icon
	}
	return setData
}

// 5 设置用户信息
func handleUserInfo(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"update user info, but msg userinfo is nil",
			map[string]string{"field": "userinfo"},
			session)
		return
	}

	user, ok := Globals.uc.GetUser(session.UserID)
	if !ok || user == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"can't find user in user cache ",
			map[string]string{"field": "userinfo"},
			session)
		return
	}

	// 更新数据库
	setData := mergeUserinfo(user, userInfo)

	_, err := Globals.mongoCli.UpdateGroupInfoPart(userInfo.UserId, setData, nil)
	if err != nil {
		Globals.Logger.Fatal("update user, mongodb err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"update user, meet database err",
			nil,
			session)
		return
	}

	// 更新redis
	err = Globals.redisCli.UpdateUserInfoPart(userInfo.UserId, setData, nil)

	// 更新内存
	user.SetBaseValue(userInfo)

	Globals.Logger.Info("recover user ok", zap.Int64("user id", userInfo.UserId),
		zap.Int64("admin id", session.UserID))

	SendBackUserOp(pbmodel.UserOperationType_SetUserInfo,
		userInfo,
		true, "ok", session)
}

// 6 验证码检查
// 3种可能，注册，登录，实名认证
func handleUserVerification(msg *pbmodel.Msg, session *Session) {
	// 如果没有没有
	if !session.HasStatus(model.UserStatusValidate) {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTStage),
			"not has a status of Verification ",
			nil,
			session)
	}
	// 检查数据
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	params := userOpMsg.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"not has params field",
			map[string]string{"field": "params"},
			session)
		return
	}

	// 用户提交的CODE
	codeUser, ok := params["code"]
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"verification code, but params.code is nil",
			map[string]string{"field": "params.code"},
			session)
		return
	}

	// 服务端存储的CODE
	codeTemp := session.GetKeyValue("code")
	if len(codeTemp) == 0 {

		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"server does not has a verification code",
			nil,
			session)
		return
	}
	// 比对CODE是否一致
	if codeUser != codeTemp {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"verification code, but params.code is not correct",
			map[string]string{"field": "params.code"},
			session)
		return
	}

	// 这里就是一致了，根据不同的阶段做出应答
	if session.HasStatus(model.UserStatusRegister) {
		// 创建用户，
		err := createUser(session.TempUserInfo, session) // 内部记录错误并回执
		if err != nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
				"create user meet error",
				nil,
				session)
			return
		}
		onRegisterSuccess(session)

	} else if session.HasStatus(model.UserStatusLogin) {
		onLoginSuccess(session)

	} else if session.HasStatus(model.UserStatusChangeInfo) {
		onChangeInfoSuccess(session)
	}

	// 取消掉当前的这个状态
	session.UnSetStatus(model.UserStatusValidate)
}

// 7 登录
func handleUserLogin(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()
	if userInfo == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"login msg userinfo is nil",
			map[string]string{"field": "userinfo"},
			session)
		return
	}

	var userInfoDb *pbmodel.UserInfo = nil

	loginMode := 1 // 1账号，2邮件，3 手机
	params := userInfo.GetParams()
	if params != nil {
		modeStr, ok := params["loginmode"]
		if ok {
			modeStr = strings.ToLower(modeStr)
			switch modeStr {
			case "id":
				loginMode = 1
			case "email":
				loginMode = 2
				ok = utils.IsValidEmail(userInfo.Email)
				if !ok {
					sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
						"email format is not correct",
						map[string]string{"field": userInfo.Email},
						session)
				}

				lst, err := Globals.mongoCli.FindUserByEmail(userInfo.Email)
				if err != nil || len(lst) < 1 {
					sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
						"email is not correct, can't find a user using it",
						map[string]string{"field": userInfo.Email},
						session)
				}
				userInfoDb = &lst[0]
				session.UserID = userInfoDb.UserId
				session.TempUserInfo = userInfoDb

				// 生成临时
				code := utils.GenerateCheckCode(5)
				code = "12345"
				session.SetKeyValue("code", code)
				// 使用后台将验证码发送给用户

				session.SetStatus(model.UserStatusLogin | model.UserStatusValidate)
				// 通知用户
				if userInfoDb.Params != nil {
					delete(userInfoDb.Params, "pwd")
				}
				SendBackUserOp(pbmodel.UserOperationType_Login,
					userInfoDb,
					true, "waitcode", session)

			case "phone":
				loginMode = 3

			default:
				loginMode = 1
			}
		}
	}

	// 用户ID和口令的匿名用户
	if loginMode == 1 {
		lst, err := Globals.mongoCli.FindUserById(userInfo.UserId)
		if err != nil || len(lst) != 1 {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"pwd is not set",
				map[string]string{"field": "userid", "value": strconv.FormatInt(userInfo.UserId, 10)},
				session)
			return
		}

		if userInfo.Params == nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"params is nil",
				map[string]string{"field": "params", "value": "nil"},
				session)
			return
		}

		pwdUser, ok := userInfo.Params["pwd"]
		if !ok {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"pwd is not set",
				map[string]string{"field": "pwd", "value": pwdUser},
				session)
			return
		}

		userInfoDb = &lst[0]
		if userInfoDb.Params == nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
				"pwd is not in db",
				map[string]string{"field": "pwd"},
				session)
			return
		}
		pwdDb, ok := userInfoDb.Params["pwd"]

		// 用户口令不对
		if pwdUser != pwdDb {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"pwd is not correct",
				map[string]string{"field": "pwd", "value": pwdUser},
				session)
			return
		}
		session.UserID = userInfoDb.UserId
		session.TempUserInfo = userInfoDb

	} else {
		// 通知用户信息不对
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"phone login is not available",
			map[string]string{"field": userInfo.Phone},
			session)
		return
	}
	onLoginSuccess(session)

	return
}

// 注册结束后，执行这一段
func onRegisterSuccess(session *Session) {
	// 加载用户信息
	err := LoadUserNew(session) // 内部记录错误并回执
	if err != nil {
		// 通知用户重新登录
		SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
			session.TempUserInfo,
			true, "needlogin", session)
		return
	}

	// 通知用户免登录
	SendBackUserOp(pbmodel.UserOperationType_RegisterUser,
		session.TempUserInfo,
		true, "loginok", session)

	session.SetStatus(model.UserStatusOk)
}

// 登录最后收尾的
func onLoginSuccess(session *Session) {
	// 加载数据
	err := LoadUserLogin(session)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"load user meet error",
			map[string]string{"error": err.Error()},
			session)
		return
	}
	OnUserLogin(session)

	// 通知用户免登录
	SendBackUserOp(pbmodel.UserOperationType_Login,
		session.TempUserInfo,
		true, "loginok", session)
	session.SetStatus(model.UserStatusOk)
}

// 验证码对了才保存到数据库
func onChangeInfoSuccess(session *Session) {

}

// 8 退出，是否需要删除缓存？不，因为可能有多用户登录
func handleUserLogout(msg *pbmodel.Msg, session *Session) {
	session.StopSession("logout")
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
