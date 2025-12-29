package db

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const DefaultUserTTL = time.Hour * 168

func userInfoToMap(userInfo *pbmodel.UserInfo) (map[string]interface{}, error) {
	return utils.AnyToMap(userInfo, nil)
}

func mapToUserInfo(data map[string]string) (*pbmodel.UserInfo, error) {
	user := pbmodel.UserInfo{}
	err := utils.FromMapString(data, &user)
	return &user, err
}

func (cli *RedisClient) RemoveUser(uid int64) error {
	keyName := GetUserInfoKey(uid)
	_, err := cli.Cmd.Del(keyName).Result()
	return err
}

// 查找用户
func (cli *RedisClient) FindUserById(uid int64) (*pbmodel.UserInfo, error) {
	keyName := GetUserInfoKey(uid)
	//fmt.Println(tblName)
	data, err := cli.Cmd.HGetAll(keyName).Result()
	if err != nil {
		//fmt.Println(err.Error())
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("not find user")
	}
	//fmt.Println("from redis get the map = ", data)
	user, err := mapToUserInfo(data)

	//fmt.Println(user)
	return user, err

}

// 保存一个完整的user
func (cli *RedisClient) SetUserInfo(user *pbmodel.UserInfo) error {
	keyName := GetUserInfoKey(user.UserId)
	mapUser, err := userInfoToMap(user)
	ret, err := cli.Cmd.HMSet(keyName, mapUser).Result()
	fmt.Println(ret)
	return err
}

// 设置部分内容
func (cli *RedisClient) UpdateUserInfoPart(id int64, setData map[string]interface{}, unsetData []string) error {
	keyName := GetUserInfoKey(id)

	tx := cli.Cmd.TxPipeline()
	// 清空集合
	if setData != nil {
		tx.HMSet(keyName, setData)
	}

	if unsetData != nil {
		tx.HDel(keyName, unsetData...)
	}

	// 执行事务
	_, err := tx.Exec()
	return err
}

func (cli *RedisClient) IncUserInfoFiledByInt(id int64, field string, num int64) (int64, error) {
	keyName := GetUserInfoKey(id)
	n, err := cli.AddHashKeyInt(keyName, field, num)
	return n, err
}

// ////////////////////////////////////////////////////////////////////////
// 设置关注列表, 这个参数与操作数据库的结构一样，保存数据库后直接添加到redis里
func (cli *RedisClient) setFriendsStores(key string, friends []model.FriendStore) error {
	aMap := make(map[string]interface{})
	for _, f := range friends {
		name := strconv.FormatInt(f.Uid2, 10)
		aMap[name] = f.Nick
	}
	return cli.SetHashMap(key, aMap)
}

func (cli *RedisClient) SetUserFollowing(uid int64, friends []model.FriendStore) error {
	keyName := GetUserFollowingKey(uid)
	return cli.setFriendsStores(keyName, friends)
}

func (cli *RedisClient) SetUserMutual(uid int64, friends []model.FriendStore) error {
	keyName := GetUserMutualFriendsKey(uid)
	return cli.setFriendsStores(keyName, friends)
}

func (cli *RedisClient) SetUserFollowingNick(uid, fid int64, nick string) error {
	keyName := GetUserFollowingKey(uid)
	field := strconv.FormatInt(fid, 10)

	_, err := cli.UpdateHashFieldIfExist(keyName, field, nick)
	return err
}

// 设置粉丝列表
func (cli *RedisClient) SetUserFans(uid int64, friends []model.FriendStore) error {
	keyName := GetUserFansKey(uid)
	return cli.setFriendsStores(keyName, friends)
}
func (cli *RedisClient) SetUserFansNick(uid, fid int64, nick string) error {
	keyName := GetUserFansKey(uid)
	field := strconv.FormatInt(fid, 10)
	_, err := cli.UpdateHashFieldIfExist(keyName, field, nick)
	return err
}

// 设置拉黑列表
func (cli *RedisClient) SetUserBlocks(uid int64, friends []model.BlockStore) error {
	key := GetUserBlockKey(uid)
	aMap := make(map[string]interface{})
	for _, f := range friends {
		name := strconv.FormatInt(f.Uid2, 10)
		aMap[name] = f.Perm
	}
	return cli.SetHashMap(key, aMap)
}

// ////////////////////////////////////////////////////
func (cli *RedisClient) addFriendsStores(key string, friends []model.FriendStore) error {
	aMap := make(map[string]interface{})
	for _, f := range friends {
		name := strconv.FormatInt(f.Uid2, 10)
		aMap[name] = f.Nick
	}
	return cli.AddHashMap(key, aMap)
}
func (cli *RedisClient) AddUserFollowing(uid int64, friends []model.FriendStore) error {
	keyName := GetUserFollowingKey(uid)
	return cli.addFriendsStores(keyName, friends)
}

