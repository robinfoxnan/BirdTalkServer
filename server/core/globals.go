package core

import "birdtalk/server/utils"

// 全局变量
type GlobalVars struct {
	maxMessageSize int64 // 最大包长

	ss  *SessionStore // 会话管理
	uc  *UserCache    // 用户内存缓存
	grc *GroupCache   // 群组信息内存缓存

	snow *utils.Snowflake // 雪花算法
}

var Globals GlobalVars

// 初始化构造函数
func init() {
	Globals = GlobalVars{}
	Globals.ss = &SessionStore{sessionMap: utils.NewConcurrentMap[int64, *Session]()}
	Globals.maxMessageSize = 10 * (1 << 20) // 10M
	Globals.snow = utils.NewSnowflake(1, 1)

}
