package db

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"fmt"
)

// 从mongoDB加载后，需要设置到redis中
func (cli *RedisClient) SetGroupInfo(grp *pbmodel.GroupInfo) error {
	keyName := GetGroupInfoKey(grp.GroupId)
	mapUser, err := utils.AnyToMap(grp, nil)
	//fmt.Println(mapUser)
	ret, err := cli.Db.HMSet(keyName, mapUser).Result()
	fmt.Println(ret)
	return err
}

func (cli *RedisClient) GetGroupInfoById(gid int64) (*pbmodel.GroupInfo, error) {
	keyName := GetGroupInfoKey(gid)
	data, err := cli.Db.HGetAll(keyName).Result()
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	//fmt.Println("from redis get the map = ", data)

	group := pbmodel.GroupInfo{}
	err = utils.FromMapString(data, &group)

	//fmt.Println(user)
	return &group, err
}

// 设置成员SET表
func (cli *RedisClient) SetGroupMembers(gid int64, members []int64) error {
	keyName := GetGroupAllMembersKey(gid)
	return cli.SetIntSet(keyName, members)
}

// 添加到成员SET表
func (cli *RedisClient) AddGroupMembers(gid int64, members []int64) (int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.AddIntSet(keyName, members)
}

// 退出去聊的的删除
func (cli *RedisClient) RemoveGroupMembers(gid int64, members []int64) (int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.RemoveIntSet(keyName, members)
}

// 获取所有的用户成员
func (cli *RedisClient) GetGroupMembers(gid int64) ([]int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.GetIntSet(keyName)
}

// ////////////////////////////////////////////////////////////////////////////////
// 某个服务器上，群组在线成员
// 设置成员SET表
func (cli *RedisClient) SetActiveGroupMembers(gid, nodeId int64, members []int64) error {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	err := cli.SetIntSet(keyName, members)
	if err != nil {
		return err
	}

	key := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)
	count := int64(len(members))
	_, err = cli.SetHashKeyInt(key, field, count)

	return err
}

// 添加到成员SET表
func (cli *RedisClient) AddActiveGroupMembers(gid, nodeId int64, members []int64) (int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	count, err := cli.AddIntSet(keyName, members)
	if err != nil {
		return count, err
	}
	if count == 0 {
		return count, err
	}

	key := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)
	n, err := cli.AddHashKeyInt(key, field, count)

	return n, err
}

// 退出去聊的的删除， 下线的也删除
func (cli *RedisClient) RemoveActiveGroupMembers(gid, nodeId int64, members []int64) (int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	count, err := cli.RemoveIntSet(keyName, members)
	if err != nil {
		return count, err
	}

	key := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)
	_, err = cli.SetHashKeyInt(key, field, count)

	return count, err
}

// 获取所有的用户成员，跨服务器转发时候有用
func (cli *RedisClient) GetActiveGroupMembers(gid, nodeId int64) ([]int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	return cli.GetIntSet(keyName)
}

// 获取群在各个服务器上活跃用户数量
func (cli *RedisClient) GetActiveGroupMemberCount(gid, nodeId int64) (int64, error) {
	key := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)
	return cli.GetHashKeyInt(key, field)
}

// 从服务器到活跃用户数量的映射
func (cli *RedisClient) GetActiveGroupMemberCountList(gid int64) (map[int16]int64, error) {
	key := GetGroupMemNumPerNodeKey(gid)
	return cli.GetHashKeyIntList(key)
}

// /////////////////////////////////////////////////////////////////////////////
// 群组消息缓存表
// 左侧插入队列，如果超过1000条，右侧弹出，使用事务
func (cli *RedisClient) PushGroupMsg(gid int64, msg string) error {
	// 获取群组消息缓存表的键名
	key := GetGroupMsgCacheKey(gid)
	// 开启 Redis 事务
	tx := cli.Db.TxPipeline()
	// 左侧插入队列
	tx.LPush(key, msg)
	// 获取列表长度
	//tx.LLen(key)
	// 如果列表长度超过 1000 条，则右侧弹出
	tx.LTrim(key, 0, 999)
	// 执行事务
	_, err := tx.Exec()
	return err

}
func (cli *RedisClient) GetGroupLatestMsg(gid, count int64) ([]string, error) {
	// 获取群组消息缓存表的键名
	key := GetGroupMsgCacheKey(gid)
	// 获取左侧的100条消息
	result, err := cli.Db.LRange(key, 0, count).Result()
	return result, err
}
