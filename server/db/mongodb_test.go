package db

import (
	"birdtalk/server/pbmodel"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//func (me *MongoDBExporter) initUsersIndexes() error {
//
//	collection := me.db.Collection(UserTableName)
//
//	index := mongo.IndexModel{
//		Keys: bson.M{
//			"username": 1, // 设置 username 字段为索引字段
//			"email":    1, // 设置 email 字段为索引字段,升序（ascending）
//			"phone":    1, // 设置 phone 字段为索引字段
//		},
//	}
//
//	_, err := collection.Indexes().CreateOne(context.Background(), index)
//
//	return err
//}

// windows version 7.0.5 community
func TestCreateIndex(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//err = MongoClient.CreateIndex(UserTableName, "userid")
	err = MongoClient.Init()
	fmt.Println(err)
}
func TestCreateUser(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
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

	tm1 := time.Now().UnixMilli()
	err = MongoClient.CreateNewUser(&userInfo)
	tm2 := time.Now().UnixMilli()
	fmt.Println("cost ms = ", tm2-tm1)
	if err != nil {
		fmt.Println(err)
	}
}

func TestFindUserByString(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userInfos, err := MongoClient.FindUserByKeyword("john_doe")
	if err == nil {
		for _, u := range userInfos {
			b, _ := json.Marshal(u)
			fmt.Println(string(b))
		}

	} else {
		fmt.Println(err)
	}
}

func TestFindUserById(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userInfos, err := MongoClient.FindUserById(10001)
	if err == nil {
		for _, u := range userInfos {
			b, _ := json.Marshal(u)
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(err)
	}
}

func TestFindUserByName(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	userInfos, err := MongoClient.FindUserByName("robin")
	if err == nil {
		for _, u := range userInfos {
			b, _ := json.Marshal(u)
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(err)
	}
}

// 更新部分信息需要在业务逻辑中控制
func TestSetUserInfo(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个用户信息对象
	userInfo := pbmodel.UserInfo{
		UserId:   10001,
		UserName: "john_doe1",
		NickName: "John1",
		Email:    "john1@example.com",
		Phone:    "1234567890",
		Gender:   "female",
		Age:      31,
		Region:   "ZH",
		Icon:     "avatar1.jpg",
		Params: map[string]string{
			"title":    "Mrs.",
			"pwd":      "password1230",
			"sid":      "session1230",
			"icon":     "avatar1.jpg",
			"phonepre": "+86",
		},
	}

	tm1 := time.Now().UnixMilli()
	n, err := MongoClient.UpdateUserInfo(&userInfo)
	tm2 := time.Now().UnixMilli()
	fmt.Println("update count = ", n)
	fmt.Println("cost ms = ", tm2-tm1)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSetUserInfoPart(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个用户信息对象
	setData := make(map[string]any)
	setData["nickname"] = "robin"
	setData["phone"] = "1100000"

	unsetData := []string{"params.phonepre"}

	tm1 := time.Now().UnixMilli()
	// n, err := MongoClient.UpdateUserInfoPart(10001, setData, unsetData)
	n, err := MongoClient.UpdateUserInfoPart(10001, setData, unsetData)
	tm2 := time.Now().UnixMilli()
	fmt.Println("update count = ", n)
	fmt.Println("cost ms = ", tm2-tm1)
	if err != nil {
		fmt.Println(err)
	}
}

//////////////////////////////////////////////////////////////////////////////

func TestCreateGroup(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	group := pbmodel.GroupInfo{
		GroupId:   10006,
		GroupType: "group", // "channel"
		GroupName: "隐藏群",
		Tags:      []string{"test1", "test3", "下棋"},
		Params: map[string]string{
			"pwd": "password123",
			//"v":    "pri", // pub, pri
			"code": "123456",
		},
	}

	err = MongoClient.CreateNewGroup(&group)
	if err != nil {
		fmt.Println(err)
	}
}

func TestFindGroupByTag(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	grps, err := MongoClient.FindGroupByKeyword("test1")
	if grps != nil {
		for _, g := range grps {
			b, _ := json.Marshal(g)
			fmt.Println(string(b))
		}
	}
}
func TestFindGroupById(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	grps, err := MongoClient.FindGroupById(10005, "123456")
	if grps != nil {
		for _, g := range grps {
			b, _ := json.Marshal(g)
			fmt.Println(string(b))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func TestSetGroupInfo(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	group := pbmodel.GroupInfo{
		GroupId:   10002,
		GroupType: "group", // "channel"
		GroupName: "测试群1",
		Tags:      []string{"越野", "徒步", "闲聊"},
		Params: map[string]string{
			"pwd": "password123",
			"q":   "public",
		},
	}

	n, err := MongoClient.UpdateGroupInfo(&group)
	fmt.Println(n)
	if err != nil {
		fmt.Println(err)
	}

}

func TestSetGroupInfoPart(t *testing.T) {
	connStr := "mongodb://admin:123456@127.0.0.1:27017"
	dbName := "birdtalk"
	err := InitMongoClient(connStr, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}

	setData := map[string]interface{}{
		"tags": []string{"学习", "聊天", "下棋", "骑行"},
	}

	n, err := MongoClient.UpdateGroupInfoPart(10003, setData, nil)
	fmt.Println(n)
	if err != nil {
		fmt.Println(err)
	}

}
