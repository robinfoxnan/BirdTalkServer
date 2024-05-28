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
			FileName:  "test.doc",
			UniqName:  "bbbbb.doc",
			Gid:       1002,
			Tm:        utils.GetTimeStamp(),
			FileSize:  12,
			UserId:    1001,
			Tags:      []string{"hello", "fox"},
		}
		err = MongoClient.SaveNewFile(&fileInfo)
		time.Sleep(time.Second)
		fmt.Println("save", err)
	}

}

func TestFindFileByString(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = MongoClient.Init()
	fmt.Println(err)

	fileInfo := model.FileInfo{
		HashCode:  "111",
		StoreType: "md5",
		FileName:  "测试1.doc",
		UniqName:  "aaaaaa.doc",
		Gid:       1001,
		Tm:        utils.GetTimeStamp(),
		FileSize:  12,
		UserId:    1001,
		Tags:      []string{"test", "robin"},
	}
	err = MongoClient.SaveNewFile(&fileInfo)
	fmt.Println("save", err)

	lst, err := MongoClient.FindFileByTag("robin")
	fmt.Println("save", lst)

}
