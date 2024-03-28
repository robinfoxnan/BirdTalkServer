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

// 内存缓存使用的用户模型
type User struct {
	pbmodel.UserInfo
	SessionId []int64           // 多用户会话ID
	Params    map[string]string // 某些状态下的附加信息都放在这里，比如验证码
	Status    uint32
	Mu        sync.Mutex

	Following map[int64]*FriendMemInfo // 关注列表
	Fans      map[int64]*FriendMemInfo // 粉丝列表
	Block     map[int64]*BlockMemInfo
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
		Status:    0,
		Mu:        sync.Mutex{}, // 初始化互斥锁

		Following: make(map[int64]*FriendMemInfo),
		Fans:      make(map[int64]*FriendMemInfo),
		Block:     make(map[int64]*BlockMemInfo),
	}
}

// 添加会话ID
func (u *User) AddSessionID(sessionID int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	u.SessionId = append(u.SessionId, sessionID)
	return len(u.SessionId)
}

// 删除会话ID
func (u *User) DeleteSessionID(sessionID int64) int {
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

// /////////////////////////////////////////////////////////////////////
// 保存到内存中
type FriendMemInfo struct {
	fid  int64
	tm   int64
	nick string
}

type BlockMemInfo struct {
	FriendMemInfo
	Perm int32
}

// 下面2个结构对应数据库中结构
// 关注和粉丝的表一样，没有权限一项
type FriendStore struct {
	Pk   int16
	Uid1 int64
	Uid2 int64
	Tm   int64
	Nick string
}

type BlockStore struct {
	FriendStore
	Perm int32
}

// //////////////////////////////////////////////////////////////////
type Group struct {
}

const (
	// 群主
	GroupOwner = 1 << iota
	// 管理员
	GroupAdmin
	// 普通用户
	GroupMember
	// 只读普通用户
	GroupMemberReadOnly
)

// 数据库中记录群组成员的结构体
type GroupMemberStore struct {
	Pk   int16
	Role int16
	Gid  int64
	Uid  int64
	Tm   int64
	Nick string
}

type UserInGStore struct {
	Pk  int16
	Uid int64
	Gid int64
}