package db

import (
	"birdtalk/server/model"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCreateFollow(t *testing.T) {

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

	id := int64(1001)

	for i := int(1002); i <= 1005; i++ {

		friend := model.FriendStore{
			Pk:   int16(id),
			Uid1: id,
			Uid2: int64(i),
			Tm:   time.Now().UTC().UnixMilli(),
			Nick: "用户" + strconv.Itoa(i),
		}
		fan := model.FriendStore{
			Pk:   int16(i),
			Uid1: int64(i),
			Uid2: id,
			Tm:   time.Now().UTC().UnixMilli(),
			Nick: "飞鸟",
		}

		err = client.InsertFollowing(&friend, &fan)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("insert follow ok =>", i)
		}

	}
}

func TestDeleteFollow(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.DeleteFollowing(1001, 1002, 1001, 1002)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("1001 delete follow ok =>", 1002)
	}
}

func TestQueryFollow(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	lst, err := client.FindFollowing(1001, 1001, 0, 10)

	fmt.Printf("size= %d: \n", len(lst))
	for _, f := range lst {

		fmt.Printf("pk: %d, uid1: %d, uid2: %d,  tm: %d, nick = %s\n",
			f.Pk, f.Uid1, f.Uid2, f.Tm, f.Nick)
	}
}

func TestCountFollowing(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	count, err := client.CountFollowing(1001, 1001)
	fmt.Println(err)
	fmt.Println(count)
}

func TestQueryFans(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	lst, err := client.FindFans(1002, 1002, 0, 10)

	fmt.Printf("size= %d: \n", len(lst))
	for _, f := range lst {

		fmt.Printf("pk: %d, uid1: %d, uid2: %d,  tm: %d, nick = %s\n",
			f.Pk, f.Uid1, f.Uid2, f.Tm, f.Nick)
	}
}

func TestCountFans(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	count, err := client.CountFans(1002, 1002)
	fmt.Println(err)
	fmt.Println(count)
}

func TestFollowingSetNick(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetFollowingNick(1001, 1001, 1002, "robinfoxnan")
	fmt.Println(err)
}

func TestFansSetNick(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.SetFansNick(1006, 1006, 1001, "飞鸟真人")
	fmt.Println(err)
}

// //////////////////////////////////////////////////////////////
func TestCreateBlock(t *testing.T) {

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

	id := int64(1001)

	for i := int(1002); i <= 1005; i++ {

		block := model.BlockStore{
			FriendStore: model.FriendStore{
				Pk:   int16(id),
				Uid1: id,
				Uid2: int64(i),
				Tm:   time.Now().UTC().UnixMilli(),
				Nick: "用户" + strconv.Itoa(i)},
			Perm: 0,
		}

		err = client.InsertBlock(&block)
		fmt.Println(err)

	}
}

func TestDeleteBlock(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.DeleteBlock(1001, 1001, 1002)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("1001 delete block ok =>", 1002)
	}
}

func TestSetBlock(t *testing.T) {
	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SetBlockPermission(1001, 1001, 1003, 5)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("1001 set block ok =>", 1002)
	}
}

func TestQueryBlock(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	lst, err := client.FindBlocks(1001, 1001, 0, 10)

	fmt.Printf("size= %d: \n", len(lst))
	for _, f := range lst {

		fmt.Printf("pk: %d, uid1: %d, uid2: %d,  tm: %d, nick = %s\n",
			f.Pk, f.Uid1, f.Uid2, f.Tm, f.Nick)
	}
}
