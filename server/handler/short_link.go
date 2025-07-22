package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bestnite/sub2clash/common"
	"github.com/bestnite/sub2clash/common/database"
	"github.com/bestnite/sub2clash/config"
	"github.com/bestnite/sub2clash/model"

	"github.com/gin-gonic/gin"
)

type shortLinkGenRequset struct {
	Url      string `form:"url" binding:"required"`
	Password string `form:"password"`
	CustomID string `form:"customId"`
}

type shortLinkUpdateRequest struct {
	Hash     string `form:"hash" binding:"required"`
	Url      string `form:"url" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func respondWithError(c *gin.Context, code int, message string) {
	c.String(code, message)
	c.Abort()
}

func GenerateLinkHandler(c *gin.Context) {
	var params shortLinkGenRequset
	if err := c.ShouldBind(&params); err != nil {
		respondWithError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if strings.TrimSpace(params.Url) == "" {
		respondWithError(c, http.StatusBadRequest, "URL 不能为空")
		return
	}

	var hash string
	var password string
	var err error

	if params.CustomID != "" {
		// 检查自定义ID是否已存在
		exists, err := database.CheckShortLinkHashExists(params.CustomID)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "数据库错误")
			return
		}
		if exists {
			respondWithError(c, http.StatusBadRequest, "短链已存在")
			return
		}
		hash = params.CustomID
		password = params.Password
	} else {
		// 自动生成短链ID和密码
		hash, err = generateUniqueHash(config.GlobalConfig.ShortLinkLength)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "生成短链接失败")
			return
		}
		if params.Password == "" {
			password = common.RandomString(8) // 生成8位随机密码
		} else {
			password = params.Password
		}
	}

	shortLink := model.ShortLink{
		Hash:     hash,
		Url:      params.Url,
		Password: password,
	}

	if err := database.SaveShortLink(&shortLink); err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	// 返回生成的短链ID和密码
	response := map[string]string{
		"hash":     hash,
		"password": password,
	}
	c.JSON(http.StatusOK, response)
}

func generateUniqueHash(length int) (string, error) {
	for {
		hash := common.RandomString(length)
		exists, err := database.CheckShortLinkHashExists(hash)
		if err != nil {
			return "", err
		}
		if !exists {
			return hash, nil
		}
	}
}

func UpdateLinkHandler(c *gin.Context) {
	var params shortLinkUpdateRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		respondWithError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 先获取原有的短链接
	existingLink, err := database.FindShortLinkByHash(params.Hash)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "未找到短链接")
		return
	}

	// 验证密码
	if existingLink.Password != params.Password {
		respondWithError(c, http.StatusUnauthorized, "密码错误")
		return
	}

	// 更新URL，但保持原密码不变
	shortLink := model.ShortLink{
		Hash:     params.Hash,
		Url:      params.Url,
		Password: existingLink.Password, // 保持原密码不变
	}

	if err := database.SaveShortLink(&shortLink); err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	c.String(http.StatusOK, "短链接更新成功")
}

func GetRawConfHandler(c *gin.Context) {
	hash := c.Param("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	shortLink.LastRequestTime = time.Now().Unix()
	err = database.SaveShortLink(shortLink)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	targetPath := strings.TrimPrefix(shortLink.Url, "/")
	requestURL := fmt.Sprintf("%s://%s/%s", scheme, host, targetPath)

	client := &http.Client{
		Timeout: 30 * time.Second, // 30秒超时
	}

	response, err := client.Get(requestURL)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "请求错误: "+err.Error())
		return
	}
	defer response.Body.Close()

	all, err := io.ReadAll(response.Body)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "读取错误: "+err.Error())
		return
	}

	c.String(http.StatusOK, string(all))
}

func GetRawConfUriHandler(c *gin.Context) {
	hash := c.Query("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	c.String(http.StatusOK, shortLink.Url)
}
