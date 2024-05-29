package db

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestSaveFiles(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = MongoClient.Init()
	fmt.Println(err)

	for i := 11; i < 20; i++ {
		fileInfo := model.FileInfo{
			HashCode:  strconv.FormatInt(int64(i), 10),
			StoreType: "md5",
			FileName:  "测试.doc",
			UniqName:  "aaaa.doc",
			Gid:       0,
			Status:    "",
			Tm:        utils.GetTimeStamp(),
			FileSize:  12,
			UserId:    1009,
			Tags:      []string{"hello", "fox"},
		}
		err = MongoClient.SaveNewFile(&fileInfo)
		time.Sleep(time.Second)
		fmt.Println("save", err)
	}

}

func TestFindFile(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = MongoClient.Init()
	//fmt.Println(err)

	//lst, err := MongoClient.FindFileByTag("robin")
	//fmt.Println("save", len(lst), lst)

	data, err := MongoClient.FindFileByHashCode("11")
	fmt.Println("save", data)

	//lst, err := MongoClient.FindFileByGroup(1002, 1716883226053)
	//fmt.Println("save", len(lst), lst)
	//
	//lst, err = MongoClient.FindFileByUser(1001, 1716883139804)
	//fmt.Println("save", len(lst), lst)

}
