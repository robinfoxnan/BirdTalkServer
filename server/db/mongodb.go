package db

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

// go get go.mongodb.org/mongo-driver
const UserTableName = "users"
const GroupTableName = "groups"
const FileTableName = "files"

// MongoDBExporter 结构体
type MongoDBExporter struct {
	client *mongo.Client
	db     *mongo.Database
}

// 全局变量
var MongoClient *MongoDBExporter = nil

// NewMongoDBExporter 创建一个新的 MongoDBExporter 实例
func NewMongoDBExporter(connectionString, dbName string) (*MongoDBExporter, error) {
	// 创建 MongoDB 连接选项
	clientOptions := options.Client().ApplyURI(connectionString)
	clientOptions.SetConnectTimeout(10 * time.Second). // 设置连接超时时间为 10 秒
								SetSocketTimeout(5 * time.Second). // 设置 Socket 超时时间为 5 秒
								SetMaxPoolSize(100)                // 设置连接池大小为 100

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// 处理连接错误
		return nil, err
	}

	if err != nil {
		// 处理连接错误
		return nil, err
	}

	// 检查连接是否成功
	err = client.Ping(ctx, nil)

	database := client.Database(dbName)

	return &MongoDBExporter{
		client: client,
		db:     database,
	}, nil
}

// "mongodb://localhost:27017"
func InitMongoClient(connectionString, dbName string) error {
	var err error
	MongoClient, err = NewMongoDBExporter(connectionString, dbName)
	return err
}

// 在日志处理中，按照日期分表
func getDateString(mt int64) string {
	t := time.UnixMilli(mt)
	// time.Now()
	return t.Format("20060102")
}

// close 关闭 MongoDBExporter
func (me *MongoDBExporter) Close() error {
	err := me.client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func printTm(tm1, tm2 int64) {
	t1 := time.UnixMilli(tm1)
	t2 := time.UnixMilli(tm2)

	layout := "2006-01-02 15:04:05"

	fmt.Printf("%s---> %s  \n", t1.Format(layout), t2.Format(layout))
}

// TODO: 把所有需要查询过滤掉属性都放在一列中，用一个seartch_tags字段，这样快一些：
// 假设有百万级群表：
//查询方式								数据量		索引命中					备注
//------------------------------------------------------------------------------------------
//单条件 + 分页							1M			走复合索引				返回前100条，毫秒级
//$or 两个复合索引 + 分页groupid>fromId	1M			两个索引扫描 + 合并		返回100条，ms → 数十ms
//$or + 非索引条件						1M			两个索引扫描 + 内存过滤		可能上百 ms

var userTableSearchFields = []string{"username", "email", "phone"}
var fileTableSearchFields = []string{"hashcode", "gid", "tm", "tags"}

// init 初始化业务相关的部分
func (me *MongoDBExporter) Init() error {
	// 用户表，四个索引
	err := me.CreateIndexUnique(UserTableName, "userid", true)
	err = me.CreateIndexes(UserTableName, []string{"nickname", "userid"})
	err = me.CreateIndexes(UserTableName, []string{"email", "userid"})
	err = me.CreateIndexes(UserTableName, []string{"phone", "userid"})

	// 1️⃣ groupid 唯一索引（必须最先）
	err = me.CreateIndexUnique(GroupTableName, "groupid", true)

	// 2️⃣ 搜索 + 游标分页索引
	err = me.CreateIndexes(GroupTableName, []string{"groupname", "groupid"})
	err = me.CreateIndexes(GroupTableName, []string{"tags", "groupid"})

	// 文件查找
	err = me.CreateIndexUnique(FileTableName, "hashcode", false)
	return err
}

// //////////////////////////////////////////////////////////////////////
func (me *MongoDBExporter) CreateIndexes(tableName string, indexFields []string) error {
	collection := me.db.Collection(tableName)

	keys := bson.D{}
	for _, field := range indexFields {
		keys = append(keys, primitive.E{
			Key:   field,
			Value: 1,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	index := mongo.IndexModel{
		Keys: keys,
		Options: options.Index().
			SetName(strings.Join(indexFields, "_") + "_idx"),
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	return err
}

func (me *MongoDBExporter) CreateIndexUnique(tableName, fieldName string, bUnique bool) error {
	collection := me.db.Collection(tableName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	indexOpts := options.Index().
		SetName(fieldName + "_idx") // 可以根据 bUnique 加后缀

	if bUnique {
		indexOpts.SetUnique(true)
		indexOpts.SetName(fieldName + "_unique_idx")
	}

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: fieldName, Value: 1}},
		Options: indexOpts,
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	return err
}

// //////////////////////////////////////////////////////////////////////
// 创建一个新的用户，如果ID重复则失败
func (me *MongoDBExporter) CreateNewUser(u *pbmodel.UserInfo) error {
	// 选择要保存数据的数据库和集合
	collection := me.db.Collection(UserTableName)

	// 将用户信息对象转换为 MongoDB 文档
	bsonData, err := bson.Marshal(u)
	_, err = collection.InsertOne(context.Background(), bsonData)
	if err != nil {
		return err
	}

	fmt.Println("User information has been saved successfully.")
	return nil
}

// 查找用户，
func (me *MongoDBExporter) FindUserById(id int64) ([]*pbmodel.UserInfo, error) {

	return me.FindUserByField("userid", id, 0, 1)
}

// 这个字段没有设置唯一索引，出于性能考虑，在更改名字的时候可以检测是否唯一，
func (me *MongoDBExporter) FindUserByName(keyword string, fromId int64, pageSize int64) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("username", keyword, fromId, pageSize)
}

func (me *MongoDBExporter) FindUserByNick(keyword string, fromId int64, pageSize int64) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("nickname", keyword, fromId, pageSize)
}

func (me *MongoDBExporter) FindUserByEmail(keyword string, fromId int64, pageSize int64) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("email", keyword, fromId, pageSize)
}

func (me *MongoDBExporter) FindUserByPhone(keyword string, fromId int64, pageSize int64) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("phone", keyword, fromId, pageSize)
}

