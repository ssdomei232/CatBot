package main

import (
	"fmt"
	"log"
	"time"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/internal"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"github.com/gin-gonic/gin"
)

func main() {
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	connectInfo := fmt.Sprintf("ws://%s:%d/?access_token=%s", Config.NapcatHost, Config.NapcatWebsocketPort, Config.NapcatToken)
	client := napcat.New(
		connectInfo,
		internal.SendGroupMsg,
		napcat.WithRetryDelay(5*time.Second),
	)

	go func() {
		client.Start(internal.SendGroupMsg)
	}()

	r := gin.Default()
	r.POST("/webhook", internal.HandleWebhook)

	if err := r.Run(":8085"); err != nil {
		log.Fatalf("Gin服务器启动失败: %v", err)
	}
}
