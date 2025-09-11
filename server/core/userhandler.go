package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"fmt"
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
	if (opCode != pbmodel.UserOperationType_Login) &&
		(opCode != pbmodel.UserOperationType_RegisterUser) &&
		(opCode != pbmodel.UserOperationType_RealNameVerification) {
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

	// 如果登录，则不能注册账号
	if session.HasStatus(model.UserStatusOk) {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTStage),
			"you are already login, do not register too much id",
			nil,
			session)
		return
	}

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
	params := userOpMsg.GetParams()
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

	if userInfo.GetUserName() == "" {
		userInfo.UserName = "momo"

	}
	if userInfo.GetNickName() == "" {
		userInfo.NickName = userInfo.UserName
	}

	// 邮件是否合法
	if regMode == 2 {
		ok := utils.IsValidEmail(userInfo.Email)
		if !ok {
			// 通知用户信息不对
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"email is not correct",
				map[string]string{"field": userInfo.Email},
				session)
			return
		}
	}

	if regMode == 2 {
		// 先检查邮箱是否合法
		lst, err := Globals.mongoCli.FindUserByEmail(userInfo.Email)
		if lst != nil && len(lst) > 0 {
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
		// 2025-05-17 修正
		return

	} else {
		// 如果设置了邮件，也需要检查，防止乱设置
		if len(userInfo.Email) > 0 {
			lst, _ := Globals.mongoCli.FindUserByEmail(userInfo.Email)
			if lst != nil && len(lst) > 0 {
				sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
					"email is already used by user",
					map[string]string{"field": userInfo.Email},
					session)
				return
			}
		}

		// 匿名注册
		err := createUser(userInfo, session) // 内部记录错误并回执
		if err != nil {
			return
		}

		// 加载并回执
		onRegisterSuccess(session)
		fmt.Println(session.TempUserInfo)
	}

}

