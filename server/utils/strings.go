package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"os"
	"path/filepath"
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
