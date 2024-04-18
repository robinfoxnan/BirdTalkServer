package db

import "time"

/*
基本的思路就是每个服务器都是hash中一个field, 值是当前的时间戳；
如果时间戳过期了，那么其他服务器发现服务器A过期了，就可以直接删除它；
如果这个服务器没有信息，则说明是下线了；

每个定时点服务器更新自己的时间戳，同时获取所有服务器在线状态，用户转发消息；
这里之所以没有保存IP，是因为目前设计中每个服务器使用一个kafka消息队列接收；
*/
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
func (cli *RedisClient) RemoveClusterServer(id int64) error {
	key := GetClusterActiveStateKey()
	field := GetServerField(id)
	return cli.Cmd.HDel(key, field).Err()
}
