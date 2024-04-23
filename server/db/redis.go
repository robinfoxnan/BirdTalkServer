package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"strings"
	"time"
)

// 需要手动调用
func InitRedisDb(host string, pwd string) error {
	var err error
	RedisCli, err = NewRedisClient(host, pwd)
	if err == nil {
		RedisCli.initData()
	}
	return err
}

// 通过该变量引用
var RedisCli *RedisClient = nil

var ctx = context.Background()

type RedisClient struct {
	Db               *redis.Client
	Dbs              *redis.ClusterClient
	Cmd              redis.Cmdable
	runningSubscribe int32
}

const dbIndex = 1

// host should be like this: "10.128.5.73:6379"
func NewRedisClient(host string, pwd string) (*RedisClient, error) {
	cli := RedisClient{
		Db:               nil,
		Dbs:              nil,
		Cmd:              nil,
		runningSubscribe: 0,
	}
	var err error
	// 将地址字符串按逗号拆分为字符串切片
	addrs := strings.Split(host, ",")
	if len(addrs) > 1 {
		err, cli.Dbs = initRedisCluster(addrs, pwd)
		if err != nil {
			fmt.Printf("connect redis failed! err : %v\n", err)
			return nil, err
		}
		cli.Cmd = cli.Dbs
	} else {
		err, cli.Db = initRedis(host, pwd, dbIndex)
		if err != nil {
			fmt.Printf("connect redis failed! err : %v\n", err)
			return nil, err
		}
		cli.Cmd = cli.Db
	}

	//atomic.StoreInt32(&cli.runningSubscribe, 0)
	return &cli, err
}

func (cli *RedisClient) Close() {
	if cli.Db != nil {
		cli.Db.Close()
		cli.Db = nil
		cli.Cmd = nil
	}

	if cli.Dbs != nil {
		cli.Dbs.Close()
		cli.Dbs = nil
		cli.Cmd = nil
	}

}

// https://blog.csdn.net/weixin_45901764/article/details/117226225
func initRedis(addr string, password string, dbIndex int) (err error, redisdb *redis.Client) {
	redisOpt := redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
		Network:  "tcp", //网络类型，tcp or unix，默认tcp

		//连接池容量及闲置连接数量
		PoolSize:     150, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,  //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		Dialer: func() (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial("tcp", addr)
		},

		//钩子函数
		OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			//fmt.Printf("conn=%v\n", conn)
			return nil
		},
	}
	redisdb = redis.NewClient(&redisOpt)
	// 判断是否能够链接到数据库
	pong, err := redisdb.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
	}

	//printRedisPool(redisdb.PoolStats())
	return err, redisdb
}

func initRedisCluster(addrs []string, password string) (error, *redis.ClusterClient) {
	// 集群配置选项
	redisOpt := redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
		// 连接池容量及闲置连接数量
		PoolSize:     150, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,  // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		// 超时
		DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		// 闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  // 闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
		// 命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		// 钩子函数
		OnConnect: func(conn *redis.Conn) error { // 仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			// fmt.Printf("conn=%v\n", conn)
			return nil
		},
	}

	// 创建 Redis 集群客户端实例
	redisdbs := redis.NewClusterClient(&redisOpt)
	// 判断是否能够连接到数据库
	pong, err := redisdbs.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
	}

	// printRedisPool(redisdb.PoolStats())
	return err, redisdbs
}

