package internal

import (
	"fmt"
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"git.mmeiblog.cn/mei/CatBot/pkg/review"
	"github.com/gorilla/websocket"
)

func ReviewText(conn *websocket.Conn, message string, groupid int, messageId int, qqNumber int) {
	if !isAdminGroup(groupid) && qqNumber != 3979567422 {
		return
	}
	if isBadMessage := review.ReviewText(message); isBadMessage {
		deleteMsg, _ := napcat.MarshalDeleteMessage(messageId)
		banMsg, _ := napcat.MarshalGroupBan(groupid, qqNumber, 60)
		replyMsg, _ := napcat.MarshalAtMsg(groupid, qqNumber, " 您的消息被检测到违规内容，禁言一分钟")

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

func ReviewImage(conn *websocket.Conn, imgUrl string, groupid int, messageId int, qqNumber int) {
	if !isAdminGroup(groupid) && qqNumber != 3979567422 {
		return
	}
	if isBadMessage, err := review.ReviewImage(imgUrl); isBadMessage {
		errMsg := fmt.Sprintf("err:%s", err)
		if err != nil {
			replyMsg, _ := napcat.MarshalAtMsg(groupid, qqNumber, errMsg)
			conn.WriteMessage(websocket.TextMessage, replyMsg)
		}
		deleteMsg, _ := napcat.MarshalDeleteMessage(messageId)
		banMsg, _ := napcat.MarshalGroupBan(groupid, qqNumber, 60)
		replyMsg, _ := napcat.MarshalAtMsg(groupid, qqNumber, " 您发送的图片中含有违规内容，禁言一分钟")
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

func isAdminGroup(groupId int) bool {
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	for _, value := range Config.AdminGroups {
		if value == groupId {
			return true
		}
	}
	return false
}
