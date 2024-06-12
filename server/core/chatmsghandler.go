package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"go.uber.org/zap"
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

	// 检查权限啊
	bFun, _ := checkFriendIsFan(fid, session.UserID)
	if !bFun {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotFriend), "not friend", map[string]string{
			"uid": strconv.FormatInt(fid, 64),
		}, session)
		sendBackChatMsgReply(false, "not friend", msgChat, session)
		return
	}

	// 社区模式
	if !Globals.Config.Server.FriendMode {
		// 要不要显示次数
	}

	bOk := checkFriendPermission(session.UserID, fid, bFun, model.PermissionMaskChat)
	if !bOk {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotPermission), "not have chat permission", map[string]string{
			"uid": strconv.FormatInt(fid, 64),
		}, session)
		sendBackChatMsgReply(false, "permission", msgChat, session)
		return
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
		sendBackChatMsgReply(false, "group id error", msgChat, session)
		return
	}

	_, b := group.HasMember(session.UserID)
	if !b {

		sendBackChatMsgReply(false, "you are not a group member", msgChat, session)
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

	sendBackChatMsgReply(true, "ok", msgChat, session)
	notifyGroupMembers(group.GroupId, msg)
}

// 6 消息应答：私聊消息需要确认
func handleChatReplyMsg(msg *pbmodel.Msg, session *Session) {
	// 直接更新数据库，并转发消息
	replyMsg := msg.GetPlainMsg().GetChatReply()
	if replyMsg != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "", nil, session)
		return
	}

	pk1 := db.ComputePk(replyMsg.UserId)
	pk2 := db.ComputePk(replyMsg.FromId)
	var err error

	if replyMsg.RecvOk > 0 && replyMsg.ReadOk > 0 {
		err = Globals.scyllaCli.SetPChatRecvReadReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
			replyMsg.MsgId, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
	} else if replyMsg.RecvOk > 0 {
		err = Globals.scyllaCli.SetPChatRecvReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
			replyMsg.MsgId, time.Now().UTC().UnixMilli())
	} else if replyMsg.ReadOk > 0 {
		err = Globals.scyllaCli.SetPChatReadReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
			replyMsg.MsgId, time.Now().UTC().UnixMilli())
	} else {
		return
	}

	// 没有找到合适的
	if err != nil {
		return
	}

	// 转发消息
	trySendMsgToUser(replyMsg.UserId, msg)

}

// 这里的查询数据，包括私聊和群聊
func onQueryChatData(queryMsg *pbmodel.MsgQuery, session *Session) {
	if queryMsg.ChatType == pbmodel.ChatType_ChatTypeP2P {
		onQueryP2PChatData(queryMsg, session)
	} else if queryMsg.ChatType == pbmodel.ChatType_ChatTypeGroup {
		onQueryGroupChatData(queryMsg, session)
	}
}

