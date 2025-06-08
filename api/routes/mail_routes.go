package routes

import (
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/internal/delivery/rest"

	"github.com/gin-gonic/gin"
)

func SetupMailRoutes(h *rest.MailHandler) *gin.Engine {
	cfg := config.LoadConfig()
	router := gin.Default()

	router.POST("mail/manual_container_report", func(c *gin.Context) {
		_ = h.SendManualContainerReport(c, *cfg)
	})

	return router
}