// 2 注销用户
func handleUserUnRegister(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	userInfo := userOpMsg.GetUser()

	_, err := Globals.mongoCli.UpdateUserInfoPart(session.UserID, map[string]interface{}{"params.status": "deleted"}, nil)

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

	// 更新redis
	Globals.redisCli.RemoveUser(session.UserID)

	// 更新用户信息
	user, ok := Globals.uc.GetUser(session.UserID)
	if ok && user != nil {
		user.SetDeleted()
	}

	// todo: 清理资源
	session.StopSession("unregok")
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
	_, err := Globals.mongoCli.UpdateUserInfoPart(userInfo.UserId, params, nil)
	if err != nil {
		Globals.Logger.Fatal("disable user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"disable user, meet database err",
			nil,
			session)
		return
	}

	params = map[string]interface{}{
		"Params.status": "disabled",
		"Params.reason": reason,
	}
	// 更新redis
	err = Globals.redisCli.UpdateUserInfoPart(userInfo.UserId, params, nil)

	// 更新内存
	user, ok := Globals.uc.GetUser(userInfo.UserId)
	if ok && user != nil {
		user.SetDisabled(reason)
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
	_, err := Globals.mongoCli.UpdateUserInfoPart(userInfo.UserId, params, nil)
	if err != nil {
		Globals.Logger.Fatal("recover user set user status err", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"recover user, meet database err",
			nil,
			session)
		return
	}

	// 更新redis
	params = map[string]interface{}{
		"Params.status": "ok",
		"Params.reason": reason,
	}
	err = Globals.redisCli.UpdateUserInfoPart(userInfo.UserId, params, nil)

	// 更新内存
	user, ok := Globals.uc.GetUser(userInfo.UserId)
	if ok && user != nil {
		user.SetRecover()
	}

	Globals.Logger.Info("recover user ok", zap.Int64("user id", userInfo.UserId),
		zap.Int64("admin id", session.UserID))

	SendBackUserOp(pbmodel.UserOperationType_DisableUser,
		userInfo,
		true, "ok", session)
}

// 计算需要更新的字段
func mergeUserinfo(user *model.User, params map[string]string, session *Session) (map[string]interface{}, map[string]interface{}, bool, bool) {

	setData := make(map[string]interface{})
	setDataRedis := make(map[string]interface{})
	hasEmail := false
	hasPhone := false

	for k, v := range params {
		key := strings.ToLower(k)

		switch key {
		case "params.status": // 防止攻击
			continue

		case "username":
			setData["username"] = v
			setDataRedis["UserName"] = v
			user.SetUserName(v)
		case "nickname":
			setData["nickname"] = v
			setDataRedis["NickName"] = v
			user.SetNickName(v)
		case "age":
			nAge, _ := strconv.Atoi(v)
			setData["age"] = nAge
			setDataRedis["Age"] = nAge
			user.SetAge(v)
		case "region":
			setData["region"] = v
			setDataRedis["Region"] = v
			user.SetRegion(v)
		case "gender":
			setData["gender"] = v
			setDataRedis["Gender"] = v
			user.SetGender(v)
		case "icon":
			setData["icon"] = v
			setDataRedis["Icon"] = v
			user.SetIcon(v)
		case "email":
			hasEmail = true
			session.SetKeyValue("changeEmail", v)
		case "phone":
			hasPhone = true
			session.SetKeyValue("changePhone", v)

		default:
			setData[key] = v
			setDataRedis[k] = v
			surfix := strings.TrimLeft(key, "params.")
			user.SetBaseKeyValue(surfix, v)
		}

	}

	return setData, setDataRedis, hasEmail, hasPhone
}

// 5 设置用户信息
func handleUserInfo(msg *pbmodel.Msg, session *Session) {
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	//userInfo := userOpMsg.GetUser()
	//if userInfo == nil {
	//	sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
	//		"update user info, but msg userinfo is nil",
	//		map[string]string{"field": "userinfo"},
	//		session)
	//	return
	//}

	// 检查要更改的字段
	params := userOpMsg.GetParams()
	if params == nil || len(params) == 0 {
		SendBackUserOp(pbmodel.UserOperationType_SetUserInfo,
			nil,
			true, "ok", session)
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
	setDataMongo, setDataRedis, hasEmail, hasPhone := mergeUserinfo(user, params, session)

	// 这里需要验证，
	if hasEmail {
		emailAddr := session.GetKeyValue("changeEmail")
		isOk := utils.IsValidEmail(emailAddr)
		if !isOk {
			// 通知用户信息不对
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"email is not correct",
				map[string]string{"field": emailAddr},
				session)
			return
		}

		// 先检查邮箱是否合法
		lst, _ := Globals.mongoCli.FindUserByEmail(emailAddr)
		if lst != nil && len(lst) > 0 {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"email is already used by user",
				map[string]string{"field": emailAddr},
				session)
			return
		}

		code := utils.GenerateCheckCode(5)
		SendEmailCode(session, emailAddr, code)
		session.SetKeyValue("code", code)
		session.SetStatus(model.UserStatusChangeInfo | model.UserStatusValidate)
	}

	// todo:
	if hasPhone {
		session.SetStatus(model.UserStatusChangeInfo | model.UserStatusValidate)
	}

	// 通过IP来检查用户地区
	if session.RemoteAddr != "" {
		fmt.Println("remote addr:", session.RemoteAddr)
		city, err := Globals.GeoHelper.GetCityByIP(session.RemoteAddr)
		if err == nil {
			setDataMongo["region"] = city.City
		}
	}

	if len(setDataMongo) > 0 {
		_, err := Globals.mongoCli.UpdateUserInfoPart(session.UserID, setDataMongo, nil)
		if err != nil {
			Globals.Logger.Fatal("update user, mongodb err", zap.Error(err))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
				"update user, meet database err",
				nil,
				session)
			return
		}
	}

	if len(setDataRedis) > 0 {
		// 更新redis
		Globals.redisCli.UpdateUserInfoPart(session.UserID, setDataRedis, nil)
	}

	// 更新内存

	Globals.Logger.Info("set user info", zap.Int64("user id", session.UserID))

	if hasEmail || hasPhone {
		SendBackUserOp(pbmodel.UserOperationType_SetUserInfo,
			user.GetUserInfo(),
			true, "waitcode", session)
	} else {
		SendBackUserOp(pbmodel.UserOperationType_SetUserInfo,
			user.GetUserInfo(),
			true, "ok", session)
	}

}

