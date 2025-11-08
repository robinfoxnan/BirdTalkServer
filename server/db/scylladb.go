package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"time"
)

// 定义 Scylla 封装结构
type Scylla struct {
	clusterConfig *gocql.ClusterConfig
	session       gocqlx.Session
}

var ScyllaClient *Scylla

// ---------------------------
// 自定义常量重连策略（新版 gocql v1.6 没有内置）
// ---------------------------
//type ConstantReconnectionPolicy struct {
//	Interval   time.Duration
//	MaxRetries int // 0 = 无限重试
//}
//
//func (p *ConstantReconnectionPolicy) Attempt(iter int) bool {
//	if p.MaxRetries == 0 {
//		return true
//	}
//	return iter < p.MaxRetries
//}
//func (p *ConstantReconnectionPolicy) GetMaxRetries() int {
//	return p.MaxRetries
//}
//
//func (p *ConstantReconnectionPolicy) GetInterval(iter int) time.Duration {
//	return p.Interval
//}

// ---------------------------
// 主连接函数
// --------------------------
func NewScyllaClient(constr []string, user, pwd string) (*Scylla, error) {
	// 连接到 ScyllaDB 集群
	client := Scylla{}
	client.clusterConfig = gocql.NewCluster(constr...)
	//cluster.Keyspace = "chatdata"
	if user != "" {
		client.clusterConfig.Authenticator = gocql.PasswordAuthenticator{
			Username: user,
			Password: pwd,
		}
	}

	// ---------------------------
	// 基本配置
	// 在 gocql v1.6 中： ClusterConfig.Timeout 定义的是 Query 的全局超时时间（也就是 Query 请求的超时时间）。
	// 这个超时时间包括：向节点发送请求,节点处理请求,接收响应
	// ---------------------------
	client.clusterConfig.ConnectTimeout = 10 * time.Second
	client.clusterConfig.Timeout = 20 * time.Second
	client.clusterConfig.NumConns = 3
	client.clusterConfig.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	// 这是操作系统内核功能，不是 gocql 自己实现的心跳。
	// 当一个 TCP 连接 空闲超过 SocketKeepalive 设置的时间时，操作系统内核会自动发送一个 小的 TCP 保活探测包（keepalive probe） 到对端。
	client.clusterConfig.SocketKeepalive = 5 * time.Second

	client.clusterConfig.MaxWaitSchemaAgreement = 15 * time.Second

	// ---------------------------
	// 自定义重连策略（gocql v1.6 不再内置 ConstantReconnectionPolicy）
	// ---------------------------
	//client.clusterConfig.ReconnectionPolicy = &ConstantReconnectionPolicy{
	//	Interval:   5 * time.Second,
	//	MaxRetries: 0, // 无限重试
	//}

	// ---------------------------
	// 创建 session
	// ---------------------------
	session, err := client.clusterConfig.CreateSession()
	if err != nil {
		fmt.Println("create scyllaDB session error", err)
		return nil, err
	}

	client.session, err = gocqlx.WrapSession(session, nil)
	if err != nil {
		session.Close()
		fmt.Println("wrap scyllaDB session error", err)
		return nil, err
	}

	return &client, nil
}

func (me *Scylla) Close(friend *model.FriendStore) {
	me.session.Close()
}

func (me *Scylla) Exec(cql string) error {
	// 使用 CQL 语句创建 Keyspace
	if err := me.session.Session.Query(cql).Exec(); err != nil {
		//fmt.Println("creating Keyspace error:", err)
		return err
	}

	//fmt.Println("Keyspace created ok")
	return nil
}

// 在后台定期执行简单查询，比如
func (me *Scylla) CheckConnection() error {

	err := me.Exec("SELECT now() FROM system.local")
	return err
}

func (me *Scylla) StartAutoReconnect(constr []string, user, pwd string, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			if err := me.CheckConnection(); err != nil {
				fmt.Println("⚠️ Scylla connection lost, trying to reconnect...:", err)
				me.session.Close()

				newClient, err := NewScyllaClient(constr, user, pwd)
				if err != nil {
					fmt.Println("❌ Scylla reconnect failed:", err)
					continue
				}
				me.session = newClient.session
				fmt.Println("✅ Scylla reconnected successfully")
			}
		}
	}()
}

// //////////////////////////////////////////////////////////////////////////
const cqlCreateKeySpaceChatUser = `CREATE KEYSPACE IF NOT EXISTS  chatuser
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}`

const cqlCreateKeySpaceChatGroup = `CREATE KEYSPACE IF NOT EXISTS  chatgroup
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}`

const cqlCreateKeySpaceChatData = `CREATE KEYSPACE IF NOT EXISTS  chatdata
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}`

const cqlCreateKeySpaceChatUserOp = `CREATE KEYSPACE IF NOT EXISTS  chatuserop
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}`

// //////////////////////////////////////////////////////////////////////////
const cqlCreateTableFollow = `CREATE TABLE IF NOT EXISTS chatuser.following (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			PRIMARY KEY (pk, uid1, uid2)
		)`

const cqlCreateTableFans = `CREATE TABLE IF NOT EXISTS chatuser.fans (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			PRIMARY KEY (pk, uid1, uid2)
		)`

