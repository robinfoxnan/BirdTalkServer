package db

import (
	"birdtalk/server/model"
	"fmt"
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

	for i := int64(1001); i <= 1300; i++ {
		f := model.FriendStore{
			Pk:   1002,
			Uid1: 1002,
			Uid2: i,
			Tm:   time.Now().UTC().UnixMilli(),
			Nick: "飞鸟",
		}
		err = client.InsertFollowing(&f)
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

	err = client.DeleteFollowing(1001, 1001, 1002)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("insert follow ok =>", 1002)
	}
}

func TestQueryFollow(t *testing.T) {

	client, err := NewScyllaClient([]string{"8.140.203.92:9042"}, "cassandra", "Tjj.31415")
	if err != nil {
		fmt.Println(err)
		return
	}

	lst, err := client.FindFollowing(1002, 1002, 1010, 10)

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
	count, err := client.CountFollowing(1002, 1002)
	fmt.Println(err)
	fmt.Println(count)
}
