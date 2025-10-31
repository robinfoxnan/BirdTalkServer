package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 结构体定义
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// NewSnowflake 创建一个新的 Snowflake 实例
func NewSnowflake(workerID, datacenterID int64) *Snowflake {
	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}
}

// GenerateID 生成唯一ID
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTimestamp := s.timeGen()

	if currentTimestamp < s.lastTimestamp {
		delta := currentTimestamp - s.lastTimestamp + 1
		fmt.Printf("Clock moved backwards. wait to %d ms to generate ID.\n", delta)
		time.Sleep(time.Millisecond * time.Duration(delta))
		currentTimestamp = s.timeGen()
		return 0
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & 0xfff
		if s.sequence == 0 {
			currentTimestamp = s.tilNextMillis()
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := ((currentTimestamp - Epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// 从ID中提取时间戳
func SnowIdtoTm(id int64) int64 {
	tm := (id >> timestampLeftShift) + Epoch
	return tm
}

// 将时间戳转为ID类似的结构，可以直接过滤比较，
func TmToSnowIdLike(tm int64) int64 {
	id := (tm - Epoch) << timestampLeftShift
	return id
}

func (s *Snowflake) timeGen() int64 {
	//return time.Now().UnixNano() / int64(time.Millisecond)
	// 在 Go 1.13 版本之后
	return time.Now().UnixMilli()
}

func (s *Snowflake) tilNextMillis() int64 {
	timestamp := s.timeGen()
	for timestamp <= s.lastTimestamp {
		timestamp = s.timeGen()
	}
	return timestamp
}

// Constants
const (
	Epoch              = 1546300800000 // 时间戳起始点，这里设置为2019-01-01 00:00:00 UTC的毫秒数
	workerIDBits       = 5             // 机器ID的位数
	datacenterIDBits   = 5             // 数据中心ID的位数
	maxWorkerID        = -1 ^ (-1 << workerIDBits)
	maxDatacenterID    = -1 ^ (-1 << datacenterIDBits)
	sequenceBits       = 12 // 序列号的位数
	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
)

//根据给出的常量，我们可以进一步解释代码中 ID 的各个部分占用的位数：
//
//时间戳部分：
//时间戳的范围是从 Epoch 开始到当前时间的毫秒数。
//根据 timestampLeftShift 的设置，时间戳部分占用的位数为 sequenceBits + workerIDBits + datacenterIDBits，即 22 位。

//数据中心 ID 部分：
//数据中心 ID 的范围是从 0 到 maxDatacenterID。
//根据 datacenterIDBits 的设置，数据中心 ID 部分占用的位数为 5 位。

//工作节点 ID 部分：
//工作节点 ID 的范围是从 0 到 maxWorkerID。
//根据 workerIDBits 的设置，工作节点 ID 部分占用的位数也是 5 位。
//序列号部分：
//
//序列号的范围是从 0 到 sequenceMask。
//根据 sequenceBits 的设置，序列号部分占用的位数为 12 位。

//func main() {
//	// 创建一个Snowflake实例，传入workerID和datacenterID
//	snowflake := NewSnowflake(1, 1)
//
//	// 生成10个唯一ID并打印
//	for i := 0; i < 10; i++ {
//		id := snowflake.GenerateID()
//		fmt.Println(id)
//	}
//}
/*
请注意，这只是一个简单的实现，真实的生产环境中需要根据实际需求进行调整和扩展。
在实际应用中，确保workerID和datacenterID的唯一性，以及适应系统的时钟回拨等问题，
是使用雪花算法时需要考虑的重要事项。
*/
