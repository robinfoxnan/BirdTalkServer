package model

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"strings"
	"sync"
)

const (
	// 群主
	RoleGroupOwner = 1 << iota
	// 管理员
	RoleGroupAdmin
	// 普通用户
	RoleGroupMember
	// 只读普通用户
	RoleGroupMemberRead
)

type GroupMember struct {
	Uid  int64
	Nick string
}

type Group struct {
	pbmodel.GroupInfo

	Owner   *GroupMember
	Admins  map[int64]*GroupMember
	Members map[int64]*GroupMember

	IsDeleted    bool
	LastActiveTm int64

	Mu sync.Mutex
}

func (g *Group) MergeGroup(other *Group) {

}

func (g *Group) MergeGroupInfo(info *pbmodel.GroupInfo) {
	if info.GroupName != "" {
		g.GroupName = info.GetGroupName()
	}
	if info.Tags != nil {
		g.Tags = info.Tags
	}

	if info.Params != nil {
		if g.Params == nil {
			g.Params = map[string]string{}
		}

		for k, v := range info.Params {
			g.Params[k] = v
		}
	}
}

// 有些敏感信息不能发给用户
func (g *Group) GetGroupInfo() *pbmodel.GroupInfo {

	params := map[string]string{}
	for k, v := range g.Params {
		key := strings.ToLower(k)
		if key == "joinanswer" {
			continue
		}
		if key == "code" {
			continue
		}
		params[key] = v
	}

	return &pbmodel.GroupInfo{
		GroupId:   g.GroupId,
		GroupName: g.GroupName,
		GroupType: g.GroupType,
		Tags:      g.Tags,
		Params:    params,
	}
}

func NewGroupFromInfo(info *pbmodel.GroupInfo) *Group {
	return &Group{
		GroupInfo:    *info,
		Owner:        nil,
		Admins:       make(map[int64]*GroupMember),
		Members:      make(map[int64]*GroupMember),
		Mu:           sync.Mutex{},
		IsDeleted:    false,
		LastActiveTm: utils.GetTimeStamp(),
	}
}

func (g *Group) SetOwner(uid int64, nick string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Owner = &GroupMember{
		Uid:  uid,
		Nick: nick,
	}
}

func (g *Group) AddMember(uid int64, nick string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	member := &GroupMember{
		Uid:  uid,
		Nick: nick,
	}

	g.Members[uid] = member
}

func (g *Group) RemoveMember(uid int64) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	delete(g.Members, uid)
}

func (g *Group) ClearMember() {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Members = make(map[int64]*GroupMember)
}

func (g *Group) SetDeleted() {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.IsDeleted = true
}

func (g *Group) CheckDeleted() bool {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	return g.IsDeleted
}

func (g *Group) GetMembers() []int64 {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	members := make([]int64, len(g.Members))
	index := 0
	for _, v := range g.Members {
		members[index] = v.Uid
		index++
	}

	return members

}

func (g *Group) IsOwner(uid int64) bool {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if g.Owner == nil {
		return false
	}

	if g.Owner.Uid == uid {
		return true
	}

	return false
}

func (g *Group) IsAdmin(uid int64) bool {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if g.Owner != nil {
		if g.Owner.Uid == uid {
			return true
		}
	}

	for k, _ := range g.Admins {
		if k == uid {
			return true
		}
	}

	return false
}

// 设置活跃的最后时间
func (g *Group) Active() {
	g.LastActiveTm = utils.GetTimeStamp()
}

func (g *Group) IsTimeout() bool {
	current := utils.GetTimeStamp()
	delta := current - g.LastActiveTm
	// 24小时的毫秒数
	// 检查时间差是否大于或等于24小时
	if delta >= utils.TwentyFourHoursInMilliseconds {
		return true
	}
	return false
}

// 查询某个值
func (g *Group) GetParamByKey(key string) string {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.GroupInfo.GetParams() == nil {
		return ""
	}

	v, ok := g.GroupInfo.GetParams()[key]
	if !ok {
		return ""
	}
	return v
}

func (g *Group) SetParamByKey(key, value string) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.GroupInfo.GetParams() == nil {
		g.GroupInfo.Params = map[string]string{
			key: value,
		}
		return
	}

	g.GroupInfo.Params[key] = value

}

func (g *Group) GetAdminMembers() []int64 {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	if g.Owner == nil {
		return nil
	}

	members := make([]int64, len(g.Admins)+1)
	members[0] = g.Owner.Uid
	index := 1
	for _, v := range g.Admins {
		members[index] = v.Uid
		index++
	}

	return members

}
