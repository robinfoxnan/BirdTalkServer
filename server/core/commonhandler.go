package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 所有消息的统一入口，这里再开始分发
// 用户与群组消息一共是6大类，
func HandleCommonMsg(msg *pbmodel.Msg, session *Session) error {
	fmt.Println(msg)
	switch msg.MsgType {
	case pbmodel.ComMsgType_MsgTUnused:
		fmt.Println("recv unused type message!!")
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
		fmt.Println("Receive error type of message ComMsgType_MsgTChatQueryResult:")
	case pbmodel.ComMsgType_MsgTUpload:
		handleFileUpload(msg, session)
	case pbmodel.ComMsgType_MsgTDownload: //下载文件的消息，文件操作分为带内和带外，这里是小文件可以这样操作
		handleFileDownload(msg, session)
	case pbmodel.ComMsgType_MsgTUploadReply:
		fmt.Println("Receive error type of messageComMsgType_MsgTUploadReply:")
	case pbmodel.ComMsgType_MsgTDownloadReply:
		fmt.Println("Receive error type of message ComMsgType_MsgTDownloadReply:")
	case pbmodel.ComMsgType_MsgTUserOp: // 所有用户相关操作的消息
		handleUserOp(msg, session)
	case pbmodel.ComMsgType_MsgTUserOpRet:
		fmt.Println("Receive error type of message ComMsgType_MsgTUserOpRet:")
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

// 检查当前客户端的协议版本号；
func checkProtoVersion(msg *pbmodel.Msg, session *Session) bool {
	ver := int(msg.GetVersion()) // 协议版本号
	if ver != ProtocolVersion {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTVersion), "version is too high", nil, session)
		return false
	}
	return true
}

// 检查秘钥指纹是否存在
func checkKeyPrint(key int64) (int, string) {
	return int(pbmodel.ErrorMsgType_ErrTKeyPrint), "should be 0"
}

// 检查是否需要重定向
func checkNeedRedirect(session *Session) (int, string) {
	//return int(pbmodel.ErrorMsgType_ErrTRedirect), "should redirect"

	return int(pbmodel.ErrorMsgType_ErrTNone), ""
}

