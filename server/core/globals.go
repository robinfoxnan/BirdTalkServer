package core

import (
	"birdtalk/server/model"
	"birdtalk/server/utils"
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
}

var Globals GlobalVars

// 初始化构造函数
func init() {
	Globals = GlobalVars{}
	Globals.ss = &SessionCache{sessionMap: utils.NewConcurrentMap[int64, *Session]()}
	Globals.maxMessageSize = 10 * (1 << 20) // 10M
	Globals.snow = utils.NewSnowflake(1, 1)

}
