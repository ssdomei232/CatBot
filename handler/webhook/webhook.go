package webhook

import (
	"fmt"
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
	napcathttp "git.mmeiblog.cn/mei/CatBot/pkg/napcat-http"
	"github.com/gin-gonic/gin"
)

type GroupMessage struct {
	GroupID string    `json:"group_id"`
	Message []Message `json:"message"`
}

type Message struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

func HandleWebhook(c *gin.Context) {
	// load Config
	Config, err := configs.GetConfig()
	if err != nil {
		log.Printf("加载配置失败: %v", err)
	}
	secret := c.DefaultQuery("secret", "")
	if secret != Config.WebhookSecret {
		c.AbortWithStatusJSON(400, gin.H{
			"code":    400,
			"message": "Webhook Secret incorrect",
		})
		return
	}

	// Get Msg&group id
	message := c.DefaultPostForm("message", "")
	if message == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"code":    400,
			"message": "message is empty",
		})
		return
	}
	groupId := c.DefaultPostForm("group_id", "")
	if groupId == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"code":    400,
			"message": "group_id is empty",
		})
		return
	}

	// request napcat
	url := fmt.Sprintf("http://%s:%d", Config.NapcatHost, Config.NapcatHttpPort)
	client := napcathttp.NewClient(Config.NapcatToken, url)
	err = client.SendGroupMsg(groupId, message)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": fmt.Sprintf("发送消息失败: %v", err),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
