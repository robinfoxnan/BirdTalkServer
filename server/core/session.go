package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

const sendQueueLimit = 128

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = pongWait / 2
)

// Session 表示单个的 WebSocket 连接。一个用户可能拥有多个会话，因为每个允许多设备登录并同步。
type Session struct {
	UserID      int64
	Sid         int64  // 会话 ID
	DeviceID    string // 客户端的设备 ID
	Platf       string // 平台：web、ios、android
	Lang        string // 客户端的语言
	CountryCode string // 客户端的国家代码
	CodeType    string // 编码支持json和protobuf
	RemoteAddr  string // 客户端的 IP 地址。

	// 客户端的协议版本：((major & 0xff) << 8) | (minor & 0xff)。
	Ver int

	// 会话起源的集群节点的引用。仅用于集群 RPC 会话。
	//clnode *ClusterNode

	// 多路复用会话的引用。仅用于代理会话。
	//multi        *Session

	// 会话接收到来自客户端的任何数据包的时间
	lastAction int64
	// WebSocket。仅用于 WebSocket 会话。
	ws *websocket.Conn

	// 输出消息，缓冲。
	// 内容必须以适合会话的格式进行序列化。
	send        chan any
	stop        chan any //// 用于关闭会话的通道，缓冲 1。
	terminating int32
}

// 创建新的会话
func NewSession(conn *websocket.Conn, sid, uid int64, code string) *Session {
	s := Session{UserID: uid, Sid: sid, ws: conn, CodeType: code}
	if sid == 0 {
		s.Sid = Globals.snow.GenerateID()
		i := 0
		for i = 0; i < 5; i++ { // 这里要防止雪花算法出问题，
			if Globals.ss.Has(s.Sid) {
				s.Sid = Globals.snow.GenerateID()
			} else {
				Globals.ss.UpSet(s.Sid, &s)
				break
			}
		}
		if i == 5 {
			return nil
		}
	}

	s.send = make(chan any, sendQueueLimit+32) // buffered
	s.stop = make(chan any, 1)                 // Buffered by 1 just to make it non-blocking

	atomic.StoreInt32(&s.terminating, 0)
	s.lastAction = time.Now().UnixMilli()
	return &s
}

func wsWrite(ws *websocket.Conn, mt int, msg any) error {
	var bits []byte
	if msg != nil {
		switch msg.(type) {
		case []byte:
			bits = msg.([]byte)
		case string:
			bits = []byte(msg.(string))
		default:
			bits = []byte{}
		}

	} else {
		bits = []byte{}
	}
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(mt, bits)
}

func (sess *Session) SendMessage(msg any) bool {
	if len(sess.send) > sendQueueLimit {
		//logs.Err.Println("ws: outbound queue limit exceeded", sess.sid)
		return false
	}

	//statsInc("OutgoingMessagesWebsockTotal", 1)
	//if err := wsWrite(sess.ws, websocket.TextMessage, msg); err != nil {
	//	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
	//		websocket.CloseNormalClosure) {
	//		logs.Err.Println("ws: writeLoop", sess.sid, err)
	//	}
	//	return false
	//}
	sess.send <- msg
	return true
}

// 写循环
func (sess *Session) WriteLoop() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		// Break readLoop.
		sess.ws.Close()
		sess.cleanUp()
		fmt.Println("write loop end here")
	}()

	for {
		select {
		case msg, ok := <-sess.send:
			if !ok {
				// Channel closed.
				return
			}
			wsWrite(sess.ws, websocket.TextMessage, msg)
			// 这里根据消息类型处理发送

		//case <-sess.bkgTimer.C:
		//	if sess.background {
		//		sess.background = false
		//		sess.onBackgroundTimer()
		//	}

		case msg := <-sess.stop:
			// Shutdown requested, don't care if the message is delivered
			if msg != nil {
				wsWrite(sess.ws, websocket.TextMessage, msg)
			}
			return

		case <-ticker.C:
			if err := wsWrite(sess.ws, websocket.PingMessage, nil); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
					websocket.CloseNormalClosure) {
					//logs.Err.Println("ws: writeLoop ping", sess.sid, err)
				}
				return
			}
		}
	}
}

func (s *Session) StopSession(data any) {
	select {
	case s.stop <- data:
		// 向通道写入数据成功
	default:
		// 通道已关闭
		fmt.Println("Channel is closed")
	}
	//s.maybeScheduleClusterWriteLoop()
}

// 读循环，
func (sess *Session) ReadLoop() {
	defer func() {
		sess.StopSession("stop")
		sess.ws.Close()
		fmt.Println("read loop end here")
	}()

	sess.ws.SetReadLimit(Globals.maxMessageSize)
	sess.ws.SetReadDeadline(time.Now().Add(pongWait))
	sess.ws.SetPongHandler(func(string) error {
		sess.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Read a ClientComMessage
		t, raw, err := sess.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				fmt.Println("ws: readLoop", sess.Sid, err)
			}
			fmt.Printf("ws: readLoop err, sid=%v, err=%v \n", sess.Sid, err)
			return
		}

		if t == websocket.CloseMessage {
			return
		} else {
			//statsInc("IncomingMessagesWebsockTotal", 1)
			sess.dispatchRaw(t, raw)
		}
	}
}

// 读循环中分发消息
func (s *Session) dispatchRaw(messageType int, msg []byte) {
	switch messageType {
	case websocket.TextMessage:
		fmt.Println("recv text message from sid=", s.Sid, string(msg))
		str := "recv msg:" + string(msg)
		s.SendMessage([]byte(str))
	case websocket.BinaryMessage:
		fmt.Println("recv bin message from sid=", s.Sid)
		encoder := BinEncoder{}
		msg, err := encoder.DecodeMsg(msg)
		if err == nil {
			fmt.Println(msg)
		}

	case websocket.PingMessage:
		fmt.Println("recv ping message from sid=", s.Sid)
	case websocket.PongMessage:
		fmt.Println("recv pong message from sid=", s.Sid)
	}
}

// 清理资源
/*
当从缓存中删除了指向对象的指针时，如果该对象没有其他指针引用，
那么该对象就会被 Go 的垃圾回收机制回收，这意味着对象所包含的所有资源，包括管道，都会被释放。
在 cleanUp 函数中，通过 Globals.ss.Remove(s.Sid) 删除了指向 Session 对象的指针，
这意味着该 Session 对象可能会被垃圾回收器回收。如果该 Session 对象不再被任何其他地方引用，
那么该对象所持有的资源，包括管道等，都会被释放。

在 Go 中，session.ReadLoop() 调用并不是一个闭包，而是一个普通的方法调用。
在 Go 中，方法调用会将对象作为接收者传递给方法，但并不创建闭包。

在 session.ReadLoop() 方法内部，虽然会引用 Session 对象的指针，但并不会导致对象的引用计数增加。
在 Go 中，不像其它语言（比如 Python），没有显式的引用计数机制。
Go 的垃圾回收器通过遍历可达对象图来确定对象的可达性，并对不可达的对象进行回收，而不是简单地根据引用计数来判断对象是否可以回收。
因此，当 session.ReadLoop() 方法执行完毕后，如果该 Session 对象没有其他地方引用，它将会被垃圾回收器回收。

这里保存一个指针，是用于在外部停止这个会话，比如服务需要优雅的退出。
*/
func (s *Session) cleanUp() {
	Globals.ss.Remove(s.Sid)
}
