package model

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
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

// 使用uint16表示针对某个人的权限设置，每1位都是设置1允许，0禁止
const (
	PermissionMP2PNone     = 0
	PermissionMaskFriend   = 1 << iota // 1是否关注 1
	PermissionMaskChat                 // 2是否可以发消息
	PermissionMaskViewInfo             // 3是否查看敏感信息，手机，邮件
	PermissionMaskViewArt              // 4是否可以查看帖子
	PermissionMaskAdd                  // 5加好友，默认是开放的, 不开放就是拉黑
	PermissionMaskViewLoc              // 6是否可以查看位置
	PermissionMask7
	PermissionMask8
	PermissionMask9
	PermissionMask10
	PermissionMask11
	PermissionMask12
	PermissionMask13
	PermissionMask14
	PermissionMaskExist
	PermissionMask16

	// 陌生人也是如此，左移16位
)

// 如果你是对方好友，即对方关注你了，你得到的权限是除了看位置信息，其他都可以！
const DefaultPermissionP2P = PermissionMaskExist | PermissionMaskAdd | PermissionMaskChat | PermissionMaskViewInfo | PermissionMaskViewArt

// 如果你是陌生人，那么默认可以添加好友，可以看帖子
const DefaultPermissionStranger = PermissionMaskExist | ((PermissionMaskAdd | PermissionMaskViewArt) << 16)

// 默认朋友权限的设置
const DefaultPermissionFriend = DefaultPermissionP2P | PermissionMaskFriend

// 相关redis表
// 1：用户基础信息表，bsui_
// 2：用户好友权限表，bsufb_
// 4：用户关注表，bsufo_
// 8：用户粉丝表，bsufa_
// 16：用户的指纹表，bsut_ 这个存储在session 中
// 32：用户所属群组，bsuing_
// 64：用户session 分布表  30分钟 bsud_

// 内存缓存使用的用户模型
type User struct {
	pbmodel.UserInfo
	SessionId []int64 // 多用户会话ID
	Status    uint32
	MaskLoad  uint32            // 目前的加载状态
	Params    map[string]string // 某些状态下的附加信息都放在这里，比如验证码

	Block      map[int64]uint64 // 权限控制列表，高位是对方给自己的权限，低32位是自己向对方的权限
	Following  map[int64]string // 关注列表
	Fans       map[int64]string // 粉丝列表
	Groups     map[int64]bool   // 群组
	SessionDis map[int64]int32  // 会话分布表

	Mu sync.Mutex
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
		Status:    UserStatusNone,
		Mu:        sync.Mutex{}, // 初始化互斥锁

		Params: make(map[string]string),

		Following: make(map[int64]string),
		Fans:      make(map[int64]string),
		Block:     make(map[int64]uint64),

		Groups:     make(map[int64]bool),
		SessionDis: make(map[int64]int32),
	}
}

// 自己给对方的权限
func (u *User) SetSelfMask(fid int64, mask uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	data, ok := u.Block[fid]
	if ok {
		// 如果好友ID已存在于屏蔽列表中，则更新其屏蔽掩码
		u.Block[fid] = data | uint64(mask)
	} else {
		// 如果好友ID不存在于屏蔽列表中，则向列表中添加新的条目
		u.Block[fid] = uint64(mask)
	}
}

func (u *User) RemoveSelfMask(fid int64, mask uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	data, ok := u.Block[fid]
	if ok {
		// 如果好友ID存在于屏蔽列表中，则从屏蔽掩码中移除指定的掩码
		u.Block[fid] = data &^ uint64(mask)
	}
	// 如果好友ID不存在于屏蔽列表中，则不执行任何操作
}

// 对方给自己的权限，对高32位设置
func (u *User) SetFriendToMeMask(fid int64, mask uint32) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	data, ok := u.Block[fid]
	if ok {
		// 如果好友ID已存在于屏蔽列表中，则更新其屏蔽掩码
		u.Block[fid] = data | (uint64(mask) << 32) | PermissionMaskExist
	} else {
		// 如果好友ID不存在于屏蔽列表中，则向列表中添加新的条目
		u.Block[fid] = uint64(mask) << 32
	}
}

// GetP2PMask 用于获取用户结构体中指定好友的 P2P 屏蔽掩码
func (u *User) GetFriendToMeMask(fid int64) (uint32, bool) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	// 尝试从屏蔽列表中获取指定好友的屏蔽掩码
	data, ok := u.Block[fid]
	if ok {
		// 如果好友ID已存在于屏蔽列表中，则返回其屏蔽掩码的高32位
		return uint32(data >> 32), true
	}
	// 如果好友ID不存在于屏蔽列表中，则返回0
	return 0, false
}

// 检查是否可以向对方发送消息，
// 第2个返回值：是否设置了，未设置应该去加载
// enable, ok := CanSendMesssageTo(1001)
func (u *User) CanSendMessageTo(fid int64) (bool, bool) {
	return u.checkP2PMask(fid, PermissionMaskChat)
}

func (u *User) CanSendAddReq(fid int64) (bool, bool) {
	return u.checkP2PMask(fid, PermissionMaskAdd)
}

func (u *User) CanSendViewInfo(fid int64) (bool, bool) {
	return u.checkP2PMask(fid, PermissionMaskViewInfo)
}

func (u *User) CanSendViewArticle(fid int64) (bool, bool) {
	return u.checkP2PMask(fid, PermissionMaskViewArt)
}

func (u *User) CanSendViewLocation(fid int64) (bool, bool) {
	return u.checkP2PMask(fid, PermissionMaskViewLoc)
}

// 检查某个权限
func (u *User) checkP2PMask(fid int64, mask uint32) (bool, bool) {
	data, ok := u.GetFriendToMeMask(fid)
	if !ok {
		return false, false
	}

	strangerMask := data >> 16

	// 还未设置
	if (data & PermissionMaskExist) == 0 {
		return false, false
	}

	if (data & PermissionMaskFriend) == 0 {
		// 对方未关注你！陌生人
		return (strangerMask & mask) > 0, true
	}

	// 好友，但是也不一定能发
	return (data & mask) > 0, true
}

// 更新时候直接将新的数据合并过来，这个比较繁琐
func (u *User) MergeUser(newUser *User) {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	// session
	u.UserInfo = newUser.UserInfo
	u.addSessionIDListNoLock(newUser.SessionId)
	u.Status = newUser.Status
	// 合并几个map
	utils.MergeMap(u.Params, newUser.Params)

	utils.MergeMap(u.Following, newUser.Following)
	utils.MergeMap(u.Fans, newUser.Fans)
	utils.MergeMapMask(u.Block, newUser.Block)

	utils.MergeMap(u.Groups, newUser.Groups)
	utils.MergeMap(u.SessionDis, newUser.SessionDis)
}

// 添加会话ID
func (u *User) AddSessionID(newSid int64) int {
	u.Mu.Lock()
	defer u.Mu.Unlock()

	for _, id := range u.SessionId {
		if id == newSid {
			return len(u.SessionId)
		}
	}
	u.SessionId = append(u.SessionId, newSid)
	return len(u.SessionId)
}

func (u *User) addSessionIDListNoLock(sidList []int64) int {

	for _, newSid := range sidList {
		bFound := false
		for _, id := range u.SessionId {
			if id == newSid {
				bFound = true
				break
			}
		}
		if !bFound {
			u.SessionId = append(u.SessionId, newSid)
		}
	}

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
