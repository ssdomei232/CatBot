package internal

import (
	"context"
	"net/http"
	"time"

	"git.mmeiblog.cn/mei/CatBot/configs"
	ha "github.com/mkelcik/go-ha-client"
)

// 获取 HA Client
func getHAClient() (client *ha.Client, err error) {
	config, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	client = ha.NewClient(ha.ClientConfig{Token: config.HAConfig.Token, Host: config.HAConfig.Token}, &http.Client{
		Timeout: 30 * time.Second,
	})

	if err = client.Ping(context.Background()); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

// 获取温度
func GetTemperature() (terperaturn string, err error) {
	client, err := getHAClient()
	if err != nil {
		return "", err
	}
	temperature, err := client.GetStateForEntity(context.Background(), "sensor.miaomiaoc_cn_blt_3_1k7hgiaj8kg00_t2_temperature_p_2_1")
	if err != nil {
		return "", err
	}
	return temperature.State, nil
}
