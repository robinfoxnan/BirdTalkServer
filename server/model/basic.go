package model

import (
	"birdtalk/server/pbmodel"
	"sync"
)

const (
	UserStatusNone     = 0         // 等待 HELLO 包
	UserStatusExchange = 1 << iota // 等待秘钥交换1
	UserStatusRegister             // 收到了注册申请，等待验证码2
	UserStatusLogin                // 收到了登录申请，需要等待验证码4
	UserStatusOk                   // 登录完成8
	UserStatusValidate             // 验证状态16，
)

type User struct {
	pbmodel.UserInfo
	sessionId []int64           // 多用户会话ID
	params    map[string]string // 某些状态下的附加信息都放在这里，比如验证码
	status    uint32
	mu        sync.Mutex
}

// New 函数用于创建一个 User 实例
func NewUser() *User {
	return &User{
		UserInfo:  pbmodel.UserInfo{}, // 初始化嵌入的 UserInfo 结构体
		sessionId: make([]int64, 0),   // 初始化 sessionId 切片
		params:    make(map[string]string),
		status:    0,
		mu:        sync.Mutex{}, // 初始化互斥锁
	}
}

// New 函数用于创建一个 User 实例
func NewUserFromInfo(userInfo *pbmodel.UserInfo) *User {
	return &User{
		UserInfo:  *userInfo, // 使用传入的 *pbmodel.UserInfo 参数
		sessionId: make([]int64, 0),
		mu:        sync.Mutex{},
	}
}

// 添加会话ID
func (u *User) AddSessionID(sessionID int64) int {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.sessionId = append(u.sessionId, sessionID)
	return len(u.sessionId)
}

// 删除会话ID
func (u *User) DeleteSessionID(sessionID int64) int {
	u.mu.Lock()
	defer u.mu.Unlock()

	for i, id := range u.sessionId {
		if id == sessionID {
			// 将要删除的会话ID与最后一个元素交换位置，然后缩减切片长度
			u.sessionId[i] = u.sessionId[len(u.sessionId)-1]
			u.sessionId = u.sessionId[:len(u.sessionId)-1]
			return len(u.sessionId)
		}
	}
	return len(u.sessionId)
}

// SetStatus 设置指定状态，同时存在的一个状态
func (u *User) SetStatus(newStatus uint32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.status |= newStatus
}

// ClearStatus 清除指定状态，取消某个同时存在的状态
func (u *User) ClearStatus(statusToClear uint32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.status &= ^statusToClear
}

// 切换到某个状态
func (u *User) ChangeToStatus(newStatus uint32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.status = newStatus
}

// HasStatus 检查是否包含指定状态
func (u *User) HasStatus(checkStatus uint32) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.status&checkStatus == checkStatus
}

///////////////////////////////////////////////////////////////////////

type Group struct {
	pbmodel.GroupInfo
}
