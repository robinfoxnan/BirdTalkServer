package core

import "birdtalk/server/utils"

// https://github.com/robinfoxnan/go_concurrent_map
// 这里使用一个支持高并发的map来管理会话
type SessionStore struct {
	//lock sync.Mutex

	sessionMap utils.ConcurrentMap[int64, *Session]
}

// 获取用户的会话
func (ss *SessionStore) Get(sid int64) (*Session, bool) {
	v, b := ss.sessionMap.Get(sid)
	return v, b
}

// 如果没有设置过就保存
func (ss *SessionStore) UpSet(sid int64, s *Session) {

	ss.sessionMap.SetIfAbsent(sid, s)
	return
}

func (ss *SessionStore) Has(sid int64) bool {
	return ss.sessionMap.Has(sid)
}

func (ss *SessionStore) Remove(sid int64) {
	ss.sessionMap.Remove(sid)
	return
}
