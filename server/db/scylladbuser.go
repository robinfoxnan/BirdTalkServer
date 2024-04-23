package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// 定义关注表，粉丝表的元数据
var metaFriend = table.Metadata{
	Columns: []string{"pk", "uid1", "uid2", "tm", "nick"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "uid2"},
}

var metaBlock = table.Metadata{
	Columns: []string{"pk", "uid1", "uid2", "tm", "nick", "perm"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "uid2"},
}

// 重复的数据会覆盖，而不是简单的报错
func (me *Scylla) InsertFollowing(friend *model.FriendStore, fan *model.FriendStore) error {

	// 同时加入粉丝表
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 设置关注
	insertFollowQry := qb.Insert(FollowingTableName).Columns(metaFriend.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertFollowQry.Release()
	batch.Query(insertFollowQry.Statement(), friend.Pk, friend.Uid1, friend.Uid2, friend.Tm, friend.Nick)

	// 设置粉丝
	insertFanQry := qb.Insert(FansTableName).Columns(metaFriend.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertFanQry.Release()
	batch.Query(insertFanQry.Statement(), fan.Pk, fan.Uid1, fan.Uid2, fan.Tm, fan.Nick)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 拉黑或者设置权限
func (me *Scylla) InsertBlock(friend *model.BlockStore) error {

	// Insert song using query builder.
	insertChat := qb.Insert(BlockTableName).Columns(metaBlock.Columns...).Query(me.session).Consistency(gocql.One)
	insertChat.BindStruct(friend)
	if err := insertChat.ExecRelease(); err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

// 同时操作2个表
func (me *Scylla) DeleteFollowing(pk1, pk2 int16, uid1, uid2 int64) error {
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 取消关注
	builder1 := qb.Delete(FollowingTableName)

	builder1.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	query1 := builder1.Query(me.session)
	defer query1.Release()
	query1.Consistency(gocql.One)
	batch.Query(query1.Statement(), pk1, uid1, uid2)

	// 取消粉丝
	builder2 := qb.Delete(FansTableName)
	builder2.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	query2 := builder2.Query(me.session)
	defer query2.Release()
	query2.Consistency(gocql.One)
	batch.Query(query2.Statement(), pk2, uid2, uid1)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (me *Scylla) DeleteBlock(pk int16, uid1, uid2 int64) error {
	// 构建删除语句
	builder := qb.Delete(BlockTableName)

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

func (me *Scylla) FindFollowing(pk, uid1, from int64, pageSize uint) ([]model.FriendStore, error) {
	return me.FindFriendStore(pk, uid1, from, pageSize, FollowingTableName)
}

func (me *Scylla) FindFans(pk, uid1, from int64, pageSize uint) ([]model.FriendStore, error) {
	return me.FindFriendStore(pk, uid1, from, pageSize, FansTableName)
}

func (me *Scylla) FindFollowingExact(pk, uid1, fid int64) (*model.FriendStore, error) {
	return me.FindFriendStoreExact(pk, uid1, fid, FollowingTableName)
}

func (me *Scylla) FindFansExact(pk, uid1, fid int64) (*model.FriendStore, error) {
	return me.FindFriendStoreExact(pk, uid1, fid, FansTableName)
}

// 精确查找
func (me *Scylla) FindFriendStoreExact(pk, uid1, uid2 int64, table string) (*model.FriendStore, error) {
	builder := qb.Select(table).Columns(metaFriend.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("uid2", qb.ASC)

	//builder.AllowFiltering()

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid1, uid2)

	var friendList []model.FriendStore

	err := q.Select(&friendList)
	if err != nil || friendList == nil || len(friendList) == 0 {
		//fmt.Println(err)
		return nil, err
	}
	return &friendList[0], nil
}

// 查询关注
// 数据库不支持偏移，但是可以按照用户id排序，当返回数量少于分页，说明都取完了，
// 否则，则需要按照最后一个ID继续
func (me *Scylla) FindFriendStore(pk, uid1, uid2 int64, pageSize uint, table string) ([]model.FriendStore, error) {
	//chatTable := table.New(pchatMetadata)
	builder := qb.Select(table).Columns(metaFriend.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Gt("uid2"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("uid2", qb.ASC)

	//builder.AllowFiltering()
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

func (me *Scylla) FindBlocksExact(pk, uid1, fid int64) (*model.BlockStore, error) {
	builder := qb.Select(BlockTableName).Columns(metaBlock.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("uid2", qb.ASC)

	//builder.AllowFiltering()
	//builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid1, fid)

	var friendList []model.BlockStore

	err := q.Select(&friendList)
	if err != nil || friendList == nil || len(friendList) == 0 {
		//fmt.Println(err)
		return nil, err
	}
	return &friendList[0], nil
}

// 查询所有的拉黑的名单
func (me *Scylla) FindBlocks(pk, uid1, from int64, pageSize uint) ([]model.BlockStore, error) {
	builder := qb.Select(BlockTableName).Columns(metaBlock.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Gt("uid2"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("uid2", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid1, from)

	var friendList []model.BlockStore

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

func (me *Scylla) CountFans(pk, uid1 int64) (int64, error) {
	return me.CountFriendStore(pk, uid1, FansTableName)
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

func (me *Scylla) SetFollowingNick(pk, uid1, uid2 int64, nick string) error {
	return me.setFriendStoreNick(pk, uid1, uid2, nick, FollowingTableName)
}

func (me *Scylla) SetFansNick(pk, uid1, uid2 int64, nick string) error {
	return me.setFriendStoreNick(pk, uid1, uid2, nick, FansTableName)
}

// 设置一个名字
func (me *Scylla) setFriendStoreNick(pk, uid1, uid2 int64, nick, table string) error {
	builder := qb.Update(table)

	builder.Set("nick")

	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(nick, pk, uid1, uid2)

	err := query.Exec()
	return err
}

func (me *Scylla) SetBlockPermission(pk, uid1, uid2 int64, perm int32) error {
	builder := qb.Update(BlockTableName)

	builder.Set("perm")

	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(perm, pk, uid1, uid2)

	err := query.Exec()
	return err
}

//////////////////////////////////////////////////////////////////////
