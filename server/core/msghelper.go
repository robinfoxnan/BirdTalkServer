package core

import "birdtalk/server/pbmodel"

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

// 将消息转发给所有的群组用户
func notifyGroupMembers(groupId int64, msg *pbmodel.Msg) {
	group, _ := Globals.grc.GetGroup(groupId)
	if group == nil {
		return
	}

	members := group.GetMembers()
	for _, mId := range members {
		trySendMsgToUser(mId, msg)
	}
}