func (me *MongoDBExporter) FindUserByField(field string, keyword interface{}, fromId, pageSize int64) ([]*pbmodel.UserInfo, error) {
	collection := me.db.Collection(UserTableName)

	// 只允许三个字段之一
	if field != "username" && field != "email" && field != "phone" {
		return nil, fmt.Errorf("unsupported search field: %s", field)
	}

	// 构建查询条件：指定字段 = keyword + userid 游标分页
	filter := bson.M{
		"$and": []bson.M{
			{field: keyword},
			{"userid": bson.M{"$gt": fromId}},
		},
	}

	// 查询选项：按 userid 升序 + 限制数量
	findOptions := options.Find().
		SetSort(bson.D{{Key: "userid", Value: 1}}).
		SetLimit(pageSize)

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*pbmodel.UserInfo
	for cursor.Next(context.Background()) {
		var user pbmodel.UserInfo
		if err := cursor.Decode(&user); err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, &user)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users, nil
}

// 根据用户名、邮件或手机号进行搜索
func (me *MongoDBExporter) FindUserByKeyword(keyword string, fromId int64, pageSize int64) ([]*pbmodel.UserInfo, error) {
	collection := me.db.Collection(UserTableName)

	// 构建查询条件
	filter := bson.M{
		"$and": []bson.M{
			{"userid": bson.M{"$gt": fromId}}, // 游标分页
			{
				"$or": []bson.M{
					{"username": keyword},
					{"email": keyword},
					{"phone": keyword},
				},
			},
		},
	}

	// 查询选项：按 userid 升序 + 限制数量
	findOptions := options.Find().
		SetSort(bson.D{{Key: "userid", Value: 1}}).
		SetLimit(pageSize)

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*pbmodel.UserInfo
	for cursor.Next(context.Background()) {
		var user pbmodel.UserInfo
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	if len(users) == 0 {
		return nil, nil
	}
	return users, nil
}

// 通过参数这个更新
func (me *MongoDBExporter) UpdateUserInfo(u *pbmodel.UserInfo) (int64, error) {
	collection := me.db.Collection(UserTableName)
	filter := bson.M{"userid": u.UserId} // 过滤条件
	// 现有的更新操作
	update := bson.M{
		"$set": bson.M{
			"username": u.UserName,
			"nickname": u.NickName,
			"email":    u.Email,
			"phone":    u.Phone,
			"gender":   u.Gender,
			"age":      u.Age, // 假设新年龄为30
			"region":   u.Region,
			"icon":     u.Icon,
		},
	}

	// 遍历param中的字段
	for k, v := range u.Params {
		key := "params." + k
		update["$set"].(bson.M)[key] = v
	}

	// "$unset" 用于删除字段
	// 后续添加的更新字段
	//update["$set"].(bson.M)["params.key3"] = "new_value3"
	//update["$set"].(bson.M)["params.key4"] = "new_value4"

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// 处理错误
		return 0, err
	}

	fmt.Println("Matched:", result.MatchedCount, "Modified:", result.ModifiedCount)
	return result.ModifiedCount, nil
}