// 设置粉丝列表
func (cli *RedisClient) AddUserFans(uid int64, friends []model.FriendStore) error {
	keyName := GetUserFansKey(uid)
	return cli.addFriendsStores(keyName, friends)
}

// 设置拉黑列表
func (cli *RedisClient) AddUserBlocks(uid int64, friends []model.BlockStore) error {
	key := GetUserBlockKey(uid)
	aMap := make(map[string]interface{})
	for _, f := range friends {
		name := strconv.FormatInt(f.Uid2, 10)
		aMap[name] = f.Perm
	}
	return cli.AddHashMap(key, aMap)
}

// //////////////////////////////////////////////////////////////////////////////////////
// 返回好友id-> 昵称的map
func (cli *RedisClient) getUserFriendStore(key string, offset uint64) (uint64, map[int64]string, error) {
	off, aMap, err := cli.ScanHashKeys(key, offset, MaxFriendBatchSize)
	if err != nil {
		return 0, nil, err
	}
	intMap := make(map[int64]string)
	for k, v := range aMap {
		i, e := strconv.ParseInt(k, 10, 64)
		if e != nil {
			continue
		}
		intMap[i] = v
	}

	return off, intMap, nil
}

func (cli *RedisClient) GetUserFollowing(uid int64, offset uint64) (uint64, map[int64]string, error) {
	keyName := GetUserFollowingKey(uid)
	return cli.getUserFriendStore(keyName, offset)
}

func (cli *RedisClient) GetUserFans(uid int64, offset uint64) (uint64, map[int64]string, error) {
	keyName := GetUserFansKey(uid)
	return cli.getUserFriendStore(keyName, offset)
}

func (cli *RedisClient) ExistFollowing(uid int64) (bool, error) {
	keyName := GetUserFollowingKey(uid)
	return cli.HasKey(keyName, DefaultUserTTL)
}

func (cli *RedisClient) ExistFans(uid int64) (bool, error) {
	keyName := GetUserFansKey(uid)
	return cli.HasKey(keyName, DefaultUserTTL)
}

func (cli *RedisClient) ExistPermission(uid int64) (bool, error) {
	keyName := GetUserBlockKey(uid)
	return cli.HasKey(keyName, DefaultUserTTL)
}

func (cli *RedisClient) ExistUserInfo(uid int64) (bool, error) {
	keyName := GetUserInfoKey(uid)
	return cli.HasKey(keyName, DefaultUserTTL)
}

// 检查是否存在好友，如果设置了空字符串，就是非好友
// 如果没有，error = redis: nil
func (cli *RedisClient) CheckUserFan(uid int64, fid int64) (bool, error) {
	keyName := GetUserFansKey(uid)
	field := strconv.FormatInt(fid, 10)
	str, err := cli.Cmd.HGet(keyName, field).Result()
	if err == nil {
		if str == "##" {
			return false, nil
		}
		return true, nil
	}

	return false, err
}

// 查看用户的好友的权限
func (cli *RedisClient) CheckUserPermission(uid int64, fid int64) (uint32, error) {
	key := GetUserBlockKey(uid)
	field := strconv.FormatInt(fid, 10)

	perm, err := cli.GetHashKeyInt(key, field)
	return uint32(perm), err
}

// 设置单个的用户权限
func (cli *RedisClient) AddUserPermission(uid int64, fid int64, perm uint32) error {
	key := GetUserBlockKey(uid)
	field := strconv.FormatInt(fid, 10)

	_, err := cli.AddHashKeyInt(key, field, int64(perm))
	return err
}

// 这里不返回昵称，直接返回掩码
func (cli *RedisClient) GetUserBLocks(uid int64, offset uint64) (uint64, map[int64]uint32, error) {
	key := GetUserBlockKey(uid)
	off, aMap, err := cli.ScanHashKeys(key, offset, MaxFriendBatchSize)
	if err != nil {
		return 0, nil, err
	}
	intMap := make(map[int64]uint32)
	for k, v := range aMap {
		intkey, e := strconv.ParseInt(k, 10, 64)
		if e != nil {
			continue
		}

		intValue, e := strconv.ParseInt(v, 10, 32)
		if e != nil {
			continue
		}
		intMap[intkey] = uint32(intValue)
	}

	return off, intMap, nil
}

