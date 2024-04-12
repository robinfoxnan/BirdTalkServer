package db

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func onStateChange(message string) {
	fmt.Println("Handler received message:", message)
}

func TestRedisCh(t *testing.T) {

	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	chName := GetStateChKey()

	err = redisCli.Subscribe(chName, onStateChange)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 1; i < 30; i++ {
		err = redisCli.Publish(chName, "message index = "+strconv.Itoa(i))
		time.Sleep(time.Second)
	}
	fmt.Println("exit here")
}
