package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

// POIResponse 高德地图POI搜索响应结构
type POIResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Infocode string `json:"infocode"`
	POIs     []struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		TypeCode string `json:"typecode"`
		Address  string `json:"address"`
		Location string `json:"location"`
		PCode    string `json:"pcode"`
		CityCode string `json:"citycode"`
		CityName string `json:"cityname"`
		AdName   string `json:"adname"`
		AdCode   string `json:"adcode"`
		Photos   []struct {
			URL string `json:"url"`
		} `json:"photos"`
	} `json:"pois"`
}

// GaodeConfig 高德地图API配置
type GaodeConfig struct {
	APIKey string
}

// FindFoodService 餐饮搜索服务
type FindFoodService struct {
	config GaodeConfig
}

// NewFindFoodService 创建新的餐饮搜索服务实例
func NewFindFoodService(apiKey string) *FindFoodService {
	return &FindFoodService{
		config: GaodeConfig{APIKey: apiKey},
	}
}

// SearchFood 搜索餐饮场所
func (s *FindFoodService) SearchFood(keyword string) (*POIResponse, error) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	// 构建请求URL
	baseURL := "https://restapi.amap.com/v3/place/text"
	params := url.Values{}
	params.Add("key", s.config.APIKey)
	params.Add("keywords", keyword)
	params.Add("city", strconv.Itoa(config.CityABCode))
	params.Add("types", "餐饮服务")
	params.Add("offset", "20")
	params.Add("page", "1")
	params.Add("extensions", "all")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 发起HTTP请求
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求高德地图API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var poiResp POIResponse
	err = json.Unmarshal(body, &poiResp)
	if err != nil {
		return nil, fmt.Errorf("解析JSON响应失败: %v", err)
	}

	// 检查API返回状态
	if poiResp.Status != "1" {
		return nil, fmt.Errorf("高德地图API返回错误，infocode: %s", poiResp.Infocode)
	}

	return &poiResp, nil
}

// FormatSinglePOIInfo 格式化单个POI信息为"店名:xxx,地址:xxx"格式
func (s *FindFoodService) FormatSinglePOIInfo(poiResp *POIResponse) (string, string) {
	// 检查是否有POI数据
	if len(poiResp.POIs) == 0 {
		return "", ""
	}

	// 随机选择一个POI
	randomIndex := rand.Intn(len(poiResp.POIs))
	selectedPOI := poiResp.POIs[randomIndex]

	// 格式化信息
	info := fmt.Sprintf("店名:%s,地址:%s", selectedPOI.Name, selectedPOI.Address)

	// 获取对应POI的第一张图片URL（如果存在）
	var photoURL string
	if len(selectedPOI.Photos) > 0 && selectedPOI.Photos[0].URL != "" {
		photoURL = selectedPOI.Photos[0].URL
	}

	return info, photoURL
}

// SearchAndFormat 搜索并格式化结果，只返回一个随机POI及其图片
func (s *FindFoodService) SearchAndFormat(keyword string) (string, string, error) {
	// 搜索餐饮场所
	poiResp, err := s.SearchFood(keyword)
	if err != nil {
		return "", "", err
	}

	// 格式化单个POI信息并获取对应的图片URL
	formattedInfo, photoURL := s.FormatSinglePOIInfo(poiResp)

	return formattedInfo, photoURL, nil
}
