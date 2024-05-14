package db

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestFindFollow(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//list, err := redisCli.GetFriendIntersect(10001)
	//fmt.Println(list)

	off := uint64(0)

	cursor, list1, err := redisCli.GetUserFollowing(10003, off)

	fmt.Println(cursor, list1, err)

	b, err := redisCli.HasKey("t10003", time.Hour)
	fmt.Println(b, err)

}
func TestUserInfoToMap(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建一个用户信息对象
	userInfo := pbmodel.UserInfo{
		UserId:   10001,
		UserName: "john_doe",
		NickName: "John",
		Email:    "john@example.com",
		Phone:    "123456789",
		Gender:   "male",
		Age:      30,
		Region:   "US",
		Icon:     "avatar.jpg",
		Params: map[string]string{
			"title": "Mr.",
			"pwd":   "password123",
			"sid":   "session123",
			"icon":  "avatar.jpg",
		},
	}
	fmt.Println(userInfo)

	userMap, err := userInfoToMap(&userInfo)
	for k, v := range userMap {
		t := reflect.TypeOf(v)
		fmt.Printf("%v (%v)= %v \n", k, t.Kind(), v)
	}

	fmt.Println(err)

	redisCli.SetUserInfo(&userInfo)

}

func TestUserInfoLoad(t *testing.T) {

	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建一个用户信息对象
	userInfo, err := redisCli.FindUserById(10005)
	fmt.Println("err = ", err)

	fmt.Println(userInfo)
}

func TestUserSetFollowing(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	id := int64(10001)
	friendList := make([]model.FriendStore, 0)
	blockList := make([]model.BlockStore, 0)

	for i := int(10002); i <= 10600; i++ {

		friend := model.FriendStore{
			Pk:   int16(id),
			Uid1: id,
			Uid2: int64(i),
			Tm:   time.Now().UTC().UnixMilli(),
			Nick: "用户" + strconv.Itoa(i),
		}
		friendList = append(friendList, friend)

		block := model.BlockStore{
			FriendStore: friend,
			Perm:        8,
		}
		blockList = append(blockList, block)
	}

	err = redisCli.SetUserFollowing(10001, friendList)
	err = redisCli.SetUserFans(10001, friendList)
	err = redisCli.SetUserBlocks(10001, blockList)

	fmt.Println(err)
}

func TestUserAddFollowing(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	id := int64(10001)
	friendList := make([]model.FriendStore, 0)
	blockList := make([]model.BlockStore, 0)

	for i := int(10040); i <= 10042; i++ {

		friend := model.FriendStore{
			Pk:   int16(id),
			Uid1: id,
			Uid2: int64(i),
			Tm:   time.Now().UTC().UnixMilli(),
			Nick: "用户" + strconv.Itoa(i),
		}
		friendList = append(friendList, friend)

		block := model.BlockStore{
			FriendStore: friend,
			Perm:        8,
		}
		blockList = append(blockList, block)
	}

	err = redisCli.AddUserFollowing(10001, friendList)
	err = redisCli.AddUserFans(10001, friendList)
	err = redisCli.AddUserBlocks(10001, blockList)

	fmt.Println(err)
}

func TestUserRemoveFollowing(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = redisCli.RemoveUserFollowing(10001, []int64{10041})
	fmt.Println(err)
	err = redisCli.RemoveUserFans(10001, []int64{10041})
	fmt.Println(err)
	err = redisCli.RemoveUserBlocks(10001, []int64{10041})
	fmt.Println(err)
}

func TestUserFriends(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//list, err := redisCli.GetFriendIntersect(10001)
	//fmt.Println(list)

	off := uint64(0)
	count := 0
	for {
		cursor, list1, _ := redisCli.GetUserFollowing(10001, off)
		count += len(list1)
		fmt.Printf("len= %d off= %d, count = %d\n", len(list1), cursor, count)

		off = cursor
		if cursor == 0 {
			break
		}
	}

	//list2, err := redisCli.GetUserFans(10001, 5)
	//fmt.Println(list2)
	//
	//list3, err := redisCli.GetUserBLocks(10001, 5)
	//fmt.Println(list3)

}

func TestUserAddGroup(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = redisCli.SetUserJoinGroup(10002, 101)
	err = redisCli.SetUserJoinGroup(10002, 102)
	fmt.Println(err)

	lst, err := redisCli.GetUsersInSameGroup(10001, 10002)
	fmt.Println(lst)

	lst, err = redisCli.GetUserInGroupAll(10002)
	fmt.Println(lst)

}
func TestIntFormat(t *testing.T) {
	tm := utils.GetTimeStamp()
	i := strconv.FormatInt(tm, 10)
	e := strconv.FormatInt(tm, 36)
	fmt.Println(i, e)
	fmt.Println(len(i), len(e))
}

// 测试查询某个值是否存在
func TestCheckUserFan(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	friend := model.FriendStore{
		Pk:   int16(10001),
		Uid1: 10001,
		Uid2: int64(10003),
		Tm:   time.Now().UTC().UnixMilli(),
		Nick: "用户1",
	}
	redisCli.AddUserFans(10001, []model.FriendStore{friend})

	ret, err := redisCli.CheckUserFan(10001, 10004)
	fmt.Println(ret, err)
}
