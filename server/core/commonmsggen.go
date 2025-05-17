package core

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func sendBackHeartMsg(session *Session) {
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
func sendBackErrorMsg(errCode int, detail string, params map[string]string, session *Session) {

	Globals.Logger.Error(detail, zap.Int64("user id", session.UserID), zap.Int("err code", errCode))
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
	session.SendMessage(msg)
}

// 应答服务端的hello消息
func sendBackHelloMsg(session *Session, stage string) {

	hello := pbmodel.MsgHello{
		ClientId: "",
		Version:  "v1.0",
		Platform: "windows",
		Stage:    stage, // "waitlogin"
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
	fmt.Println("tm string is: ", tmStr)
	fmt.Println("tm temp data is: ", checkData)

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

// 通知客户端秘钥交换完毕
func sendBackExchange4(session *Session) {

	exMsg := pbmodel.MsgKeyExchange{
		KeyPrint: session.KeyEx.SharedKeyPrint,
		RsaPrint: 0,
		Stage:    4,
		TempKey:  nil, // 这里发送共享密钥使用对称算法加密时间戳的密文
		PubKey:   nil,
		EncType:  session.KeyEx.EncType,
		Status:   "needlogin",
		Detail:   "check data ok, key print ok",
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_KeyEx{
			KeyEx: &exMsg,
		},
	}

	tm := utils.GetTimeStamp()

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
