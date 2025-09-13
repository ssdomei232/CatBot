package internal

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"git.mmeiblog.cn/mei/aiComplain/pkg/ai"
	"git.mmeiblog.cn/mei/aiComplain/pkg/napcat"
	"git.mmeiblog.cn/mei/aiComplain/tools"
	"github.com/gorilla/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var writeMutex sync.Mutex

func SendGroupMsg(conn *websocket.Conn, messageType int, message []byte) {
	var err error
	var GroupMsg *napcat.Message
	GroupMsg, err = napcat.Parse(message)
	if err != nil {
		return
	}
	if len(GroupMsg.Message) == 0 || GroupMsg.Message[0].Data.Text == "" {
		return
	}

	// 功能部分
	var returnMessage string
	if strings.Contains(GroupMsg.Message[0].Data.Text, "天钿") || strings.Contains(GroupMsg.Message[0].Data.Text, "/chat") {
		log.Println("触发关键词")
		returnMessage, err = ai.SendComplain(GroupMsg.Message[0].Data.Text[5:])
		if err != nil {
			log.Printf("ai处理失败: %v", err)
			return
		}
	} else if strings.Contains(GroupMsg.Message[0].Data.Text, "/ping") {
		ip := GroupMsg.Message[0].Data.Text[6:]
		returnMessage, err = tools.Ping(ip)
		if err != nil {
			log.Println(err)
			returnMessage = fmt.Sprintf("Ping失败:%s", err)
		}
	} else if rand.Float64() < 0.5 {
		log.Println("随机决定不回复此消息")
		return
	} else {
		returnMessage, err = ai.SendComplain(GroupMsg.Message[0].Data.Text)
		if err != nil {
			log.Printf("ai处理失败: %v", err)
			return
		}
	}

	sendMessage, err := napcat.Marshal("send_group_msg", fmt.Sprint(GroupMsg.GroupID), "text", returnMessage)
	if err != nil {
		log.Printf("生成群组消息失败: %v", err)
		return
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()
	if err = conn.WriteMessage(websocket.TextMessage, sendMessage); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
}
