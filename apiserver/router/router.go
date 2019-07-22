/**
 * 功能描述: 路由管理
 * @Date: 2019-04-14
 * @author: lixiaoming
 */
package router

import (
	_ "apiserver/docs"

	"apiserver/controllers/sd"
	"apiserver/controllers/user"
	"apiserver/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

// 加载中间件, 路由,返回gin引擎
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// 安装swag命令行工具 go get -u github.com/swaggo/swag/cmd/swag
	// 加载swagger api 文档
	// 重新生成文档  swag init
	// 访问地址/swagger/index.html
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.POST("/login", user.Login)

	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
	}

	return g
}
