package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/utils"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
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

	Logger   *zap.Logger
	logLevel zap.AtomicLevel

	Config      *LocalConfig
	GeoHelper   *GeoIPHelper
	maxPageSize uint
}

var Globals GlobalVars

// 初始化构造函数
func init() {
	Globals = GlobalVars{}
	Globals.ss = NewSessionCache()
	Globals.uc = model.NewUserCache()
	Globals.grc = model.NewGroupCache()
	Globals.logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	Globals.Logger = utils.CreateLogger(&Globals.logLevel)
	Globals.segment = utils.NewSegment()
	Globals.GeoHelper = nil
	Globals.maxPageSize = 100

}

// 加载配置
func (g *GlobalVars) LoadConfig(fileName string) error {
	var err error
	g.Config, err = LoadConfig(fileName)
	return err
}

func (g *GlobalVars) InitWithConfig() error {

	key, err := hex.DecodeString(g.Config.Server.TokenHex)
	if err != nil {
		panic(err)
	}
	g.Config.Server.TokenKey = key
	switch g.Config.Server.LogLevel {
	case "debug":
		g.logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		g.logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		g.logLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		g.logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		g.logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		g.logLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "disabled":
		g.logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		g.logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	g.maxMessageSize = 10 * (1 << 20) // 10M
	g.snow = utils.NewSnowflake(1, 1)
	n := int64(g.Config.Email.Workers)
	if n < 2 {
		n = 2
	}
	g.emailWorkerManager = NewEmailWorkerManager(n)

	err = utils.InitFont(Globals.Config.Server.AvatarFont)
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
	g.scyllaCli.Init()

	// 2025-01-12 added by robin
	g.scyllaCli.StartAutoReconnect(hosts, g.Config.ScyllaDb.User, g.Config.ScyllaDb.Pwd, 10*time.Second)

	g.mongoCli, err = db.NewMongoDBExporter(g.Config.MongoDb.MongoHost, g.Config.MongoDb.DbName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 通过指纹来查看用户的基本信息
func LoadUserByKeyPrint(keyPrint int64) (int64, *utils.KeyExchange, error) {
	return Globals.redisCli.LoadToken(keyPrint)
}
