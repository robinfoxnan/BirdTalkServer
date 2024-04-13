package core

import "birdtalk/server/pbmodel"

// 5) 交换秘钥消息
func handleChatMsg(msg *pbmodel.Msg, session *Session) {

}

// 6 消息应答：私聊消息需要确认
func handleChatReplyMsg(msg *pbmodel.Msg, session *Session) {

}

// 7 这里的查询包括：私聊消息同步，群消息同步，回执查询；好友请求查询，群用户操作消息
func handleCommonQuery(msg *pbmodel.Msg, session *Session) {

}
