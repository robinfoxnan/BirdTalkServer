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
	"strings"
	"time"
)

// go get go.mongodb.org/mongo-driver
const UserTableName = "users"
const UserTableIndex = "userid"
const GroupTableName = "groups"
const GroupTableIndex = "groupid"

// 内部使用
var userTableSearchFields = []string{"username", "email", "phone"}

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

// init 初始化业务相关的部分
func (me *MongoDBExporter) Init() error {
	// 在此进行初始化操作，如果有需要的话
	err := me.CreateIndex(UserTableName, UserTableIndex)
	err = me.CreateIndex(GroupTableName, GroupTableIndex)
	err = me.CreateIndexes(UserTableName, userTableSearchFields)
	return err
}

func (me *MongoDBExporter) CreateIndexes(tableName string, indexFields []string) error {
	collection := me.db.Collection(tableName)

	// 构建索引键的映射
	keys := bson.D{}
	for _, field := range indexFields {
		keys = append(keys, primitive.E{Key: field, Value: 1})
	}

	// 创建索引模型
	index := mongo.IndexModel{
		Keys: keys,
	}

	// 创建索引
	_, err := collection.Indexes().CreateOne(context.Background(), index)
	return err
}

// //////////////////////////////////////////////////////////////////////
func (me *MongoDBExporter) CreateIndex(tableName, fieldName string) error {
	collection := me.db.Collection(tableName)
	// 创建唯一索引
	indexOptions := options.Index().SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bson.M{fieldName: 1}, // 设置 userName 字段为唯一索引
		Options: indexOptions,
	}
	_, err := collection.Indexes().CreateOne(context.Background(), index)
	// 检查错误类型
	if err != nil {
		// 检查是否是索引重复错误，MongoDB 7.0.3 Community版本没有遇到
		if mongo.IsDuplicateKeyError(err) {
			e := fmt.Errorf("unique index on field '%s' already exists", fieldName)
			fmt.Println(e)
			return nil
		}
		// 处理其他类型的错误
		return err
	}

	return err
}

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

	return me.FindUserByField("userid", id)
}

// 这个字段没有设置唯一索引，出于性能考虑，在更改名字的时候可以检测是否唯一，
func (me *MongoDBExporter) FindUserByName(keyword string) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("username", keyword)
}

func (me *MongoDBExporter) FindUserByEmail(keyword string) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("email", keyword)
}

func (me *MongoDBExporter) FindUserByPhone(keyword string) ([]*pbmodel.UserInfo, error) {
	return me.FindUserByField("phone", keyword)
}

func (me *MongoDBExporter) FindUserByField(field string, keyword interface{}) ([]*pbmodel.UserInfo, error) {
	collection := me.db.Collection(UserTableName)

	//var filter bson.M
	//switch keyword.(type) {
	//case int64:
	//	// 如果关键字是 int64 类型，则直接返回
	//	i := keyword.(int64)
	//	// 构建查询条件
	//	filter = bson.M{field: i}
	//case string:
	//	str := keyword.(string)
	//	filter = bson.M{field: str}
	//default:
	//	return nil, errors.New("keyword type err")
	//}

	filter := bson.M{field: keyword}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*pbmodel.UserInfo
	//if err = cursor.All(context.Background(), &users); err != nil {
	//	return nil, err
	//}
	count := 0
	for cursor.Next(context.Background()) {
		var user pbmodel.UserInfo
		if err = cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, &user)
		count++

		// 检查结果数量是否达到限制
		if count >= limit {
			break
		}
	}

	if count == 0 {
		return nil, nil
	}

	return users, nil
}

// 根据用户名、邮件或手机号进行搜索
func (me *MongoDBExporter) FindUserByKeyword(keyword string) ([]*pbmodel.UserInfo, error) {
	collection := me.db.Collection(UserTableName)

	// 构建查询条件
	// 构建正则表达式查询条件
	//regexName := primitive.Regex{Pattern: keyword, Options: "i"} // "i" 表示不区分大小写
	// 构建正则表达式查询条件
	//regexMail := primitive.Regex{Pattern: "^" + keyword, Options: "i"} // "i" 表示不区分大小写
	filter := bson.M{
		"$or": []bson.M{
			//{"username": bson.M{"$regex": regexName}},
			//{"email": bson.M{"$regex": regexMail}},
			{"username": keyword},
			{"email": keyword},
			{"phone": keyword},
		},
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []*pbmodel.UserInfo
	//if err := cursor.All(context.Background(), &users); err != nil {
	//	return nil, err
	//}

	count := 0
	for cursor.Next(context.Background()) {
		var user pbmodel.UserInfo
		if err = cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, &user)
		count++

		// 检查结果数量是否达到限制
		if count >= limit {
			break
		}
	}

	if count == 0 {
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
func (me *MongoDBExporter) UpdateUserInfoPart(id int64, setData map[string]interface{}, unsetData []string) (int64, error) {
	collection := me.db.Collection(UserTableName)
	filter := bson.M{"userid": id} // 过滤条件

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

	fmt.Println("Group information has been saved successfully.")
	return nil

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
			"grouptags": g.Tags,
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
const limit = 100

func (me *MongoDBExporter) FindGroupByKeyword(key string, bFilter bool) ([]*pbmodel.GroupInfo, error) {
	collection := me.db.Collection(GroupTableName)

	// 构建查询条件，不区分大小写了，影响性能，精确查找
	// bson.M{"$regex": key, "$options": "i"}

	filter := bson.M{
		"$or": []bson.M{
			bson.M{"groupname": key}, // 精确匹配 groupname
			bson.M{"tags": key},      // 精确匹配 tags
		},
		"$and": []bson.M{
			bson.M{"params.v": bson.M{"$ne": "pri"}}, // 未设置为pri 私有
		},
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	cursor.SetMaxTime(time.Second * 10)
	var groups []*pbmodel.GroupInfo
	//if err = cursor.All(context.Background(), &groups); err != nil {
	//	return nil, err
	//}

	count := 0
	for cursor.Next(context.Background()) {
		var group pbmodel.GroupInfo
		if err = cursor.Decode(&group); err != nil {
			continue
		}
		// 过滤掉私有的
		if bFilter {
			if model.CheckGroupInfoIsPrivate(&group) {
				continue
			}
		}
		groups = append(groups, &group)
		count++

		// 检查结果数量是否达到限制
		if count >= limit {
			break
		}
	}

	if groups == nil || len(groups) == 0 {
		return nil, nil
	}

	return groups, nil
}
