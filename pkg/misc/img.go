package misc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sapphire-server/internal/conf"
	"strconv"
	"strings"
	"time"
)

func UploadImage(src []byte, fileName string) (string, error) {
	var err error
	info := conf.GetImgConfig()
	re := regexp.MustCompile(`svrUrl: (.*?); directUrl: (.*?); auth string: (.*?);`)
	matches := re.FindStringSubmatch(info)
	var svrUrl, directUrl, auth string
	if len(matches) == 4 {
		svrUrl = matches[1]
		directUrl = matches[2]
		auth = matches[3]
	} else {
		fmt.Println("String format is not valid.")
	}

	filename := "sapphire_" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + fileName

	fileBytes := src

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("PUT", svrUrl+filename, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", auth)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return directUrl + filename, nil
}

func getExtension(contentType string) string {
	fileTypeMap := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}

	parts := strings.Split(contentType, "/")
	if len(parts) != 2 {
		return ""
	}

	extension, ok := fileTypeMap[contentType]
	if !ok {
		return ""
	}

	return extension
}
