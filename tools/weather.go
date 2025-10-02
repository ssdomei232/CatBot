package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

type WeatherResponse struct {
	Status    string     `json:"status"`
	Count     string     `json:"count"`
	Info      string     `json:"info"`
	Infocode  string     `json:"infocode"`
	Forecasts []Forecast `json:"forecasts"`
}

type Forecast struct {
	City       string `json:"city"`
	Adcode     string `json:"adcode"`
	Province   string `json:"province"`
	Reporttime string `json:"reporttime"`
	Casts      []Cast `json:"casts"`
}

type Cast struct {
	Date           string `json:"date"`
	Week           string `json:"week"`
	DayWeather     string `json:"dayweather"`
	NightWeather   string `json:"nightweather"`
	DayTemp        string `json:"daytemp"`
	NightTemp      string `json:"nighttemp"`
	DayWind        string `json:"daywind"`
	NightWind      string `json:"nightwind"`
	DayPower       string `json:"daypower"`
	NightPower     string `json:"nightpower"`
	DayTempFloat   string `json:"daytemp_float"`
	NightTempFloat string `json:"nighttemp_float"`
}

// GetWeather 获取并格式化天气信息
func GetWeather() string {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 构建请求URL
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?parameters&key=%s&city=%d&extensions=all",
		config.GDKey, config.CityABCode)
	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("请求天气API失败: %v", err)
		return "获取天气信息失败"
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据失败: %v", err)
		return "读取天气信息失败"
	}

	// 解析JSON数据
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Printf("解析JSON数据失败: %v", err)
		return "解析天气信息失败"
	}

	// 检查API返回状态
	if weatherResp.Status != "1" {
		log.Printf("API返回错误: %s, info: %s, infocode: %s", weatherResp.Status, weatherResp.Info, weatherResp.Infocode)
		return "天气API调用失败"
	}

	// 格式化天气信息
	return formatWeatherInfo(weatherResp)
}

// formatWeatherInfo 格式化天气信息为指定格式
func formatWeatherInfo(resp WeatherResponse) string {
	if len(resp.Forecasts) == 0 || len(resp.Forecasts[0].Casts) < 2 {
		return "天气数据不足"
	}

	casts := resp.Forecasts[0].Casts

	// 只取今天和明天的数据
	today := casts[0]
	tomorrow := casts[1]

	// 构建格式化字符串
	result := fmt.Sprintf("今天上午:%s\n今天下午:%s\n今天晚上:%s\n----------------\n明天上午:%s\n明天下午:%s\n明天晚上:%s",
		today.DayWeather, today.DayWeather, today.NightWeather,
		tomorrow.DayWeather, tomorrow.DayWeather, tomorrow.NightWeather)

	return result
}
