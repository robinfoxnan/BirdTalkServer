package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"strconv"
	"time"
)

// 5 聊天的信息，这里有3类，给自己的，私聊一对一的，群聊的；
func handleChatMsg(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	if session.UserID != msgChat.FromId {
		sendBackChatMsgReply(true, "from id is not same as session user id", msgChat, session)
		return
	}

	msgChat.MsgId = Globals.snow.GenerateID()
	// 发给自己的消息
	if msgChat.FromId == msgChat.ToId {
		onSelfChatMessage(msg, session)
		return
	}

	// 私聊消息
	if msgChat.ChatType == pbmodel.ChatType_ChatTypeP2P {
		onP2pChatMessage(msg, session)
		return
	} else if msgChat.ChatType == pbmodel.ChatType_ChatTypeGroup {
		onGroupChatMessage(msg, session)
		return
	} else {
		// 消息类型错误
	}

}

// 发给自己的消息
func onSelfChatMessage(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	// 保存消息
	pk1 := db.ComputePk(msgChat.FromId)
	msgData := model.PChatDataStore{
		Pk:    pk1,
		Uid1:  msgChat.FromId,
		Uid2:  msgChat.ToId,
		Id:    msgChat.MsgId,
		Usid:  msgChat.SendId,
		Tm:    time.Now().UTC().UnixMilli(),
		Tm1:   0,
		Tm2:   0,
		Io:    0,
		St:    0,
		Ct:    0,
		Mt:    int8(msgChat.MsgType.Number()),
		Print: 0,
		Ref:   msgChat.RefMessageId,
		Draf:  msgChat.Data,
	}
	err := Globals.scyllaCli.SavePChatSelfData(&msgData)
	if err != nil {
		sendBackChatMsgReply(false, "save data err", msgChat, session)
		return
	}

	trySendMsgToMe(msgChat.ToId, msg, session)
	sendBackChatMsgReply(true, "ok", msgChat, session)
}

// 私聊
func onP2pChatMessage(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	// 检查权限啊
	fid := msgChat.ToId
	if Globals.Config.Server.FriendMode {
		bFun, _ := checkFriendIsFan(fid, session.UserID)
		if !bFun {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotFriend), "not friend", map[string]string{
				"uid": strconv.FormatInt(fid, 64),
			}, session)
			sendBackChatMsgReply(false, "not friend", msgChat, session)
			return
		}
		bOk := checkFriendPermission(session.UserID, fid, bFun, model.PermissionMaskChat)
		if !bOk {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "not have chat permission", map[string]string{
				"uid": strconv.FormatInt(fid, 64),
			}, session)
			sendBackChatMsgReply(false, "permission", msgChat, session)
			return
		}
	}

	// 保存消息
	pk1 := db.ComputePk(msgChat.FromId)
	pk2 := db.ComputePk(msgChat.ToId)
	tm := time.Now().UTC().UnixMilli()
	msgData := model.PChatDataStore{
		Pk:    pk1,
		Uid1:  msgChat.FromId,
		Uid2:  msgChat.ToId,
		Id:    msgChat.MsgId,
		Usid:  msgChat.SendId,
		Tm:    tm,
		Tm1:   0,
		Tm2:   0,
		Io:    0,
		St:    0,
		Ct:    0,
		Mt:    int8(msgChat.MsgType.Number()),
		Print: 0,
		Ref:   msgChat.RefMessageId,
		Draf:  msgChat.Data,
	}

	err := Globals.scyllaCli.SavePChatData(&msgData, pk2)
	if err != nil {
		sendBackChatMsgReply(false, "error when save data", msgChat, session)
		return
	}

	msgChat.Tm = tm
	sendBackChatMsgReply(true, "save data ok", msgChat, session)
	trySendMsgToUser(msgChat.ToId, msg)
	trySendMsgToMe(msgChat.ToId, msg, session)

}

// 应答保存完毕
func sendBackChatMsgReply(ok bool, detail string, msgChat *pbmodel.MsgChat, session *Session) {

	msgChatReply := pbmodel.MsgChatReply{
		MsgId:    msgChat.MsgId,
		SendId:   msgChat.SendId,
		SendOk:   msgChat.Tm,
		RecvOk:   0,
		ReadOk:   0,
		ExtraMsg: "save ok",
		UserId:   session.UserID, //
		FromId:   0,              // 从服务器处得到的应答
		Params:   nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_ChatReply{

			ChatReply: &msgChatReply,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTChatReply,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(msg)
}

// 群消息
func onGroupChatMessage(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	// 检查权限啊
	group, err := findGroup(msgChat.ToId)
	if group == nil {
		sendBackChatMsgReply(true, "group id error", msgChat, session)
		return
	}

	_, b := group.HasMember(session.UserID)
	if !b {

		sendBackChatMsgReply(true, "you are not a group member", msgChat, session)
		return
	}

	// 保存消息
	pk1 := db.ComputePk(msgChat.ToId)
	msgData := model.GChatDataStore{
		Pk:    pk1,
		Gid:   msgChat.ToId,
		Uid:   msgChat.FromId,
		Id:    msgChat.MsgId,
		Usid:  msgChat.SendId,
		Tm:    time.Now().UTC().UnixMilli(),
		St:    0,
		Ct:    0,
		Mt:    int8(msgChat.MsgType.Number()),
		Print: 0,
		Ref:   msgChat.RefMessageId,
		Draf:  msgChat.Data,
	}
	err = Globals.scyllaCli.SaveGChatData(&msgData)
	if err != nil {
		sendBackChatMsgReply(false, "save data error", msgChat, session)
		return
	}

	notifyGroupMembers(group.GroupId, msg)
}

// 6 消息应答：私聊消息需要确认
func handleChatReplyMsg(msg *pbmodel.Msg, session *Session) {

}

func onQueryChatData(queryMsg *pbmodel.MsgQuery, session *Session) {

}

func onQueryChatReply(queryMsg *pbmodel.MsgQuery, session *Session) {

}
