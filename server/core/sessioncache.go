package core

import (
	"birdtalk/server/utils"
)

// https://github.com/robinfoxnan/go_concurrent_map
// 这里使用一个支持高并发的map来管理会话
type SessionCache struct {
	//lock sync.Mutex

	sessionMap utils.ConcurrentMap[int64, *Session]
}

// 获取用户的会话
func (ss *SessionCache) Get(sid int64) (*Session, bool) {
	v, b := ss.sessionMap.Get(sid)
	return v, b
}

// 如果没有设置过就保存
func (ss *SessionCache) UpSet(sid int64, s *Session) {

	ss.sessionMap.SetIfAbsent(sid, s)
	return
}

func (ss *SessionCache) Has(sid int64) bool {
	return ss.sessionMap.Has(sid)
}

func (ss *SessionCache) Remove(sid int64, uid int64) {

	ss.sessionMap.Remove(sid)

	// 从用户的信息中删除这个会话
	if uid != 0 {
		user, b := Globals.uc.GetUser(uid)
		if b {
			user.RemoveSessionID(sid)
		}
	}

	return
}
