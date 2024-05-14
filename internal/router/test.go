package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"sapphire-server/internal/conf"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"strconv"
	"strings"
	"time"
)

type TestRouter struct {
}

func NewTestRouter(engine *gin.Engine) *TestRouter {
	router := &TestRouter{}
	testGroup := engine.Group("/test")
	testGroup.GET("/hw", router.HandleHw)
	testGroup.POST("/upload", router.HandleUpload)
	testGroup.POST("/changeRole", router.HandleChangeRole)
	return router
}

func (t *TestRouter) HandleHw(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

// HandleChangeRole 根据用户 id 更改权限
func (t *TestRouter) HandleChangeRole(ctx *gin.Context) {
	userId := ctx.Query("userId")
	role := ctx.PostForm("role")
	//fmt.Println("userId:", userId, "role:", role)
	if role == "" {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("role is required"))
		return
	}
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("userId is required"))
		return
	}
	user, err := dao.First[domain.User]("id = ?", userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	} else if user == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("user not found"))
		return
	}
	dao.Modify(user, "role", role)
}

// HandleUpload Upload a picture
func (t *TestRouter) HandleUpload(ctx *gin.Context) {
	// Read from form-data
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}
	// Judge if the file is a picture
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("only support jpg and png"))
		return
	}

	// Start uploading
	directUrl, err := t.Upload(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"url": directUrl,
	})
}

// Upload Implement the upload function
func (t *TestRouter) Upload(fileHeader *multipart.FileHeader) (string, error) {
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

	filename := "sapphire_" + strconv.FormatInt(time.Now().Unix(), 10) + t.getExtension(fileHeader.Header.Get("Content-Type"))

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

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
	defer resp.Body.Close()

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

func (t *TestRouter) getExtension(contentType string) string {
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
