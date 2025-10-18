package main

import (
	_ "embed"
	"io"

	"github.com/bestnite/sub2clash/common"
	"github.com/bestnite/sub2clash/config"
	"github.com/bestnite/sub2clash/logger"
	"github.com/bestnite/sub2clash/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	var err error

	err = common.MkEssentialDir()
	if err != nil {
		logger.Logger.Panic("create essential dir failed", zap.Error(err))
	}

	err = config.LoadConfig()

	logger.InitLogger(config.GlobalConfig.LogLevel)
	if err != nil {
		logger.Logger.Panic("load config failed", zap.Error(err))
	}

	logger.Logger.Info("database connect success")
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard

	r := gin.Default()

	server.SetRoute(r)
	logger.Logger.Info("server is running at " + config.GlobalConfig.Address)
	err := r.Run(config.GlobalConfig.Address)
	if err != nil {
		logger.Logger.Error("server running failed", zap.Error(err))
		return
	}
}
