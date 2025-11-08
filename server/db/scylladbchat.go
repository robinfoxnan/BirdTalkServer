package db

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

// 定义表的元数据
var metaPrivateChatData = table.Metadata{
	Name:    "pchat",
	Columns: []string{"pk", "uid1", "uid2", "id", "usid", "tm", "tm1", "tm2", "io", "st", "ct", "mt", "pr", "ref", "draf"},
	PartKey: []string{"pk"},
	SortKey: []string{"uid1", "id"},
}

var metaGroupChatData = table.Metadata{
	Name:    "gchat",
	Columns: []string{"pk", "gid", "uid", "id", "usid", "tm", "res", "st", "ct", "mt", "pr", "ref", "draf"},
	PartKey: []string{"pk"},
	SortKey: []string{"gid", "id"},
}

// 写2次，首先是发方A，然后是收方B
func (me *Scylla) SavePChatData(msg *model.PChatDataStore, pk2 int16) error {
	// 同时加入粉丝表
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 发方的IO = 0:OUT
	insertFirst := qb.Insert(PrivateChatTableName).Columns(metaPrivateChatData.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertFirst.Release()
	batch.Query(insertFirst.Statement(), msg.Pk, msg.Uid1, msg.Uid2,
		msg.Id, msg.Usid, msg.Tm, msg.Tm1, msg.Tm2,
		model.ChatDataIOOut, msg.St, msg.Ct, msg.Mt,
		msg.Print, msg.Ref, msg.Draf)

	// 收方的IO = 1:IN
	insertSecond := qb.Insert(PrivateChatTableName).Columns(metaPrivateChatData.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertSecond.Release()
	batch.Query(insertSecond.Statement(), pk2, msg.Uid2, msg.Uid1,
		msg.Id, msg.Usid, msg.Tm, msg.Tm1, msg.Tm2,
		model.ChatDataIOIn, msg.St, msg.Ct, msg.Mt,
		msg.Print, msg.Ref, msg.Draf)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

func (me *Scylla) SavePChatSelfData(msg *model.PChatDataStore) error {
	// 同时加入粉丝表
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 收方
	insertFirst := qb.Insert(PrivateChatTableName).Columns(metaPrivateChatData.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertFirst.Release()
	batch.Query(insertFirst.Statement(), msg.Pk, msg.Uid1, msg.Uid2,
		msg.Id, msg.Usid, msg.Tm, msg.Tm1, msg.Tm2,
		model.ChatDataIOIn, msg.St, msg.Ct, msg.Mt,
		msg.Print, msg.Ref, msg.Draf)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 系统发给A的系统通知，A作为收方，直接写一次数据库， uid2这里是系统账户号码，默认你1000以下的都是
func (me *Scylla) SavePChatDataSystem(msg *model.PChatDataStore) error {
	// 同时加入粉丝表
	// 创建 Batch
	batch := me.session.Session.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	// 收方
	insertFirst := qb.Insert(PrivateChatTableName).Columns(metaPrivateChatData.Columns...).Query(me.session).Consistency(gocql.One)
	defer insertFirst.Release()
	batch.Query(insertFirst.Statement(), msg.Pk, msg.Uid1, msg.Uid2,
		msg.Id, msg.Usid, msg.Tm, msg.Tm1, msg.Tm2,
		model.ChatDataIOIn, msg.St, msg.Ct, msg.Mt,
		msg.Print, msg.Ref, msg.Draf)

	if err := me.session.ExecuteBatch(batch); err != nil {
		return err
	}
	return nil
}

// 对发送方设置回执，收方不需要设置
func (me *Scylla) SetPChatRecvReply(pk1, pk2 int16, uid1, uid2, msgId, tm1 int64) error {
	builder := qb.Update(PrivateChatTableName)

	builder.Set("tm1")

	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id")).Existing()

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(tm1, pk1, uid1, msgId)

	// 执行查询并检查是否应用更新
	applied, err := query.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}

	return nil

	//err := query.Exec()
	//return err
}

func (me *Scylla) SetPChatReadReply(pk1, pk2 int16, uid1, uid2, msgId, tm2 int64) error {
	builder := qb.Update(PrivateChatTableName)

	builder.Set("tm2")

	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id")).Existing()

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(tm2, pk1, uid1, msgId)

	// 执行查询并检查是否应用更新
	applied, err := query.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}

	return nil
}

func (me *Scylla) SetPChatRecvReadReply(pk1, pk2 int16, uid1, uid2, msgId, tm1, tm2 int64) error {
	builder := qb.Update(PrivateChatTableName)

	builder.Set("tm1", "tm2")

	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id")).Existing()

	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(tm1, tm2, pk1, uid1, msgId)

	// 执行查询并检查是否应用更新
	applied, err := query.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}

	return nil
}

// 设置删除，不可逆
func (me *Scylla) SetPChatMsgDeleted(pk1, pk2 int16, uid1, uid2, msgId int64) error {
	// Update 发方的 DrafStateDel
	builder1 := qb.Update(PrivateChatTableName).
		Set("st").
		Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("id")).
		Existing()

	query1 := builder1.Query(me.session)
	defer query1.Release()

	query1.Consistency(gocql.One)
	query1.Bind(pbmodel.ChatMsgStatus_DELETED, pk1, uid1, msgId)

	// 执行查询并检查是否应用更新
	applied, err := query1.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}
	////////////////////////////////////////////////////////////
	query2 := builder1.Query(me.session)
	defer query2.Release()

	query2.Consistency(gocql.One)
	query2.Bind(pbmodel.ChatMsgStatus_DELETED, pk2, uid2, msgId)

	// 执行查询并检查是否应用更新
	applied, err = query2.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}

	return nil

}

