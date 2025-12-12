package internal

import (
	"fmt"
	"log"

	"git.mmeiblog.cn/mei/CatBot/configs"
	napcathttp "git.mmeiblog.cn/mei/CatBot/pkg/napcat-http"
	"git.mmeiblog.cn/mei/CatBot/tools"
)

func CheckWeatherPerDays() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("加载配置失败: %v", err)
	}

	msg := tools.CheckRainForecast()

	url := fmt.Sprintf("http://%s:%d", config.NapcatHost, config.NapcatHttpPort)
	napcatClient := napcathttp.NewClient(config.NapcatToken, url)
	err = napcatClient.SendGroupMsg(fmt.Sprint(726833553), msg)
	if err != nil {
		log.Printf("发送天气预报请求失败: %v", err)
		return
	}
}
