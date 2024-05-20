package model

import (
	"birdtalk/server/pbmodel"
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

	Mu sync.Mutex
}

func (g *Group) MergeGroup(other *Group) {

}

func NewGroupFromInfo(info *pbmodel.GroupInfo) *Group {
	return &Group{GroupInfo: *info}
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
