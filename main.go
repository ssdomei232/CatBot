package main

import (
	"fmt"
	"log"
	"time"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/internal"
	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	// config
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// CRON
	c := cron.New()
	c.AddFunc("0 5 * * *", internal.CheckWeatherPerDays)
	log.Println("cron job started")

	// napcat
	connectInfo := fmt.Sprintf("ws://%s:%d/?access_token=%s", Config.NapcatHost, Config.NapcatWebsocketPort, Config.NapcatToken)
	client := napcat.New(
		connectInfo,
		internal.HandleMsg,
		napcat.WithRetryDelay(5*time.Second),
	)
	go func() {
		client.Start(internal.HandleMsg)
	}()
	log.Println("napcat client started")

	// webhook
	r := gin.Default()
	r.POST("/webhook", internal.HandleWebhook)
	r.Run(":8085")
}
