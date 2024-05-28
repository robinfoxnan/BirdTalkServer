package db

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (me *MongoDBExporter) SaveNewFile(f *model.FileInfo) error {

	// 选择要保存数据的数据库和集合
	collection := me.db.Collection(FileTableName)

	// 将用户信息对象转换为 MongoDB 文档
	bsonData, err := bson.Marshal(f)
	_, err = collection.InsertOne(context.Background(), bsonData)
	if err != nil {
		return err
	}

	return nil
}

// 查找

func (me *MongoDBExporter) FindFileByHash(keyword string) (*model.FileInfo, error) {
	fileLst, err := me.findFileByFieldAndTmOrder("hashcode", keyword, utils.GetTimeStamp())

	if fileLst == nil || len(fileLst) == 0 {
		return nil, err
	}

	return &fileLst[0], nil
}

func (me *MongoDBExporter) findFileByFieldAndTmOrder(field string, keyword interface{}, bigTm int64) ([]model.FileInfo, error) {
	collection := me.db.Collection(FileTableName)

	filter := bson.M{field: keyword, "tm": bson.M{"$lt": bigTm}}

	// 按照 tm 倒序排列
	// 设置排序选项
	findOptions := options.Find().SetSort(bson.D{{"tm", -1}}).SetLimit(100).SetMaxTime(time.Second * 10)

	// 执行查询并排序
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var fileLst []model.FileInfo
	if err = cursor.All(context.Background(), &fileLst); err != nil {
		return nil, err
	}

	return fileLst, nil
}

func (me *MongoDBExporter) FindFileByGroup(gid, tm int64) ([]model.FileInfo, error) {
	return me.findFileByFieldAndTmOrder("gid", gid, tm)
}

func (me *MongoDBExporter) FindFileByUser(uid, tm int64) ([]model.FileInfo, error) {
	return me.findFileByFieldAndTmOrder("userid", uid, tm)
}

func (me *MongoDBExporter) FindFileByTag(keyword string) ([]model.FileInfo, error) {
	return me.findFileByFieldAndTmOrder("tags", keyword, utils.GetTimeStamp())
}

func (me *MongoDBExporter) FindFileById(hashcode string) (*model.FileInfo, error) {
	collection := me.db.Collection(FileTableName)

	filter := bson.M{"hashcode": hashcode}

	// 按照 tm 倒序排列
	// 设置排序选项
	findOptions := options.Find().SetLimit(1).SetMaxTime(time.Second * 10)

	// 执行查询并排序
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var fileLst []model.FileInfo
	if err = cursor.All(context.Background(), &fileLst); err != nil {
		return nil, err
	}

	if len(fileLst) == 0 {
		return nil, nil
	}

	return &fileLst[0], nil
}
