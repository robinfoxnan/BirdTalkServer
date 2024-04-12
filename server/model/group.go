package model

import "birdtalk/server/pbmodel"

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

type Group struct {
	pbmodel.GroupInfo
}
