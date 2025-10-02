package internal

import (
	"log"
	"strconv"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"git.mmeiblog.cn/mei/CatBot/tools"
	"github.com/gorilla/websocket"
)

func findfood(conn *websocket.Conn, poi string, groupId int) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatal("获取配置文件错误: ", err)
	}
	foodService := tools.NewFindFoodService(config.GDKey)
	info, photoURL, err := foodService.SearchAndFormat(poi)
	if err != nil {
		log.Printf("搜索失败: %v\n", err)
		return
	}
	sendMessage, _ := napcat.Marshal("send_group_msg", strconv.Itoa(groupId), "text", info)
	sendImgMsg, err := napcat.MarshalGroupImgMsg(groupId, photoURL)

	writeMutex.Lock()
	defer writeMutex.Unlock()
	if err = conn.WriteMessage(websocket.TextMessage, sendMessage); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, sendImgMsg); err != nil {
		log.Printf("发送响应失败: %v", err)
	}
}
