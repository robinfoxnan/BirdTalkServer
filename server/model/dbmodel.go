package model

// 下面2个结构对应数据库中结构
// 关注和粉丝的表一样，没有权限一项
type FriendStore struct {
	Pk   int16
	Uid1 int64
	Uid2 int64
	Tm   int64
	Nick string
}

type BlockStore struct {
	FriendStore
	Perm int64
}

// 数据库中记录群组成员的结构体
type GroupMemberStore struct {
	Pk   int16
	Role int16
	Gid  int64
	Uid  int64
	Tm   int64
	Nick string
}

type UserInGStore struct {
	Pk  int16
	Uid int64
	Gid int64
}

const (
	DrafStateNormal = 1 << iota // 常规 1
	DrafStateArrive             // 送达 2
	DrafStateRead               // 已读 4
	DrafStateDel                // 撤回 8

)

const ChatDataIOOut = 0
const ChatDataIOIn = 1
const UserOpResultRefuse = 2
const UserOpResultOk = 1

// 私聊消息存储结构
type PChatDataStore struct {
	Pk   int16 `db:"pk"`
	Uid1 int64 `db:"uid1"`
	Uid2 int64 `db:"uid2"`
	Id   int64 `db:"id"`
	Usid int64 `db:"usid"`
	Tm   int64 `db:"tm"`
	Tm1  int64 `db:"tm1"`
	Tm2  int64 `db:"tm2"`

	Io    int8   `db:"io"`  // 0=out, 1=in
	St    int8   `db:"st"`  // 0=normal, 1=送达,2阅读，
	Ct    int8   `db:"ct"`  // 0=p2p_plain, 1=system, 2=p2_encrypted,
	Mt    int8   `db:"mt"`  // 0=text, 1=pic, 2=
	Print int64  `db:"pr"`  // 秘钥哈希的低8字节作为指纹
	Ref   int64  `db:"ref"` // 引用
	Draf  []byte `db:"draf"`
}

// 群聊消息存储结构
// 群成员的操作记录作为消息记录直接保存，
type GChatDataStore struct {
	Pk   int16 `db:"pk"`
	Gid  int64 `db:"gid"`
	Uid  int64 `db:"uid"`
	Id   int64 `db:"id"`
	Usid int64 `db:"usid"`
	Tm   int64 `db:"tm"`
	Res  int8  `db:"res"` // 保留
	St   int8  `db:"st"`  // 0=normal, 1=送达,2阅读，
	Ct   int8  `db:"ct"`  // 0=普通，1=广播
	Mt   int8  `db:"mt"`  // 0=text, 1=pic, 2=

	Print int64  `db:"pr"`  // 秘钥哈希的低8字节作为指纹
	Ref   int64  `db:"ref"` // 引用
	Draf  []byte `db:"draf"`
}

// 用户好友相关记录的存储，
// 群组操作相关的记录存储
type CommonOpStore struct {
	Pk   int16 `db:"pk"`
	Uid1 int64 `db:"uid1"`
	Uid2 int64 `db:"uid2"`
	Gid  int64 `db:"gid"`
	Id   int64 `db:"id"`
	Usid int64 `db:"usid"`
	Tm   int64 `db:"tm"`
	Tm1  int64 `db:"tm1"`
	Tm2  int64 `db:"tm2"`

	Io   int8   `db:"io"`   // 0=out, 1=in
	St   int8   `db:"st"`   // 0=normal, 1=送达,2阅读，
	Cmd  int8   `db:"cmd"`  // 0=删除好友; 1=申请好友;2= 设置权限; 10=邀请好友；11=申请加入；
	Ret  int8   `db:"ret"`  // 2=拒绝， 1=同意
	Mask int32  `db:"mask"` // 权限操作的掩码
	Ref  int64  `db:"ref"`  // 引用
	Draf []byte `db:"draf"` // 附加消息
}
