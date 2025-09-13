package napcat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// MessageHandler 定义消息处理函数签名
type MessageHandler func(conn *websocket.Conn, messageType int, message []byte)

// Client WebSocket客户端配置
type Client struct {
	url        string
	retryDelay time.Duration
}

// New 创建客户端实例
func New(url string, handler MessageHandler, opts ...Option) *Client {
	client := &Client{
		url:        url,
		retryDelay: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Start 启动客户端
func (c *Client) Start(handler MessageHandler) {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			log.Printf("连接失败: %v，%v后重试...", err, c.retryDelay)
			time.Sleep(c.retryDelay)
			continue
		}

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("连接断开: %v，尝试重连...", err)
				conn.Close()
				break
			}

			// 使用goroutine处理消息
			go func() {
				if !isHeartbeat(message) {
					handler(conn, messageType, message)
				}
			}()
		}
	}
}

// isHeartbeat 判断是否为心跳包
func isHeartbeat(message []byte) bool {
	var msgMap map[string]any
	if err := json.Unmarshal(message, &msgMap); err != nil {
		return false
	}

	if metaType, ok := msgMap["meta_event_type"].(string); ok {
		return metaType == "heartbeat"
	}
	return false
}

// Option 自定义客户端配置
type Option func(*Client)

// WithRetryDelay 设置重试间隔
func WithRetryDelay(delay time.Duration) Option {
	return func(c *Client) {
		c.retryDelay = delay
	}
}
