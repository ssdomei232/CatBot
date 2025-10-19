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
	var returnMessage string
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

	// 原神拦截器
	if strings.Contains(GroupMsg.RawMessage, "gamecenter.qq.com") {
		returnMessage = "以上消息存在欺诈行为(点击可能会下载某种'热门'游戏)，请勿相信"
	}

	// 每次消息都需要执行的部分
	Record(*GroupMsg)
	ReviewText(conn, GroupMsg.RawMessage, GroupMsg.GroupID, GroupMsg.MessageID, GroupMsg.UserID)

	// 功能部分
	if strings.Contains(commandText, ".chat") {
		if promptWaf(commandText) {
			returnMessage = "消息被 Prompt WAF 拦截"
		} else {
			returnMessage, err = ai.SendComplain(commandText[5:]) // 去掉".chat"前缀
			if err != nil {
				log.Printf("ai处理失败: %v", err)
				return
			}
			if llmwaf(returnMessage) {
				returnMessage = "消息被 LLM WAF 拦截"
			}
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
	} else if strings.Contains(commandText, "杜奕") || strings.Contains(commandText, "杜伊") || strings.Contains(commandText, "喵") {
		returnMessage = "👀"
	} else if strings.Contains(commandText, ".help") {
		sendHelp(conn, GroupMsg.GroupID)
		return
	} else if strings.Contains(commandText, ".weather") {
		returnMessage = tools.GetWeather()
	} else if strings.Contains(commandText, ".findfood") {
		go findfood(conn, commandText[10:], GroupMsg.GroupID)
		returnMessage = "正在搜索..."
	} else if strings.Contains(commandText, ".bus") {
		returnMessage = tools.FindBus(commandText[5:])
	} else if strings.Contains(commandText, ".zanwo") {
		sendLike(conn, GroupMsg.GroupID, GroupMsg.UserID)
		return
	} else if strings.Contains(commandText, "wcnm") || strings.Contains(commandText, "cnmd") {
		returnMessage = "中国国家导弹防御系统"
	} else if strings.Contains(commandText, ".tp") {
		rconMsg, err := sendRconCmd(commandText[4:])
		if err != nil {
			log.Printf("RCON执行命令失败: %v", err)
			returnMessage = "RCON执行命令失败"
		} else {
			returnMessage = rconMsg
		}
	} else if strings.Contains(commandText, ".temperature") {
		temperature, err := GetTemperature()
		if err != nil {
			log.Printf("获取温度失败: %v", err)
			returnMessage = "获取温度失败"
		} else {
			returnMessage = fmt.Sprintf("当前室外温度为: %s", temperature)
		}
	}

	sendMessage, err := napcat.MarshalGroupTextMsg(GroupMsg.GroupID, returnMessage)
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
