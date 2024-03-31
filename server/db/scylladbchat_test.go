package db

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"fmt"
	"testing"
	"time"
)

func TestAddPMsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	snow := utils.NewSnowflake(1, 1)
	msgId := snow.GenerateID()

	msg := model.PChatDataStore{
		Pk:    1001,
		Uid1:  1001,
		Uid2:  1002,
		Id:    msgId,
		Usid:  1,
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

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetPChatMsgDeleted(1001, 1002, 1001, 1002, 693649456985411584)
	fmt.Println(err)
}

func TestSetPMsgReply(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetPChatRecvReply(1001, 1002, 1001, 1002, 693649456985411584, time.Now().UTC().UnixMilli())
	fmt.Println(err)

	err = client.SetPChatReadReply(1001, 1002, 1001, 1002, 693649456985411584, time.Now().UTC().UnixMilli())
	fmt.Println(err)
}

func TestSetPMsgReply1(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetPChatRecvReadReply(1001, 1002, 1001, 1002, 693649502762045440, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
	fmt.Println(err)
}

func TestFindPChat(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	lst, err := client.FindPChatMsg(1002, 1002, 693649502762045440, 100)
	for _, item := range lst {

		tm := utils.SnowIdtoTm(item.Id)
		str2 := utils.TmToLocalString(tm)
		fmt.Println(item, utils.TmToLocalString(item.Tm), str2)
		//fmt.Println(utils.UtcTm2LocalString(item.Tm1))
	}

}

// //////////////////////////////////////////////////
func TestAddGroupMsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	snow := utils.NewSnowflake(1, 1)
	msgId := snow.GenerateID()

	msg := model.GChatDataStore{
		Pk:    101,
		Gid:   101,
		Uid:   10001,
		Id:    msgId,
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

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetGChatMsgDeleted(101, 101, 321313213)
	fmt.Println(err)
}

func TestFindGroupmsg(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	//tm := 0
	//id := utils.TmToSnowIdLike(int64(tm))
	//fmt.Println(id)

	var id int64 = 693645496463527935

	list, err := client.FindGChatMsg(101, 101, id, 100)
	for _, item := range list {

		str1 := utils.TmToLocalString(item.Tm)
		tm := utils.SnowIdtoTm(item.Id)
		str2 := utils.TmToLocalString(tm)

		fmt.Println(item, str1, str2)
	}

}
