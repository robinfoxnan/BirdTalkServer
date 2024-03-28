package db

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
)

// 定义表的元数据
var pchatMetadata = table.Metadata{
	Name:    "pchat",
	Columns: []string{"pk", "uid1", "uid2", "id", "usid", "tm", "tm1", "tm2", "draf", "io", "del", "t"},
	PartKey: []string{"pk"},
	SortKey: []string{"tm", "uid1", "id"},
}

// 创建表对象
//var pchatTable = table.New(pchatMetadata)

// 定义数据结构
type PchatData struct {
	Pk   int       `db:"pk"`
	Uid1 int64     `db:"uid1"`
	Uid2 int64     `db:"uid2"`
	Id   int64     `db:"id"`
	Usid int64     `db:"usid"`
	Tm   time.Time `db:"tm"`
	Tm1  time.Time `db:"tm1"`
	Tm2  time.Time `db:"tm2"`
	Draf string    `db:"draf"`
	Io   bool      `db:"io"`
	Del  bool      `db:"del"`
	T    int       `db:"t"`
}

func PchatDataToSlice(data PchatData) []interface{} {
	return []interface{}{
		data.Pk,
		data.Uid1,
		data.Uid2,
		data.Id,
		data.Usid,
		data.Tm,
		data.Tm1,
		data.Tm2,
		data.Draf,
		data.Io,
		data.Del,
		data.T,
	}
}

func TestDb() {
	// 连接到 ScyllaDB 集群
	cluster := gocql.NewCluster("8.140.203.92:9042")
	cluster.Keyspace = "chatdata"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "Tjj.31415",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("创建会话时发生错误:", err)
		return
	}
	defer session.Close()

	sessionx, err := gocqlx.WrapSession(session, nil)
	if err != nil {
	}
	defer sessionx.Close()

	// 插入数据
	//if err := insertData(&sessionx); err != nil {
	//	fmt.Println("插入数据时发生错误:", err)
	//	return
	//}

	//记录程序启动的时间
	//start := time.Now()
	//if err := insertBatch(&sessionx); err != nil {
	//	fmt.Println("批量插入数据时发生错误:", err)
	//	return
	//}
	//duration := time.Since(start)
	//fmt.Printf("Program execution time: %s\n", duration)

	// 查询数据
	//if err := queryData(&sessionx); err != nil {
	//	fmt.Println("查询数据时发生错误:", err)
	//	return
	//}

	//err = queryDataByPage(&sessionx)
	//if err != nil {
	//	fmt.Println("查询数据时发生错误:", err)
	//}

	//err = queryDataByIdPage(&sessionx)
	//if err != nil {
	//	fmt.Println("查询数据时发生错误:", err)
	//}

	//update(&sessionx)
	batchUpdate(&sessionx)

	err = queryDataBytmPage(&sessionx)
}

func insertData(session *gocqlx.Session) error {
	data := PchatData{
		Pk:   1,
		Uid1: 123456,
		Uid2: 789012,
		Id:   987654,
		Usid: 654321,
		Tm:   time.Now(),
		Tm1:  time.UnixMilli(0),
		Tm2:  time.UnixMilli(0),
		Draf: "你的草稿内容",
		Io:   true,
		Del:  false,
		T:    42,
	}

	// Insert song using query builder.
	insertChat := qb.Insert("chatdata.pchat").Columns(pchatMetadata.Columns...).Query(*session).Consistency(gocql.One)

	insertChat.BindStruct(data)
	if err := insertChat.ExecRelease(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func insertBatch(session *gocqlx.Session) error {

	// 创建 Batch
	batch := session.Session.NewBatch(gocql.LoggedBatch)
	// 创建 Batch
	//batch := gocql.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	var index int64 = 1
	// 构建多个插入语句
	for i := index; i < index+1000; i++ {
		data := PchatData{
			Pk:   1,
			Uid1: 1001,
			Uid2: 1005,
			Id:   i,
			Usid: i,
			Tm:   time.Now(),
			Tm1:  time.UnixMilli(0),
			Tm2:  time.UnixMilli(0),
			Draf: "你的草稿内容",
			Io:   true,
			Del:  false,
			T:    1,
		}

		insertChatQry := qb.Insert("chatdata.pchat").Columns(pchatMetadata.Columns...).Query(*session).Consistency(gocql.One)
		batch.Query(insertChatQry.Statement(),
			PchatDataToSlice(data)...)

	}

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}

func queryData(session *gocqlx.Session) error {

	var dataList []PchatData

	q := qb.Select("chatdata.pchat").Columns(pchatMetadata.Columns...).Query(*session).Consistency(gocql.One)
	if err := q.Select(&dataList); err != nil {
		return err
	}

	//for _, c := range dataList {
	//	fmt.Printf("%+v \n", c)
	//}

	for _, d := range dataList {
		fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d, usid: %d, tm: %v, tm1: %v, tm2: %v, draf: %s, io: %t, del: %t, t: %d\n",
			d.Pk, d.Uid1, d.Uid2, d.Id, d.Usid, d.Tm, d.Tm1, d.Tm2, d.Draf, d.Io, d.Del, d.T)
	}
	return nil
}

func queryDataByPage(session *gocqlx.Session) error {

	var pageSize = 10

	//chatTable := table.New(pchatMetadata)
	builder := qb.Select("chatdata.pchat").Columns(pchatMetadata.Columns...)
	builder.Where(qb.Eq("uid1"))
	builder.AllowFiltering()

	q := builder.Query(*session)
	defer q.Release()
	q.PageSize(pageSize)
	q.Consistency(gocql.One)
	q.Bind(1001)

	getUserChatFunc := func(userID int64, page []byte) (chats []PchatData, nextPage []byte, err error) {
		if len(page) > 0 {
			q.PageState(page)
		}
		iter := q.Iter()
		return chats, iter.PageState(), iter.Select(&chats)
	}

	var (
		dataList []PchatData
		nextPage []byte
		err      error
	)

	for i := 1; ; i++ {
		dataList, nextPage, err = getUserChatFunc(1001, nextPage)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Page %d: \n", i)
		for _, d := range dataList {
			//fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d, usid: %d, tm: %v, tm1: %v, tm2: %v, draf: %s, io: %t, del: %t, t: %d\n",
			//	d.Pk, d.Uid1, d.Uid2, d.Id, d.Usid, d.Tm, d.Tm1, d.Tm2, d.Draf, d.Io, d.Del, d.T)

			fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d \n", d.Pk, d.Uid1, d.Uid2, d.Id)
		}
		if len(nextPage) == 0 {
			break
		}
	}

	return nil
}

func queryDataByIdPage(session *gocqlx.Session) error {

	var pageSize uint = 10

	//chatTable := table.New(pchatMetadata)
	builder := qb.Select("chatdata.pchat").Columns(pchatMetadata.Columns...)
	builder.Where(qb.Eq("uid1"), qb.Gt("id"))

	builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(*session)
	defer q.Release()
	q.Consistency(gocql.One)
	q.Bind(1002, 900)

	var dataList []PchatData

	err := q.Select(&dataList)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("size= %d: \n", len(dataList))
	for _, d := range dataList {
		//fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d, usid: %d, tm: %v, tm1: %v, tm2: %v, draf: %s, io: %t, del: %t, t: %d\n",
		//	d.Pk, d.Uid1, d.Uid2, d.Id, d.Usid, d.Tm, d.Tm1, d.Tm2, d.Draf, d.Io, d.Del, d.T)

		fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d tm: %v \n", d.Pk, d.Uid1, d.Uid2, d.Id, d.Tm)
	}

	return nil
}

func string2time(dateString string) (time.Time, error) {

	// 注意日期格式必须与提供的字符串匹配，否则会出错
	parsedTime, err := time.Parse("2006-01-02 15:04:05", dateString)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return time.Now(), err
	}

	return parsedTime, nil
}

