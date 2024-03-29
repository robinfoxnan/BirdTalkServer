package db

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"fmt"
	"testing"
	"time"
)

func TestAddPMsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := model.PChatDataStore{
		Pk:    1001,
		Uid1:  1001,
		Uid2:  1002,
		Id:    123456,
		Usid:  010101222,
		Tm:    time.Now().UTC().UnixMilli(),
		Tm1:   0,
		Tm2:   0,
		Io:    0,
		St:    0,
		Ct:    0,
		Mt:    0,
		Print: 0,
		Ref:   0,
		Draf:  []byte("it is a test"),
	}
	err = client.SavePChatData(&msg, 1002)
	fmt.Println(err)
}

func TestSetPMsgDelete(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetPChatMsgDeleted(1001, 1002, 1001, 1002, 1711672998994, 123456)
	fmt.Println(err)
}

func TestSetPMsgReply(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetPChatRecvReply(1001, 1002, 1001, 1002, 1711672998994, 123456, time.Now().UTC().UnixMilli())
	fmt.Println(err)

	err = client.SetPChatReadReply(1001, 1002, 1001, 1002, 1711672998994, 123456, time.Now().UTC().UnixMilli())
	fmt.Println(err)
}

func TestSetPMsgReply1(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetPChatRecvReadReply(1001, 1002, 1001, 1002, 1711672998994, 123456, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
	fmt.Println(err)
}

func TestFindPChat(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	lst, err := client.FindPChatMsg(1001, 1001, 1711672998994, 100)
	for _, item := range lst {
		fmt.Println(item, utils.UtcTm2LocalString(item.Tm))
		//fmt.Println(utils.UtcTm2LocalString(item.Tm1))
	}

}

// //////////////////////////////////////////////////
func TestAddGroupMsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := model.GChatDataStore{
		Pk:    101,
		Gid:   101,
		Uid:   10001,
		Id:    321313213,
		Usid:  68666,
		Tm:    time.Now().UTC().UnixMilli(),
		Res:   0,
		St:    0,
		Ct:    0,
		Mt:    0,
		Print: 0,
		Ref:   0,
		Draf:  []byte("it is a test"),
	}

	err = client.SaveGChatData(&msg)
	fmt.Println(err)
}

func TestSetGroupMsgDel(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetGChatMsgDeleted(101, 101, 1711671994002, 321313213)
	fmt.Println(err)
}

func TestFindGroupmsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	list, err := client.FindGChatMsg(101, 101, 0, 100)
	fmt.Println(list)

}
