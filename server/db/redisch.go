package db

import (
	"fmt"
	"sync/atomic"
)

// 订阅指定的频道，需要手动调用
func (cli *RedisClient) Subscribe(channelName string, handler func(message string)) error {
	i := atomic.LoadInt32(&cli.runningSubscribe)
	if i == 1 {
		return nil
	}
	// 创建订阅对象
	pubSub := cli.Db.Subscribe(channelName)

	// 从通道接收消息
	ch := pubSub.Channel()

	// 启动一个 goroutine 处理接收到的消息
	go func() {
		atomic.StoreInt32(&cli.runningSubscribe, 1)
		defer func() {
			atomic.StoreInt32(&cli.runningSubscribe, 0)
			pubSub.Close()
		}() // 在 goroutine 外部延迟关闭订阅对象

		for msg := range ch {
			//fmt.Printf("Received message from channel %s: %s\n", channelName, msg.Payload)
			handler(msg.Payload) // 调用处理函数处理接收到的消息
		}
	}()

	return nil
}

// 向指定的频道发布消息
func (cli *RedisClient) Publish(channelName, message string) error {
	err := cli.Db.Publish(channelName, message).Err()
	if err != nil {
		fmt.Printf("Error publishing message to channel %s: %s\n", channelName, err)
		return err
	}
	//fmt.Printf("Published message '%s' to channel %s\n", message, channelName)
	return nil
}