// 查询私聊
func onQueryP2PChatData(queryMsg *pbmodel.MsgQuery, session *Session) {
	pk := db.ComputePk(session.UserID)
	var lst []model.PChatDataStore
	var err error
	switch queryMsg.SynType {
	case pbmodel.SynType_SynTypeForward:
		lst, err = Globals.scyllaCli.FindPChatMsgForward(pk, session.UserID, queryMsg.LittleId, 100)
	case pbmodel.SynType_SynTypeBackward:
		lst, err = Globals.scyllaCli.FindPChatMsgBackwardFrom(pk, session.UserID, queryMsg.BigId, 100)
	case pbmodel.SynType_SynTypeBetween:
		lst, err = Globals.scyllaCli.FindPChatMsgForwardBetween(pk, session.UserID, queryMsg.LittleId, queryMsg.BigId, 100)
	}
	if err != nil {
		Globals.Logger.Error("find p2p chat msg error", zap.Error(err))
	}

	var littleId, bigId int64 = 0, 0
	var fromId int64
	var toId int64
	var chatDataList []*pbmodel.MsgChat = nil

	if lst != nil && len(lst) > 0 {
		chatDataList = make([]*pbmodel.MsgChat, len(lst))
		littleId = lst[0].Id
		bigId = lst[len(lst)-1].Id

		for index, item := range lst {
			// 自己发出去的
			if item.Io == model.ChatDataIOOut {
				fromId = item.Uid1
				toId = item.Uid2
			} else {
				fromId = item.Uid2
				toId = item.Uid1
			}

			data := pbmodel.MsgChat{
				MsgId:        item.Id,
				UserId:       session.UserID,
				FromId:       fromId,
				ToId:         toId,
				Tm:           item.Tm,
				DevId:        "",
				SendId:       item.Usid,
				MsgType:      pbmodel.ChatMsgType(item.Mt),
				Data:         item.Draf,
				Priority:     0,
				RefMessageId: item.Ref,
				Status:       0,
				SendReply:    item.Tm,
				RecvReply:    item.Tm1,
				ReadReply:    item.Tm2,
				EncType:      0,
				ChatType:     pbmodel.ChatType_ChatTypeP2P,
				SubMsgType:   0,
				KeyPrint:     0,
				Params:       nil,
			}
			chatDataList[index] = &data
		}
	}

	chatDataRet := pbmodel.MsgQueryResult{
		UserId:          session.UserID,
		GroupId:         0,
		BigId:           bigId,
		LittleId:        littleId,
		SynType:         queryMsg.SynType,
		Tm:              utils.GetTimeStamp(),
		ChatType:        pbmodel.ChatType_ChatTypeP2P,
		QueryType:       pbmodel.QueryDataType_QueryDataTypeChatData,
		ChatDataList:    chatDataList,
		ChatReplyList:   nil,
		FriendOpRetList: nil,
		GroupOpRetList:  nil,
		Result:          "ok",
		Params:          nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_CommonQueryRet{
			CommonQueryRet: &chatDataRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTQueryResult,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(msg)

}

func sendErrQueryChatDataResult(detail string, queryMsg *pbmodel.MsgQuery, session *Session) {
	chatDataRet := pbmodel.MsgQueryResult{
		UserId:          session.UserID,
		GroupId:         0,
		BigId:           0,
		LittleId:        0,
		SynType:         queryMsg.SynType,
		Tm:              utils.GetTimeStamp(),
		ChatType:        queryMsg.ChatType,
		QueryType:       queryMsg.QueryType,
		ChatDataList:    nil,
		ChatReplyList:   nil,
		FriendOpRetList: nil,
		GroupOpRetList:  nil,
		Result:          "fail",
		Detail:          detail,
		Params:          nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_CommonQueryRet{
			CommonQueryRet: &chatDataRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTQueryResult,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(msg)
}

// 查询群聊
func onQueryGroupChatData(queryMsg *pbmodel.MsgQuery, session *Session) {
	group, _ := findGroup(queryMsg.GroupId)
	if group == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group id is not correct", nil, session)
		return
	}

	_, bMember := group.HasMember(session.UserID)
	if !bMember {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "not a member", nil, session)
		sendErrQueryChatDataResult("not member of group", queryMsg, session)
		return
	}

	pk := db.ComputePk(queryMsg.GroupId)
	var lst []model.GChatDataStore
	var err error
	switch queryMsg.SynType {
	case pbmodel.SynType_SynTypeForward:
		lst, err = Globals.scyllaCli.FindGChatMsgForward(pk, group.GroupId, queryMsg.LittleId, 100)
	case pbmodel.SynType_SynTypeBackward:
		lst, err = Globals.scyllaCli.FindGChatMsgBackwardFrom(pk, group.GroupId, queryMsg.BigId, 100)
	case pbmodel.SynType_SynTypeBetween:
		lst, err = Globals.scyllaCli.FindGChatMsgForwardBetween(pk, session.UserID, queryMsg.LittleId, queryMsg.BigId, 100)
	}
	if err != nil {
		Globals.Logger.Error("find p2p chat msg error", zap.Error(err))
	}

	var littleId, bigId int64 = 0, 0
	var chatDataList []*pbmodel.MsgChat = nil

	if lst != nil && len(lst) > 0 {
		chatDataList = make([]*pbmodel.MsgChat, len(lst))
		littleId = lst[0].Id
		bigId = lst[len(lst)-1].Id

		for index, item := range lst {
			// 自己发出去的

			data := pbmodel.MsgChat{
				MsgId:        item.Id,
				UserId:       session.UserID,
				FromId:       item.Uid,
				ToId:         item.Gid,
				Tm:           item.Tm,
				DevId:        "",
				SendId:       item.Usid,
				MsgType:      pbmodel.ChatMsgType(item.Mt),
				Data:         item.Draf,
				Priority:     0,
				RefMessageId: item.Ref,
				Status:       0,
				SendReply:    item.Tm,
				EncType:      0,
				ChatType:     pbmodel.ChatType_ChatTypeGroup,
				SubMsgType:   0,
				KeyPrint:     0,
				Params:       nil,
			}
			chatDataList[index] = &data
		}
	}

	chatDataRet := pbmodel.MsgQueryResult{
		UserId:          session.UserID,
		GroupId:         0,
		BigId:           bigId,
		LittleId:        littleId,
		SynType:         queryMsg.SynType,
		Tm:              utils.GetTimeStamp(),
		ChatType:        pbmodel.ChatType_ChatTypeGroup,
		QueryType:       pbmodel.QueryDataType_QueryDataTypeChatData,
		ChatDataList:    chatDataList,
		ChatReplyList:   nil,
		FriendOpRetList: nil,
		GroupOpRetList:  nil,
		Params:          nil,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_CommonQueryRet{
			CommonQueryRet: &chatDataRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTQueryResult,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(msg)

}

func onQueryChatReply(queryMsg *pbmodel.MsgQuery, session *Session) {

}
