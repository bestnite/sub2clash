package server

import (
	"embed"
	"net/http"

	"github.com/bestnite/sub2clash/server/handler"
	"github.com/bestnite/sub2clash/server/middleware"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var staticFiles embed.FS

func SetRoute(r *gin.Engine) {
	r.GET("/convert/:config", middleware.ZapLogger(), handler.ConvertHandler())
	r.GET("/s/:id", middleware.ZapLogger(), handler.GetRawConfHandler)
	r.POST("/short", middleware.ZapLogger(), handler.GenerateLinkHandler)
	r.PUT("/short", middleware.ZapLogger(), handler.UpdateLinkHandler)
	r.GET("/short/:id", middleware.ZapLogger(), handler.GetRawConfUriHandler)
	r.DELETE("/short/:id", middleware.ZapLogger(), handler.DeleteShortLinkHandler)

	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("frontend/dist/", http.FS(staticFiles))
	})
	r.GET(
		"/assets/*filepath", func(c *gin.Context) {
			c.FileFromFS("frontend/dist/assets/"+c.Param("filepath"), http.FS(staticFiles))
		},
	)
}
