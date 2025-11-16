package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.mmeiblog.cn/mei/CatBot/configs"
	"git.mmeiblog.cn/mei/CatBot/tools"
)

func CheckWeatherPerDays() {
	Config, err := configs.GetConfig()
	if err != nil {
		log.Printf("加载配置失败: %v", err)
	}

	msg := tools.CheckRainForecast()

	url := fmt.Sprintf("http://%s:%d/send_group_msg", Config.NapcatHost, Config.NapcatHttpPort)
	method := "POST"
	paylooadStruct := GroupMessage{
		GroupID: "726833553",
		Message: []Message{
			{
				Type: "text",
				Data: struct {
					Text string `json:"text"`
				}{Text: msg},
			},
		},
	}
	payloadString, _ := json.Marshal(paylooadStruct)
	payload := bytes.NewBuffer(payloadString)
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Content-Type", "application/json")
	authorization := fmt.Sprintf("Bearer %s", Config.NapcatToken)
	req.Header.Add("Authorization", authorization)

	res, _ := client.Do(req)
	defer res.Body.Close()
	return
}
