package db

import (
	"birdtalk/server/pbmodel"
	"fmt"
	"testing"
)

func TestGroupSave(t *testing.T) {

	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建一个用户信息对象
	group := pbmodel.GroupInfo{
		GroupId:   101,
		Tags:      []string{"test", "聊天"},
		GroupName: "测试群聊",
		GroupType: "group",
		Params: map[string]string{
			"title": "GroupChat",
		},
	}
	err = redisCli.SetGroupInfo(&group)
	fmt.Println("err = ", err)

	fmt.Println(err)
}

func TestGroupGet(t *testing.T) {

	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	grp, err := redisCli.GetGroupInfoById(101)
	fmt.Println(*grp)
}

func TestAddGroupMember(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	count, err := redisCli.AddGroupMembers(101, []int64{10001, 10002, 10003})
	fmt.Println(count)

	//count, err = redisCli.AddGroupMembers(101, []int64{10004, 10005, 10006})
	//count, err = redisCli.RemoveGroupMembers(101, []int64{10004, 10005, 10006})
	//redisCli.SetGroupMembers(101, []int64{10004, 10005, 10006})
	lst, err := redisCli.GetGroupMembers(101)
	fmt.Println(lst, err)
}

func TestAddActiveGroupMember(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	count, err := redisCli.AddActiveGroupMembers(101, 1, []int64{10012})
	fmt.Println("server 1 users:", count, err)
	count, err = redisCli.AddActiveGroupMembers(101, 2, []int64{10013, 10014})
	fmt.Println("server 2 users:", count, err)

	//lst, err := redisCli.GetActiveGroupMembers(101, 1)
	//fmt.Println(lst, err)
	//
	//lst, err = redisCli.GetActiveGroupMembers(101, 2)
	//fmt.Println(lst, err)
	//
	//count, err = redisCli.GetActiveGroupMemberCount(101, 1)
	//fmt.Println(count)
	//countMap, err := redisCli.GetActiveGroupMemberCountList(101)
	//fmt.Println(countMap, err)
}
