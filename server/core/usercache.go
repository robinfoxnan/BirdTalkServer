package core

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
)

// https://github.com/robinfoxnan/go_concurrent_map
// 这里使用一个支持高并发的map来管理会话
type UserCache struct {
	//lock sync.Mutex

	userMap utils.ConcurrentMap[int64, *model.User]
}

// 如果没有，则从redis中查找，如果redis中没有，则从数据库中查找
func (u *UserCache) Get(uid int64) (*model.User, error) {

	user := model.NewUser()
	user.UserId = uid
	return user, nil
}

func (u *UserCache) Set(uid int64, user *model.User) {
	u.userMap.SetIfAbsent(uid, user)
}

// 超时不使用则删除
func (u *UserCache) Remove(uid int64) {
	u.userMap.Remove(uid)
}
