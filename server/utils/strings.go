package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

func BytesToHexStr(data []byte) string {
	return hex.EncodeToString(data)
}

func HexStrToBytes(hexStr string) ([]byte, error) {
	// 解码十六进制字符串为字节数组
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAbsolutePath 获取相对路径的绝对路径
func GetAbsolutePath(relativePath string) (string, error) {
	// 检查输入路径是否是绝对路径
	if filepath.IsAbs(relativePath) {
		return relativePath, nil
	}

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 判断操作系统类型
	var absPath string
	if runtime.GOOS == "windows" {
		// 如果是 Windows，使用 filepath.Join 拼接路径
		absPath = filepath.Join(filepath.Dir(exePath), relativePath)
	} else {
		// 如果是 Linux，直接拼接路径
		absPath = filepath.Join(filepath.Dir(exePath), relativePath)
	}

	// 返回绝对路径
	return absPath, nil
}

func BytesToInt64(data []byte) (int64, error) {
	if len(data) < 8 {
		return 0, errors.New("insufficient bytes to convert to int64")
	}
	return int64(binary.LittleEndian.Uint64(data[:8])), nil
}

func Int64ToInt16(i int64) int16 {
	module := math.MaxInt16
	return int16(i % int64(module))
}

func DecodeBase64(str string) ([]byte, error) {
	// 解码base64字符串为字节切片
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	return decodedBytes, err
}

// go:inline
func EncodeBase64(data []byte) string {
	// 解码base64字符串为字节切片
	return base64.StdEncoding.EncodeToString(data)

}

// IsValidEmail 检查电子邮件地址是否合法
func IsValidEmail(email string) bool {
	// 使用正则表达式匹配电子邮件地址的模式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 使用正则表达式判断是否匹配
	return reg.MatchString(email)
}

// 合并表
// mergeMap 合并两个map，并将结果存回第一个map
func MergeMap[K comparable, V any](m1, m2 map[K]V) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func MergeMapMask(m1, m2 map[int64]uint64) {
	for k, v := range m2 {
		v1, ok := m1[k]
		if ok {
			m1[k] = v1 | v
		} else {
			m1[k] = v
		}
	}
}
