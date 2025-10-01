package review

import (
	"github.com/kirklin/go-swd"
)

func ReviewText(text string) (isBadMessage bool) {
	detector, _ := swd.New()
	return detector.Detect(text)
}