// 注意，这里需要自己拼接params.key这样的关键字
// 例如删除一个"params.ip"
// 2025-09-11 updata ,这样的处理方式可以避免当prarams为null时候报错。
// 要实现「params 为空时自动替换成 {} 再合并子字段」，必须用 update pipeline（MongoDB 4.2+ 支持）。
func (me *MongoDBExporter) UpdateUserInfoPart(id int64, setData map[string]interface{}, unsetData []string) (int64, error) {
	collection := me.db.Collection(UserTableName)
	filter := bson.M{"userid": id} // 过滤条件

	update := bson.M{
		"$set":   bson.M{},
		"$unset": bson.M{},
	}

	for k, v := range setData {
		if strings.HasPrefix(k, "params.") {
			// 直接使用 "params.xxx" 形式
			update["$set"].(bson.M)[k] = v
		} else {
			update["$set"].(bson.M)[k] = v
		}
	}

	for _, k := range unsetData {
		update["$unset"].(bson.M)[k] = ""
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return 0, err
	}

	fmt.Println("Matched:", result.MatchedCount, "Modified:", result.ModifiedCount)
	return result.ModifiedCount, nil
}

// 更新用户的关注的个数
func (me *MongoDBExporter) UpdateUserFieldIncNum(field string, num, id int64) bool {
	collection := me.db.Collection(UserTableName)

	// 创建一个上下文，设置操作超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 自增操作
	filter := bson.M{"userid": id} // 过滤条件
	update := bson.M{"$inc": bson.M{field: num}}

	// 更新文档，如果文档不存在则创建它
	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount > 0 {
		fmt.Println("Matched an existing document and incremented the count.")

	} else if result.UpsertedCount > 0 {
		fmt.Printf("Inserted a new document with ID %v and set the initial count.\n", result.UpsertedID)
	}
	fmt.Println(result)

	return true
}

// ////////////////////////////////////////////////////////////////////////////////////
// 保存新
func (me *MongoDBExporter) CreateNewGroup(g *pbmodel.GroupInfo) error {

	collection := me.db.Collection(GroupTableName)

	// 将用户信息对象转换为 MongoDB 文档
	bsonData, err := bson.Marshal(g)
	_, err = collection.InsertOne(context.Background(), bsonData)
	if err != nil {
		return err
	}

	//fmt.Println("Group information has been saved successfully.")
	return nil

}

// 更新群基础信息，用的不会太多
func (me *MongoDBExporter) UpdateGroupInfo(g *pbmodel.GroupInfo) (int64, error) {
	collection := me.db.Collection(GroupTableName)
	filter := bson.M{"groupid": g.GroupId} // 过滤条件
	// 现有的更新操作
	update := bson.M{
		"$set": bson.M{
			"groupname": g.GroupName,
			//"grouptype": g.GroupType,
			"tags": g.Tags,
		},
	}

	// 遍历param中的字段
	for k, v := range g.Params {
		key := "params." + strings.ToLower(k)
		update["$set"].(bson.M)[key] = v
	}

	// "$unset" 用于删除字段
	// 后续添加的更新字段
	//update["$set"].(bson.M)["params.key3"] = "new_value3"
	//update["$set"].(bson.M)["params.key4"] = "new_value4"

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// 处理错误
		return 0, err
	}

	fmt.Println("Matched:", result.MatchedCount, "Modified:", result.ModifiedCount)
	return result.ModifiedCount, nil
	return 1, nil
}

