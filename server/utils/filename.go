package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// 雪花算法一般不会出现小于4字节，雪花算法按照36进制序列化，
// 文件名取前2个为1级目录，那么就是1600个目录，二级也是2个字符，所以2级就是100万个目录；
func fileNameExt2FilePath(base, mainName, extName string, bCreate bool) (string, error) {

	if base == "" {
		// 获取当前工作目录
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return "", errors.New("")
		}
		base = filepath.Join(currentDir, "web/filestore")
	}

	if len(mainName) < 4 {
		newPath := filepath.Join(base, "less4")
		// 创建目录, 下载时候不需要只需要检查是否存在
		if bCreate {
			err := os.MkdirAll(newPath, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directories:", err)
				return "", err
			}
		}
		return filepath.Join(newPath, mainName+extName), nil
	}

	// 获取文件名的前两个字节
	firstTwoBytes := mainName[:2]
	nextTwoBytes := mainName[2:4]

	newPath := filepath.Join(base, firstTwoBytes, nextTwoBytes)
	// 创建目录, 下载时候不需要只需要检查是否存在
	if bCreate {
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directories:", err)
			return "", err
		}
	}

	return filepath.Join(newPath, mainName+extName), nil
}

// 从当前文件名，计算出新文件该放在哪里
func FileName2FilePath(base, fileName string, bCreate bool) (string, error) {

	baseName := filepath.Base(fileName)
	mainName := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	ext := filepath.Ext(baseName)

	return fileNameExt2FilePath(base, mainName, ext, bCreate)
}

// 文件名分析为主文件名和扩展文件名
func DepartFileName(fileName string) (string, string) {
	baseName := filepath.Base(fileName)
	mainName := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	ext := filepath.Ext(baseName)
	return mainName, ext
}
