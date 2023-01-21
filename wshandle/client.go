package wshandle

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	Connecting   bool
	Connected    bool            // 连接状态
	MsgSendCh    chan string     // 消息发送通道
	MsgReceiveCh chan string     // 消息接收通道
	Done         chan struct{}   // 发送结束通道
	Exit         chan bool       // 程序退出信号
	Interrupt    chan os.Signal  // 终端中止信号
	Conn         *websocket.Conn // 主连接
	ConnMux      sync.Mutex      // 连接互斥锁
}

var wsconn *WsConn

// 读取hostp环境变量
var hostP = GetenvDefault("NEXTTRACE_HOSTPORT", "api.leo.moe")

func (c *WsConn) keepAlive() {
	go func() {
		// 开启一个定时器
		for {
			<-time.After(time.Second * 54)
			if c.Connected {
				c.Conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			}
		}
	}()
	for {
		if !c.Connected && !c.Connecting {
			c.Connecting = true
			c.recreateWsConn()
			// log.Println("WebSocket 连接意外断开，正在尝试重连...")
			// return
		}
		// 降低检测频率，优化 CPU 占用情况
		<-time.After(200 * time.Millisecond)
	}
}

func (c *WsConn) messageReceiveHandler() {
	// defer close(c.Done)
	for {
		if c.Connected {
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				// 读取信息出错，连接已经意外断开
				// log.Println(err)
				c.Connected = false
				return
			}
			if string(msg) != "pong" {
				c.MsgReceiveCh <- string(msg)
			}
		}
	}
}

func (c *WsConn) messageSendHandler() {
	for {
		// 循环监听发送
		select {
		case <-c.Done:
			log.Println("发送协程已经退出")
			return
		case t := <-c.MsgSendCh:
			// log.Println(t)
			if !c.Connected {
				c.MsgReceiveCh <- `{"ip":"` + t + `", "asnumber":"API服务端异常"}`
			} else {
				err := c.Conn.WriteMessage(websocket.TextMessage, []byte(t))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		// 来自终端的中断运行请求
		case <-c.Interrupt:
			// 向 websocket 发起关闭连接任务
			err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				// log.Println("write close:", err)
				os.Exit(1)
			}
			select {
			// 等到了结果，直接退出
			case <-c.Done:
			// 如果等待 1s 还是拿不到结果，不再等待，超时退出
			case <-time.After(time.Second):
			}
			os.Exit(1)
			// return
		}
	}
}

func (c *WsConn) recreateWsConn() {
	u := url.URL{Scheme: "wss", Host: hostP, Path: "/v2/ipGeoWs"}
	// log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c.Conn = ws
	if err != nil {
		log.Println("dial:", err)
		// <-time.After(time.Second * 1)
		c.Connected = false
		c.Connecting = false
		return
	} else {
		c.Connected = true
	}
	c.Connecting = false

	c.Done = make(chan struct{})
	go c.messageReceiveHandler()
}

func createWsConn() *WsConn {
	// 设置终端中断通道
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: hostP, Path: "/v2/ipGeoWs"}
	// log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	wsconn = &WsConn{
		Conn:         c,
		Connected:    true,
		Connecting:   false,
		MsgSendCh:    make(chan string, 10),
		MsgReceiveCh: make(chan string, 10),
	}

	if err != nil {
		log.Println("dial:", err)
		// <-time.After(time.Second * 1)
		wsconn.Connected = false
		wsconn.Done = make(chan struct{})
		go wsconn.keepAlive()
		go wsconn.messageSendHandler()
		return wsconn
	}
	// defer c.Close()
	// 将连接写入WsConn，方便随时可取
	wsconn.Done = make(chan struct{})
	go wsconn.keepAlive()
	go wsconn.messageReceiveHandler()
	go wsconn.messageSendHandler()
	return wsconn
}

func New() *WsConn {
	return createWsConn()
}

func GetWsConn() *WsConn {
	return wsconn
}

func GetenvDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defVal
}
