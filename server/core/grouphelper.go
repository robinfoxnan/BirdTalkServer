package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"go.uber.org/zap"
)

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