// //////////////////////////////////////////////////////////////////////////
// 移除关注
func (cli *RedisClient) RemoveUserFollowing(uid int64, friends []int64) error {
	keyName := GetUserFollowingKey(uid)
	return cli.RemoveHashMapWithIntFields(keyName, friends)
}

// 粉丝列表
func (cli *RedisClient) RemoveUserFans(uid int64, friends []int64) error {
	keyName := GetUserFansKey(uid)
	return cli.RemoveHashMapWithIntFields(keyName, friends)
}

// 拉黑列表
func (cli *RedisClient) RemoveUserBlocks(uid int64, friends []int64) error {
	keyName := GetUserBlockKey(uid)
	return cli.RemoveHashMapWithIntFields(keyName, friends)
}

// 求粉丝和关注的交集，那么就是双向好友了
func (cli *RedisClient) GetFriendIntersect(uid int64) ([]int64, error) {
	key1 := GetUserFollowingKey(uid)
	key2 := GetUserFansKey(uid)
	friends, err := cli.GetHashIntersect(key1, key2)
	if err != nil {
		return nil, err
	}

	intList := make([]int64, len(friends))
	for i, item := range friends {
		intValue, err1 := strconv.ParseInt(item, 10, 64)
		if err1 != nil {
			continue
		}
		intList[i] = intValue
	}

	return intList, nil
}

// ////////////////////////////////////////////////////////////////////
// 用户登录初始化加载时候使用
func (cli *RedisClient) SetUserInGroup(uid int64, gidList []int64) error {
	keyUserInG := GetUseringKey(uid)
	return cli.SetIntSet(keyUserInG, gidList)
}

func (cli *RedisClient) HasUserInGroup(uid int64) (bool, error) {
	keyUserInG := GetUseringKey(uid)
	return cli.HasKey(keyUserInG, DefaultUserTTL)
}

// 用户加入的群组个数
func (cli *RedisClient) GetUserInGroupCount(uid int64) (int64, error) {
	keyUserInG := GetUseringKey(uid)
	return cli.GetSetLen(keyUserInG)
}

// 直接返回用户所在的所有的群组
func (cli *RedisClient) GetUserInGroupAll(uid int64) ([]int64, error) {
	keyUserInG := GetUseringKey(uid)
	return cli.GetIntSet(keyUserInG)
}

// 加入有机器人用户，可能会加入太多的群组，一般应该在应用层限制，防止过多
func (cli *RedisClient) GetUserInGroupPage(uid, offset, pageSize int64) (uint64, []int64, error) {
	keyUserInG := GetUseringKey(uid)
	return cli.ScanIntSet(keyUserInG, uint64(offset), pageSize)
}

// 用户加入群组，处理2处
// 1） 用户参加的群列表
// 2) 群全员用户的hash表
func (cli *RedisClient) SetUserJoinGroup(uid, gid int64, nick string) error {

	keyUserInG := GetUseringKey(uid)
	keyGroupMem := GetGroupAllMembersKey(gid)
	idStr := strconv.FormatInt(uid, 10)
	// 创建事务
	tx := cli.Cmd.TxPipeline()
	// 清空集合
	tx.SAdd(keyUserInG, gid)
	tx.Expire(keyUserInG, DefaultUserTTL) // 如果不登录，7天后消失
	//tx.SAdd(keyGroupMem, uid)
	tx.HSet(keyGroupMem, idStr, nick)
	// 执行事务
	_, err := tx.Exec()

	return err
}

// 用户退出群组
func (cli *RedisClient) SetUserLeaveGroup(uid, gid int64) error {

	keyUserInG := GetUseringKey(uid)
	keyGroupMem := GetGroupAllMembersKey(gid)
	idStr := strconv.FormatInt(uid, 10)
	// 创建事务
	tx := cli.Cmd.TxPipeline()
	// 清空集合
	tx.SRem(keyUserInG, gid)
	//tx.SRem(keyGroupMem, uid)
	tx.HDel(keyGroupMem, idStr)
	// 执行事务
	_, err := tx.Exec()

	return err
}

// 求2个用户的共同的所在的群组
func (cli *RedisClient) GetUsersInSameGroup(uid1, uid2 int64) ([]int64, error) {

	key1 := GetUseringKey(uid1)
	key2 := GetUseringKey(uid2)
	return cli.IntersectIntSets(key1, key2)
}

