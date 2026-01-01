package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sort"
)

// 这个文件主要是处理三级缓存相关的逻辑
// 当内存没有数据时候，群在加载时候就初始化了成员列表，

// 从内存或者redis数据库找群成员
func findGroupMembersFromId(group *model.Group, fromId int64, pageSize int) []*pbmodel.GroupMember {
	members := group.GetRawMembers()

	// 2. 自定义比较函数：按 ID 升序排列
	sort.Slice(members, func(i, j int) bool {
		// 返回 true 表示 groups[i] 应该排在 groups[j] 前面
		iId := int64(0)
		if members[i].U != nil {
			iId = members[i].U.UserId
		}
		jId := int64(0)
		if members[j].U != nil {
			jId = members[j].U.UserId
		}
		return iId < jId
	})

	Globals.Logger.Debug("findGroupMembersFromId() ", zap.Int64("groupId", group.GroupId))
	// 3. 遍历输出排序后的结果
	for i, m := range members {
		if m.U != nil {
			Globals.Logger.Debug(fmt.Sprintf("gid=%d, index=%d, ID: %d, Name: %s\n", group.GroupId, i, m.U.UserId, m.Nick))
		} else {
			Globals.Logger.Debug(fmt.Sprintf("gid=%d, index=%d, ID: %d, Name: %s\n", group.GroupId, i, 0, m.Nick))
		}
	}

	// 格式转换
	sz := pageSize
	if len(members) < pageSize {
		sz = len(members)
	}
	groupMembers := make([]*pbmodel.GroupMember, 0, sz)
	count := int(0)
	for _, m := range members {
		if (m.U == nil) || (m.U.UserId < fromId) {
			continue
		}

		data := &pbmodel.GroupMember{
			UserId:  m.U.UserId,
			Nick:    m.Nick,
			Icon:    m.U.Icon,
			Role:    group.GetMemberRoleString(m.U.UserId),
			GroupId: 0,
			Params:  nil,
		}
		groupMembers = append(groupMembers, data)
		count++
		if count >= pageSize {
			break
		}

	}
	return groupMembers
}

// 根据数据库或者redis中记录的群用户列表，查找用户信息
// 这里保证不为空
func getUsersFromMemberStores(lst []model.GroupMemberStore) []*model.User {
	users := make([]*model.User, len(lst))
	for i, m := range lst {
		u, _, err := findUser(m.Uid)
		if err != nil {
			Globals.Logger.Error("getUsersFromMemberStores() load user but nil", zap.Error(err), zap.Int64("uid", m.Uid))
			users[i] = model.NewUser()
			users[i].UserId = m.Uid
			users[i].NickName = m.Nick
		} else {
			users[i] = u
		}

	}
	return users
}

// 当群内没有成员信息时候，加载群的用户信息; 群的数量为2000上限，所以这里直接从数据库加载所有群用户，放到redis中
func loadGroupMembers(group *model.Group) {
	// 从redis中查找，如果找到了，
	lst, err := Globals.redisCli.GetGroupMembers(group.GroupId)
	if lst != nil && len(lst) > 0 {
		Globals.Logger.Debug("loadGroupMembers() find group members in redis", zap.Int64("group_id", group.GroupId), zap.Any("lst", lst))
		users := getUsersFromMemberStores(lst)
		err = group.SetMembers(lst, users)
		if err != nil {
			Globals.Logger.Error("loadGroupMembers() set group members error", zap.Error(err))
		}
		Globals.Logger.Debug("loadGroupMembers() re-print group member in memory", zap.Int64("group_id", group.GroupId), zap.Any("group", group.Members))
		return
	}

	// 从数据库中加载
	memlist, err := Globals.scyllaCli.FindGroupMembers(db.ComputePk(group.GroupId), group.GroupId, 0, 2000)
	if err != nil {
		Globals.Logger.Fatal("loadGroupMembers() db error", zap.Error(err))
		return
	}

	Globals.Logger.Debug("loadGroupMembers() find group members in scyllaDb", zap.Int64("group_id", group.GroupId), zap.Any("lst", memlist))
	// 同步到redis
	err = Globals.redisCli.SetGroupMembers(group.GroupId, memlist)
	// 如果错了，大不了下次再加载

	// 同步到内存在中
	if memlist != nil && len(memlist) > 0 {
		users := getUsersFromMemberStores(lst)
		group.SetMembers(memlist, users)
	}

}

// 从数据库中加载Group基础信息
func findGroupAndLoad(gid int64) (*model.Group, error) {
	group, ok := Globals.grc.GetGroup(gid)
	if ok && group != nil {
		// 这里需要检查用户是否为空，如果是空，则应该加载，因为还没有初始化
		if len(group.Members) == 0 {
			Globals.Logger.Debug("find group in memory, but need ->loadGroupMembers()", zap.Int64("group_id", gid))
			loadGroupMembers(group)
		}
		return group, nil
	}

	// 从redis中查找群
	groupInfo, err := Globals.redisCli.GetGroupInfoById(gid)
	if groupInfo != nil {
		Globals.Logger.Debug("find group in redis", zap.Int64("group_id", gid))
		group = model.NewGroupFromInfo(groupInfo)
		loadGroupMembers(group)

		Globals.grc.InsertGroup(gid, group)

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
	loadGroupMembers(group)
	Globals.grc.InsertGroup(gid, group)

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
		g, _ := findGroupAndLoad(item)
		ret[index] = g.GetGroupInfo()
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
		g, _ := findGroupAndLoad(item.Gid)
		ret[index] = g.GetGroupInfo()
	}

	return ret
}
