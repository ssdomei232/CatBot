package review

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"git.mmeiblog.cn/mei/CatBot/configs"
)

type ImageReviewInfo struct {
	Porn    float64 `json:"porn"`
	Sexy    float64 `json:"sexy"`
	Hentai  float64 `json:"hentai"`
	Neutral float64 `json:"neutral"`
	Drawing float64 `json:"drawing"`
}

const (
	maxFileSize = 10 << 30 // 最大支持10GB
	cacheDir    = "cache"
)

// CacheImg 将图片下载到 cache/ 下并重命名为基于URL和时间戳的SHA1哈希值
func CacheImg(imgUrl string) (filename string, err error) {
	// 校验URL合法性
	if imgUrl == "" {
		return "", fmt.Errorf("invalid image url")
	}

	// 创建缓存目录
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache dir: %v", err)
	}

	// 发起GET请求获取图片数据
	resp, err := http.Get(imgUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	// 检查状态码是否成功
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}

	// 验证Content-Type是否为图片格式
	contentType := resp.Header.Get("Content-Type")
	if !isImageMIME(contentType) {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	// 检查文件大小是否超出限制
	if resp.ContentLength > maxFileSize || resp.ContentLength <= 0 {
		return "", fmt.Errorf("file size exceeds limit or invalid: %d bytes", resp.ContentLength)
	}

	// 构造目标文件名
	hashName := fmt.Sprintf("%x", sha1.Sum([]byte(imgUrl+time.Now().String())))[:16]
	ext := getFileExtension(contentType)
	if ext != "" {
		hashName = hashName + ext
	}
	filename = filepath.Join(cacheDir, hashName)

	// 打开本地文件准备写入
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	f, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	// 写入文件内容
	buf := make([]byte, 16*1024)
	_, err = io.CopyBuffer(f, resp.Body, buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	// 这里的 filename 就是路径
	return filename, nil
}

// isImageMIME 判断MIME类型是否属于常见图像格式
func isImageMIME(mime string) bool {
	switch mime {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

func ReviewImage(imgUrl string) (isBadImage bool, err error) {
	config, err := configs.GetConfig()
	if err != nil {
		return false, fmt.Errorf("加载配置失败: %w", err)
	}

	filename, err := CacheImg(imgUrl)
	if err != nil {
		return false, fmt.Errorf("缓存图片失败: %w", err)
	}

	payload := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("打开缓存文件失败: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("关闭文件失败: %v", closeErr)
		}
	}()
	part, err := writer.CreateFormFile("image", filepath.Base(filename))
	if err != nil {
		return false, fmt.Errorf("创建表单文件失败: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return false, fmt.Errorf("复制文件内容失败: %w", err)
	}
	if err = writer.Close(); err != nil {
		return false, fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.NsfwApiUrl, payload)
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("发送请求失败: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应体失败: %w", err)
	}
	var result ImageReviewInfo
	if err = json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("解析响应体失败: %w", err)
	}

	if result.Hentai > 0.5 || result.Porn > 0.5 || result.Sexy > 0.5 {
		return true, nil
	}
	return false, nil
}
