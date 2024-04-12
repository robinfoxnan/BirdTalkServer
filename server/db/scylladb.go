package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type Scylla struct {
	clusterConfig *gocql.ClusterConfig
	session       gocqlx.Session
}

var ScyllaClient *Scylla

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

// ////////////////////////////////////////////////////////
// 用户关系
const FollowingTableName = "chatuser.following"
const FansTableName = "chatuser.fans"
const BlockTableName = "chatuser.block"

// 群组关系
const GroupMemberTableName = "chatgroup.members"
const UserInGroupTableName = "chatgroup.uing"

// 聊天消息
const PrivateChatTableName = "chatdata.pchat"
const GroupChatTableName = "chatdata.gchat"

// 好友申请，群申请记录
const UserOpTableName = "chatuserop.userop"

var initCqlList = []string{cqlCreateKeySpaceChatUser, cqlCreateKeySpaceChatGroup, cqlCreateKeySpaceChatData, cqlCreateKeySpaceChatUserOp,
	cqlCreateTableFollow, cqlCreateTableFans, cqlCreateTableBlock,
	cqlCreateTableGroupMem, cqlCreateTableUinG,
	cqlCreateTablePChat, cqlCreateTableGChat,
	cqlCreateTableUserOp,
}

func (me *Scylla) Init() error {
	for _, cql := range initCqlList {
		err := me.Exec(cql)
		if err != nil {
			return err
		}
	}
	fmt.Println("init keyspace and tables ok")
	return nil
}
