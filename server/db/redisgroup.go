package db

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"strconv"
)

// 从mongoDB加载后，需要设置到redis中
func (cli *RedisClient) SetGroupInfo(grp *pbmodel.GroupInfo) error {
	keyName := GetGroupInfoKey(grp.GroupId)
	mapUser, err := utils.AnyToMap(grp, nil)
	//fmt.Println(mapUser)
	ret, err := cli.Cmd.HMSet(keyName, mapUser).Result()
	fmt.Println(ret)
	return err
}

func (cli *RedisClient) SetGroupInfoPart(gid int64, field, value string) error {
	keyName := GetGroupInfoKey(gid)
	_, err := cli.Cmd.HSet(keyName, field, value).Result()
	return err
}

func (cli *RedisClient) GetGroupInfoById(gid int64) (*pbmodel.GroupInfo, error) {
	keyName := GetGroupInfoKey(gid)
	data, err := cli.Cmd.HGetAll(keyName).Result()
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

// ////////////////////////////////////////////////////////////////////////////////
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

// 解散
func (cli *RedisClient) RemoveAllUserOfGroup(gid int64) error {
	keyName := GetGroupAllMembersKey(gid)
	return cli.RemoveKey(keyName)
}

// 退出群聊的的删除
func (cli *RedisClient) RemoveGroupMembers(gid int64, members []int64) (int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.RemoveIntSet(keyName, members)
}

// 获取所有的用户成员
func (cli *RedisClient) GetGroupMembers(gid int64) ([]int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.GetIntSet(keyName)
}

// 计算成员个数
func (cli *RedisClient) GetGroupMembersCount(gid int64) (int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.GetSetLen(keyName)
}

func (cli *RedisClient) GetGroupMembersPage(gid, offset, pageSize int64) (uint64, []int64, error) {
	keyName := GetGroupAllMembersKey(gid)
	return cli.ScanIntSet(keyName, uint64(offset), pageSize)
}

// ////////////////////////////////////////////////////////////////////////////////
// 某个服务器上，群组在线成员
// 设置成员SET表, 这个函数也只有创建群的时候使用，其他时候都是上线加入
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

// 添加到成员SET表，用户上线时候使用
// 首先设置bsgdi_1001_1 set中添加用户，如果返回正数，说明是新用户，需要计数器加N
func (cli *RedisClient) AddActiveGroupMembers(gid, nodeId int64, members []int64) (int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	hashKey := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)

	count, err := cli.AddIntSet(keyName, members)
	if err != nil {
		return count, err
	}
	if count == 0 {
		count, err = cli.GetSetLen(keyName)
		cli.SetHashKeyInt(hashKey, field, count)
		return count, err
	}

	n, err := cli.AddHashKeyInt(hashKey, field, count)
	return n, err
}

// 使用LUA脚本一次性添加到set 和hash中，网络开销小
func (cli *RedisClient) AddActiveGroupMembersLua(gid, nodeId int64, members []int64) (int64, error) {
	lua := `local setKey = KEYS[1]
	local hashKey = KEYS[2]
	local field = KEYS[3]
	for _, member in ipairs(ARGV) do
	redis.call("SADD", setKey, member)
	end

	local length = redis.call("SCARD", setKey)
	redis.call("HSET", hashKey, field, length)

	return length`

	setKey := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	hashKey := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)

	strMem := make([]string, len(members))
	for i, m := range members {
		strMem[i] = strconv.FormatInt(m, 10)
	}

	// 执行 Lua 脚本
	ret, err := cli.Cmd.Eval(lua, []string{setKey, hashKey, field}, strMem).Result()
	if err != nil {
		return 0, err
	}

	if value, ok := ret.(int64); ok {
		return value, nil
	}
	return 0, errors.New("return value error")

}

// 退出去聊的的删除， 下线的也删除
func (cli *RedisClient) RemoveActiveGroupMembers(gid, nodeId int64, members []int64) (int64, error) {
	setKey := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	hashKey := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)

	count, err := cli.RemoveIntSet(setKey, members)
	if err != nil {
		count, err = cli.GetSetLen(setKey)
		cli.SetHashKeyInt(hashKey, field, count)
		return count, err
	}

	count, err = cli.AddHashKeyInt(hashKey, field, 0-count)

	return count, err
}

// 解散时候，记得清楚相关的群组成员分布数据
func (cli *RedisClient) RemoveActiveGroupRelated(gid int64) error {
	setKey := GetGroupActiveMemsPerNodeKey(gid, 1)
	cli.RemoveKey(setKey)

	hashKey := GetGroupMemNumPerNodeKey(gid)
	cli.RemoveKey(hashKey)

	return nil
}

// bsgdi_1001_1 set中删除相关用户，同时重新计算计数
func (cli *RedisClient) RemoveActiveGroupMembersLua(gid, nodeId int64, members []int64) (int64, error) {
	lua := `local setKey = KEYS[1]
	local hashKey = KEYS[2]
	local field = KEYS[3]
	for _, member in ipairs(ARGV) do
	redis.call("SREM", setKey, member)
	end

	local length = redis.call("SCARD", setKey)
	redis.call("HSET", hashKey, field, length)

	return length`

	setKey := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	hashKey := GetGroupMemNumPerNodeKey(gid)
	field := GetServerField(nodeId)

	strMem := make([]string, len(members))
	for i, m := range members {
		strMem[i] = strconv.FormatInt(m, 10)
	}

	// 执行 Lua 脚本
	ret, err := cli.Cmd.Eval(lua, []string{setKey, hashKey, field}, strMem).Result()
	if err != nil {
		return 0, err
	}

	if value, ok := ret.(int64); ok {
		return value, nil
	}
	return 0, errors.New("return value error")
}

// 获取所有的用户成员，跨服务器转发时候有用
func (cli *RedisClient) GetActiveGroupMembers(gid, nodeId int64) ([]int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	return cli.GetIntSet(keyName)
}

// 获取群在各个服务器上活跃用户数量，直接根据SET计算长度
func (cli *RedisClient) GetActiveGroupMemberCount(gid, nodeId int64) (int64, error) {
	keyName := GetGroupActiveMemsPerNodeKey(gid, nodeId)
	return cli.GetSetLen(keyName)
}

// 从服务器到活跃用户数量的映射，这个应该不用分页了，服务器数量不会太多，即使10000也可以处理
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
	tx := cli.Cmd.TxPipeline()
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
	result, err := cli.Cmd.LRange(key, 0, count).Result()
	return result, err
}

func (cli *RedisClient) GetGroupLatestMsgPage(gid, offset, count int64) ([]string, error) {
	// 获取群组消息缓存表的键名
	key := GetGroupMsgCacheKey(gid)
	// 获取左侧的100条消息
	result, err := cli.Cmd.LRange(key, offset, count).Result()
	return result, err
}

// 查看当前缓存有多少条，一个持续运行的群，消息正常都是满的
func (cli *RedisClient) GetGroupLatestMsgCount(gid, count int64) (int64, error) {
	key := GetGroupMsgCacheKey(gid)
	cmd := cli.Cmd.LLen(key)
	return cmd.Result()
}
