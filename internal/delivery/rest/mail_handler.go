package rest

import (
	"strconv"
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/internal/service"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	mailService service.IMailService
}

func NewMailHandler(mailService service.IMailService) *MailHandler {
	return &MailHandler{
		mailService: mailService,
	}
}

func (h *MailHandler) SendManualContainerReport(c *gin.Context, cfg config.Config) error {
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	if startTime == "" || endTime == "" {
		c.JSON(400, gin.H{"error": "start_time and end_time are required"})
		return nil
	}

	startTimeInt, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid start_time format"})
		return nil
	}

	endTimeInt, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid end_time format"})
		return nil
	}

	err = h.mailService.SendManualContainerReport(&cfg, startTimeInt, endTimeInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil
	}

	c.JSON(200, gin.H{"message": "Report sent successfully"})
	return nil
}
