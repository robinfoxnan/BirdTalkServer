package db

import (
	"birdtalk/server/model"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// 定义关注表，粉丝表的元数据
var metaUserOp = table.Metadata{
	Columns: []string{"pk", "uid1", "uid2", "gid", "id", "usid", "tm", "tm1", "tm2", "io", "st", "cmd", "ret", "mask", "ref", "draf"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "id"},
}

// 插入好友申请,
func (me *Scylla) SaveUserOp(record *model.CommonOpStore, pk2 int16) error {

	// 发出申请的人保存一条，io=out,   收到的人写一条io=in
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	insertQry := qb.Insert(UserOpTableName).Columns(metaUserOp.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertQry.Release()

	// 添加发出申请的人,io=out,  uid1为主角
	batch.Query(insertQry.Statement(), record.Pk, record.Uid1, record.Uid2, record.Gid, record.Id, record.Usid,
		record.Tm, record.Tm1, record.Tm2, model.ChatDataIOOut, record.St, record.Cmd, record.Ret, record.Mask, record.Ref, record.Draf)

	// 收到申请的人，io=in, uid2为主角
	batch.Query(insertQry.Statement(), pk2, record.Uid2, record.Uid1, record.Gid, record.Id, record.Usid,
		record.Tm, record.Tm1, record.Tm2, model.ChatDataIOIn, record.St, record.Cmd, record.Ret, record.Mask, record.Ref, record.Draf)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 设置送达、阅读回执
func (me *Scylla) SetUserOpRecvReply(pk1, pk2, uid1, uid2, msgId, tm1 int64) error {

	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	builder := qb.Update(UserOpTableName)
	builder.Set("tm1")
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id"))

	query := builder.Query(me.session)
	defer query.Release()
	query.Consistency(gocql.One)

	batch.Query(query.Statement(), tm1, pk1, uid1, msgId)
	batch.Query(query.Statement(), tm1, pk2, uid2, msgId)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (me *Scylla) SetUserOpReadReply(pk1, pk2, uid1, uid2, msgId, tm2 int64) error {

	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	builder := qb.Update(UserOpTableName)
	builder.Set("tm2")
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id"))

	query := builder.Query(me.session)
	defer query.Release()
	query.Consistency(gocql.One)

	batch.Query(query.Statement(), tm2, pk1, uid1, msgId)
	batch.Query(query.Statement(), tm2, pk2, uid2, msgId)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (me *Scylla) SetUserOpRecvReadReply(pk1, pk2, uid1, uid2, msgId, tm1, tm2 int64) error {
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	builder := qb.Update(UserOpTableName)
	builder.Set("tm1", "tm2")
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id"))

	query := builder.Query(me.session)
	defer query.Release()
	query.Consistency(gocql.One)

	batch.Query(query.Statement(), tm1, tm2, pk1, uid1, msgId)
	batch.Query(query.Statement(), tm1, tm2, pk2, uid2, msgId)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 设置收方对请求的同意或者拒绝
// const UserOpResultRefuse = 2
// const UserOpResultOk = 1
func (me *Scylla) SetUserOpResult(pk1, pk2, uid1, uid2, msgId int64, result int) error {
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	builder := qb.Update(UserOpTableName)
	builder.Set("ret")
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id"))

	query := builder.Query(me.session)
	defer query.Release()
	query.Consistency(gocql.One)

	batch.Query(query.Statement(), result, pk1, uid1, msgId)
	batch.Query(query.Statement(), result, pk2, uid2, msgId)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 正向查找，如果从头开始查找，那么设置为littleId = 0
func (me *Scylla) FindUserOpForward(pk, uid, littleId int64, pageSize uint) ([]model.CommonOpStore, error) {

	builder := qb.Select(UserOpTableName).Columns(metaUserOp.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.GtOrEq("id"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid, littleId)

	var lst []model.CommonOpStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}
