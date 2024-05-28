package db

import (
	"birdtalk/server/model"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

/////////////////////////////////////////////////////////////
// 群用户的操作记录

// 定义关注表，粉丝表的元数据
var metaGroupOp = table.Metadata{
	Columns: []string{"pk", "gid", "id", "uid1", "uid2", "usid", "tm", "tm1", "tm2", "io", "st", "cmd", "ret", "mask", "ref", "draf"},
	PartKey: []string{"pk"},
	SortKey: []string{"gid", "id"},
}

// 插入群组操作的记录
func (me *Scylla) SaveGroupOp(record *model.CommonOpStore) error {

	// uid1为申请人，uid2是附加的操作记录内容，比如uid2是被删除的，被设置为管理员的，同意通过申请的；
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	insertQry := qb.Insert(GroupOpTableName).Columns(metaGroupOp.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertQry.Release()

	// 添加发出申请的人,io=out,  uid1为主角
	batch.Query(insertQry.Statement(), record.Pk, record.Gid, record.Id, record.Uid1, record.Uid2, record.Usid,
		record.Tm, record.Tm1, record.Tm2, model.ChatDataIOOut, record.St, record.Cmd, record.Ret, record.Mask, record.Ref, record.Draf)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 更新操作的结果，管理员的应答
// 需要注意的是，如果没有找到主键，会插入一条新的记录
func (me *Scylla) SetGroupOpResult(pk int16, gid, logId int64, adminId int64, ok bool) error {

	// uid1为申请人，uid2是附加的操作记录内容，比如uid2是被删除的，被设置为管理员的，同意通过申请的；
	// 创建 Batch

	builder := qb.Update(GroupOpTableName)

	builder.Set("uid2", "ret")
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.Eq("id"))

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)

	ret := model.UserOpResultRefuse
	if ok {
		ret = model.UserOpResultOk
	}
	query.Bind(adminId, ret, pk, gid, logId)

	err := query.Exec()
	return err
}

// 精确的查找一条记录
func (me *Scylla) FindGroupOpExact(pk int16, gid, logId int64) (*model.CommonOpStore, error) {
	builder := qb.Select(GroupOpTableName).Columns(metaGroupOp.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.Eq("id"))

	builder.OrderBy("gid", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, logId)

	var lst []model.CommonOpStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if lst == nil || len(lst) < 1 {
		return nil, errors.New("not find group op record")
	}

	return &lst[0], nil
}

// 正向查找，如果从头开始查找，那么设置为littleId = 0
func (me *Scylla) FindGroupOpForward(pk int16, gid, littleId int64, pageSize uint) ([]model.CommonOpStore, error) {

	builder := qb.Select(GroupOpTableName).Columns(metaGroupOp.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.GtOrEq("id"))

	builder.OrderBy("gid", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, littleId)

	var lst []model.CommonOpStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}
