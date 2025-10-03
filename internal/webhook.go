package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.mmeiblog.cn/mei/CatBot/configs"
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
	Config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	secret := c.DefaultQuery("secret", "")
	if secret != Config.WebhookSecret {
		c.AbortWithStatusJSON(400, gin.H{
			"code":    400,
			"message": "Webhook Secret incorrect",
		})
		return
	}

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

	url := fmt.Sprintf("http://%s:%d/send_group_msg", Config.NapcatHost, Config.NapcatHttpPort)
	method := "POST"
	paylooadStruct := GroupMessage{
		GroupID: groupId,
		Message: []Message{
			{
				Type: "text",
				Data: struct {
					Text string `json:"text"`
				}{Text: message},
			},
		},
	}
	payloadString, err := json.Marshal(paylooadStruct)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"code":    400,
			"message": "message is empty",
		})
	}
	payload := bytes.NewBuffer(payloadString)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	authorization := fmt.Sprintf("Bearer %s", Config.NapcatToken)
	req.Header.Add("Authorization", authorization)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
