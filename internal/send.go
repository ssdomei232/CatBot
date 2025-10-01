package internal

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"git.mmeiblog.cn/mei/CatBot/pkg/ai"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"git.mmeiblog.cn/mei/CatBot/pkg/review"
	"git.mmeiblog.cn/mei/CatBot/tools"
	"github.com/gorilla/websocket"
)

var writeMutex sync.Mutex

func SendGroupMsg(conn *websocket.Conn, messageType int, message []byte) {
	var err error
	var GroupMsg *napcat.Message
	GroupMsg, err = napcat.Parse(message)
	if err != nil {
		return
	}

	// 解析消息项
	messageItems, err := GroupMsg.GetMessageItems()
	if err != nil || len(messageItems) == 0 {
		return
	}

	// 查找第一个文本消息作为命令输入
	var commandText string
	for _, item := range messageItems {
		if textData, ok := item.Data.(napcat.TextData); ok && commandText == "" {
			commandText = textData.Text
		}
		if imgData, ok := item.Data.(napcat.ImageData); ok {
			review.CacheImg(imgData.URL)
			ReviewImage(conn, imgData.URL, GroupMsg.GroupID, GroupMsg.MessageID, GroupMsg.UserID)
		}
	}

	if commandText == "" {
		return
	}

	// 每次消息都需要执行的部分
	Record(*GroupMsg)
	ReviewText(conn, commandText, GroupMsg.GroupID, GroupMsg.MessageID, GroupMsg.UserID)

	// 功能部分
	var returnMessage string
	if strings.Contains(commandText, ".chat") {
		log.Println("触发关键词")
		returnMessage, err = ai.SendComplain(commandText[5:]) // 去掉".chat"前缀
		if err != nil {
			log.Printf("ai处理失败: %v", err)
			return
		}
		if review.ReviewText(returnMessage) {
			return
		}
	} else if strings.Contains(commandText, ".ping") {
		ip := commandText[6:] // 去掉".ping "前缀
		returnMessage, err = tools.Ping(ip)
		if err != nil {
			log.Println(err)
			returnMessage = fmt.Sprintf("Ping失败:%s", err)
		}
	} else if strings.Contains(commandText, ".nc") {
		returnMessage = "你猜"
	} else if strings.Contains(commandText, "xmsl") {
		returnMessage = "羡慕死了"
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
