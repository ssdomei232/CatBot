package tools

import (
	"fmt"

	"github.com/go-ping/ping"
)

func Ping(ip string) (response string, err error) {
	pinger, err := ping.NewPinger(ip) // 修复：使用传入的ip参数
	if err != nil {
		return "", err
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return "", err
	}
	stats := pinger.Statistics()
	// 修复：使用fmt.Sprintf格式化结构体
	return fmt.Sprintf("PacketLoss: %.2f%%, AvgRtt: %v",
		stats.PacketLoss, stats.AvgRtt), nil
}
