package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	weatherResp, _ := GetWeatherRaw()
	// 格式化天气信息
	return formatWeatherInfo(weatherResp)
}

func GetWeatherStruct() SimpleWeatherData {
	weatherResp, _ := GetWeatherRaw()
	return FormatWeatherToStruct(weatherResp)
}

func GetWeatherRaw() (WeatherResponse, error) {
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
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据失败: %v", err)
		return WeatherResponse{}, err
	}

	// 解析JSON数据
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		log.Printf("解析JSON数据失败: %v", err)
		return WeatherResponse{}, err
	}

	// 检查API返回状态
	if weatherResp.Status != "1" {
		log.Printf("API返回错误: %s, info: %s, infocode: %s", weatherResp.Status, weatherResp.Info, weatherResp.Infocode)
		return WeatherResponse{}, err
	}

	return weatherResp, nil
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

// SimpleWeatherData 简单天气数据结构
type SimpleWeatherData struct {
	City     string             `json:"city"`
	Province string             `json:"province"`
	Weather  []SimpleDayWeather `json:"weather"`
}

// SimpleDayWeather 日天气数据
type SimpleDayWeather struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"day_weather"`
	NightWeather string `json:"night_weather"`
	DayTemp      string `json:"day_temp"`
	NightTemp    string `json:"night_temp"`
}

// FormatWeatherToStruct 将天气数据格式化为简单结构体
func FormatWeatherToStruct(resp WeatherResponse) SimpleWeatherData {
	result := SimpleWeatherData{}

	if len(resp.Forecasts) == 0 {
		return result
	}

	forecast := resp.Forecasts[0]
	result.City = forecast.City
	result.Province = forecast.Province

	// 处理最多三天的数据
	daysCount := min(len(forecast.Casts), 3)

	for i := 0; i < daysCount; i++ {
		cast := forecast.Casts[i]
		dayWeather := SimpleDayWeather{
			Date:         cast.Date,
			Week:         cast.Week,
			DayWeather:   cast.DayWeather,
			NightWeather: cast.NightWeather,
			DayTemp:      cast.DayTemp,
			NightTemp:    cast.NightTemp,
		}
		result.Weather = append(result.Weather, dayWeather)
	}

	return result
}

// CheckRainForecast 检查天气预报中是否有雨，如果有则返回提醒信息
func CheckRainForecast() string {
	weatherData, _ := GetWeatherRaw()

	if len(weatherData.Forecasts) == 0 || len(weatherData.Forecasts[0].Casts) == 0 {
		return ""
	}

	forecast := weatherData.Forecasts[0]

	// 检查未来几天的天气
	for _, cast := range forecast.Casts {
		// 检查白天或夜间是否有雨
		if hasRain(cast.DayWeather) || hasRain(cast.NightWeather) {
			return fmt.Sprintf("%s(%s)有雨，记得带伞", cast.Date, cast.Week)
		}
	}

	return ""
}

// hasRain 判断天气描述中是否包含雨
func hasRain(weather string) bool {
	rainKeywords := []string{"雨", "雷阵雨", "小雨", "中雨", "大雨", "暴雨"}
	for _, keyword := range rainKeywords {
		if strings.Contains(weather, keyword) {
			return true
		}
	}
	return false
}
