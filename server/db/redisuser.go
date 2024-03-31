package db

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"fmt"
	"strconv"
)

func userInfoToMap(userInfo *pbmodel.UserInfo) (map[string]interface{}, error) {
	return utils.AnyToMap(userInfo, nil)
}

func mapToUserInfo(data map[string]string) (*pbmodel.UserInfo, error) {
	user := pbmodel.UserInfo{}
	err := utils.FromMapString(data, &user)
	return &user, err
}

func (cli *RedisClient) FindUserById(uid int64) (*pbmodel.UserInfo, error) {
	keyName := GetUserInfoKey(uid)
	//fmt.Println(tblName)
	data, err := cli.Db.HGetAll(keyName).Result()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
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
	ret, err := cli.Db.HMSet(keyName, mapUser).Result()
	fmt.Println(ret)
	return err
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

// 设置粉丝列表
func (cli *RedisClient) SetUserFans(uid int64, friends []model.FriendStore) error {
	keyName := GetUserFansKey(uid)
	return cli.setFriendsStores(keyName, friends)
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

// 这里不返回昵称，直接返回掩码
func (cli *RedisClient) GetUserBLocks(uid int64, offset uint64) (uint64, map[int64]int32, error) {
	key := GetUserBlockKey(uid)
	off, aMap, err := cli.ScanHashKeys(key, offset, MaxFriendBatchSize)
	if err != nil {
		return 0, nil, err
	}
	intMap := make(map[int64]int32)
	for k, v := range aMap {
		intkey, e := strconv.ParseInt(k, 10, 64)
		if e != nil {
			continue
		}

		intValue, e := strconv.ParseInt(v, 10, 32)
		if e != nil {
			continue
		}
		intMap[intkey] = int32(intValue)
	}

	return off, intMap, nil
}

// //////////////////////////////////////////////////////////////////////////
// 移除关注
func (cli *RedisClient) RemoveUserFollowing(uid int64, friends []int64) error {
	keyName := GetUserFollowingKey(uid)
	return cli.RemoveHashMapWithIntFields(keyName, friends)
}

// 设置粉丝列表
func (cli *RedisClient) RemoveUserFans(uid int64, friends []int64) error {
	keyName := GetUserFansKey(uid)
	return cli.RemoveHashMapWithIntFields(keyName, friends)
}

// 设置拉黑列表
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
