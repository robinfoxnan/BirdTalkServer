package model

import (
	"birdtalk/server/pbmodel"
	"sync"
)

// 最多支持32个状态
const (
	UserStatusNone       = 0         // 等待 HELLO 包
	UserStatusExchange   = 1 << iota // 等待秘钥交换1
	UserWaitLogin                    // 可以继续
	UserStatusRegister               // 收到了注册申请，等待验证码2,不需要验证则继续UserReady
	UserStatusLogin                  // 收到了登录申请，需要等待验证码4,不需要验证则继续UserReady
	UserStatusOk                     // 登录完成8
	UserStatusValidate               // 验证状态16，
	UserStatusChangeInfo             // 32
)

// 内存缓存使用的用户模型
type User struct {
	pbmodel.UserInfo
	SessionId []int64           // 多用户会话ID
	Params    map[string]string // 某些状态下的附加信息都放在这里，比如验证码
	Status    uint32
	Mu        sync.Mutex

	Following map[int64]string // 关注列表
	Fans      map[int64]string // 粉丝列表
	Block     map[int64]int32  // 权限控制列表
	Groups    map[int64]bool   // 群组
}

// New 函数用于创建一个 User 实例
func NewUser() *User {
	userInfo := pbmodel.UserInfo{}
	return NewUserFromInfo(&userInfo)
}

// New 函数用于创建一个 User 实例
func NewUserFromInfo(userInfo *pbmodel.UserInfo) *User {
	return &User{
		UserInfo:  *userInfo,        // 使用传入的 *pbmodel.UserInfo 参数
		SessionId: make([]int64, 0), // 初始化 sessionId 切片
		Params:    make(map[string]string),
		Status:    UserStatusNone,
		Mu:        sync.Mutex{}, // 初始化互斥锁

		Following: make(map[int64]string),
		Fans:      make(map[int64]string),
		Block:     make(map[int64]int32),
		Groups:    make(map[int64]bool),
	}
}

// 添加会话ID
func (u *User) AddSessionID(sid int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	for _, id := range u.SessionId {
		if id == sid {
			return len(u.SessionId)
		}
	}
	u.SessionId = append(u.SessionId, sid)
	return len(u.SessionId)
}

// 删除会话ID
func (u *User) RemoveSessionID(sessionID int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	for i, id := range u.SessionId {
		if id == sessionID {
			// 将要删除的会话ID与最后一个元素交换位置，然后缩减切片长度
			u.SessionId[i] = u.SessionId[len(u.SessionId)-1]
			u.SessionId = u.SessionId[:len(u.SessionId)-1]
			return len(u.SessionId)
		}
	}
	return len(u.SessionId)
}

// SetStatus 设置指定状态，同时存在的一个状态
func (u *User) SetStatus(newStatus uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Status |= newStatus
}

// ClearStatus 清除指定状态，取消某个同时存在的状态
func (u *User) ClearStatus(statusToClear uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Status &= ^statusToClear
}

// 切换到某个状态
func (u *User) ChangeToStatus(newStatus uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Status = newStatus
}

// HasStatus 检查是否包含指定状态
func (u *User) HasStatus(checkStatus uint32) bool {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	return u.Status&checkStatus == checkStatus
}

func (u *User) SetFollow(fid int64, nick string) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Following[fid] = nick
	return len(u.Following)
}

func (u *User) DelFollow(fid int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	delete(u.Following, fid)
	return len(u.Following)
}

func (u *User) SetFan(fid int64, nick string) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Fans[fid] = nick
	return len(u.Fans)
}

func (u *User) DelFan(fid int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	delete(u.Fans, fid)
	return len(u.Fans)

}

func (u *User) SetPermission(fid int64, mask int32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.Block[fid] = mask
}
func (u *User) GetPermission(fid int64) (int32, bool) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	ret, b := u.Block[fid]
	return ret, b
}

func (u *User) SetExtraKeyValue(key, value string) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	u.Params[key] = value
}

func (u *User) GetExtraKeyValue(key string) (string, bool) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	ret, b := u.Params[key]
	return ret, b
}

func (u *User) DelExtraKeyValue(key string) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	delete(u.Params, key)
}
