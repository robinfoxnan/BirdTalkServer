package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// String 方法实现 fmt.Stringer 接口，用于格式化输出 MsgChat 结构体信息
func msg2String(m *pbmodel.MsgChat) string {
	if m == nil {
		return "MsgChat(nil)"
	}

	// 构建字符串缓冲区，提高拼接效率
	var buf strings.Builder

	// 基础字段
	buf.WriteString(fmt.Sprintf("MsgChat{msgId: %d, userId: %d, fromId: %d, toId: %d, tm: %d (time: %s), ",
		m.MsgId,
		m.UserId,
		m.FromId,
		m.ToId,
		m.Tm,
		time.Unix(m.Tm/1000, 0).Format("2006-01-02 15:04:05"), // 假设时间戳是毫秒级
	))

	// 设备和发送相关字段
	buf.WriteString(fmt.Sprintf("devId: %q, sendId: %d, ", m.DevId, m.SendId))

	// 枚举类型字段（假设这些枚举都实现了 String() 方法）
	buf.WriteString(fmt.Sprintf("msgType: %s, ", m.MsgType))

	// 消息内容（处理二进制数据）
	if m.MsgType != pbmodel.ChatMsgType_VOICE {
		dataStr := "[]byte"
		if len(m.Data) > 0 {
			// 尝试将 bytes 解析为字符串，否则显示长度
			if str := string(m.Data); isPrintable(str) {
				dataStr = fmt.Sprintf("%q", str)
			} else {
				dataStr = fmt.Sprintf("bytes(len=%d)", len(m.Data))
			}
		}
		buf.WriteString(fmt.Sprintf("data: %s, ", dataStr))
	}

	// 优先级和引用消息
	buf.WriteString(fmt.Sprintf("priority: %s, refMessageId: %d, ", m.Priority, m.RefMessageId))

	// 状态相关字段
	buf.WriteString(fmt.Sprintf("status: %s, sendReply: %d, recvReply: %d, readReply: %d, ",
		m.Status, m.SendReply, m.RecvReply, m.ReadReply,
	))

	// 加密和聊天类型
	buf.WriteString(fmt.Sprintf("encType: %s, chatType: %s, subMsgType: %d, keyPrint: %d, ",
		m.EncType, m.ChatType, m.SubMsgType, m.KeyPrint,
	))

	// 处理 params 映射
	buf.WriteString("params: {")
	first := true
	for k, v := range m.Params {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%q: %q", k, v))
		first = false
	}
	buf.WriteString("}}")

	return buf.String()
}

// 辅助函数：判断字符串是否可打印
func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// 5 聊天的信息，这里有3类，给自己的，私聊一对一的，群聊的；
func handleChatMsg(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	if session.UserID != msgChat.FromId {
		sendBackChatMsgReply(false, "from id is not same as session user id", msgChat, session, 0)
		return
	}

	msgChat.MsgId = Globals.snow.GenerateID()
	// 发给自己的消息
	if msgChat.FromId == msgChat.ToId {
		onSelfChatMessage(msg, session)
		return
	}

	// 单独的打印函数
	txtInfo := fmt.Sprintf("receive chat msg(%d) from %d -> %d", msgChat.SendId, msgChat.FromId, msgChat.ToId)
	Globals.Logger.Debug(txtInfo, zap.String("msg", msg2String(msgChat)))

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
		sendBackChatMsgReply(false, "save data err", msgChat, session, 0)
		return
	}

	trySendMsgToMe(msgChat.ToId, msg, session)
	sendBackChatMsgReply(true, "send to self", msgChat, session, 0)
}

// 私聊
func onP2pChatMessage(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()

	if msgChat.FromId != session.UserID {
		msgChat.FromId = session.UserID
	}

	// 检查权限啊
	fid := msgChat.ToId
	pk1 := db.ComputePk(msgChat.FromId)
	pk2 := db.ComputePk(msgChat.ToId)
	tm := time.Now().UTC().UnixMilli()

	// 单独处理撤回消息
	if msgChat.MsgType == pbmodel.ChatMsgType_DELETE {

		err := Globals.scyllaCli.SetPChatMsgDeleted(pk1, pk2, session.UserID, msgChat.ToId, msgChat.RefMessageId)
		if err != nil {
			sendBackChatMsgReply(false, "ref error", msgChat, session, 0)
			return
		}
	} else {
		// 检查权限啊
		bFun, _ := checkFriendIsFan(fid, session.UserID)
		if !bFun {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotFriend), "not friend", map[string]string{
				"uid": strconv.FormatInt(fid, 10),
			}, session)
			sendBackChatMsgReply(false, "not friend", msgChat, session, 0)
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
			sendBackChatMsgReply(false, "permission", msgChat, session, 0)
			return
		}
	}

	// 保存消息

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
		sendBackChatMsgReply(false, "error when save data", msgChat, session, 0)
		return
	}

	msgChat.Tm = tm
	sendBackChatMsgReply(true, "save data ok", msgChat, session, 0)
	Globals.Logger.Debug("send back reply 'save data ok'")

	// 先放到队列中
	tryPushToUserMsgCache(msgChat.ToId, msgChat.MsgId, msg)
	// 转发用户
	trySendMsgToUser(msgChat.ToId, msg)
	trySendMsgToMe(session.UserID, msg, session)

}