// 6 验证码检查
// 3种可能，注册，登录，实名认证
func handleUserVerification(msg *pbmodel.Msg, session *Session) {
	// 如果没有没有
	if !session.HasStatus(model.UserStatusValidate) {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTStage),
			"current user session not has a status of Verification ",
			nil,
			session)

		return
	}
	// 检查数据
	userOpMsg := msg.GetPlainMsg().GetUserOp()
	//fmt.Println(userOpMsg)
	params := userOpMsg.GetParams()
	if params == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"not has params field",
			map[string]string{"field": "params"},
			session)
		//fmt.Println("----------------------------")
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
			map[string]string{"field": "params.code", "value": codeUser},
			session)

		Globals.Logger.Info("verification code is not correct", zap.Int64("userid", session.UserID),
			zap.String("code", codeTemp), zap.String("user code", codeUser))
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
		onLoginSuccess(session, true)

	} else if session.HasStatus(model.UserStatusChangeInfo) {
		onChangeInfoSuccess(session)

	}

	// 取消掉当前的这个状态
	session.UnSetStatus(model.UserStatusValidate)
	session.RemoveKeyValue("code")
}

// 如果用户可以用，返回true
func checkUserUsable(userInfo *pbmodel.UserInfo, session *Session) bool {
	// 检查状态，是否被禁用了
	if userInfo.GetParams() != nil {
		status, ok := userInfo.GetParams()["status"]
		if ok {
			if status == "deleted" {
				sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTDeleted),
					"user is deleted, can't be used again",
					map[string]string{"field": "userInfo.Params.status"},
					session)
				return false
			} else if status == "disabled" {
				sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTDisabled),
					"user is disabled, contact with admin",
					map[string]string{"field": "userInfo.Params.status"},
					session)
				return false
			}
		}
	}

	return true
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
	params := userOpMsg.GetParams()
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
					return
				}

				lst, err := Globals.mongoCli.FindUserByEmail(userInfo.Email)
				if err != nil || len(lst) < 1 {
					//sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
					//	"email is not correct, can't find a user using it",
					//	map[string]string{"field": userInfo.Email},
					//	session)
					// 2025-05-17 改为如果使用邮件登录，直接注册
					handleUserRegister(msg, session)
					return
				}
				userInfoDb = lst[0]
				session.UserID = userInfoDb.UserId
				session.TempUserInfo = userInfoDb
				// 检查是否可以用
				ret := checkUserUsable(userInfoDb, session)
				if !ret {
					return
				}

				// 生成临时
				code := utils.GenerateCheckCode(5)

				session.SetKeyValue("code", code)
				// 使用后台将验证码发送给用户
				SendEmailCode(session, userInfoDb.Email, code)

				session.SetStatus(model.UserStatusLogin | model.UserStatusValidate)
				// 通知用户, 这里还不能通知用户真实数据
				//if userInfoDb.Params != nil {
				//	delete(userInfoDb.Params, "pwd")
				//}
				SendBackUserOp(pbmodel.UserOperationType_Login,
					userInfo,
					true, "waitcode", session)
				return
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
			Globals.Logger.Error("can't find user in mongodb",
				zap.Int64("user id ", userInfo.UserId),
				zap.Error(err))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
				"can't find user by id",
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
				"pwd is not set in your login data",
				map[string]string{"field": "params.pwd", "value": pwdUser},
				session)
			return
		}

		userInfoDb = lst[0]
		// 检查是否可以用
		ret := checkUserUsable(userInfoDb, session)
		if !ret {
			return
		}

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

	} else if loginMode == 3 {
		// 通知用户信息不对
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent),
			"phone login is not available",
			map[string]string{"field": userInfo.Phone},
			session)
		return
	}
	onLoginSuccess(session, true)

	return
}

