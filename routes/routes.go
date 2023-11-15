package routes

import (
	"ddns/handlers"
	"ddns/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func DnspodRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("cymjj123"))
	r.Use(sessions.Sessions("ddv_ddns", store))

	r.GET("/", handlers.LoginView)
	r.POST("/loginauth", handlers.LoginAuth)

	r.GET("/ddns", handlers.DDNS)

	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())
	authRoutes.GET("/domainlist", handlers.DomainList)
	authRoutes.GET("/recordlist", handlers.Recordlist)
	authRoutes.POST("/domaincreate", handlers.Domaincreate)
	authRoutes.GET("/domainremove", handlers.Domainremove)
	authRoutes.GET("/domainstatus", handlers.DomainStatus)
	authRoutes.GET("/recordeditf", handlers.RecordEditf)
	authRoutes.POST("/modifyrecord", handlers.ModifyRecord)
	authRoutes.GET("/recordstatus", handlers.RecordStatus)
	authRoutes.GET("/recordremove", handlers.RecordRemove)

	return r
}
