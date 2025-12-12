package internal

import (
	"fmt"
	"strings"
	"sync"

	"git.mmeiblog.cn/mei/CatBot/pkg/review"
	"git.mmeiblog.cn/mei/CatBot/tools"
	"github.com/gorilla/websocket"
	napcat "github.com/ssdomei232/go-napcat-ws"
)

var writeMutex sync.Mutex

func HandleMsg(conn *websocket.Conn, messageType int, message []byte) {
	var err error
	var groupMsg *napcat.Message
	groupMsg, err = napcat.Parse(message)
	if err != nil {
		return
	}

	// 解析消息项
	messageItems, err := groupMsg.GetMessageItems()
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
			if isAdminGroup(groupMsg.GroupID) {
				review.CacheImg(imgData.URL)
				ReviewImage(conn, imgData.URL, groupMsg.GroupID, groupMsg.MessageID, groupMsg.UserID)
			}
		}
	}

	// 原神拦截器
	if strings.Contains(groupMsg.RawMessage, "gamecenter.qq.com") {
		reply, _ := napcat.MarshalGroupReplyMsg(groupMsg.GroupID, groupMsg.MessageID, "以上消息存在欺诈行为(点击可能会下载某种'热门'游戏)，请勿相信")
		writeMutex.Lock()
		defer writeMutex.Unlock()
		napcat.SendMsg(conn, reply)
	}

	// 每次消息都需要执行的部分
	Record(*groupMsg)
	ReviewText(conn, groupMsg.RawMessage, groupMsg.GroupID, groupMsg.MessageID, groupMsg.UserID)

	// Split command
	cmdList := strings.Split(commandText, " ")
	if len(cmdList) < 1 {
		return
	}

	switch cmdList[0] {
	case ".chat":
		handleAIChat(conn, cmdList, groupMsg)
	case ".ping":
		handlePing(conn, cmdList, groupMsg)
	case ".tq":
		msg := tools.GetWeather()
		reply, _ := napcat.MarshalGroupTextMsg(groupMsg.GroupID, msg)
		writeMutex.Lock()
		defer writeMutex.Unlock()
		napcat.SendMsg(conn, reply)
	case ".eat":
		handleEat(conn, cmdList, groupMsg.GroupID)
	case ".temp":
		var msg string
		temperature, err := GetTemperature()
		if err != nil {
			return
		} else {
			msg = fmt.Sprintf("当前温度为: %s°C", temperature)
		}
		reply, _ := napcat.MarshalGroupTextMsg(groupMsg.GroupID, msg)
		writeMutex.Lock()
		defer writeMutex.Unlock()
		napcat.SendMsg(conn, reply)
	case ".mc":
		handleMc(conn, cmdList, groupMsg)
	case ".bus":
		handleFindBus(conn, cmdList, groupMsg.GroupID)
	case ".like":
		handleLike(conn, groupMsg)
	}
}