// 更新用户相关的TTL
// 用户基础信息表7天bsui_
// 用户好友权限表7天bsufb_
// 用户关注表7天bsufo_
// 用户粉丝表7天bsufa_
// 用户的指纹表7天bsut_
// 用户所属群组7天bsuing_
// 用户session 分布表  30分钟 bsud_
func (cli *RedisClient) UpdateUserTTL(uid int64) (int, error) {
	keyUserInfo := GetUserInfoKey(uid)
	keyUserFollow := GetUserFollowingKey(uid)
	keyUserFan := GetUserFansKey(uid)
	keyUserPermission := GetUserBlockKey(uid)
	keyUserToken := GetUserTokenKey(uid)
	keyUserInGroup := GetUseringKey(uid)
	keys := []string{keyUserInfo, keyUserFollow, keyUserFan, keyUserPermission, keyUserToken, keyUserInGroup}

	count, err := cli.SetKeysExpire(keys, DefaultUserTTL) // 7 天

	//keyUserDistribution := GetUserDistributionKey(uid)
	//err = cli.SetKeyExpire(keyUserDistribution, time.Minute*30)
	//if err != nil {
	//	count++
	//}

	return count, err
}

// 登录时候，设置用户会话在某个服务器上，设置超时时间为30分钟
func (cli *RedisClient) SetUserSessionOnServer(uid, sid, serverIndex int64) error {
	keyUserDistribution := GetUserDistributionKey(uid)
	field := strconv.FormatInt(sid, 10)

	// 使用管道
	pipe := cli.Cmd.Pipeline()

	pipe.HSet(keyUserDistribution, field, serverIndex)
	pipe.Expire(keyUserDistribution, time.Minute*30)

	// 执行管道操作
	_, err := pipe.Exec()

	return err
}

// 会话断开时候，删除会话标记
func (cli *RedisClient) RemoveUserSessionOnServer(uid, sid int64) error {
	keyUserDistribution := GetUserDistributionKey(uid)
	field := strconv.FormatInt(sid, 10)
	return cli.RemoveHashInt(keyUserDistribution, field)
}

func (cli *RedisClient) GetUserSessionOnServer(uid int64) (map[int64]int32, error) {
	keyUserDistribution := GetUserDistributionKey(uid)
	return cli.GetHashKeyInt64List(keyUserDistribution)
}

// 设置好友的关注和粉丝的个数，永久有效，用户注册时候就加载
func (cli *RedisClient) SetUserFriendNum(uid, numFollow, numFans, numFriends int64) error {

	key := GetUserFriendNumKey(uid)
	// 使用管道
	pipe := cli.Cmd.Pipeline()

	pipe.HSet(key, "follows", numFollow)
	pipe.HSet(key, "fans", numFans)
	pipe.HSet(key, "friends", numFriends)
	// 执行管道操作
	_, err := pipe.Exec()
	return err
}

func (cli *RedisClient) GetUserFriendNum(uid int64) (int64, int64, int64, error) {
	key := GetUserFriendNumKey(uid)
	mapRet, err := cli.Cmd.HGetAll(key).Result()
	if err != nil {
		return 0, 0, 0, err
	}
	numFollow := int64(0)
	numFans := int64(0)
	numFriends := int64(0)
	str, ok := mapRet["follows"]
	if ok {
		numFollow, err = strconv.ParseInt(str, 10, 64)
	}

	str, ok = mapRet["fans"]
	if ok {
		numFans, err = strconv.ParseInt(str, 10, 64)
	}

	str, ok = mapRet["friends"]
	if ok {
		numFriends, err = strconv.ParseInt(str, 10, 64)
	}

	return numFollow, numFans, numFriends, err
}

// 不增加的就设置为0
func (cli *RedisClient) AddUserFriendNum(uid, numFollow, numFans int64) error {
	key := GetUserFriendNumKey(uid)
	if numFollow > 0 {
		return cli.Cmd.HIncrBy(key, "follows", numFollow).Err()
	}
	if numFans > 0 {
		return cli.Cmd.HIncrBy(key, "fans", numFollow).Err()
	}
	return nil
}

// 单独的增加或者减少，可以设置1或者-1
func (cli *RedisClient) AddUserFollowsNum(uid, num int64) error {
	key := GetUserFriendNumKey(uid)

	return cli.Cmd.HIncrBy(key, "follows", num).Err()
}

func (cli *RedisClient) AddUserFansNum(uid, num int64) error {
	key := GetUserFriendNumKey(uid)
	return cli.Cmd.HIncrBy(key, "fans", num).Err()

}

func (cli *RedisClient) AddUserFriendsNum(uid, num int64) error {
	key := GetUserFriendNumKey(uid)
	return cli.Cmd.HIncrBy(key, "friends", num).Err()
}
