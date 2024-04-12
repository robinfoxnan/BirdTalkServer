package model

import (
	"birdtalk/server/utils"
)

// https://github.com/robinfoxnan/go_concurrent_map
// 这里使用一个支持高并发的map来管理会话
type UserCache struct {
	//lock sync.Mutex

	userMap utils.ConcurrentMap[int64, *User]
}

// 如果没有，则从redis中查找，如果redis中没有，则从数据库中查找
func (uc *UserCache) GetUser(uid int64) (*User, bool) {

	user, err := uc.userMap.Get(uid)
	return user, err
}

func (uc *UserCache) SetUser(uid int64, user *User) {
	uc.userMap.SetIfAbsent(uid, user)
}

// 超时不使用则删除
func (uc *UserCache) RemoveUser(uid int64) {
	uc.userMap.Remove(uid)
}
