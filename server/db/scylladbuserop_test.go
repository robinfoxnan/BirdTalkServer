package db

import (
	"birdtalk/server/model"
	"fmt"
	"testing"
	"time"
)

func TestSaveUserOp(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Init()

	record := model.CommonOpStore{
		Pk:   1001,
		Uid1: 1001,
		Uid2: 1003,
		Gid:  2,
		Id:   111111000002,
		Usid: 500002,
		Tm:   time.Now().UnixMilli(),
		Tm1:  0,
		Tm2:  0,
		Io:   0,
		St:   0,
		Cmd:  1,
		Ret:  0,
		Mask: 0,
		Ref:  0,
		Draf: nil,
	}

	err = client.SaveUserOp(&record, 1002)
	fmt.Println(err)
}

func TestUpdateUserOp(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	//client.Init()

	//err = client.SetUserOpRecvReply(1001, 1002, 1001, 1002, 111111000001, time.Now().UnixMilli())
	//err = client.SetUserOpReadReply(1001, 1002, 1001, 1002, 111111000001, time.Now().UnixMilli())
	//err = client.SetUserOpRecvReadReply(1001, 1002, 1001, 1002, 111111000001, time.Now().UnixMilli(), time.Now().UnixMilli()+1)

	err = client.SetUserOpResult(1001, 1002, 1001, 1002, 111111000001, 2)
	fmt.Println(err)
}

func TestFindUserOp(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	//client.Init()

	lst, err := client.FindUserOpForward(1001, 1001, 111111000002, 100)
	for index, item := range lst {
		fmt.Printf("index =%d, record= %v \n", index, item)
	}
	fmt.Println(err)
}
