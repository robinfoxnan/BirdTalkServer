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

func NewUserCache() *UserCache {
	return &UserCache{userMap: utils.NewConcurrentMap[int64, *User]()}
}

// 如果没有，则从redis中查找，如果redis中没有，则从数据库中查找
func (uc *UserCache) GetUser(uid int64) (*User, bool) {

	user, ok := uc.userMap.Get(uid)
	return user, ok
}

// 更新时候的回调函数，如果未设置，则
func updateInsertUser(exist bool, oldUser *User, newUser *User) *User {
	if exist == false {
		return newUser
	} else {
		oldUser.MergeUser(newUser)
		return oldUser
	}
}

// 这里可能会有并发冲突，需要解决的就是session列表需要合并
func (uc *UserCache) SetOrUpdateUser(uid int64, user *User) *User {
	res := uc.userMap.Upsert(uid, user, updateInsertUser)
	return res
}

// 超时不使用则删除
func (uc *UserCache) RemoveUser(uid int64) {
	uc.userMap.Remove(uid)
}
