package handler

import (
	"birdtalk/server/core"
	"birdtalk/server/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 扩展下载
// 这个是根据文件名来直接下载，
func FileDownloadExHandler(c *gin.Context) {

	//设置默认值
	//user := c.DefaultQuery("username", "test")
	//sid := c.Query("sid")
	//fmt.Println("sid = ", sid)

	filename := c.Param("filename")

	baseDir := core.Globals.Config.Server.FileBasePath
	filePath, err := utils.FileName2FilePath(baseDir, filename, false)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// 检查文件是否存在
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// 提供文件下载
	c.File(filePath)
}

// Content-Disposition: attachment; filename="test.pdf"
// Content-Disposition: attachment; filename*=UTF-8”%E6%B5%8B%E8%AF%95.pdf
func getFileName(c *gin.Context) (string, error) {
	cd := c.GetHeader("Content-Disposition")
	if cd == "" {
		c.String(http.StatusBadRequest, "Content-Disposition not found")
		return "", errors.New("Content-Disposition not found")
	}

	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		c.String(http.StatusBadRequest, "parse Content-Disposition failed")
		return "", errors.New("Content-Disposition not found")
	}

	var filename string = ""

	// 优先 RFC5987 filename*
	if v, ok := params["filename*"]; ok {
		// UTF-8''xxxx
		if strings.HasPrefix(v, "UTF-8''") {
			filename, _ = url.QueryUnescape(strings.TrimPrefix(v, "UTF-8''"))
		} else {
			filename = v
		}
	} else if v, ok := params["filename"]; ok {
		filename = v
	}
	return filename, nil
}

// /download?filename=test.pdf
//filename := c.Query("filename")
//if filename == "" {
//	// filename: 测试文件.zip
//	filename = c.GetHeader("filename")
//
//}
//if filename == "" {
//	filename, _ = getFileName(c)
//}
//
//if filename == "" {
//
//}

// https://127.0.0.1:7817/filestore/download?u=123345664555&sign=
func FileDownloadExHandler1(c *gin.Context) {
	str := c.Query("u")
	keyPrint, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		c.String(http.StatusNotFound, "user key print error")
		return
	}

	_, keyEx, err := core.LoadUserByKeyPrint(keyPrint)
	if err != nil || keyEx == nil {
		core.Globals.Logger.Info("try download by key print, but not found: ", zap.Int64("keyPrint", keyPrint))
		c.String(http.StatusNotFound, "user key print error")
		return
	}

	token := c.Query("sign")
	// 服务器密钥
	//expireTm, filename, err := core.DecryptFileToken(core.Globals.Config.Server.TokenKey, token)

	// 使用用户的密钥
	expireTm, filename, err := core.DecryptFileToken(keyEx.SharedKey, token)
	if err != nil {
		c.String(http.StatusNotFound, "decrypt file token failed")
		//print(err.Error())
		return
	}

	// 校验时间戳是否过期
	if time.Now().Unix() > expireTm {
		// t := time.UnixMilli(expireTm)
		c.String(http.StatusNotFound, "sign is expired")

		t := time.Unix(expireTm, 0) // 秒级时间戳转换
		core.Globals.Logger.Debug("download?", zap.String("filename", filename),
			zap.String("expire", t.Format("2006-01-02 15:04:05")))
		return
	}

	if filename == "" {
		c.String(400, "filename not found")
		return
	}

	baseDir := core.Globals.Config.Server.FileBasePath
	filePath, err := utils.FileName2FilePath(baseDir, filename, false)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}

	// 检查文件是否存在
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	core.Globals.Logger.Debug("download?", zap.String("filename", filename))

	// 提供文件下载
	c.File(filePath)
}