// 应答保存完毕
func sendBackChatMsgReply(ok bool, detail string, msgChat *pbmodel.MsgChat, session *Session, gid int64) {

	tm := utils.GetTimeStamp()
	params := map[string]string{
		"gid":    fmt.Sprintf("%d", gid),
		"detail": detail,
	}

	var extraMsg string
	if ok {
		extraMsg = "ok"
	} else {
		extraMsg = "fail"
	}

	msgChatReply := pbmodel.MsgChatReply{
		MsgId:    msgChat.MsgId,
		SendId:   msgChat.SendId,
		SendOk:   tm,
		RecvOk:   0,
		ReadOk:   0,
		ExtraMsg: extraMsg,
		UserId:   session.UserID, // 给这个用户应答
		FromId:   msgChat.ToId,   // 从服务器处得到的应答
		Params:   params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_ChatReply{

			ChatReply: &msgChatReply,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       tm,
		MsgType:  pbmodel.ComMsgType_MsgTChatReply,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	session.SendMessage(&msg)
}

// 群消息
func onGroupChatMessage(msg *pbmodel.Msg, session *Session) {
	msgPlain := msg.GetPlainMsg()
	msgChat := msgPlain.GetChatData()
	// 检查权限啊
	group, err := findGroup(msgChat.ToId)
	if group == nil {
		sendBackChatMsgReply(false, "group id error", msgChat, session, msgChat.ToId)
		return
	}

	_, b := group.HasMember(session.UserID)
	if !b {

		sendBackChatMsgReply(false, "you are not a group member", msgChat, session, msgChat.ToId)
		return
	}

	// 保存消息
	pk1 := db.ComputePk(msgChat.ToId)

	// 单独处理删除消息
	if msgChat.MsgType == pbmodel.ChatMsgType_DELETE {

		err = Globals.scyllaCli.SetGChatMsgDeleted(pk1, msgChat.ToId, msgChat.RefMessageId)
		if err != nil {
			sendBackChatMsgReply(false, "ref error", msgChat, session, msgChat.ToId)
			return
		}
	}

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
		sendBackChatMsgReply(false, "save data error", msgChat, session, msgChat.ToId)
		return
	}

	sendBackChatMsgReply(true, "ok", msgChat, session, msgChat.ToId)
	sendToGroupMembersExceptMe(group.GroupId, msg, session)
}

// 6 消息应答：私聊消息需要确认
func handleChatReplyMsg(msg *pbmodel.Msg, session *Session) {
	// 直接更新数据库，并转发消息
	replyMsg := msg.GetPlainMsg().GetChatReply()
	if replyMsg == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "reply is null", nil, session)
		return
	}
	Globals.Logger.Debug("Recv User ChatMsgReply:", zap.Int64("from", replyMsg.FromId),
		zap.Int64("to", replyMsg.UserId), zap.Any("msg", replyMsg))

	replyMsg.Params["gid"] = "0"
	batchIds, bOk := replyMsg.Params["batch"]
	if bOk {
		handleChatReplyMsgBatch(msg, session, batchIds)
		return
	}

	pk1 := db.ComputePk(replyMsg.UserId)
	pk2 := db.ComputePk(replyMsg.FromId)
	var err error

	user := session.GetUser()
	if user == nil {
		Globals.Logger.Error("user is nil, return here", zap.Int64("from", replyMsg.FromId))
		return
	}

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
		Globals.Logger.Debug("Recv User ChatMsgReply:", zap.String("error", "recv and read time is 0"))
	}

	// 收到回执后，从发送列表中删除
	user.PopMsgInCache(replyMsg.MsgId)

	// 没有找到合适的
	if err != nil {
		Globals.Logger.Error("Recv User ChatMsgReply:", zap.Error(err))
		return
	}

	// 转发消息
	trySendMsgToUser(replyMsg.UserId, msg)

}

// SplitToInt64Array 将 "1, 2, 3, 4" 格式字符串转为 []int64
// 若输入为空、格式错误（含非数字），返回错误信息
func SplitToInt64Array(s string) ([]int64, error) {
	// 1. 处理空输入
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil, errors.New("输入字符串不能为空")
	}

	// 2. 分割字符串：处理逗号+任意空格（如 "1,2" "1,  2" 均兼容）
	strSlice := strings.Split(trimmed, ",")
	int64Slice := make([]int64, 0, len(strSlice)) // 预分配容量，提升性能

	// 3. 遍历转换：逐个字符串转 int64（base=10，bitSize=64 对应 long）
	for _, str := range strSlice {
		str = strings.TrimSpace(str) // 去除每个元素的前后空格（如 " 3 " → "3"）
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			// 明确报错位置，便于调试（如第 2 个元素是 "a"）
			//return nil, fmt.Errorf("第 %d 个元素格式错误：%s（需为整数）", i+1, str)
		}
		int64Slice = append(int64Slice, num)
	}

	return int64Slice, nil
}

