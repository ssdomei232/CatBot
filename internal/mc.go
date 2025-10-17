package internal

import (
	"fmt"
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
	rcon "git.mmeiblog.cn/mei/CatBot/pkg/Rcon"
)

func sendRconCmd(arg string) (Msg string, err error) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	connectInfo := fmt.Sprintf("%s:%d", config.MCSConfig.Host, config.MCSConfig.Port)
	rcon, err := rcon.Dial(connectInfo, config.MCSConfig.Password)
	if err != nil {
		log.Fatalf("RCON连接失败: %v", err)
		return "", err
	}
	cmd := fmt.Sprintf("tp %s", arg)
	Msg, err = rcon.Execute(cmd)
	if err != nil {
		log.Fatalf("RCON执行命令失败: %v", err)
		return "", err
	}
	return Msg, nil
}
