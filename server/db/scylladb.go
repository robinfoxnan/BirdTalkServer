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

// //////////////////////////////////////////////////////////////////////////
const cqlCrateTableFollow = `CREATE TABLE IF NOT EXISTS chatuser.following (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			PRIMARY KEY (pk, uid1, uid2)
		)`

const cqlCrateTableFans = `CREATE TABLE IF NOT EXISTS chatuser.fans (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
            perm int,
			PRIMARY KEY (pk, uid1, uid2)
		)`

const cqlCrateTableBlock = `CREATE TABLE IF NOT EXISTS chatuser.block (
			pk smallint,
			uid1 bigint,
			uid2 bigint,
			tm bigint,
			nick text,
			PRIMARY KEY (pk, uid1, uid2)
		)`

// ////////////////////////////////////////////////////////
const FollowingTableName = "chatuser.following"
const FansTableName = "chatuser.fans"
const BlockTableName = "chatuser.block"
const GroupMemberTableName = "chatgroup.members"
const PrivateChatTableName = "chatdata.pchat"
const GroupChatTableName = "chatdata.gchat"

var initCqlList = []string{cqlCreateKeySpaceChatUser, cqlCreateKeySpaceChatGroup, cqlCreateKeySpaceChatData,
	cqlCrateTableFollow, cqlCrateTableFans, cqlCrateTableBlock}

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