// 正向查找，如果从头开始查找，那么设置为littleId = 0
func (me *Scylla) FindPChatMsgForward(pk int16, uid, littleId int64, pageSize uint) ([]model.PChatDataStore, error) {

	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.GtOrEq("id"))

	// 如果你只排序 id，会报错或者被忽略，因为排序的逻辑是按聚簇键顺序 (uid1, id)
	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)
	q.Bind(pk, uid, littleId)

	var lst []model.PChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

func (me *Scylla) FindPChatMsgForwardTopic(pk int16, uid, fid, littleId int64, pageSize uint) ([]model.PChatDataStore, error) {

	builder := qb.Select(PrivateChatViewName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"), qb.GtOrEq("id"))
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid, fid, littleId)

	var lst []model.PChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 正序查找，设置边界范围
func (me *Scylla) FindPChatMsgForwardBetween(pk int16, uid, littleId, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {

	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.GtOrEq("id"), qb.LtOrEq("id"))

	builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)
	q.Bind(pk, uid, littleId, bigId)

	var lst []model.PChatDataStore
	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

func (me *Scylla) FindPChatMsgForwardBetweenTopic(pk int16, uid, fid, littleId, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {

	builder := qb.Select(PrivateChatViewName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"), qb.GtOrEq("id"), qb.LtOrEq("id"))

	//builder.OrderBy("uid1", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)
	q.Bind(pk, uid, fid, littleId, bigId)

	var lst []model.PChatDataStore
	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 从最新的数据向前倒序查若干条
func (me *Scylla) FindPChatMsgBackward(pk int16, uid, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {
	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.LtOrEq("id"))

	builder.OrderBy("uid1", qb.DESC)
	builder.OrderBy("id", qb.DESC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid, bigId)
	var lst []model.PChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

func (me *Scylla) FindPChatMsgBackwardTopic(pk int16, uid, fid, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {
	builder := qb.Select(PrivateChatViewName).Columns(metaPrivateChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.Eq("uid2"), qb.LtOrEq("id"))

	builder.OrderBy("id", qb.DESC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, uid, fid, bigId)
	var lst []model.PChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 从某一点开始向之前的历史数据反向查找,即 所有小于bigId 的
//func (me *Scylla) FindPChatMsgBackwardFrom(pk int16, uid, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {
//	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
//	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.LtOrEq("id"))
//
//	builder.OrderBy("uid1", qb.DESC)
//	builder.OrderBy("id", qb.DESC)
//
//	//builder.AllowFiltering()
//	builder.Limit(pageSize)
//
//	q := builder.Query(me.session)
//	defer q.Release()
//
//	q.Consistency(gocql.One)
//
//	q.Bind(pk, uid, bigId)
//
//	var lst []model.PChatDataStore
//
//	err := q.Select(&lst)
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	return lst, nil
//}

// 从当前最新开始向之前的历史数据反向查找，即 所有大于littlId 的
//func (me *Scylla) FindPChatMsgBackwardTo(pk, uid, littleId int64, pageSize uint) ([]model.PChatDataStore, error) {
//	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
//	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.GtOrEq("id"))
//
//	builder.OrderBy("uid1", qb.DESC)
//	builder.OrderBy("id", qb.DESC)
//
//	//builder.AllowFiltering()
//	builder.Limit(pageSize)
//
//	q := builder.Query(me.session)
//	defer q.Release()
//
//	q.Consistency(gocql.One)
//
//	q.Bind(pk, uid, littleId)
//
//	var lst []model.PChatDataStore
//
//	err := q.Select(&lst)
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	return lst, nil
//}

// 向之前的历史数据反向查找
//func (me *Scylla) FindPChatMsgBackwardBetween(pk int16, uid, littleId, bigId int64, pageSize uint) ([]model.PChatDataStore, error) {
//	builder := qb.Select(PrivateChatTableName).Columns(metaPrivateChatData.Columns...)
//	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.LtOrEq("id"), qb.GtOrEq("id"))
//
//	builder.OrderBy("uid1", qb.DESC)
//	builder.OrderBy("id", qb.DESC)
//
//	//builder.AllowFiltering()
//	builder.Limit(pageSize)
//
//	q := builder.Query(me.session)
//	defer q.Release()
//
//	q.Consistency(gocql.One)
//
//	q.Bind(pk, uid, bigId, littleId)
//
//	var lst []model.PChatDataStore
//
//	err := q.Select(&lst)
//	if err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//	return lst, nil
//}

// //////////////////////////////////////////////////////////////////////////////////////////
// 写1次，
func (me *Scylla) SaveGChatData(msg *model.GChatDataStore) error {
	insertChat := qb.Insert(GroupChatTableName).Columns(metaGroupChatData.Columns...).Query(me.session).Consistency(gocql.One)
	insertChat.BindStruct(msg)
	if err := insertChat.ExecRelease(); err != nil {
		//fmt.Println(err)
		return err
	}
	return nil
}

// 设置删除，不可逆
func (me *Scylla) SetGChatMsgDeleted(pk int16, gid, msgId int64) error {

	builder := qb.Update(GroupChatTableName)
	builder.Set("st").Where(qb.Eq("pk"), qb.Eq("gid"), qb.Eq("id")).Existing()
	query := builder.Query(me.session)
	defer query.Release()

	query.Consistency(gocql.One)
	query.Bind(pbmodel.ChatMsgStatus_DELETED, pk, gid, msgId)
	applied, err := query.ExecCAS()
	if err != nil {
		return fmt.Errorf("error executing update: %w", err)
	}

	// Check if the update was applied, if not, return an error
	if !applied {
		return errors.New("primary key not found, update not applied")
	}
	return nil
}

// 正向查找
func (me *Scylla) FindGChatMsgForward(pk int16, gid, littleId int64, pageSize uint) ([]model.GChatDataStore, error) {
	builder := qb.Select(GroupChatTableName).Columns(metaGroupChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.GtOrEq("id"))

	builder.OrderBy("gid", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, littleId)

	var lst []model.GChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 正向查找，设置某个边界范围内
func (me *Scylla) FindGChatMsgForwardBetween(pk int16, gid, littleId, bigId int64, pageSize uint) ([]model.GChatDataStore, error) {
	builder := qb.Select(GroupChatTableName).Columns(metaGroupChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.GtOrEq("id"), qb.LtOrEq("id"))

	builder.OrderBy("gid", qb.ASC)
	builder.OrderBy("id", qb.ASC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, littleId, bigId)

	var lst []model.GChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 倒序，反向历史数据方向查找
func (me *Scylla) FindGChatMsgBackwardFrom(pk int16, gid, bigId int64, pageSize uint) ([]model.GChatDataStore, error) {
	builder := qb.Select(GroupChatTableName).Columns(metaGroupChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.LtOrEq("id"))

	builder.OrderBy("gid", qb.DESC)
	builder.OrderBy("id", qb.DESC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, bigId)

	var lst []model.GChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 倒序，反向历史数据方向查找，从最新的数据开始向前加载
func (me *Scylla) FindGChatMsgBackwardTo(pk, gid, littleId int64, pageSize uint) ([]model.GChatDataStore, error) {
	builder := qb.Select(GroupChatTableName).Columns(metaGroupChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.GtOrEq("id"))

	builder.OrderBy("gid", qb.DESC)
	builder.OrderBy("id", qb.DESC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, littleId)

	var lst []model.GChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// 倒序，从bigId 向littleId方向去查找，限定一定的个数，如果无法覆盖边界，再来一次
func (me *Scylla) FindGChatMsgBackwardBetween(pk, gid, littleId, bigId int64, pageSize uint) ([]model.GChatDataStore, error) {
	builder := qb.Select(GroupChatTableName).Columns(metaGroupChatData.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("gid"), qb.LtOrEq("id"), qb.GtOrEq("id"))

	builder.OrderBy("gid", qb.DESC)
	builder.OrderBy("id", qb.DESC)

	//builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(me.session)
	defer q.Release()

	q.Consistency(gocql.One)

	q.Bind(pk, gid, bigId, littleId)

	var lst []model.GChatDataStore

	err := q.Select(&lst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return lst, nil
}

// todo: 是否需要添加批量写入多条消息，暂时不做，因为得知写入出错的条目，就需要逐条处理；
// 在集群模式下可以尝试，从消息队列读取后批量处理；
