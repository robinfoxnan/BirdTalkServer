package main

import (
	"birdtalk/server/core"
	"birdtalk/server/debugapi"
	"birdtalk/server/handler"
	"birdtalk/server/ws"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"sort"
)

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// mapHTMLFiles 将指定目录下的所有HTML文件映射为主文件名的路由
func mapHTMLFiles(router *gin.Engine, dir string) {
	// 获取目录下的所有HTML文件
	files, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 对文件名进行排序
	sort.Strings(files)
	//fmt.Println(files)

	// 遍历文件并映射路由
	for _, file := range files {
		if file == "" {
			continue
		}
		//fmt.Println(file)
		// 在循环内部创建一个局部变量，以确保每个匿名函数引用的是正确的文件
		localFile := file
		// 使用完整文件名（包括扩展名）作为路由
		router.GET("/"+filepath.Base(localFile), func(c *gin.Context) {
			// 渲染相应的HTML页面

			c.HTML(http.StatusOK, filepath.Base(localFile), nil)
		})
	}
}

// https://127.0.0.1:7817/index
func startServer() {
	router := gin.Default()

	// 使用 GinLogger 中间件处理日志记录
	//router.Use(utils.GinLogger(utils.Logger))

	// 使用 GinRecovery 中间件处理恢复
	//router.Use(utils.GinRecovery(utils.Logger, true))

	router.LoadHTMLGlob("page/*.html") // 加载page目录下的所有HTML文件
	router.Static("/js", "./js")       // 设置静态文件目录

	// 自动映射路由
	mapHTMLFiles(router, "page")

	// 添加下载路由

	router.GET("/index", Index)
	router.GET("/", Index)
	router.GET("/ws", ws.HandleWebSocket)
	router.GET("/filestore/:filename", handler.FileDownloadExHandler)
	router.GET("/filestore/download", handler.FileDownloadExHandler1)

	// 调试模式下使用，
	router.GET("/debug", debugapi.DebugHandler)

	fmt.Println("Server is running on port ...")
	//err := router.Run(":80")

	// 启动HTTP/2服务器

	ip := core.Globals.Config.Server.Host
	port := core.Globals.Config.Server.Port
	host := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println(host)
	server := &http.Server{
		Addr:    host, // 0.0.0.0:7817
		Handler: router,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // 不验证客户端证书
		},
	}
	core.Globals.Logger.Info("server started here...")
	//err = server.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem")
	err := server.ListenAndServeTLS(core.Globals.Config.Server.CertFile, core.Globals.Config.Server.KeyFile)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	// load config
	err := core.Globals.LoadConfig("config.yaml")
	if err != nil {

		fmt.Println("load config err!")
		return
	}

	//fmt.Printf("%v", core.Globals.Config)
	core.Globals.InitWithConfig()

	// test
	//city, err := core.Globals.GeoHelper.GetCityByIP("8.140.203.92")
	//fmt.Println(city)

	// init db
	err = core.Globals.InitDb()
	if err != nil {

		fmt.Println("init db err! ", err.Error())
		return
	}

	//core.TestEmailWorkers1()
	startServer()

}
