package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

// FindBus 查询公交信息并返回格式化的字符串
func FindBus(line string) string {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatal("获取配置文件错误: ", err)
	}

	// 构建API请求URL
	apiURL := fmt.Sprintf("https://api.lolimi.cn/API/che/api.php?type=text&city=%s&line=%s",
		config.CityName, line)

	// 发送HTTP请求
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Sprintf("请求公交API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取响应数据失败: %v", err)
	}

	// 直接返回API返回的数据，因为已经是所需格式
	return string(body)
}