// 批量处理
func handleChatReplyMsgBatch(msg *pbmodel.Msg, session *Session, ids string) {
	var err error
	replyMsg := msg.GetPlainMsg().GetChatReply()

	idLst, err := SplitToInt64Array(ids)
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "reply ids format err", nil, session)
		return
	}

	pk1 := db.ComputePk(replyMsg.UserId)
	pk2 := db.ComputePk(replyMsg.FromId)
	user := session.GetUser()

	for _, msgId := range idLst {

		if replyMsg.RecvOk > 0 && replyMsg.ReadOk > 0 {
			err = Globals.scyllaCli.SetPChatRecvReadReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
				msgId, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
		} else if replyMsg.RecvOk > 0 {
			err = Globals.scyllaCli.SetPChatRecvReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
				msgId, time.Now().UTC().UnixMilli())
		} else if replyMsg.ReadOk > 0 {
			err = Globals.scyllaCli.SetPChatReadReply(pk1, pk2, replyMsg.UserId, replyMsg.FromId,
				msgId, time.Now().UTC().UnixMilli())
		} else {
			Globals.Logger.Debug("Recv User ChatMsgReply:", zap.String("error", "recv and read time is 0"))
		}

		// 收到回执后，从发送列表中删除
		user.PopMsgInCache(msgId)

		// 没有找到合适的
		if err != nil {
			Globals.Logger.Error("Recv User ChatMsgReply:", zap.Error(err))
			return
		}
	}

	// 转发消息
	trySendMsgToUser(replyMsg.UserId, msg)
}

// 这里的查询数据，包括私聊和群聊
func onQueryChatData(queryMsg *pbmodel.MsgQuery, session *Session) {
	Globals.Logger.Debug("onQueryChatData:", zap.String("queryMsg.QueryType", queryMsg.QueryType.String()), zap.Any("query", queryMsg.String()))
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
	if queryMsg.BigId == 0 {
		queryMsg.BigId = Globals.snow.GenerateID()
	}

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
				Status:       pbmodel.ChatMsgStatus(item.St),
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

	Globals.Logger.Debug("onQueryP2PChatData:", zap.Any("msg count:", len(chatDataList)))

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
	session.SendMessage(&msg)

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
	session.SendMessage(&msg)
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

	if queryMsg.BigId == 0 {
		queryMsg.BigId = Globals.snow.GenerateID()
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
				Status:       pbmodel.ChatMsgStatus(item.St),
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
	session.SendMessage(&msg)

}

// 查询消息是否送达的回执
// 对于没有给接收回执的消息需要重送，并使用发缓存跟踪
func onQueryChatReply(queryMsg *pbmodel.MsgQuery, session *Session) {
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
	var userId int64
	var chatReplyList []*pbmodel.MsgChatReply = nil

	if lst != nil && len(lst) > 0 {
		chatReplyList = make([]*pbmodel.MsgChatReply, len(lst))
		littleId = lst[0].Id
		bigId = lst[len(lst)-1].Id

		for index, item := range lst {
			// 自己发出去的
			if item.Io == model.ChatDataIOOut {
				fromId = item.Uid2
				userId = item.Uid1
			} else {
				continue
			}
			// 检查是否需要重发
			resendChatMsg(&item, session)

			data := pbmodel.MsgChatReply{
				MsgId:  item.Id,
				SendId: item.Usid,

				UserId: userId,
				FromId: fromId,

				SendOk:   item.Tm,
				RecvOk:   item.Tm1,
				ReadOk:   item.Tm2,
				ExtraMsg: "",
				Params:   nil,
			}
			chatReplyList[index] = &data
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
		QueryType:       pbmodel.QueryDataType_QueryDataTypeChatReply,
		ChatDataList:    nil,
		ChatReplyList:   chatReplyList,
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
	session.SendMessage(&msg)
}

// 重发消息
func resendChatMsg(msgStore *model.PChatDataStore, session *Session) {
	if msgStore.Tm1 > 0 || msgStore.Tm2 > 0 {
		return
	}

	if msgStore.Io == model.ChatDataIOIn {
		return
	}

	if msgStore.Uid1 != session.UserID {
		Globals.Logger.Fatal("resendChatMsg()  uid1 not equal with session user",
			zap.Int64("userid", session.UserID),
			zap.Int64("uid1", msgStore.Uid1))
		return
	}

	fid := msgStore.Uid2
	item := msgStore

	msgChat := pbmodel.MsgChat{
		MsgId:        item.Id,
		UserId:       session.UserID,
		FromId:       item.Uid1,
		ToId:         fid,
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

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_ChatData{
			ChatData: &msgChat,
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

	tryPushToUserMsgCache(fid, msgStore.Id, &msg)
	trySendMsgToUser(fid, &msg)

}
