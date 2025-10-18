package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bestnite/sub2clash/common"
	"github.com/bestnite/sub2clash/common/database"
	"github.com/bestnite/sub2clash/config"
	"github.com/bestnite/sub2clash/model"
	M "github.com/bestnite/sub2clash/model"
	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"
)

type shortLinkGenRequset struct {
	Config   model.ConvertConfig `form:"config" binding:"required"`
	Password string              `form:"password"`
	ID       string              `form:"id"`
}

type shortLinkUpdateRequest struct {
	Config   model.ConvertConfig `form:"config" binding:"required"`
	Password string              `form:"password" binding:"required"`
	ID       string              `form:"id" binding:"required"`
}

var DB *database.Database

func init() {
	var err error
	DB, err = database.ConnectDB()
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		os.Exit(1)
	}
}

func GenerateLinkHandler(c *gin.Context) {
	var params shortLinkGenRequset
	if err := c.ShouldBind(&params); err != nil {
		c.String(http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var id string
	var password string
	var err error

	if params.ID != "" {
		// 检查自定义ID是否已存在
		exists, err := DB.CheckShortLinkIDExists(params.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "数据库错误")
			return
		}
		if exists {
			c.String(http.StatusBadRequest, "短链已存在")
			return
		}
		id = params.ID
		password = params.Password
	} else {
		// 自动生成短链ID和密码
		id, err = generateUniqueHash(config.GlobalConfig.ShortLinkLength)
		if err != nil {
			c.String(http.StatusInternalServerError, "生成短链失败")
			return
		}
		if params.Password == "" {
			password = common.RandomString(8) // 生成8位随机密码
		} else {
			password = params.Password
		}
	}

	shortLink := model.ShortLink{
		ID:       id,
		Config:   params.Config,
		Password: password,
	}

	if err := DB.CreateShortLink(&shortLink); err != nil {
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}

	// 返回生成的短链ID和密码
	response := map[string]string{
		"id":       id,
		"password": password,
	}
	c.JSON(http.StatusOK, response)
}

func generateUniqueHash(length int) (string, error) {
	for {
		hash := common.RandomString(length)
		exists, err := DB.CheckShortLinkIDExists(hash)
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
		c.String(http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 先获取原有的短链
	existingLink, err := DB.FindShortLinkByID(params.ID)
	if err != nil {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	// 验证密码
	if existingLink.Password != params.Password {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	jsonData, err := json.Marshal(params.Config)
	if err != nil {
		c.String(http.StatusBadRequest, "配置格式错误")
		return
	}
	if err := DB.UpdataShortLink(params.ID, "config", jsonData); err != nil {
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}

	c.String(http.StatusOK, "短链更新成功")
}

func GetRawConfHandler(c *gin.Context) {
	id := c.Param("id")
	password := c.Query("password")

	if strings.TrimSpace(id) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := DB.FindShortLinkByID(id)
	if err != nil {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	err = DB.UpdataShortLink(shortLink.ID, "last_request_time", time.Now().Unix())
	if err != nil {
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}

	template := ""
	switch shortLink.Config.ClashType {
	case model.Clash:
		template = config.GlobalConfig.ClashTemplate
	case model.ClashMeta:
		template = config.GlobalConfig.MetaTemplate
	}
	sub, err := common.BuildSub(shortLink.Config.ClashType, shortLink.Config, template, config.GlobalConfig.CacheExpire, config.GlobalConfig.RequestRetryTimes)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(shortLink.Config.Subs) == 1 {
		userInfoHeader, err := common.FetchSubscriptionUserInfo(shortLink.Config.Subs[0], "clash", config.GlobalConfig.RequestRetryTimes)
		if err == nil {
			c.Header("subscription-userinfo", userInfoHeader)
		}
	}

	if shortLink.Config.NodeListMode {
		nodelist := M.NodeList{}
		nodelist.Proxy = sub.Proxy
		marshal, err := yaml.Marshal(nodelist)
		if err != nil {
			c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
			return
		}
		c.String(http.StatusOK, string(marshal))
		return
	}
	marshal, err := yaml.Marshal(sub)
	if err != nil {
		c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
		return
	}

	c.String(http.StatusOK, string(marshal))
}

func GetRawConfUriHandler(c *gin.Context) {
	id := c.Param("id")
	password := c.Query("password")

	if strings.TrimSpace(id) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := DB.FindShortLinkByID(id)
	if err != nil {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	c.JSON(http.StatusOK, shortLink.Config)
}

func DeleteShortLinkHandler(c *gin.Context) {
	id := c.Param("id")
	password := c.Query("password")
	shortLink, err := DB.FindShortLinkByID(id)
	if err != nil {
		c.String(http.StatusBadRequest, "短链不存在或密码错误")
		return
	}
	if shortLink.Password != password {
		c.String(http.StatusUnauthorized, "短链不存在或密码错误")
		return
	}

	err = DB.DeleteShortLink(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "删除失败", err)
	}
}