func printRedisPool(stats *redis.PoolStats) {
	fmt.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n",
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

func printRedisOption(opt *redis.Options) {
	fmt.Printf("Network=%v\n", opt.Network)
	fmt.Printf("Addr=%v\n", opt.Addr)
	fmt.Printf("Password=%v\n", opt.Password)
	fmt.Printf("DB=%v\n", opt.DB)
	fmt.Printf("MaxRetries=%v\n", opt.MaxRetries)
	fmt.Printf("MinRetryBackoff=%v\n", opt.MinRetryBackoff)
	fmt.Printf("MaxRetryBackoff=%v\n", opt.MaxRetryBackoff)
	fmt.Printf("DialTimeout=%v\n", opt.DialTimeout)
	fmt.Printf("ReadTimeout=%v\n", opt.ReadTimeout)
	fmt.Printf("WriteTimeout=%v\n", opt.WriteTimeout)
	fmt.Printf("PoolSize=%v\n", opt.PoolSize)
	fmt.Printf("MinIdleConns=%v\n", opt.MinIdleConns)
	fmt.Printf("MaxConnAge=%v\n", opt.MaxConnAge)
	fmt.Printf("PoolTimeout=%v\n", opt.PoolTimeout)
	fmt.Printf("IdleTimeout=%v\n", opt.IdleTimeout)
	fmt.Printf("IdleCheckFrequency=%v\n", opt.IdleCheckFrequency)
	fmt.Printf("TLSConfig=%v\n", opt.TLSConfig)

}

func (cli *RedisClient) CountFields(key string) (int64, error) {
	count, err := cli.Cmd.SCard(key).Result()
	if err != nil {
		// 处理错误
	}
	//fmt.Println("Set 成员个数:", count)
	return count, err
}

func (cli *RedisClient) SetIntSet(key string, intArray []int64) error {
	friendStrs := make([]string, len(intArray))
	// 将 friends 切片中的 int64 转换为字符串切片
	for i, friend := range intArray {
		friendStrs[i] = strconv.FormatInt(friend, 10)
	}

	// 创建事务
	tx := cli.Cmd.TxPipeline()
	// 清空集合
	tx.Del(key)
	tx.SAdd(key, friendStrs)
	// 执行事务
	_, err := tx.Exec()
	if err != nil {
		return err
	}

	return nil
}

// 这里返回的是真正加入的个数，因为重复就不算在添加范围
func (cli *RedisClient) AddIntSet(key string, intArray []int64) (int64, error) {
	friendStrs := make([]string, len(intArray))
	// 将 intArray 切片中的 int64 转换为字符串切片
	for i, friend := range intArray {
		friendStrs[i] = strconv.FormatInt(friend, 10)
	}

	count, err := cli.Cmd.SAdd(key, friendStrs).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 加载set，并转换为
func (cli *RedisClient) GetIntSet(key string) ([]int64, error) {

	members, err := cli.Cmd.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	data := make([]int64, len(members))
	// 将 friends 切片中的 int64 转换为字符串切片
	for i, str := range members {
		data[i], err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (cli *RedisClient) RemoveIntSet(key string, intArray []int64) (int64, error) {
	// 将 intArray 切片中的 int64 转换为字符串切片
	friendStrs := make([]string, len(intArray))
	for i, friend := range intArray {
		friendStrs[i] = strconv.FormatInt(friend, 10)
	}

	count, err := cli.Cmd.SRem(key, friendStrs).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (cli *RedisClient) IntersectSets(set1, set2 string) ([]string, error) {
	// 执行 SINTER 命令来求两个集合的交集
	result, err := cli.Cmd.SInter(set1, set2).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (cli *RedisClient) IntersectIntSets(set1, set2 string) ([]int64, error) {
	// 执行 SINTER 命令来求两个集合的交集
	strList, err := cli.IntersectSets(set1, set2)
	if err != nil {
		return nil, err
	}
	data := make([]int64, len(strList))
	// 将 friends 切片中的 int64 转换为字符串切片
	for i, str := range strList {
		data[i], err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return data, nil

}

// 计算元素集合中的数据个数
func (cli *RedisClient) GetSetLen(key string) (int64, error) {
	count, err := cli.Cmd.SCard(key).Result()
	return count, err
}

// 分页获取群组等set类型的结构
func (cli *RedisClient) ScanIntSet(key string, cursor uint64, pageSize int64) (uint64, []int64, error) {

	// 使用 HSCAN 命令扫描哈希键
	members, nextCursor, err := cli.Cmd.SScan(key, cursor, "", pageSize).Result()
	if err != nil {
		return 0, nil, err
	}

	result := make([]int64, len(members))
	for index, member := range members {
		//fmt.Println("Member:", member)
		result[index], err = strconv.ParseInt(member, 10, 64)
		if err != nil {
			continue
		}
	}

	return nextCursor, result, nil
}

// ///////////////////////////////////////////////////////////////////////////////////
// 哈希表中存储int，是为了统计各个服务器上群组用户分布而使用的
func (cli *RedisClient) SetHashKeyInt(key, field string, value int64) (bool, error) {
	cmd := cli.Cmd.HSet(key, field, value)
	return cmd.Result()
}

func (cli *RedisClient) AddHashKeyInt(key, field string, value int64) (int64, error) {
	cmd := cli.Cmd.HIncrBy(key, field, value)
	return cmd.Result()
}

func (cli *RedisClient) GetHashKeyInt(key, field string) (int64, error) {
	cmd := cli.Cmd.HGet(key, field)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}

	i, err := strconv.ParseInt(cmd.Val(), 10, 64)
	return i, err
}

func (cli *RedisClient) RemoveHashInt(key, field string) error {
	cmd := cli.Cmd.HDel(key, field)
	return cmd.Err()
}

// /////////////////////////////////////////////////////////////////////////////////
// 这部分是用户的好友相关内容使用的
// 使用脚本求2个hash的交集
func (cli *RedisClient) GetHashIntersect(key1, key2 string) ([]string, error) {
	// Lua 脚本代码
	script := `
local hash1 = KEYS[1]
local hash2 = KEYS[2]
local result = {}

local fields1 = redis.call('HKEYS', hash1)

for _, field in ipairs(fields1) do
    local exists = redis.call('HEXISTS', hash2, field)
    if exists == 1 then
        table.insert(result, field)
    end
end

return result
`

	// 执行 Lua 脚本
	result, err := cli.Cmd.Eval(script, []string{key1, key2}).Result()
	if err != nil {
		return nil, err
	}

	//fmt.Println(result, err)
	resultSlice, ok := result.([]interface{})
	if !ok {
		// 如果 result 不是 []interface{} 类型，返回错误
		return nil, errors.New("unexpected type for result")
	}

	strList := make([]string, len(resultSlice))
	for i, item := range resultSlice { // 正确使用 result 而不是 item
		strList[i] = item.(string) // 转换为字符串类型
	}

	// 处理结果
	//fmt.Println("交集字段:", result)
	return strList, err
}

// 将hash表中字段提取出来，然后转为map
func (cli *RedisClient) GetHashKeyIntList(key string) (map[int16]int64, error) {
	// 获取哈希表中所有字段和值
	cmd := cli.Cmd.HGetAll(key)
	// 检查命令执行是否出错
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	// 提取结果并转换为 map[int16]int64 类型
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	hashMap := make(map[int16]int64)
	for field, value := range result {
		// 将哈希表的键转换为 int16 类型
		keyInt, _ := strconv.ParseInt(field, 10, 16)
		// 将哈希表的值转换为 int64 类型
		valInt, _ := strconv.ParseInt(value, 10, 64)
		hashMap[int16(keyInt)] = valInt
	}
	return hashMap, nil

}

func (cli *RedisClient) GetHashKeyInt64List(key string) (map[int64]int32, error) {
	// 获取哈希表中所有字段和值
	cmd := cli.Cmd.HGetAll(key)
	// 检查命令执行是否出错
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	// 提取结果并转换为 map[int16]int64 类型
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	hashMap := make(map[int64]int32)
	for field, value := range result {
		// 将哈希表的键转换为 int16 类型
		keyInt, _ := strconv.ParseInt(field, 10, 64)
		// 将哈希表的值转换为 int64 类型
		valInt, _ := strconv.ParseInt(value, 10, 32)
		hashMap[int64(keyInt)] = int32(valInt)
	}
	return hashMap, nil

}

// 哈希表设置与添加，如果元素超过限制值，则不再添加了
func (cli *RedisClient) SetHashMap(key string, aMap map[string]interface{}) error {
	// 添加 hash 元素
	cmd := cli.Cmd.HMSet(key, aMap)
	return cmd.Err()
}

func (cli *RedisClient) AddHashMap(key string, aMap map[string]interface{}) error {

	cmd := cli.Cmd.HLen(key)
	count, err := cmd.Result()
	if err != nil {
		return err
	}
	if count > MaxFriendCacheSize {
		return errors.New("too mush fields in hash")
	}
	// 添加 hash 元素
	cmd1 := cli.Cmd.HMSet(key, aMap)
	return cmd1.Err()
}

// 哈希表删除
func (cli *RedisClient) RemoveHashMap(key string, fields []string) error {
	cmd := cli.Cmd.HDel(key, fields...)
	return cmd.Err()
}

func (cli *RedisClient) RemoveHashMapWithIntFields(key string, fields []int64) error {

	strArray := make([]string, len(fields))
	for i, v := range fields {
		strArray[i] = strconv.FormatInt(v, 10)
	}
	cmd := cli.Cmd.HDel(key, strArray...)
	return cmd.Err()
}

func (cli *RedisClient) GetHashMap(key string) (map[string]string, error) {
	cmd := cli.Cmd.HGetAll(key)
	return cmd.Result()
}

// 测试和网上说的一致，少于512条数据，不会发生分页
func (cli *RedisClient) ScanHashKeys(key string, cursor uint64, pageSize int64) (uint64, map[string]string, error) {
	result := make(map[string]string)

	// 使用 HSCAN 命令扫描哈希键
	vals, nextCursor, err := cli.Cmd.HScan(key, cursor, "", pageSize).Result()
	if err != nil {
		return 0, nil, err
	}

	// 将扫描结果中的字段和对应的值存储到结果映射中
	for i := 0; i < len(vals); i += 2 {
		field := vals[i]
		value := vals[i+1]
		result[string(field)] = string(value)
	}

	// 更新游标
	cursor = nextCursor
	//fmt.Println("next cursor=", cursor)

	return nextCursor, result, nil
}

func (cli *RedisClient) SetKeyExpire(key string, span time.Duration) error {
	// 使用Redis客户端设置键的过期时间为10秒
	expirationTime := 10 * time.Second
	if span > expirationTime {
		expirationTime = span
	}
	err := cli.Cmd.Expire("key", expirationTime).Err()

	return err
}

func (cli *RedisClient) SetKeysExpire(keys []string, span time.Duration) (int, error) {
	// 使用Redis客户端设置键的过期时间为10秒
	expirationTime := 10 * time.Second
	if span > expirationTime {
		expirationTime = span
	}
	err := cli.Cmd.Expire("key", expirationTime).Err()

	// 使用管道为每个键设置超时时间
	pipe := cli.Cmd.Pipeline()
	for _, key := range keys {
		pipe.Expire(key, expirationTime)
	}

	// 执行管道操作
	cmders, err := pipe.Exec()
	count := 0
	for _, cmd := range cmders {
		if cmd.Err() != nil {
			count++
		}
	}
	return count, err
}
