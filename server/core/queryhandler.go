package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
)

// 7 这里的查询包括：私聊消息同步，群消息同步，回执查询；好友请求查询，群用户操作消息
func handleCommonQuery(msg *pbmodel.Msg, session *Session) {

	// 检查数据指针是否为空
	queryMsg := msg.GetPlainMsg().GetCommonQuery()
	if queryMsg == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "query message  is null", nil, session)
		return
	}

	ok := checkUserLogin(session) // 内部发送错误通知
	if !ok {
		return
	}

	switch queryMsg.QueryType {
	case pbmodel.QueryDataType_QueryDataTypeChatData:
		onQueryChatData(queryMsg, session)
	case pbmodel.QueryDataType_QueryDataTypeChatReply:
		onQueryChatReply(queryMsg, session)
	case pbmodel.QueryDataType_QueryDataTypeFriendOP:
		onQueryFriendOp(queryMsg, session)
	case pbmodel.QueryDataType_QueryDataTypeGroupOP:
		onQueryGroupOp(queryMsg, session)
	default:
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "query message, query type is unknown", nil, session)
	}
}

// 格式转换
func FriendOpRecords2Msg(records []model.CommonOpStore) ([]*pbmodel.FriendOpResult, []*pbmodel.GroupOpResult) {
	if records == nil || len(records) == 0 {
		return nil, nil
	}

	lstFriend := make([]*pbmodel.FriendOpResult, 0)
	lstGroup := make([]*pbmodel.GroupOpResult, 0)

	var friendItem *pbmodel.FriendOpResult = nil
	var groupItem *pbmodel.GroupOpResult = nil

	for _, record := range records {
		userReq := int64(0)
		userParam := int64(0)
		friendItem = nil
		groupItem = nil

		result := ""
		if record.Ret == model.UserOpResultOk {
			result = "accept"
		} else if record.Ret == model.UserOpResultRefuse {
			result = "refuse"
		}

		// uid1是自己，但是需要根据数据流向判断自己发出请求，还是对方请求自己
		if record.Io == model.ChatDataIOOut {
			userReq = record.Uid1
			userParam = record.Uid2
		} else {
			userReq = record.Uid2
			userParam = record.Uid1
		}

		switch record.Cmd {

		// 这个还没有记录
		//case model.CommonGroupOpJoinRequest:

		case model.CommonGroupOpInviteRequest: // 邀请入群
			groupItem = &pbmodel.GroupOpResult{
				Result:    result,
				Operation: pbmodel.GroupOperationType_GroupInviteRequest,
				MsgId:     record.Id,
				SendId:    record.Usid,
				Detail:    "",
				Group:     &pbmodel.GroupInfo{GroupId: record.Gid},
				ReqMem: &pbmodel.GroupMember{
					UserId: userReq,
				},
				Members: []*pbmodel.GroupMember{
					&pbmodel.GroupMember{
						UserId: userParam,
					},
				},
			}

		case model.CommonUserOpAddRequest: // 请求或者被请求添加好友

			friendItem = &pbmodel.FriendOpResult{
				Result:    result,
				Operation: pbmodel.UserOperationType_AddFriend,
				MsgId:     record.Id,
				SendId:    record.Usid,
				User: &pbmodel.UserInfo{
					UserId: userReq,
				},
				Users: []*pbmodel.UserInfo{
					&pbmodel.UserInfo{
						UserId: userParam,
					},
				},
			}
		}

		if friendItem != nil {
			lstFriend = append(lstFriend, friendItem)
		}

		if groupItem != nil {
			lstGroup = append(lstGroup, groupItem)
		}

	}

	return lstFriend, lstGroup
}

// 查询好友加好友以及被加好友的操作记录，以及群邀请记录
func onQueryFriendOp(queryMsg *pbmodel.MsgQuery, session *Session) {

	pk := db.ComputePk(session.UserID)
	var err error
	var lst []model.CommonOpStore

	switch queryMsg.GetSynType() {
	case pbmodel.SynType_SynTypeForward:
		lst, err = Globals.scyllaCli.FindUserOpForward(pk, session.UserID, queryMsg.LittleId, 100)
	default:
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "query message, syn type is not accepted", nil, session)
		return
	}
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "query db meet some err", nil, session)
		return
	}

	//params := map[string]string{
	//	"status": status,
	//}
	friendLst, groupLst := FriendOpRecords2Msg(lst)

	queryRet := pbmodel.MsgQueryResult{
		UserId:          session.UserID,
		GroupId:         0,
		LittleId:        0,
		BigId:           0,
		SynType:         queryMsg.GetSynType(),
		Tm:              utils.GetTimeStamp(),
		ChatType:        pbmodel.ChatType_ChatTypeNone,
		QueryType:       pbmodel.QueryDataType_QueryDataTypeFriendOP,
		FriendOpRetList: friendLst,
		GroupOpRetList:  groupLst,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_CommonQueryRet{
			CommonQueryRet: &queryRet,
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

// 群管理员才需要同步这个数据，仅仅是用户请求加入群的消息
func onQueryGroupOp(queryMsg *pbmodel.MsgQuery, session *Session) {

}
