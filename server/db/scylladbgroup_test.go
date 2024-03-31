package db

import (
	"birdtalk/server/model"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestAddMemToGroup(t *testing.T) {

	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	gid := int64(100)

	for i := int(1001); i <= 1005; i++ {

		mem := model.GroupMemberStore{
			Pk:   int16(gid),
			Gid:  gid,
			Uid:  int64(i),
			Tm:   time.Now().UTC().UnixMilli(),
			Role: 1,
			Nick: "用户" + strconv.Itoa(i),
		}

		item := model.UserInGStore{
			Pk:  int16(i),
			Uid: int64(i),
			Gid: gid,
		}
		err = client.InsertGroupMember(&mem, &item)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("insert user to group ok =>", i)
		}

	}
}

func TestFindGroupMem(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	list, err := client.FindGroupMembers(100, 100, 0, 100)
	fmt.Println(err)
	fmt.Println(list)
}

func TestUserInG(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	list, err := client.FindUserInGroups(1003, 1003, 0, 100)
	fmt.Println(err)
	fmt.Println(list)
}

func TestSetGroupMemberInfo(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetGroupMemberNickRole(100, 100, 1003, "飞鸟真人", 5)
	fmt.Println(err)
}

func TestDeleteGroupMember(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.DeleteGroupMember(100, 1002, 100, 1002)
	fmt.Println(err)
}

func TestGroupDisovle(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.DissolveGroupAllMember(100, 100)
	fmt.Println(err)
}

func TestDeleteUing(t *testing.T) {
	client, err := NewScyllaClient([]string{"127.0.0.1:9042"}, "cassandra", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.DeleteUserInG(1004, 1004, 100)
	fmt.Println(err)
}
