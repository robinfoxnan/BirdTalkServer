package db

import "time"

// 每个服务器一个节点号，正整数，从1开始，
func (cli *RedisClient) SetClusterServerActive(id int64) (bool, error) {
	key := GetClusterActiveStateKey()
	field := GetServerField(id)
	return cli.SetHashKeyInt(key, field, time.Now().UnixMilli())
}

// 服务器最后活跃时间戳
func (cli *RedisClient) GetClusterServerActiveTime() (map[int16]int64, error) {
	key := GetClusterActiveStateKey()
	return cli.GetHashKeyIntList(key)
}

// 如果检测到长时间为活跃就应该删除这个节点
//
func (cli *RedisClient) ReomoveClusterServer(id int64) error {
	key := GetClusterActiveStateKey()
	field := GetServerField(id)
	return cli.Db.HDel(key, field).Err()
}
