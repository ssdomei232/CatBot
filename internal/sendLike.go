package internal

import (
	"log"
	"strconv"

	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"github.com/gorilla/websocket"
)

func sendLike(conn *websocket.Conn, groupId int, qqNumber int) {
	likeMsg, _ := napcat.MarshalLikeMsg(qqNumber, 10)
	returnMsg, _ := napcat.Marshal("send_group_msg", strconv.Itoa(groupId), "text", "赞了你10下")
	writeMutex.Lock()
	defer writeMutex.Unlock()
	if err := conn.WriteMessage(websocket.TextMessage, likeMsg); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, returnMsg); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
}
