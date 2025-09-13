package main

import (
	"time"

	"git.mmeiblog.cn/mei/aiComplain/internal"
	"git.mmeiblog.cn/mei/aiComplain/pkg/napcat"
)

func main() {
	client := napcat.New(
		"ws://127.0.0.1:3001/?access_token=^l^}BOdE[8s<k@g@",
		internal.SendGroupMsg,
		napcat.WithRetryDelay(5*time.Second),
	)
	client.Start(internal.SendGroupMsg)
}
