package core

import (
	"birdtalk/server/pbmodel"
	"fmt"
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

// 1) 握手消息
func handleHelloMsg(msg *pbmodel.Msg, session *Session) {

}

// 2) 心跳消息
func handleHeartMsg(msg *pbmodel.Msg, session *Session) {

}

// 3) 错误消息
func handleErrorMsg(msg *pbmodel.Msg, session *Session) {

}

// 4) 交换秘钥消息
func handleKeyExchange(msg *pbmodel.Msg, session *Session) {

}
