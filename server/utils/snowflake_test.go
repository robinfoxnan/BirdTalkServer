package utils

import (
	"fmt"
	"testing"
)

func TestSnow(t *testing.T) {
	snowflake := NewSnowflake(1, 1)

	// 生成10个唯一ID并打印
	for i := 0; i < 10; i++ {
		id := snowflake.GenerateID()
		tm := snowflake.Id2Tm(id)
		str := UtcTm2LocalString(tm)
		fmt.Println(id, str)
	}
}
