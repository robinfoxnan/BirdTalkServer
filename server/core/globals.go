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

	snow    *utils.Snowflake // 雪花算法
	segment *utils.SegmentJieba

	scyllaCli *db.Scylla
	redisCli  *db.RedisClient
	mongoCli  *db.MongoDBExporter

	emailWorkerManager *Manager[Task, *EmailWorker]

	Logger    *zap.Logger
	Config    *LocalConfig
	GeoHelper *GeoIPHelper
}

var Globals GlobalVars

// 初始化构造函数
func init() {
	Globals = GlobalVars{}
	Globals.ss = NewSessionCache()
	Globals.uc = model.NewUserCache()
	Globals.grc = model.NewGroupCache()
	Globals.Logger = utils.CreateLogger()
	Globals.segment = utils.NewSegment()
	Globals.GeoHelper = nil

}

// 加载配置
func (g *GlobalVars) LoadConfig(fileName string) error {
	var err error
	g.Config, err = LoadConfig(fileName)
	return err
}

func (g *GlobalVars) InitWithConfig() error {
	g.maxMessageSize = 10 * (1 << 20) // 10M
	g.snow = utils.NewSnowflake(1, 1)
	n := int64(g.Config.Email.Workers)
	if n < 2 {
		n = 2
	}
	g.emailWorkerManager = NewEmailWorkerManager(n)

	err := utils.InitFont(Globals.Config.Server.AvatarFont)
	if err != nil {
		Globals.Logger.Error("load font error", zap.Error(err))
	}

	g.GeoHelper, err = NewGeoIPHelper(Globals.Config.Server.GeoLite2Path)
	if err != nil {
		Globals.Logger.Error("load geolite2 error", zap.Error(err))
	}
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
