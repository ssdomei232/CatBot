package internal

import (
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/pkg/ai"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"git.mmeiblog.cn/mei/CatBot/tools"
	"github.com/gorilla/websocket"
)

func handleAIChat(conn *websocket.Conn, cmdList []string, groupMsg *napcat.Message) {
	var err error
	var msg string
	var reply []byte

	if len(cmdList) < 2 {
		return
	}

	if promptWaf(cmdList[1]) {
		return
	}

	msg, err = ai.SendComplain(cmdList[1])
	if err != nil {
		return
	}
	if llmwaf(msg) {
		return
	}

	reply, _ = napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, msg)
	writeMutex.Lock()
	defer writeMutex.Unlock()
	napcat.SendMsg(conn, reply)

	return
}

func handlePing(conn *websocket.Conn, cmdList []string, groupMsg *napcat.Message) {
	var err error
	var msg string
	var reply []byte

	if len(cmdList) < 2 {
		return
	}

	msg, err = tools.Ping(cmdList[1])
	if err != nil {
		return
	}

	reply, _ = napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, msg)
	writeMutex.Lock()
	defer writeMutex.Unlock()
	napcat.SendMsg(conn, reply)

	return
}

func handleEat(conn *websocket.Conn, cmdList []string, groupId int) {
	var err error

	if len(cmdList) < 2 {
		return
	}

	config, err := configs.GetConfig()
	if err != nil {
		log.Fatal("获取配置文件错误: ", err)
	}

	foodService := tools.NewFindFoodService(config.GDKey)
	info, photoURL, _ := foodService.SearchAndFormat(cmdList[1])
	sendMessage, _ := napcat.MarshalGroupTextMsg(groupId, info)
	sendImgMsg, err := napcat.MarshalGroupImgMsg(groupId, photoURL)

	writeMutex.Lock()
	defer writeMutex.Unlock()
	napcat.SendMsg(conn, sendMessage)
	napcat.SendMsg(conn, sendImgMsg)
}

func handleFindBus(conn *websocket.Conn, cmdList []string, groupId int) {
	var msg string
	var reply []byte

	if len(cmdList) < 2 {
		return
	}

	msg = tools.FindBus(cmdList[1])
	reply, _ = napcat.MarshalGroupTextMsg(groupId, msg)
	writeMutex.Lock()
	defer writeMutex.Unlock()
	napcat.SendMsg(conn, reply)
}

func handleLike(conn *websocket.Conn, groupMsg *napcat.Message) {
	likeMsg, _ := napcat.MarshalLikeMsg(groupMsg.UserID, 10)
	reply, _ := napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, "赞了你10下")
	writeMutex.Lock()
	defer writeMutex.Unlock()
	napcat.SendMsg(conn, likeMsg)
	napcat.SendMsg(conn, reply)
}

func handleMc(conn *websocket.Conn, cmdList []string, groupMsg *napcat.Message) {
	var err error
	var msg string
	var reply []byte

	if len(cmdList) < 2 {
		return
	}

	switch cmdList[1] {
	case "bind":
		msg, err = bindMCSGamer(cmdList, groupMsg)
		if err != nil {
			msg = "绑定失败"
		}
		writeMutex.Lock()
		defer writeMutex.Unlock()
		reply, _ = napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, msg)
		napcat.SendMsg(conn, reply)
	case "tp":
		msg = sendRconTpCmd(groupMsg)
		reply, _ = napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, msg)
		writeMutex.Lock()
		defer writeMutex.Unlock()
		napcat.SendMsg(conn, reply)
	}

}
