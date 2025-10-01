package internal

import (
	"log"

	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"git.mmeiblog.cn/mei/CatBot/pkg/review"
	"github.com/gorilla/websocket"
)

func ReviewText(conn *websocket.Conn, message string, groupid int, messageId int, qqNumber int) {
	if groupid != 726833553 && groupid != 945592981 {
		return
	}
	if isBadMessage := review.ReviewText(message); isBadMessage {
		deleteMsg, _ := napcat.MarshalDeleteMessage(messageId)
		banMsg, _ := napcat.MarshalGroupBan(groupid, qqNumber, 600)
		replyMsg, _ := napcat.MarshalAtMsg(groupid, qqNumber, " 您的消息被检测到违规内容，禁言十分钟")

		writeMutex.Lock()
		defer writeMutex.Unlock()
		if err := conn.WriteMessage(websocket.TextMessage, deleteMsg); err != nil {
			log.Printf("发送响应失败: %v", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, banMsg); err != nil {
			log.Printf("发送响应失败: %v", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, replyMsg); err != nil {
			log.Printf("发送响应失败: %v", err)
		}
	}
}

// func ReviewImage()
