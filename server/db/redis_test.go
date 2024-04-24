package db

import (
	"birdtalk/server/model"
	"fmt"
	"testing"
)

func TestGetHashInt(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	// model.DefaultPermissionP2P | model.PermissionMaskFriend
	redisCli.AddUserPermission(10001, 10002, model.DefaultPermissionP2P|model.PermissionMaskFriend)
	mask, err := redisCli.CheckUserPermission(10001, 10008)
	//mask, err := redisCli.GetHashKeyInt("p10001", "10002")
	if err == nil {
		fmt.Println("mask", mask)
	} else {
		fmt.Println(err.Error())
	}

}
