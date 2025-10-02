package internal

import (
	"log"

	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"github.com/gorilla/websocket"
)

func sendHelp(conn *websocket.Conn, groupId int) {
	returnMessage, _ := napcat.MarshalGroupImgMsg(groupId, "https://git.mmeiblog.cn/mei/CatBot/raw/branch/main/.gitea/tool.png")
	writeMutex.Lock()
	defer writeMutex.Unlock()
	if err := conn.WriteMessage(websocket.TextMessage, returnMessage); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
}
