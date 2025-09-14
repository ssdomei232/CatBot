package tools

import (
	"fmt"

	"github.com/go-ping/ping"
)

func Ping(ip string) (response string, err error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return "", err
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return "", err
	}
	stats := pinger.Statistics()
	return fmt.Sprintf("PacketLoss: %.2f%%, AvgRtt: %v",
		stats.PacketLoss, stats.AvgRtt), nil
}