// 使用指纹登录
func LoginWithPrint(session *Session, keyPrint int64, checkTokenData string, tm int64) bool {

	uid, keyEx, err := Globals.redisCli.LoadToken(keyPrint)
	if err != nil || keyEx == nil {
		Globals.Logger.Info("try login with keyprint, but not found: ", zap.Int64("keyPrint", keyPrint))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTKeyPrint), "key print is not exist", nil, session)
		return false
	}

	session.UserID = uid
	session.KeyEx = keyEx

	cipher, err := utils.DecodeBase64(checkTokenData)
	if err != nil {
		//fmt.Println("check data error", err)
		Globals.Logger.Info("decrypt base64 check data error: ", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error: "+err.Error(), nil, session)

		session.UserID = 0
		session.KeyEx = nil
		return false
	}

	tmData, err := decryptDataAuto(cipher, session)
	if err != nil {
		//fmt.Println("check data error", err)
		Globals.Logger.Info("decrypt hello check data error: ", zap.ByteString("data", tmData), zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check hello check data  decrypt error", nil, session)
		session.UserID = 0
		session.KeyEx = nil
		return false
	}

	tmStr := strconv.FormatInt(tm, 10)
	decryptedTm := string(tmData)

	fmt.Printf("receive tm = %v \n", decryptedTm)

	if tmStr != decryptedTm {
		Globals.Logger.Info("decrypt check data error: ", zap.Error(err))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error: "+string(checkTokenData), nil, session)
		session.UserID = 0
		session.KeyEx = nil
		return false
	}

	// 曾经交换过秘钥，但是当前显示未登录
	if uid == 0 {
		sendBackHelloMsg(session, "needlogin")
		return false
	}

	// 加载用户信息等
	onLoginSuccess(session, false)

	return true
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
func onLoginSuccess(session *Session, bSaveToken bool) {
	// 前面比对密码或者验证码成功了，这里应该绑定指纹和用户，7天
	if bSaveToken {
		if session.KeyEx != nil {
			Globals.redisCli.SaveToken(session.UserID, session.KeyEx)
		}
	}

	// 加载数据
	err := LoadUserLogin(session)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"onLoginSuccess() load user meet error",
			map[string]string{"error": err.Error()},
			session)
		return
	}
	session.updateTTL()

	// 通知用户免登录
	SendBackUserOp(pbmodel.UserOperationType_Login,
		session.GetUser().GetUserInfo(),
		true, "loginok", session)
	session.SetStatus(model.UserStatusOk)
}

// 验证码对了才保存到数据库
func onChangeInfoSuccess(session *Session) {
	setDataMongo := make(map[string]interface{})
	setDataRedis := make(map[string]interface{})

	user, ok := Globals.uc.GetUser(session.UserID)
	if !ok || user == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "user in not in cache, please login again", nil, session)
		return
	}

	emailAddr := session.GetKeyValue("changeEmail")
	if emailAddr != "" {
		setDataMongo["email"] = emailAddr
		setDataRedis["Email"] = emailAddr
		user.SetEmail(emailAddr)
		Globals.Logger.Info("set user email", zap.Int64("user id", session.UserID),
			zap.String("email", emailAddr))
		session.RemoveKeyValue("changeEmail")
	}

	phoneStr := session.GetKeyValue("changePhone")
	if phoneStr != "" {
		setDataMongo["phone"] = phoneStr
		setDataRedis["Phone"] = phoneStr
		user.SetPhone(phoneStr)
		Globals.Logger.Info("set user phone", zap.Int64("user id", session.UserID),
			zap.String("phone", phoneStr))
		session.RemoveKeyValue("changePhone")
	}

	if len(setDataMongo) > 0 {
		_, err := Globals.mongoCli.UpdateUserInfoPart(session.UserID, setDataMongo, nil)
		if err != nil {
			Globals.Logger.Fatal("update user, mongodb err", zap.Error(err))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
				"update user, meet database err",
				nil,
				session)
			return
		}
	} else {
		Globals.Logger.Fatal("update user, email and phone not found in session params")
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"update user, email and phone not found in session params",
			nil,
			session)
		return
	}

	if len(setDataRedis) > 0 {
		// 更新redis
		Globals.redisCli.UpdateUserInfoPart(session.UserID, setDataRedis, nil)
	}

	SendBackUserOp(pbmodel.UserOperationType_SetUserInfo,
		user.GetUserInfo(),
		true, "ok", session)

	session.SetStatus(model.UserStatusChangeInfo)
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
		Status:    status,
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
	if session.UserID < 100 {
		return true
	}
	return false
}
