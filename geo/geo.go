package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"net/http"
	"strings"
	"time"
)

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

// 逆地理地址结构
type Address struct {
	Address  string `json:"address"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Street   string `json:"street"`
}

// ============================
// 逆地理缓存（Redis GEO 核心）
// ============================

// Geo 缓存 key（固定）
const geoCacheKey = "regeo:address:cache"

// 查询 100 米内是否已有缓存地址
func (cli *RedisClient) GetNearbyAddress(lng, lat float64) (*Address, error) {
	// 查询 100 米内最近的一个点
	res, err := cli.Cmd.GeoRadius(geoCacheKey, lng, lat, &redis.GeoRadiusQuery{
		Radius: 100, // 100 米内算同一位置
		Unit:   "m",
		Count:  1,
		Sort:   "ASC",
	}).Result()

	if err != nil || len(res) == 0 {
		return nil, err
	}

	// 拿到缓存 key
	key := res[0].Name

	// 获取地址 JSON
	jsonStr, err := cli.Cmd.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var addr Address
	if err := json.Unmarshal([]byte(jsonStr), &addr); err != nil {
		return nil, err
	}

	return &addr, nil
}

// 把经纬度 + 地址 存入 GEO 缓存
func (cli *RedisClient) SaveAddressToGeo(lng, lat float64, addr Address) error {
	// 生成唯一 key
	key := fmt.Sprintf("regeo:addr:%f:%f", lng, lat)

	// 1. 存入 GEO 位置
	_, err := cli.Cmd.GeoAdd(geoCacheKey, &redis.GeoLocation{
		Name:      key,
		Longitude: lng,
		Latitude:  lat,
	}).Result()
	if err != nil {
		return err
	}

	// 2. 存入地址详情（7天过期）
	jsonStr, _ := json.Marshal(addr)
	//_, err = cli.Cmd.Set(key, jsonStr, 7*24*time.Hour).Result()
	_, err = cli.Cmd.Set(key, jsonStr, 0).Result()
	return err
}

// ============================
// 调用高德逆地理 API（你直接替换 KEY 即可）
// ============================
func (cli *RedisClient) RequestAmapRegeo(lng, lat float64) (Address, error) {
	apiKey := "你的高德Web服务Key"
	url := fmt.Sprintf(
		"https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%f,%f&output=json",
		apiKey, lng, lat,
	)

	resp, err := http.Get(url)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Status    string `json:"status"`
		Regeocode struct {
			FormattedAddress string `json:"formatted_address"`
			AddressComponent struct {
				Province string `json:"province"`
				City     string `json:"city"`
				District string `json:"district"`
				Township string `json:"township"`
				Street   string `json:"street"`
			} `json:"addressComponent"`
		} `json:"regeocode"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Address{}, err
	}

	if result.Status != "1" {
		return Address{}, fmt.Errorf("amap api error")
	}

	return Address{
		Address:  result.Regeocode.FormattedAddress,
		Province: result.Regeocode.AddressComponent.Province,
		City:     result.Regeocode.AddressComponent.City,
		District: result.Regeocode.AddressComponent.District,
		Street:   result.Regeocode.AddressComponent.Street,
	}, nil
}

// ============================
// 最终对外接口：缓存优先
// ============================
func (cli *RedisClient) GetRegeoAddress(lng, lat float64) (addr *Address, from string, err error) {
	// 1. 先查 100 米内缓存
	cacheAddr, err := cli.GetNearbyAddress(lng, lat)
	if err == nil && cacheAddr != nil {
		return cacheAddr, "redis_geo_cache", nil
	}

	// 2. 无缓存 → 请求高德
	newAddr, err := cli.RequestAmapRegeo(lng, lat)
	if err != nil {
		return nil, "", err
	}

	// 3. 存入 GEO 缓存
	_ = cli.SaveAddressToGeo(lng, lat, newAddr)

	return &newAddr, "amap_api", nil
}

func main() {

}
