package model

import "birdtalk/server/utils"

type GroupCache struct {
	groupMap utils.ConcurrentMap[int64, *Group]
}

// 新建一个
func NewGroupCache() *GroupCache {
	return &GroupCache{groupMap: utils.NewConcurrentMap[int64, *Group]()}
}

func (cache *GroupCache) GetGroup(id int64) (*Group, bool) {
	return cache.groupMap.Get(id)
}

// 更新时候的回调函数，
// 这一唯一被调用的可能性就是，服务中没有这个group, 多个用户的操作同时引发了加载操作，
func updateInsertGroup(exist bool, oldGroup *Group, newGroup *Group) *Group {
	if exist == false {
		return newGroup
	} else {
		oldGroup.MergeGroup(newGroup)
		return oldGroup
	}
}

// 这里可能会有并发冲突，需要解决的就是session列表需要合并
func (cache *GroupCache) InsertGroup(id int64, g *Group) *Group {
	ret := cache.groupMap.Upsert(id, g, updateInsertGroup)
	return ret
}
