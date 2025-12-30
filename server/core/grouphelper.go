package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"errors"
	"go.uber.org/zap"
)

// 这个文件主要是处理三级缓存相关的逻辑

// 加载群的用户信息; 群的数量为2000上限，所以这里直接从数据库加载所有群用户，放到redis中
func loadGroupMembers(group *model.Group) {
	// 从redis中查找，如果找到了，
	lst, err := Globals.redisCli.GetGroupMembers(group.GroupId)
	if lst != nil && len(lst) > 0 {
		Globals.Logger.Debug("find group members in redis", zap.Int64("group_id", group.GroupId), zap.Any("lst", lst))
		group.SetMembers(lst)
		return
	}

	// 从数据库中加载
	memlist, err := Globals.scyllaCli.FindGroupMembers(db.ComputePk(group.GroupId), group.GroupId, 0, 2000)
	if err != nil {
		Globals.Logger.Fatal("loadGroupMembers db error", zap.Error(err))
		return
	}

	Globals.Logger.Debug("find group members in scyllaDb", zap.Int64("group_id", group.GroupId), zap.Any("lst", memlist))
	// 同步到redis
	err = Globals.redisCli.SetGroupMembers(group.GroupId, memlist)
	// 如果错了，大不了下次再加载

	// 同步到内存在中
	if memlist != nil && len(memlist) > 0 {
		group.SetMembers(memlist)
	}

}

// 从数据库中加载Group基础信息
func findGroup(gid int64) (*model.Group, error) {
	group, ok := Globals.grc.GetGroup(gid)
	if ok && group != nil {
		// 这里需要检查用户是否为空，如果是空，则应该加载，因为还没有初始化
		if len(group.Members) == 0 {
			Globals.Logger.Debug("find group in memory", zap.Int64("group_id", gid))
			loadGroupMembers(group)
		}
		return group, nil
	}

	// 从redis中查找群
	groupInfo, err := Globals.redisCli.GetGroupInfoById(gid)
	if groupInfo != nil {
		Globals.Logger.Debug("find group in redis", zap.Int64("group_id", gid))
		group = model.NewGroupFromInfo(groupInfo)
		Globals.grc.InsertGroup(gid, group)
		loadGroupMembers(group)

		return group, nil
	}

	lst, err := Globals.mongoCli.FindGroupById(gid)
	if lst == nil || len(lst) == 0 {

		Globals.Logger.Error("findGroup() ->FindGroupById() err", zap.Int64("gid", gid), zap.Error(err))
		return nil, err
	}

	if len(lst) > 1 {
		Globals.Logger.Error("find more than one group by id", zap.Int64("gid", gid))
		return nil, errors.New("find more than one group by id")
	}

	Globals.Logger.Debug("find group in scyllaDb", zap.Int64("group_id", gid))
	groupInfo = &lst[0]

	// 保存到当前的群基础信息到redis中
	Globals.redisCli.SetGroupInfo(groupInfo)

	// 保存到内存
	group = model.NewGroupFromInfo(groupInfo)
	Globals.grc.InsertGroup(gid, group)
	loadGroupMembers(group)

	return group, nil

}

// 计算用户所在的群的信息列表，先获取一个数组，再填充
// 重启后内存没有数据，redis重启，里面可能也没有数据，所以尝试逐级加载
func LoadUserInGroupList(uid int64, fromId int64) ([]*pbmodel.GroupInfo, error) {

	u, _ := Globals.uc.GetUser(uid)
	if u != nil {
		gidLst := u.GetInGroups(fromId, Globals.maxPageSize)
		if gidLst != nil && len(gidLst) == 0 {
			Globals.Logger.Debug("LoadUserInGroupList from memory", zap.Any("uid", uid), zap.Any("gid", gidLst))
			return intList2GroupList(gidLst), nil
		}

	}

	_, gidLst, err := Globals.redisCli.GetUserInGroupPage(uid, fromId, int64(Globals.maxPageSize))
	if err == nil && gidLst != nil && len(gidLst) > 0 {
		Globals.Logger.Debug("LoadUserInGroupList from redis", zap.Any("uid", uid), zap.Any("gid", gidLst))
		if u != nil {
			u.SetInGroup(gidLst)
		}
		return intList2GroupList(gidLst), nil
	}

	// 从数据库查询
	pk := db.ComputePk(uid)
	lst, err := Globals.scyllaCli.FindUserInGroups(pk, uid, fromId, Globals.maxPageSize)
	if err != nil {
		return nil, err
	}
	if len(lst) == 0 {
		return nil, nil
	}

	// 这里应该保存到redis
	gidLst = make([]int64, len(lst))
	for index, item := range lst {
		gidLst[index] = item.Gid
	}
	if u != nil {
		u.SetInGroup(gidLst)
	}
	Globals.redisCli.SetUserInGroup(uid, gidLst)

	return gStoreList2GroupList(lst), nil

}

// 数字解析为群列表
func intList2GroupList(lst []int64) []*pbmodel.GroupInfo {
	if lst == nil || len(lst) == 0 {
		return nil
	}

	ret := make([]*pbmodel.GroupInfo, len(lst))
	for index, item := range lst {
		g, _ := LoadGroupInfoByGId(item)
		ret[index] = g
	}

	return ret
}

// 目前存储的仅仅是一组数字
func gStoreList2GroupList(lst []model.UserInGStore) []*pbmodel.GroupInfo {
	if lst == nil || len(lst) == 0 {
		return nil
	}

	ret := make([]*pbmodel.GroupInfo, len(lst))
	for index, item := range lst {
		g, _ := LoadGroupInfoByGId(item.Gid)
		ret[index] = g
	}

	return ret
}

// 尝试逐级的加载群组信息
func LoadGroupInfoByGId(gid int64) (*pbmodel.GroupInfo, error) {
	// 尝试内存加载
	group, ok := Globals.grc.GetGroup(gid)
	if ok {
		Globals.Logger.Debug("LoadGroupInfoById ok from memory", zap.Int64("gid", gid))
		return group.GetGroupInfo(), nil
	}

	// 尝试redis加载
	g, err := Globals.redisCli.GetGroupInfoById(gid)
	if err == nil && g != nil {
		group = model.NewGroupFromInfo(g)
		Globals.grc.InsertGroup(gid, group)
		Globals.Logger.Debug("LoadGroupInfoById ok from redis", zap.Int64("gid", gid))
		return g, nil
	}
	// 尝试mongodb加载
	gl, err := Globals.mongoCli.FindGroupById(gid)
	if err == nil && gl != nil && len(gl) > 0 {
		Globals.Logger.Debug("LoadGroupInfoById ok from mongodb", zap.Int64("gid", gid))
		group = model.NewGroupFromInfo(g)
		Globals.grc.InsertGroup(gid, group)
		Globals.redisCli.SetGroupInfo(g)
		return g, nil
	}

	g = &pbmodel.GroupInfo{
		GroupId: gid}
	Globals.Logger.Error("LoadGroupInfoById error, can't find in mongodb", zap.Int64("gid", gid))

	return g, nil
}