func (me *MongoDBExporter) UpdateGroupInfoPart(id int64, setData map[string]interface{}, unsetData []string) (int64, error) {

	collection := me.db.Collection(GroupTableName)
	filter := bson.M{"groupid": id} // 过滤条件

	// 初始化 update 变量
	update := bson.M{
		"$set":   bson.M{},
		"$unset": bson.M{},
	}
	for k, v := range setData {
		update["$set"].(bson.M)[k] = v
	}

	for _, k := range unsetData {
		update["$unset"].(bson.M)[k] = nil
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// 处理错误
		return 0, err
	}

	fmt.Println("Matched:", result.MatchedCount, "Modified:", result.ModifiedCount)
	return result.ModifiedCount, nil
	return 1, nil
}

// 通过关键字直接找到
func (me *MongoDBExporter) FindGroupById(id int64) ([]pbmodel.GroupInfo, error) {
	collection := me.db.Collection(GroupTableName)

	filter := bson.M{"groupid": id}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var groups []pbmodel.GroupInfo
	if err := cursor.All(context.Background(), &groups); err != nil {
		return nil, err
	}

	//if len(groups) == 0 {
	//	return nil, err
	//}
	//
	//// 只要不是0个，开始检查公开属性
	//g := groups[0]
	//v, ok := g.Params["v"] // 如果未设置，或者设置为公开
	//if !ok || v == "pub" {
	//	return groups, nil
	//}
	//
	//c, ok := g.Params["code"] // 必须存在
	//if !ok || c != code {
	//	return nil, errors.New("validate code is not correct")
	//}

	return groups, nil
}

// 通过名字或者TAG字段来查找

// 备注
// 错误的做法：{ groupid: 1, groupname: 1 }
// 顺序反了，搜索用不上
// db.group.createIndex({
// groupname: 1,
// groupid: 1
// })
//
// db.group.createIndex({
// tags: 1,
// groupid: 1
// })
// MongoDB 在遇到 $or 时：
// 每个 $or 分支单独走索引
// 然后 merge 结果
// ⚠️ 一个索引不能同时服务两个 $or 分支
func (me *MongoDBExporter) FindGroupByKeyword(key string, fromId int64, bFilter bool, pageSize int64) ([]*pbmodel.GroupInfo, error) {

	collection := me.db.Collection(GroupTableName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"groupname": key},
			{"tags": key}, // tags 为数组时是 OK 的
		},
		"groupid": bson.M{"$gt": fromId},
	}

	if bFilter {
		filter["params.visibility"] = bson.M{"$ne": "private"}
	}

	opts := options.Find().
		SetSort(bson.M{"groupid": 1}). // 升序
		SetLimit(pageSize).
		SetMaxTime(10 * time.Second)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	groups := make([]*pbmodel.GroupInfo, 0)

	for cursor.Next(ctx) {
		var group pbmodel.GroupInfo
		if err := cursor.Decode(&group); err != nil {
			continue
		}
		groups = append(groups, &group)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// 2025-11-08 for debug
func (me *MongoDBExporter) LoadPbUsersFromId(startId int64, limit int64) ([]*pbmodel.UserInfo, error) {
	collection := me.db.Collection(UserTableName)

	// 构建查询条件：userid >= startId
	filter := bson.M{
		"userid": bson.M{"$gte": startId},
	}

	// 设置查询选项
	opts := options.Find().
		SetSort(bson.D{{Key: "userid", Value: 1}}). // 按 userid 升序排序
		SetLimit(limit)                             // 限制返回数量

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*pbmodel.UserInfo
	for cursor.Next(context.Background()) {
		var user pbmodel.UserInfo
		if err := cursor.Decode(&user); err != nil {
			fmt.Println("decode error:", err)
			continue
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (me *MongoDBExporter) LoadUsersFromId(startId int64, limit int64) ([]*model.User, error) {
	collection := me.db.Collection(UserTableName)

	// 查询条件
	filter := bson.M{
		"userid": bson.M{"$gte": startId},
	}

	// 排序 + 限制
	opts := options.Find().
		SetSort(bson.D{{Key: "userid", Value: 1}}).
		SetLimit(limit)

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*model.User

	for cursor.Next(context.Background()) {
		var info pbmodel.UserInfo
		if err := cursor.Decode(&info); err != nil {
			fmt.Println("decode userinfo error:", err)
			continue
		}

		// 创建内存 User 对象
		u := model.NewUserFromInfo(&info)
		users = append(users, u)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
