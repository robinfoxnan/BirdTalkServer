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
		fmt.Println("Clock moved backwards. Refusing to generate ID.")
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

func (s *Snowflake) timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
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