func sendBackErrorMsg(errCode int, detail string, params map[string]string, session *Session) {

	// 创建一个 MsgError 消息
	errorMsg := pbmodel.MsgError{
		Code:   int32(errCode),
		Detail: detail,
		Params: params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_ErrorMsg{
			ErrorMsg: &errorMsg,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTError,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	//// 序列化消息
	//data, err := proto.Marshal(&msg)
	//if err != nil {
	//	fmt.Println("Error marshaling message:", err)
	//	return
	//}
	//
	//// 打印序列化后的消息
	//fmt.Println("Serialized message:", data)

	session.SendMessage(msg)
}

// 应答服务端的hello消息
func sendBackHelloMsg(session *Session) {

	hello := pbmodel.MsgHello{
		ClientId: "",
		Version:  "v1.0",
		Platform: "windows",
		Stage:    "waitlogin",
		KeyPrint: 0,
		RsaPrint: 0,
		Params:   nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_Hello{
			Hello: &hello,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTHello,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	session.SendMessage(&msg)
}

// 1) 握手消息
// 检查是否需要协商秘钥；检查keyprint是否存在；检查协议版本号；
func handleHelloMsg(msg *pbmodel.Msg, session *Session) {

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

	msgHello := msg.GetPlainMsg().GetHello()
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

	// 检查是否需要重定向
	//errCode, errStr = checkNeedRedirect(session)
	//if errCode != int(pbmodel.ErrorMsgType_ErrTNone) {
	//	redirectParam := map[string]string{
	//		"host": "127.0.0.1:8080",
	//	}
	//	sendBackErrorMsg(errCode, errStr, redirectParam, session)
	//	return
	//}

	session.Status = model.UserWaitLogin
	fmt.Println(&session)
	sendBackHelloMsg(session)
}

// 2) 心跳消息
func handleHeartMsg(heartHello *pbmodel.Msg, session *Session) {
	ping := heartHello.GetPlainMsg().GetHeartBeat()
	tmStr := utils.TmToLocalString(ping.Tm)
	fmt.Printf("tm=%s, userid=%d \n", tmStr, ping.UserId)

	heart := pbmodel.MsgHeartBeat{
		Tm:     utils.GetTimeStamp(),
		UserId: session.UserID,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_HeartBeat{
			HeartBeat: &heart,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTHello,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	session.SendMessage(msg)
}

// 自动选择算法执行加密
func encryptDataAuto(data []byte, session *Session) ([]byte, error) {
	encType := strings.ToLower(session.KeyEx.EncType)
	switch encType {
	case "chacha20":
		return utils.EncryptChaCha20(data, session.KeyEx.SharedKeyHash)
	case "aes-ctr":
		return utils.EncryptAES_CTR(data, session.KeyEx.SharedKeyHash)
	case "twofish":
		return utils.EncryptTwofish(data, session.KeyEx.SharedKeyHash)
	}

	return nil, errors.New("not supported encrypt algorithm")
}

func decryptDataAuto(data []byte, session *Session) ([]byte, error) {
	encType := strings.ToLower(session.KeyEx.EncType)
	switch encType {
	case "chacha20":
		return utils.DecryptChaCha20(data, session.KeyEx.SharedKeyHash)
	case "aes-ctr":
		return utils.DecryptAES_CTR(data, session.KeyEx.SharedKeyHash)
	case "twofish":
		return utils.DecryptTwofish(data, session.KeyEx.SharedKeyHash)
	}

	return nil, errors.New("not supported encrypt algorithm")
}

// 回复秘钥交换的阶段2
func sendBackExchange2(session *Session) {

	tm := utils.GetTimeStamp()
	tmStr := strconv.FormatInt(tm, 10)
	checkData, err := encryptDataAuto([]byte(tmStr), session)
	if err != nil {
		return
	}

	exMsg := pbmodel.MsgKeyExchange{
		KeyPrint: session.KeyEx.SharedKeyPrint,
		RsaPrint: 0,
		Stage:    2,
		TempKey:  checkData, // 这里发送共享密钥使用对称算法加密时间戳的密文
		PubKey:   session.KeyEx.PublicKey,
		EncType:  session.KeyEx.EncType,
		Status:   "ready",
		Detail:   "reply public key and key print",
	}

	fmt.Println("local public key is:", string(session.KeyEx.PublicKey))
	fmt.Println("share key is: ", session.KeyEx.SharedKeyHash)
	fmt.Println("share key print is: ", session.KeyEx.SharedKeyPrint)

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_KeyEx{
			KeyEx: &exMsg,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       tm,
		MsgType:  pbmodel.ComMsgType_MsgTKeyExchange,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(msg)
}

// 发送验证结果
func sendBackExchange4(session *Session) {

}

// 3) 错误消息
func handleErrorMsg(msg *pbmodel.Msg, session *Session) {
	errMsg := msg.GetPlainMsg().GetErrorMsg()
	fmt.Printf("client error = %d, detail = %s \n", errMsg.Code, errMsg.Detail)
}

// 解码对方的公钥，有可能是加密，根据rasPrint区分
func decodeRemotePublicKey(exMsg *pbmodel.MsgKeyExchange, session *Session) ([]byte, error) {
	rsaPrint := exMsg.GetRsaPrint()
	if rsaPrint == 0 {

		publicKey := exMsg.GetPubKey()
		fmt.Println("remote public key=", string(publicKey))
		if len(publicKey) < 65 {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTPublicKey), "", nil, session)
			return nil, errors.New("public key len less bytes")
		}
		return publicKey, nil
	}
	// 使用RAS解码对称密钥，使用RSA解码对方公钥；

	sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTRsaPrint), "", nil, session)
	return nil, errors.New("decode remote public key error")
}

// 4) 交换秘钥消息
func handleKeyExchange(msg *pbmodel.Msg, session *Session) {
	ok := checkProtoVersion(msg, session) // 错误会自动应答
	if !ok {
		return
	}

	exMsg := msg.GetPlainMsg().GetKeyEx()
	//fmt.Printf("收到KeyExchange %v", exMsg)

	stage := exMsg.GetStage()
	if stage == 1 { // 开始阶段
		encType := exMsg.EncType
		if len(encType) == 0 {
			encType = "chacha20"
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
			fmt.Println("calculate share key error", err)
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTPublicKey), "calculate share key error: "+err.Error(), nil, session)
			return
		}
		sendBackExchange2(session)

	} else if stage == 3 {
		cipher := exMsg.GetTempKey()
		tmData, err := decryptDataAuto(cipher, session)
		if err != nil {
			fmt.Println("check data error", err)
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error: "+err.Error(), nil, session)
		}

		remoteKeyPrint := exMsg.GetKeyPrint()

		tmStr := string(tmData)
		tm := strconv.FormatInt(msg.GetTm(), 10)
		// 如果指纹不一样，或者
		if remoteKeyPrint != session.KeyEx.SharedKeyPrint || tmStr != tm {
			fmt.Println("check data error:", " key print or data not same")
			fmt.Printf("tm = %v, tmstr= %v \n", tm, tmStr)
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTCheckData), "check data error: key print or data not same", nil, session)
		}

		// 保存指纹到redis

		sendBackExchange4(session)
	} else {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "", nil, session)
		return
	}
	return
}
