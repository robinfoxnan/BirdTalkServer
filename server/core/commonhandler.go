package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

// 所有消息的统一入口，这里再开始分发
// 用户与群组消息一共是6大类，
func HandleCommonMsg(msg *pbmodel.Msg, session *Session) error {

	//fmt.Println(msg)
	keyPrint := msg.GetKeyPrint() // 加密指纹，明文传输需要先检查类型是否正确
	if keyPrint == 0 {
		msgPlain := msg.GetPlainMsg()
		if msgPlain == nil {
			Globals.Logger.Debug("receive wrong heart msg",
				zap.Int64("sid", session.Sid),
				zap.Int64("uid", session.UserID))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "msg.plainMsg is null", nil, session)
			return errors.New("plain msg, but msg.plainmsg pointer is null.")
		}
	} else { // 目前阶段不支持密文传输，后期如果支持密文，这里需要在这里解密，然后再分发
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "current version, msg.keyPrint must be 0", nil, session)
		return errors.New("plain msg, but msg.plainmsg pointer is null.")
	}

	switch msg.MsgType {
	case pbmodel.ComMsgType_MsgTUnused:
		//fmt.Println("recv unused type message!!")
		Globals.Logger.Info("recv unused type message!!",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		break

	case pbmodel.ComMsgType_MsgTHello: // 用于握手的消息
		handleHelloMsg(msg, session)
	case pbmodel.ComMsgType_MsgTHeartBeat: // 用于保持连接的心跳消息
		handleHeartMsg(msg, session)
	case pbmodel.ComMsgType_MsgTError: // 用于传递错误信息的消息
		handleErrorMsg(msg, session)
	case pbmodel.ComMsgType_MsgTKeyExchange: // DH密钥交换的消息
		handleKeyExchange(msg, session)
	case pbmodel.ComMsgType_MsgTChatMsg:
		handleChatMsg(msg, session)
	case pbmodel.ComMsgType_MsgTChatReply:
		handleChatReplyMsg(msg, session)
	case pbmodel.ComMsgType_MsgTQuery:
		handleCommonQuery(msg, session)
	case pbmodel.ComMsgType_MsgTQueryResult:
		//fmt.Println("Receive error type of message ComMsgType_MsgTChatQueryResult:")
		Globals.Logger.Info("Receive error type of message ComMsgType_MsgTChatQueryResult:",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
	case pbmodel.ComMsgType_MsgTUpload:
		handleFileUpload(msg, session)
	case pbmodel.ComMsgType_MsgTDownload: //下载文件的消息，文件操作分为带内和带外，这里是小文件可以这样操作
		handleFileDownload(msg, session)
	case pbmodel.ComMsgType_MsgTUploadReply:
		//fmt.Println("Receive error type of messageComMsgType_MsgTUploadReply:")
		Globals.Logger.Info("Receive error type of messageComMsgType_MsgTUploadReply:",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
	case pbmodel.ComMsgType_MsgTDownloadReply:
		//fmt.Println("Receive error type of message ComMsgType_MsgTDownloadReply:")
		Globals.Logger.Info("Receive error type of message ComMsgType_MsgTDownloadReply:",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
	case pbmodel.ComMsgType_MsgTUserOp: // 所有用户相关操作的消息
		handleUserOp(msg, session)
	case pbmodel.ComMsgType_MsgTUserOpRet:
		//fmt.Println("Receive error type of message ComMsgType_MsgTUserOpRet:")
		Globals.Logger.Info("Receive error type of message ComMsgType_MsgTUserOpRet:",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
	case pbmodel.ComMsgType_MsgTFriendOp:
		handleFriendOp(msg, session)
	case pbmodel.ComMsgType_MsgTFriendOpRet:
		handleFriendOpRet(msg, session)
	case pbmodel.ComMsgType_MsgTGroupOp: // 所有群组相关的操作
		handleGroupOp(msg, session)
	case pbmodel.ComMsgType_MsgTGroupOpRet:
		handleGroupOpRet(msg, session)
	case pbmodel.ComMsgType_MsgTOther: // 转发给其他的扩展模块的
		handleOther(msg, session)
	}

	return nil
}

// 1) 握手消息
// 检查是否需要协商秘钥；检查keyprint是否存在；检查协议版本号；
func handleHelloMsg(msg *pbmodel.Msg, session *Session) {

	session.Ver = int(msg.GetVersion())   // 协议版本号
	ok := checkProtoVersion(msg, session) // 错误会自动应答
	if !ok {
		return
	}

	// hello 不能使用密文传输，
	msgHello := msg.GetPlainMsg().GetHello()
	if msgHello == nil {
		Globals.Logger.Debug("receive wrong hello msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "hello msg is null", nil, session)
		return
	}

	session.DeviceID = msgHello.GetClientId()     // 设备唯一编号
	session.Platf = msgHello.GetPlatform()        // "web","android"
	session.Params["Stage"] = msgHello.GetStage() // "clienthello"

	params := msgHello.GetParams()
	session.Params["ClientVersion"] = msgHello.GetVersion()
	str, b := params["Lang"]
	if b {
		session.Params["Lang"] = str
	}
	str, b = params["CountryCode"]
	if b {
		session.Params["CountryCode"] = str
	}

	str, b = params["CodeType"]
	if b {
		session.Params["CodeType"] = str
	}

	checkTokenData, b := params["checkTokenData"]

	//fmt.Println(&session)
	Globals.Logger.Debug("handle client hello msg", zap.Any("session", session))

	if msgHello.GetKeyPrint() != 0 {
		fmt.Println("key print=", msgHello.GetKeyPrint())
		ok = LoginWithPrint(session, msgHello.GetKeyPrint(), checkTokenData, msg.GetTm())
		if ok {
			// pass
		} else {

			session.SetStatus(model.UserWaitLogin)
		}
	} else {
		sendBackHelloMsg(session, "waitlogin")
		session.SetStatus(model.UserWaitLogin)
	}

}

// 2) 心跳消息
func handleHeartMsg(msg *pbmodel.Msg, session *Session) {

	ping := msg.GetPlainMsg().GetHeartBeat()
	if ping == nil {
		Globals.Logger.Debug("receive wrong heart msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "heart msg ping is null", nil, session)
		return
	}

	tmStr := utils.TmToLocalString(ping.Tm)
	//fmt.Printf("tm=%s, userid=%d \n", tmStr, ping.UserId)
	Globals.Logger.Debug("receive heart beat",
		zap.Int64("sid", session.Sid),
		zap.Int64("uid", session.UserID),
		zap.String("tm", tmStr))

	sendBackHeartMsg(session)
}

// 3) 错误消息
func handleErrorMsg(msg *pbmodel.Msg, session *Session) {
	errMsg := msg.GetPlainMsg().GetErrorMsg()
	if errMsg == nil {
		Globals.Logger.Debug("receive wrong error msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "error msg error part is null", nil, session)
		return
	}
	//fmt.Printf("client error = %d, detail = %s \n", errMsg.Code, errMsg.Detail)
	Globals.Logger.Info("client report error",
		zap.Int64("sid", session.Sid),
		zap.Int64("uid", session.UserID),
		zap.Int32("code", errMsg.Code),
		zap.String("detail", errMsg.Detail),
	)
}

// 4) 交换秘钥消息
func handleKeyExchange(msg *pbmodel.Msg, session *Session) {
	ok := checkProtoVersion(msg, session) // 错误会自动应答
	if !ok {
		return
	}

	exMsg := msg.GetPlainMsg().GetKeyEx()
	//fmt.Printf("收到KeyExchange %v", exMsg)
	if exMsg == nil {
		Globals.Logger.Debug("receive wrong key exchange msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "key exchange part is null", nil, session)
		return
	}

	stage := exMsg.GetStage()
	if stage == 1 { // 开始阶段
		encType := exMsg.EncType
		if len(encType) == 0 {
			encType = "AES-CTR"
		}
		exChangeData, err := utils.NewKeyExchange(encType)
		if err != nil {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "create key pair error: "+err.Error(), nil, session)
			return
		}
		session.KeyEx = exChangeData // 设置

		publicKeyRemote, err := decodeRemotePublicKey(exMsg, session) // 出错直接应答
		if err != nil {
			return
		}

		exChangeData.Stage = utils.KeyExchangeStagePublicKey
		// 计算共享密钥，并计算指纹
		_, err = exChangeData.GenShareKey(publicKeyRemote)
		if err != nil {
			//fmt.Println("calculate share key error", err)
			Globals.Logger.Info("calculate share key error: ", zap.Error(err))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTPublicKey), "calculate share key error: "+err.Error(), nil, session)
			return
		}
		sendBackExchange2(session)
		session.SetStatus(model.UserStatusExchange)

	} else if stage == 3 {
		cipher := exMsg.GetTempKey()
		tmData, err := decryptDataAuto(cipher, session)
		if err != nil {
			//fmt.Println("check data error", err)
			Globals.Logger.Info("decrypt check data error: ", zap.Error(err))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error: "+err.Error(), nil, session)
		}

		remoteKeyPrint := exMsg.GetKeyPrint()

		tmStr := string(tmData)
		tm := strconv.FormatInt(msg.GetTm(), 10)
		// 如果指纹不一样，或者
		if remoteKeyPrint != session.KeyEx.SharedKeyPrint {
			//fmt.Println("check data error:", " key print or data not same")
			//fmt.Printf("tm = %v, tmstr= %v \n", tm, tmStr)
			Globals.Logger.Info("check data error: ", zap.Int64("print", session.KeyEx.SharedKeyPrint), zap.Int64("check print", remoteKeyPrint))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check print error: key print not same", nil, session)
			return
		}

		if tmStr != tm {
			//fmt.Println("check data error:", " key print or data not same")
			//fmt.Printf("tm = %v, tmstr= %v \n", tm, tmStr)
			Globals.Logger.Info("check data error: ", zap.String("tm", tm), zap.String("check tm", tmStr))
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error:  data not same", nil, session)
			return
		}

		// 保存指纹到redis
		//fmt.Printf("check data ok! tm = %v, tmstr= %v \n", tm, tmStr)
		Globals.Logger.Debug("check data ok!", zap.String("tm", tm), zap.String("check tm", tmStr))
		err = Globals.redisCli.SaveToken(0, session.KeyEx)
		if err != nil {
			Globals.Logger.Error("exchange stage3 save token err.")
		}

		sendBackExchange4(session)
		session.SetStatus(model.UserWaitLogin)
	} else {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "", nil, session)
		return
	}
	return
}
