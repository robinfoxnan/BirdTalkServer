package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// 定义关注表，拉黑表的元数据
var pchatMetaFriend = table.Metadata{
	Columns: []string{"pk", "uid1", "uid2", "tm", "nick"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "uid2"},
}

// 重复的数据会覆盖，而不是简单的报错
func (me *Scylla) InsertFollowing(friend *model.FriendStore) error {
	return me.InsertFriendStore(friend, FollowingTableName)
}

func (me *Scylla) InsertBlock(friend *model.FriendStore) error {
	return me.InsertFriendStore(friend, BlockTableName)
}

func (me *Scylla) InsertFriendStore(friend *model.FriendStore, table string) error {

	// Insert song using query builder.
	insertChat := qb.Insert(table).Columns(pchatMetaFriend.Columns...).Query(me.session).Consistency(gocql.One)
	insertChat.BindStruct(friend)
	if err := insertChat.ExecRelease(); err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

func (me *Scylla) DeleteFollowing(pk int16, uid1, uid2 int64) error {
	return me.DeleteFriendStore(pk, uid1, uid2, FollowingTableName)
}

func (me *Scylla) DeleteBlock(pk int16, uid1, uid2 int64) error {
	return me.DeleteFriendStore(pk, uid1, uid2, BlockTableName)
}

func (me *Scylla) DeleteFriendStore(pk int16, uid1, uid2 int64, table string) error {
	// 构建删除语句
	builder := qb.Delete(table)

	builder.Where(qb.Eq("pk"))
	builder.Where(qb.Eq("uid1"))
	builder.Where(qb.Eq("uid2"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(pk, uid1, uid2)

	err := query.Exec()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

func (me *Scylla) FindFollowing(pk, uid1, uid2 int64, pageSize uint) ([]model.FriendStore, error) {
	return me.FindFriendStore(pk, uid1, uid2, pageSize, FollowingTableName)
}

func (me *Scylla) FindBlock(pk, uid1, uid2 int64, pageSize uint) ([]model.FriendStore, error) {
	return me.FindFriendStore(pk, uid1, uid2, pageSize, BlockTableName)
}

// 查询关注
// 数据库不支持偏移，但是可以按照用户id排序，当返回数量少于分页，说明都取完了，
// 否则，则需要按照最后一个ID继续
func (me *Scylla) FindFriendStore(pk, uid1, uid2 int64, pageSize uint, table string) ([]model.FriendStore, error) {
	//chatTable := table.New(pchatMetadata)
	builder := qb.Select(table).Columns(pchatMetaFriend.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Gt("uid2"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("uid2", qb.ASC)

	builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid1, uid2)

	var friendList []model.FriendStore

	err := q.Select(&friendList)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return friendList, nil
}
func (me *Scylla) CountFollowing(pk, uid1 int64) (int64, error) {
	return me.CountFriendStore(pk, uid1, FollowingTableName)
}

func (me *Scylla) CountBlock(pk, uid1 int64) (int64, error) {
	return me.CountFriendStore(pk, uid1, BlockTableName)
}

func (me *Scylla) CountFriendStore(pk, uid1 int64, table string) (int64, error) {
	builder := qb.Select(table)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"))
	builder.Count("uid1")

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid1)

	var count int64

	// 通过 Scan 函数将结果存储到 count 变量中
	if err := q.Scan(&count); err != nil {
		fmt.Println(err)
		return 0, err
	}
	//fmt.Println(count)

	return count, nil
}

// ////////////////////////////////////////////////////////////////////////
var pchatMetaFans = table.Metadata{
	Columns: []string{"pk", "uid1", "uid2", "tm", "nick", "perm"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "uid2"},
}
