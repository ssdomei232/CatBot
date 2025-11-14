package internal

import (
	"fmt"
	"log"
	"strings"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/handler"
	rcon "git.mmeiblog.cn/mei/CatBot/pkg/Rcon"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
)

func sendRconTpCmd(groupMsg *napcat.Message) (msg string) {
	if len(groupMsg.RawMessage) < 8 {
		return "请输入正确的指令"
	}

	if !isGamerExist(groupMsg.UserID) {
		return "请先绑定游戏名"
	}
	gamerName, _ := getGamerName(groupMsg.Sender.UserID)

	cmd := fmt.Sprintf("tp %s %s", gamerName, groupMsg.RawMessage[4:])
	msg, err := runRconCmd(cmd)
	if err != nil {
		log.Printf("RCON执行命令失败: %v", err)
		return "RCON执行命令失败"
	}

	return "传送成功"
}

func bindMCSGamer(cmdList []string, groupMsg *napcat.Message) (msg string, err error) {
	if len(cmdList) < 3 {
		return "请输入正确的指令", nil
	}

	db, err := handler.GetDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	gameName := cmdList[2]
	if strings.Contains(gameName, " ") {
		return "游戏名不能包含空格", nil
	}

	if isGameNameUsed(gameName) {
		return "游戏名已存在", nil
	}

	if isGamerExist(groupMsg.UserID) {
		// 如果玩家已存在，将旧游戏名从白名单删除
		oldGameName, _ := getGamerName(groupMsg.UserID)
		_, err = runRconCmd(fmt.Sprintf("whitelist remove %s", oldGameName))
		if err != nil {
			log.Printf("删除白名单失败: %v", err)
			return "", err
		}

		// update db
		_, err = db.Exec("UPDATE mcs SET game_name = ? WHERE qq = ?", gameName, groupMsg.UserID)
		if err != nil {
			log.Printf("更新数据失败: %v", err)
			return "", err
		}

		return "更新数据成功,您的旧游戏名已从白名单移除,请放心游玩", nil
	}

	_, err = db.Exec("INSERT INTO mcs (qq, qq_nickname, game_name) VALUES (?, ?, ?)", groupMsg.UserID, groupMsg.Sender.Nickname, gameName)
	if err != nil {
		log.Printf("插入数据失败: %v", err)
		return "", err
	}

	return "绑定成功\n欢迎来mc.mei.lv玩(1.21.8原版)", nil
}

func isGamerExist(qqNumber int) bool {
	db, err := handler.GetDB()
	if err != nil {
		return false
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM mcs WHERE qq = ?", qqNumber).Scan(&count)
	if err != nil {
		log.Printf("查询数据失败: %v", err)
		return false
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}

func isGameNameUsed(gameName string) bool {
	db, err := handler.GetDB()
	if err != nil {
		return false
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM mcs WHERE game_name = ?", gameName).Scan(&count)
	if err != nil {
		log.Printf("查询数据失败: %v", err)
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func getGamerName(qqNumber int) (gameName string, err error) {
	db, err := handler.GetDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	err = db.QueryRow("SELECT game_name FROM mcs WHERE qq = ?", qqNumber).Scan(&gameName)
	if err != nil {
		log.Printf("查询数据失败: %v", err)
		return "", err
	}

	return gameName, nil
}

func runRconCmd(cmd string) (Msg string, err error) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("加载配置失败: %v", err)
	}
	connectInfo := fmt.Sprintf("%s:%d", config.MCSConfig.Host, config.MCSConfig.Port)
	rcon, err := rcon.Dial(connectInfo, config.MCSConfig.Password)
	defer rcon.Close()
	if err != nil {
		log.Printf("RCON连接失败: %v", err)
		return "", err
	}
	Msg, err = rcon.Execute(cmd)
	if err != nil {
		log.Printf("RCON执行命令失败: %v", err)
		return "", err
	}
	return Msg, nil
}
