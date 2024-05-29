package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// 一次性计算file的md5
func CalculateFileMD5Small(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 一次性计算file的md5
func CalculateFileMD5Chunk(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()

	// 逐个添加 chunk 并计算散列值
	const chunkSize = 8192
	n := 0
	buffer := make([]byte, chunkSize)
	for {
		n, err = file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		hash.Write(buffer[:n])
	}

	// 计算最终的 MD5 散列值
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
