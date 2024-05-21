package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// 用于管理群组中的用户信息
var metaGroupMembers = table.Metadata{
	Columns: []string{"pk", "gid", "uid", "tm", "role", "nick"},
	PartKey: []string{"pk"},
	SortKey: []string{"gid", "uid"},
}

var metaUserInGroups = table.Metadata{
	Columns: []string{"pk", "uid", "gid"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid", "gid"},
}

// 用户加入组时候，在组成员表中加一条，在用户所属组中加一条
func (me *Scylla) InsertGroupMember(gmem *model.GroupMemberStore, uing *model.UserInGStore) error {

	// 同时加入粉丝表
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 设置关注
	insertGroupMemQry := qb.Insert(GroupMemberTableName).Columns(metaGroupMembers.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertGroupMemQry.Release()
	batch.Query(insertGroupMemQry.Statement(), gmem.Pk, gmem.Gid, gmem.Uid, gmem.Tm, gmem.Role, gmem.Nick)

	// 设置粉丝
	insertUseringQry := qb.Insert(UserInGroupTableName).Columns(metaUserInGroups.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertUseringQry.Release()
	batch.Query(insertUseringQry.Statement(), uing.Pk, uing.Uid, uing.Gid)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 设置成员在群众的昵称和角色
func (me *Scylla) SetGroupMemberNickRole(pk, gid, uid int64, nick string, role int32) error {
	builder := qb.Update(GroupMemberTableName)

	builder.Set("nick", "role")
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.Eq("uid"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(nick, role, pk, gid, uid)

	err := query.Exec()
	return err
}

// from 是查询起点，第一页应该从0开始
// 查询用户所在的群列表
func (me *Scylla) FindUserInGroups(pk, uid, from int64, pageSize uint) ([]model.UserInGStore, error) {
	builder := qb.Select(UserInGroupTableName).Columns(metaUserInGroups.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid"), qb.Gt("gid"))

	builder.OrderBy("uid", qb.ASC)
	builder.OrderBy("gid", qb.ASC)

	builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid, from)

	var lst []model.UserInGStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 查询组内成员列表
func (me *Scylla) FindGroupMembers(pk, gid, from int64, pageSize uint) ([]model.GroupMemberStore, error) {
	builder := qb.Select(GroupMemberTableName).Columns(metaGroupMembers.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.Gt("uid"))

	builder.OrderBy("gid", qb.ASC)
	builder.OrderBy("uid", qb.ASC)

	builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, from)

	var lst []model.GroupMemberStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 从群内移除成员时，需要操作2个表
func (me *Scylla) DeleteGroupMember(pk1, pk2 int16, gid, uid int64) error {
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 删除成员
	builder1 := qb.Delete(GroupMemberTableName)

	builder1.Where(qb.Eq("pk"), qb.Eq("gid"), qb.Eq("uid"))

	query1 := builder1.Query(me.session)
	defer query1.Release()
	query1.Consistency(gocql.One)
	batch.Query(query1.Statement(), pk1, gid, uid)

	// 取消在群
	builder2 := qb.Delete(UserInGroupTableName)
	builder2.Where(qb.Eq("pk"), qb.Eq("uid"), qb.Eq("gid"))

	query2 := builder2.Query(me.session)
	defer query2.Release()
	query2.Consistency(gocql.One)
	batch.Query(query2.Statement(), pk2, uid, gid)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 解散时候，删除所有的成员即可，用户再发送消息出错再处理所在群的表
func (me *Scylla) DissolveGroupAllMember(pk int16, gid int64) error {
	// 构建删除语句
	builder := qb.Delete(GroupMemberTableName)

	builder.Where(qb.Eq("pk"), qb.Eq("gid"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(pk, gid)

	err := query.Exec()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

// 删除用户在某一个表，主动删除，
func (me *Scylla) DeleteUserInG(pk int16, uid, gid int64) error {
	// 构建删除语句
	builder := qb.Delete(UserInGroupTableName)

	builder.Where(qb.Eq("pk"), qb.Eq("uid"), qb.Eq("gid"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(pk, uid, gid)

	err := query.Exec()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}
