package routes

import (
	"thanhnt208/mail-service/api/middlewares"
	"thanhnt208/mail-service/internal/delivery/rest"

	"github.com/gin-gonic/gin"
)

func SetupMailRoutes(h *rest.MailHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/send_uptime_ratio",
		middlewares.JWTAuthMiddleware(),
		middlewares.CheckScopeMiddleware("mail:send_uptime_ratio"),
		func(c *gin.Context) {
			_ = h.SendManualContainerReport(c)
		})

	router.POST("/send_uptime_duration",
		middlewares.JWTAuthMiddleware(),
		middlewares.CheckScopeMiddleware("mail:send_uptime_duration"),
		func(c *gin.Context) {
			_ = h.SendUptimeReport(c)
		})

	return router
}
