package service

import (
	"fmt"
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/external/client"
	"thanhnt208/mail-service/pkg/logger"
	"time"

	"gopkg.in/gomail.v2"
)

type IMailService interface {
	SendManualContainerReport(cfg *config.Config, startTime, endTime int64) error
}

type mailService struct {
	grpcClient client.IGetContainerInfoClient
	logger     logger.ILogger
}

func NewMailService(grpcClient client.IGetContainerInfoClient, logger logger.ILogger) IMailService {
	return &mailService{
		grpcClient: grpcClient,
		logger:     logger,
	}
}

func (s *mailService) SendManualContainerReport(cfg *config.Config, startTime, endTime int64) error {
	resp, err := s.grpcClient.GetContainerInformation(startTime, endTime)
	if err != nil {
		s.logger.Error("Failed to get container info", "error", err)
		return fmt.Errorf("failed to get container info: %w", err)
	}

	numContainers := int(resp.NumContainers)
	numRunningContainers := int(resp.NumRunningContainers)
	numStoppedContainers := int(resp.NumStoppedContainers)
	meanUptimeRatio := float64(resp.MeanUptimeRatio)

	to := cfg.AdminEmailAddr
	subject := "Docker Container Report for " + time.Unix(startTime, 0).Format("2006-01-02") + " to " + time.Unix(endTime, 0).Format("2006-01-02")
	body := fmt.Sprintf(
		"Hello,\n\nHere is the Docker container report for the period from %s to %s:\n\n"+
			"Total number of containers: %d\n"+
			"Number of running containers: %d\n"+
			"Number of stopped containers: %d\n"+
			"Mean uptime ratio: %.2f%%\n\n"+
			"Best regards,\nYour Docker Monitoring Service",
		time.Unix(startTime, 0).Format("2006-01-02"),
		time.Unix(endTime, 0).Format("2006-01-02"),
		numContainers,
		numRunningContainers,
		numStoppedContainers,
		meanUptimeRatio*100.0,
	)

	senderEmail := cfg.SenderEmailAddr
	senderPassword := cfg.SenderEmailPassword

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(m); err != nil {
		s.logger.Error("Failed to send email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}
