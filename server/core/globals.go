package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

// 当前协议版本
const ProtocolVersion int = 1

// 全局变量
type GlobalVars struct {
	maxMessageSize int64 // 最大包长

	ss  *SessionCache     // 会话管理
	uc  *model.UserCache  // 用户内存缓存
	grc *model.GroupCache // 群组信息内存缓存

	snow *utils.Snowflake // 雪花算法

	scyllaCli *db.Scylla
	redisCli  *db.RedisClient
	mongoCli  *db.MongoDBExporter

	Logger *zap.Logger
	Config *LocalConfig
}

var Globals GlobalVars

// 初始化构造函数
func init() {
	Globals = GlobalVars{}
	Globals.ss = NewSessionCache()
	Globals.uc = model.NewUserCache()
	Globals.Logger = utils.CreateLogger()

}

// 加载配置
func (g *GlobalVars) LoadConfig(fileName string) error {
	var err error
	g.Config, err = LoadConfig(fileName)
	return err
}

func (g *GlobalVars) InitWithConfig() error {
	Globals.maxMessageSize = 10 * (1 << 20) // 10M
	Globals.snow = utils.NewSnowflake(1, 1)
	return nil
}

func (g *GlobalVars) InitDb() error {
	var err error
	g.redisCli, err = db.NewRedisClient(g.Config.Redis.RedisHost, g.Config.Redis.RedisPwd)
	if err != nil {
		fmt.Println(err)
		return err
	}

	g.redisCli.InitData() // 初始一些数据

	hosts := strings.Split(g.Config.ScyllaDb.Host, ",")
	g.scyllaCli, err = db.NewScyllaClient(hosts, g.Config.ScyllaDb.User, g.Config.ScyllaDb.Pwd)
	if err != nil {
		fmt.Println(err)
		return err
	}

	g.mongoCli, err = db.NewMongoDBExporter(g.Config.MongoDb.MongoHost, g.Config.MongoDb.DbName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
