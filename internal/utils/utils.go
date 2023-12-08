package utils

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type DownloadResult struct {
	Success bool
	URL     string
	Error   error
	Image   []byte
}

func DownloadImage(url string) DownloadResult {
	resp, err := http.Get(url)
	var ret DownloadResult
	if err != nil {
		ret = DownloadResult{
			Success: false,
			URL:     url,
			Error:   err,
		}
		return ret
	}
	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		ret = DownloadResult{
			Success: false,
			URL:     url,
			Error:   err,
		}
		return ret
	}
	ret = DownloadResult{
		Success: false,
		URL:     url,
		Error:   err,
		Image:   image,
	}
	return ret
}

func DownloadImages(urls []string) ([]DownloadResult, error) {
	n := len(urls)
	ret := make([]DownloadResult, n)
	var wg sync.WaitGroup
	for idx, url := range urls {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()
			ret[i] = DownloadImage(url)
		}(url, idx)
	}
	wg.Wait()
	return ret, nil
}

func GenerateShortKey() string {
	rand.Seed(time.Now().UnixNano())

	// 生成基于时间戳的部分
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestampStr := strconv.FormatInt(timestamp, 36) // 转换为36进制

	// 截取长度在4到7位之间的时间戳部分，确保总长度为5到8位
	prefixLength := rand.Intn(4) + 4
	if len(timestampStr) > prefixLength {
		timestampStr = timestampStr[:prefixLength]
	}

	// 生成随机数部分
	randomPart := rand.Intn(36) // 生成一个0到35之间的随机数
	randomStr := strconv.FormatInt(int64(randomPart), 36)

	// 合并两部分
	return timestampStr + randomStr
}
