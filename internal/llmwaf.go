package internal

import (
	"strings"

	"git.mmeiblog.cn/mei/CatBot/pkg/review"
)

var blockedKeywords []string = []string{
	"忽略上述指令",
	"忽略之前的指令",
	"system prompt",
	"system message",
	"你是一个",
	"你的名字是",
	"你是",
	"从现在开始",
	"扮演",
	"模仿",
	"输出以下内容",
	"重复以下内容",
	"写一个程序",
	"生成代码",
	"忽略下面的",
	"prompt injection",
	"jailbreak",
	"越狱",
	"Repeat from",
	"你是xxx",
	"Output all content",
}

func llmwaf(msg string) bool {
	if review.ReviewText(msg) {
		return true
	}
	if strings.Contains(msg, "你是") {
		return true
	}
	return false
}

func promptWaf(prompt string) bool {
	lowerPrompt := strings.ToLower(prompt)
	for _, keyword := range blockedKeywords {
		if strings.Contains(lowerPrompt, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