func string2timeLoc(dateString string) (time.Time, error) {
	// 设置东八区（中国标准时间）的地理位置
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("加载地理位置错误:", err)
		return time.Now(), err
	}

	// 使用地理位置信息进行日期解析
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateString, loc)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return time.Now(), err
	}

	return parsedTime, nil
}
func queryDataBytmPage(session *gocqlx.Session) error {

	var pageSize uint = 15

	//chatTable := table.New(pchatMetadata)
	builder := qb.Select("chatdata.pchat").Columns(pchatMetadata.Columns...)
	builder.Where(qb.Eq("pk"), qb.Eq("uid1"), qb.GtOrEq("tm"), qb.LtOrEq("tm"))

	builder.OrderBy("uid1", qb.DESC)
	//builder.OrderBy("tm", qb.DESC)
	//builder.OrderBy("id", qb.DESC)

	builder.AllowFiltering()
	builder.Limit(pageSize)

	q := builder.Query(*session)
	defer q.Release()
	q.Consistency(gocql.One)
	tm1, _ := string2timeLoc("2024-01-27 13:24:00")
	tm2, _ := string2timeLoc("2024-01-27 13:25:56")
	q.Bind(1, 1001, tm1, tm2)

	var dataList []PchatData

	err := q.Select(&dataList)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("size= %d: \n", len(dataList))
	for _, d := range dataList {
		//fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d, usid: %d, tm: %v, tm1: %v, tm2: %v, draf: %s, io: %t, del: %t, t: %d\n",
		//	d.Pk, d.Uid1, d.Uid2, d.Id, d.Usid, d.Tm, d.Tm1, d.Tm2, d.Draf, d.Io, d.Del, d.T)

		//fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d tm: %v \n", d.Pk, d.Uid1, d.Uid2, d.Id, d.Tm)
		fmt.Printf("pk: %d, uid1: %d, uid2: %d, id: %d tm: %d, tm1 %v = %d\n",
			d.Pk, d.Uid1, d.Uid2, d.Id, d.Tm.UnixMilli(), d.Tm1, d.Tm1.UnixMilli())
	}

	return nil
}

func update(session *gocqlx.Session) error {
	builder := qb.Update("chatdata.pchat")

	builder.Set("tm1")

	builder.Where(qb.Eq("pk"))
	builder.Where(qb.Eq("uid1"))
	builder.Where(qb.Eq("tm"))
	builder.Where(qb.Eq("id"))

	query := builder.Query(*session)
	defer query.Release()
	query.Consistency(gocql.One)

	tm := time.UnixMilli(1706333060752)
	query.Bind(time.Now(), 1, 1001, tm, 1000)

	err := query.Exec()
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func batchUpdate(session *gocqlx.Session) error {
	// 创建 Batch
	batch := session.Session.NewBatch(gocql.LoggedBatch)
	// 创建 Batch
	//batch := gocql.NewBatch(gocql.LoggedBatch)
	batch.Cons = gocql.LocalOne

	for i := 0; i < 10; i++ {
		builder := qb.Update("chatdata.pchat")

		builder.Set("tm1")

		builder.Where(qb.Eq("pk"))
		builder.Where(qb.Eq("uid1"))
		builder.Where(qb.Eq("tm"))
		builder.Where(qb.Eq("id"))

		query := builder.Query(*session)
		defer query.Release()
		query.Consistency(gocql.One)

		tm := time.UnixMilli(1706333060752)
		//query.Bind(time.Now(), 1, 1001, tm, 1000-i)
		batch.Query(query.Statement(), time.Now(), 1, 1001, tm, 1000-i)

	}

	if err := session.ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}
