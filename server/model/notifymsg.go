package model

type NotifyType int16

/*
用户或者群组的用户改变状态时候使用广播方式通知其他的
*/
const (
	NotifyTypeNone         NotifyType = iota // 0
	NotifyUserOnline                         // 1
	NotifyUserOffline                        // 2
	NotifyUserSetInfo                        // 3
	NotifyUserAddFollow                      // 4
	NotifyUserDelFollow                      // 5
	NotifyUserPermission                     // 6
	NotifyGroupCreate                        // 7
	NotifyGroupDis                           // 8
	NotifyGroupMemJoin                       // 9
	NotifyGroupMemLeave                      // 10
	NotifyGroupSetInfo                       // 11
	NotifyGroupMemOnline                     // 12
	NotifyGroupMemOffline                    // 13
	NotifyGroupChangeAdmin                   // 14
	NotifyGroupChangeOwner                   // 15
)

// 使用redis通道来广播，用户与群组的变化，同步到内存的缓存上
type CacheNotifyMsg struct {
	NType  NotifyType `json:"nt" ` // 消息类型
	FromId int64      `json:"fromId" `
	Action string     `json:"action" `
	Uid    int64      `json:"uid" `
	Gid    int64      `json:"gid" `
}
