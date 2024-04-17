package db

import (
	"errors"
	"fmt"
)

// 用于获取几个流水号的键值
const BirdServerUserId = "birds_uid"
const BirdServerGroupId = "birds_gid"
const BirdServerNewsId = "birds_nid"
const BirdServerCommentId = "birds_cid"
const BirdServerUserPrefix = "bsui_%d"      // 用户基础信息以及动态信息 hash
const BirdServerUFolPrefix = "bsufo_%d"     // 关注表 set
const BirdServerUFanPrefix = "bsufa_%d"     // 粉丝表 set
const BirdServerUBloPrefix = "bsufb_%d"     // 拉黑表 hash表
const BirdServerUserInGPrefix = "bsuing_%d" // 用户所属群的集合

const BirdServerGroupPrefix = "bsgi_%d"         // 组基础信息  hash 存储各种属性
const BirdServerGrpUsers = "bsgu_%d"            // 组内所有成员表 hash， 每个成员包括昵称，权限，设置
const BirdServerGrpDistribution = "bsgd_%d"     // hash, 一个群在各个服务器上登录的用户个数 服务器器号->计数
const BirdServerGrpDistriDetail = "bsgdi_%d_%d" // 某一个群，在某个服务器上的成员列表 set

const BirdServerStateCh = "bssch"              //用户、群组、服务器状态广播频道
const BirdServerClusterPrefix = "bscs_%d"      // 集群每个服务器状态hash
const BirdServerCluSerStaPrefix = "bscs_state" // 集群信息hash表

const BirdServerGroupMsgCache = "bsgmsg_%d" //群组数据缓存，如果1000条
const MaxFriendBatchSize = 512              // 最大的批处理的个数
const MaxFriendCacheSize = 1024             // 缓存中粉丝之类的最大数量

// 秘钥存储
const BirdServerTokenPrefix = "bsut_%d" // 秘钥的keyPrint

// 某些值必须要有的，确保初始化
func (cli *RedisClient) initData() {
	if cli == nil {
		return
	}

	cli.makeSureInt(BirdServerUserId, 10000)
	cli.makeSureInt(BirdServerGroupId, 1000)
	cli.makeSureInt(BirdServerNewsId, 0)
	cli.makeSureInt(BirdServerCommentId, 0)

}

func (cli *RedisClient) makeSureInt(key string, def int) error {
	idCmd := cli.Db.Get(key)
	_, err := idCmd.Result()
	if err != nil {
		statusCmd := cli.Db.Set(key, 1000, 0)
		if statusCmd.Err() != nil {
			return statusCmd.Err()
		}
	}

	return nil
}
func (cli *RedisClient) GetNextKeyId(key string) (int64, error) {
	if cli == nil {
		return -1, errors.New("not connected")
	}
	idCmd := cli.Db.Incr(key)
	return idCmd.Val(), idCmd.Err()
}

// 直接获取一段范围的KEY，放在服务器内存中缓存，至少取10个
func (cli *RedisClient) GetNextKeyIdRange(key string, span int64) (int64, int64, error) {
	if cli == nil {
		return -1, -1, errors.New("not connected")
	}
	if span < 10 {
		span = 10
	}
	idCmd := cli.Db.IncrBy(key, span)
	newValue, err := idCmd.Result()

	return newValue - span + 1, newValue, err
}

///////////////////////////////////////////////////////////////////

// 使用普通的key保存一个值，每次都增加一个值，用于计算新增用户的ID
func (cli *RedisClient) GetNextUserId() (int64, error) {
	return cli.GetNextKeyId(BirdServerUserId)
}

func (cli *RedisClient) GetNextGroupId() (int64, error) {
	return cli.GetNextKeyId(BirdServerNewsId)
}

func (cli *RedisClient) GetNextNewsId() (int64, error) {
	return cli.GetNextKeyId(BirdServerCommentId)
}

func (cli *RedisClient) GetNextCommentId() (int64, error) {
	return cli.GetNextKeyId(BirdServerGroupId)
}

/////////////////////////////////////////////////////////////////////////
// 各种键值的名字拼接

// 用户基础信息以及动态信息 hash   bsui_10001
//
//go:inline
func GetUserInfoKey(id int64) string {
	return fmt.Sprintf(BirdServerUserPrefix, id)
}

// 关注表 set "bsufo_"
//
//go:inline
func GetUserFollowingKey(id int64) string {
	return fmt.Sprintf(BirdServerUFolPrefix, id)
}

// 粉丝表 set
//
//go:inline
func GetUserFansKey(id int64) string {
	return fmt.Sprintf(BirdServerUFanPrefix, id)
}

// 拉黑表 hash表
//
//go:inline
func GetUserBlockKey(id int64) string {
	return fmt.Sprintf(BirdServerUBloPrefix, id)
}

// 组基础信息  hash 存储各种属性,
//
//go:inline
func GetGroupInfoKey(id int64) string {
	return fmt.Sprintf(BirdServerGroupPrefix, id)
}

// 组内所有成员表 hash， 每个成员包括昵称，权限，设置"bsgu_"
//
//go:inline
func GetGroupAllMembersKey(id int64) string {
	return fmt.Sprintf(BirdServerGrpUsers, id)
}

// hash, 一个群在各个服务器上登录的用户个数 服务器器号->计数   "bsgd_"
//
//go:inline
func GetGroupMemNumPerNodeKey(id int64) string {
	return fmt.Sprintf(BirdServerGrpDistribution, id)
}

// 某一个群，在某个服务器上的成员列表 set      "bsgdi_%d_%d"
//
//go:inline
func GetGroupActiveMemsPerNodeKey(gid, nodeId int64) string {
	return fmt.Sprintf(BirdServerGrpDistriDetail, gid, nodeId)
}

// 群组数据缓存，如果1000条 "bsgmsg_"
//
//go:inline
func GetGroupMsgCacheKey(id int64) string {
	return fmt.Sprintf(BirdServerGroupMsgCache, id)
}

//go:inline
func GetServerField(id int64) string {
	field := fmt.Sprintf("%d", id)
	return field
}

// 用户所属群
//
//go:inline
func GetUseringKey(id int64) string {

	return fmt.Sprintf(BirdServerUserInGPrefix, id)
}

// 所有的状态广播都使用这一个频道
//
//go:inline
func GetStateChKey() string {
	return BirdServerStateCh
}

// 集群的信息表 "bscs_state", 保存各个服务器的活动时间戳
//
//go:inline
func GetClusterActiveStateKey() string {
	return BirdServerCluSerStaPrefix
}

// 集群每个服务器状态hash
// "bscs_%d"
//
//go:inline
func GetClusterServerStateKey(id int64) string {

	return fmt.Sprintf(BirdServerClusterPrefix, id)
}

// bsut_11122222
//
//go:inline
func GetUserTokenKey(id int64) string {
	return fmt.Sprintf(BirdServerTokenPrefix, id)
}
