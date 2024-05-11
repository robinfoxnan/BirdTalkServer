package db

// 计算分区的值
func ComputePk(id int64) int16 {
	return int16(id)
}
