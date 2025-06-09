package rest

import (
	"fmt"
	"net/http"
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/internal/dto"
	"thanhnt208/mail-service/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	mailService service.IMailService
	cfg         *config.Config
}

func NewMailHandler(mailService service.IMailService, cfg *config.Config) *MailHandler {
	return &MailHandler{
		mailService: mailService,
		cfg:         cfg,
	}
}

func (h *MailHandler) SendManualContainerReport(c *gin.Context) error {
	var req *dto.ReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
	}

	layout := "2006-01-02"
	startTime, startErr := time.Parse(layout, req.StartDate)
	endTime, endErr := time.Parse(layout, req.EndDate)
	if startErr != nil || endErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format. Use YYYY-MM-DD.",
		})
		return nil
	}

	endTime = endTime.Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
	adminAddrConf := *h.cfg
	adminAddrConf.AdminEmailAddr = req.AdminEmail

	if err := h.mailService.SendManualContainerReport(&adminAddrConf, startTime.Unix(), endTime.Unix()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send report: " + err.Error(),
		})
		return nil
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Report for period %s to %s has been sent to %s.", req.StartDate, req.EndDate, req.AdminEmail),
	})

	return nil
}
