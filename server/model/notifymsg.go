package model

type NotifyType int16

const (
	add    NotifyType = iota // 0
	remove                   // 1
	update                   // 2
	expire                   // 3
	hset
	hupdate
	hremove
)

// 使用redis通道来广播，用户与群组的变化，同步到内存的缓存上
type CacheNotifyMsg struct {
	Key    int64      `json:"key,omitempty"` // 用户或者群的ID
	Field  int64      `json:"field,omitempty"`
	Value  string     `json:"value,omitempty"`
	Action NotifyType `json:"action,omitempty"` // 消息类型
	Type   int16      `json:"type,omitempty"`   // user 或者group
}