const cqlCreateTableBlock = `CREATE TABLE IF NOT EXISTS chatuser.block (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			perm int,
			PRIMARY KEY (pk, uid1, uid2)
		)`

// 2025-11-05 add
const cqlCreateTableMutual = `CREATE TABLE IF NOT EXISTS chatuser.friends (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			label text,
			perm int,
			PRIMARY KEY (pk, uid1, uid2)
		)`

const cqlCreateTableGroupMem = `CREATE TABLE IF NOT EXISTS chatgroup.members (
			pk smallint,
			role smallint,
			gid bigint,
			uid bigint,
			tm bigint,
			nick text,
			PRIMARY KEY (pk, gid, uid)
		)`

const cqlCreateTableUinG = `CREATE TABLE IF NOT EXISTS chatgroup.uing (
			pk smallint,
			uid bigint,
			gid bigint,
			PRIMARY KEY (pk, uid, gid)
		)`

const cqlCreateTablePChat = `CREATE TABLE IF NOT EXISTS  chatdata.pchat (
			pk smallint,
			uid1 bigint, 
			uid2 bigint,
			id bigint,
			usid bigint,
			tm bigint,
			tm1 bigint,
			tm2 bigint,
			io tinyint,
			st tinyint,
			ct tinyint,
			mt tinyint,
			draf blob,
			pr  varint,
			ref varint,
			PRIMARY KEY (pk, uid1, id)
		)`

// robin add 2025-11-07
const cqlCreateTablePChatTopic = `CREATE MATERIALIZED VIEW IF NOT EXISTS chatdata.pchat_topic AS
		SELECT pk, uid1, uid2, id, usid, tm, tm1, tm2, io, st, ct, mt, draf, pr, ref
		FROM chatdata.pchat
		WHERE pk IS NOT NULL
		  AND uid1 IS NOT NULL
		  AND uid2 IS NOT NULL
		  AND id IS NOT NULL
		PRIMARY KEY ((pk, uid1, uid2), id)
		WITH CLUSTERING ORDER BY (id DESC)`

const cqlCreateTableGChat = `CREATE TABLE IF NOT EXISTS  chatdata.gchat (
			pk smallint,
			gid bigint,
			uid bigint, 
			id bigint,
			usid bigint,
			tm bigint,
			res tinyint,
			st tinyint,
			ct tinyint,
			mt tinyint,
			draf blob,
			pr  varint,
			ref varint,
			PRIMARY KEY (pk, gid, id)
		)`

const cqlCreateTableUserOp = `CREATE TABLE IF NOT EXISTS  chatuserop.userop (
			pk SMALLINT,
			uid1 BIGINT,
			uid2 BIGINT,
			gid BIGINT,
			id BIGINT,
			usid BIGINT,
			tm BIGINT,
			tm1 BIGINT,
			tm2 BIGINT,
			io TINYINT,
			st TINYINT,
			cmd TINYINT,
			ret TINYINT,
			mask INT,
			ref BIGINT,
			draf BLOB,
			PRIMARY KEY (pk, uid1, id)
		)`

const cqlCreateTableGroupOp = `CREATE TABLE IF NOT EXISTS  chatuserop.groupop (
			pk SMALLINT,
			gid BIGINT,
			id BIGINT,
			uid1 BIGINT,
			uid2 BIGINT,
			usid BIGINT,
			tm BIGINT,
			tm1 BIGINT,
			tm2 BIGINT,
			io TINYINT,
			st TINYINT,
			cmd TINYINT,
			ret TINYINT,
			mask INT,
			ref BIGINT,
			draf BLOB,
			PRIMARY KEY (pk, gid, id)
		)`

// ////////////////////////////////////////////////////////
// 用户关系
const FollowingTableName = "chatuser.following"
const FansTableName = "chatuser.fans"
const BlockTableName = "chatuser.block"
const MutualFriendTableName = "chatuser.friends"

// 群组关系
const GroupMemberTableName = "chatgroup.members"
const UserInGroupTableName = "chatgroup.uing"

// 聊天消息
const PrivateChatTableName = "chatdata.pchat"
const GroupChatTableName = "chatdata.gchat"
const PrivateChatViewName = "chatdata.pchat_topic"

// 好友申请，群申请记录
const UserOpTableName = "chatuserop.userop"

// 群组的操作记录
const GroupOpTableName = "chatuserop.groupop"

var initCqlList = []string{cqlCreateKeySpaceChatUser, cqlCreateKeySpaceChatGroup, cqlCreateKeySpaceChatData, cqlCreateKeySpaceChatUserOp,
	cqlCreateTableFollow, cqlCreateTableFans, cqlCreateTableBlock, cqlCreateTableMutual,
	cqlCreateTableGroupMem, cqlCreateTableUinG,
	cqlCreateTablePChat, cqlCreateTableGChat, cqlCreateTablePChatTopic,
	cqlCreateTableUserOp, cqlCreateTableGroupOp,
}

func (me *Scylla) Init() error {
	for _, cql := range initCqlList {
		//fmt.Println(cql)
		err := me.Exec(cql)
		if err != nil {
			return err
		}
	}
	fmt.Println("init keyspace and tables ok")
	return nil
}
