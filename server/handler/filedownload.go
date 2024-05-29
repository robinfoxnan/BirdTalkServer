package handler

import (
	"birdtalk/server/core"
	"birdtalk/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// 扩展下载
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
