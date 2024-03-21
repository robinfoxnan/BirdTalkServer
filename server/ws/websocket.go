package ws

import (
	"birdtalk/server/core"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

func HandleWebSocket(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//fmt.Println("upgrade protocal to websocket", r.Header["User-Agent"])
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	codetype := c.Query("code")
	fmt.Println("code:", codetype)

	sess := core.NewSession(conn, 0, 0, codetype)
	if sess == nil {
		// 如果生成的SID有问题，那么属于重大错误，防止覆盖应该重来，不然消息会发错用户
		conn.Close()
		return
	}

	sess.RemoteAddr = c.ClientIP()

	userAgent := c.Request.UserAgent()
	var platform string
	if strings.Contains(userAgent, "Windows") {
		platform = "Windows"
	} else if strings.Contains(userAgent, "Android") {
		platform = "Android"
	} else {
		platform = "Unknown"
	}
	sess.Platf = platform

	// 启动2个协程处理读和写；
	go sess.WriteLoop()
	go sess.ReadLoop()

}
