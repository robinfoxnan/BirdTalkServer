package db

import (
	"birdtalk/server/model"
	"fmt"
	"testing"
	"time"
)

func TestSaveGroupOp(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Init()

	record := model.CommonOpStore{
		Pk:   2,
		Gid:  2,
		Uid1: 1001,
		Uid2: 0,
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
	fmt.Println(err)

	err = client.SaveGroupOp(&record)

	log, err := client.FindGroupOpExact(2, 2, 111111000002)

	fmt.Println(log, err)
}

func TestUpdateGroupOp(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetGroupOpResult(2, 2, 1111110000021, 110, true)
	fmt.Println("Update ret=", err)

	log, err := client.FindGroupOpExact(2, 2, 1111110000021)

	fmt.Println(log, err)

}
