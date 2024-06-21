package core

import "birdtalk/server/pbmodel"

func tryPushToUserMsgCache(uid int64, msgId int64, msg *pbmodel.Msg) {
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		user.PushMsgInCache(msgId, msg)
	}
}

// 尝试发送到对方用户
func trySendMsgToUser(uid int64, msg *pbmodel.Msg) {
	// 现在本机查找
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		for _, sid := range user.SessionId {

			sess, b := Globals.ss.Get(sid)
			if b && sess != nil {
				sess.SendMessage(msg)
			}
		}
		return
	}

	// TODO: 如果是集群模式，则应该查找并转发
	if Globals.Config.Server.ClusterMode {

		// 将对方的登录情况写到自己的的User中，如果没有，去redis中查一下
	}
}

func trySendMsgToMe(uid int64, msg *pbmodel.Msg, session *Session) {
	// 现在本机查找
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		for _, sid := range user.SessionId {

			if session.Sid == sid {
				continue
			}
			sess, b := Globals.ss.Get(sid)
			if b && sess != nil {
				sess.SendMessage(msg)
			}
		}
		return
	}

	// TODO: 如果是集群模式，则应该查找并转发
	if Globals.Config.Server.ClusterMode {

		// 将对方的登录情况写到自己的的User中，如果没有，去redis中查一下
	}
}

func trySendMsgToUserList(uidList []int64, msg *pbmodel.Msg) {
	for _, uid := range uidList {
		trySendMsgToUser(uid, msg)
	}
}

// 将消息转发给所有的群组用户，除了自己的会话，防止
func sendToGroupMembersExceptMe(groupId int64, msg *pbmodel.Msg, session *Session) {
	group, _ := Globals.grc.GetGroup(groupId)
	if group == nil {
		return
	}

	// 集群模式需要检查各个用户的分布情况
	if Globals.Config.Server.ClusterMode {

	} else {
		// 单机模式直接转发在线的用户
		members := group.GetMembers()
		for _, mId := range members {
			if mId == session.UserID {
				trySendMsgToMe(mId, msg, session)
			} else {
				trySendMsgToUser(mId, msg)
			}

		}
	}

}

// 群回执，每个用户都有份
func notifyGroupMembers(groupId int64, msg *pbmodel.Msg) {
	group, _ := Globals.grc.GetGroup(groupId)
	if group == nil {
		return
	}

	// 集群模式需要检查各个用户的分布情况
	if Globals.Config.Server.ClusterMode {

	} else {
		// 单机模式直接转发在线的用户
		members := group.GetMembers()
		for _, mId := range members {
			trySendMsgToUser(mId, msg)
		}
	}

}
