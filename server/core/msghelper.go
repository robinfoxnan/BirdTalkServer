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

	// 如果是集群模式，则应该查找并转发
	if Globals.Config.Server.ClusterMode {

	}
}