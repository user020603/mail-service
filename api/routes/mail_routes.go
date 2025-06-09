package routes

import (
	"thanhnt208/mail-service/internal/delivery/rest"

	"github.com/gin-gonic/gin"
)

func SetupMailRoutes(h *rest.MailHandler) *gin.Engine {
	router := gin.Default()

	router.POST("mail/send_period_report", func(c *gin.Context) {
		_ = h.SendManualContainerReport(c)
	})

	return router
}
