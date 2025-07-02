package service

import (
	"fmt"
	"os"
	"path/filepath"
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/external/client"
	"thanhnt208/mail-service/pkg/logger"
	"time"

	"github.com/go-pdf/fpdf" 
	"gopkg.in/gomail.v2"
)

type IMailService interface {
	SendManualContainerReport(cfg *config.Config, startTime, endTime int64) error
	SendUptimeReport(cfg *config.Config, startTime, endTime int64) error
	RunDailyReportJob(cfg *config.Config) error
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

func generateEmailHTMLBody(startTime, endTime time.Time, numContainers, numRunningContainers, numStoppedContainers int, meanUptimeRatio float64) string {
	return fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; color: #222;">
        <p>Dear System Administrator,</p>
        <p>
            Please find below a summary of the Docker container health report for the period 
            <strong>%s</strong> to <strong>%s</strong>.
        </p>
        <p>
            The detailed report is attached as a PDF for your reference.
        </p>
        <h3 style="color: #1565c0;">Report Summary</h3>
        <ul>
            <li><strong>Total Containers:</strong> %d</li>
            <li><strong>Running Containers:</strong> %d</li>
            <li><strong>Stopped Containers:</strong> %d</li>
            <li><strong>Mean Uptime Ratio:</strong> %.2f%%</li>
        </ul>
        <p>
            If you have any questions or require further details, please let us know.
        </p>
        <p>
            Best regards,<br>
            <strong>Docker Monitoring Service</strong>
        </p>
        <img src="https://miro.medium.com/v2/resize:fit:1100/format:webp/1*QQk-kwU6qwPlIkR_rzxrYQ.gif" 
             alt="Docker Monitoring Service" style="width:320px; margin-top:16px;" />
    </div>
    `,
		startTime.Format("January 2, 2006"),
		endTime.Format("January 2, 2006"),
		numContainers,
		numRunningContainers,
		numStoppedContainers,
		meanUptimeRatio*100.0,
	)
}

func generateReportPDF(startTime, endTime time.Time, numContainers, numRunningContainers, numStoppedContainers int, meanUptimeRatio float64) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)	

	pdf.Cell(40, 10, "Docker Container Health Report")
	pdf.Ln(20) 

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Report Period:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 10, "Metric")
	pdf.Cell(40, 10, "Value")
	pdf.Ln(10)
	pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+100, pdf.GetY()) // Separator line
	pdf.Ln(5)

	reportData := map[string]string{
		"Total number of containers":   fmt.Sprintf("%d", numContainers),
		"Number of running containers": fmt.Sprintf("%d", numRunningContainers),
		"Number of stopped containers": fmt.Sprintf("%d", numStoppedContainers),
		"Mean uptime ratio":            fmt.Sprintf("%.2f%%", meanUptimeRatio*100.0),
	}

	pdf.SetFont("Arial", "", 12)
	for key, value := range reportData {
		pdf.Cell(60, 10, key)
		pdf.Cell(40, 10, value)
		pdf.Ln(10)
	}

	pdfPath := filepath.Join(os.TempDir(), fmt.Sprintf("container-report-%d.pdf", time.Now().Unix()))
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfPath, nil
}

func (s *mailService) SendManualContainerReport(cfg *config.Config, startTimeInt, endTimeInt int64) error {
	resp, err := s.grpcClient.GetContainerInformation(startTimeInt, endTimeInt)
	if err != nil {
		s.logger.Error("Failed to get container info", "error", err)
		return fmt.Errorf("failed to get container info: %w", err)
	}

	numContainers := int(resp.NumContainers)
	numRunningContainers := int(resp.NumRunningContainers)
	numStoppedContainers := int(resp.NumStoppedContainers)
	meanUptimeRatio := float64(resp.MeanUptimeRatio)

	startTime := time.Unix(startTimeInt, 0)
	endTime := time.Unix(endTimeInt, 0)

	pdfPath, err := generateReportPDF(startTime, endTime, numContainers, numRunningContainers, numStoppedContainers, meanUptimeRatio)
	if err != nil {
		s.logger.Error("Failed to create PDF report", "error", err)
		return err 
	}

	defer func() {
		if err := os.Remove(pdfPath); err != nil {
			s.logger.Error("Failed to remove temporary PDF file", "path", pdfPath, "error", err)
		}
	}()

	subject := fmt.Sprintf("[Monitoring] Daily Container Health Report: %s", startTime.Format("2006-01-02"))
	htmlBody := generateEmailHTMLBody(startTime, endTime, numContainers, numRunningContainers, numStoppedContainers, meanUptimeRatio)

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SenderEmailAddr)
	m.SetHeader("To", cfg.AdminEmailAddr)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", htmlBody)

	m.Attach(pdfPath)

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.SenderEmailAddr, cfg.SenderEmailPassword)

	if err := d.DialAndSend(m); err != nil {
		s.logger.Error("Failed to send email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Email sent successfully with PDF attachment", "to", cfg.AdminEmailAddr, "subject", subject)
	return nil
}

func (s *mailService) RunDailyReportJob(cfg *config.Config) error {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	startTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, yesterday.Location())

	err := s.SendManualContainerReport(cfg, startTime.Unix(), endTime.Unix())
	if err != nil {
		s.logger.Error("Failed to run daily report job", "error", err)
		return fmt.Errorf("failed to run daily report job: %w", err)
	} else {
		s.logger.Info("Daily report job completed successfully", "date", yesterday.Format("2006-01-02"))
		return nil
	}
}

func generateUptimeEmailHTMLBody(startTime, endTime time.Time, numContainers, numRunningContainers, numStoppedContainers int, totalUptime time.Duration, perContainerUptime map[string]time.Duration) string {
    return fmt.Sprintf(`
    <div style="font-family: Arial, sans-serif; color: #222;">
        <p>Dear System Administrator,</p>
        <p>
            Please find below a summary of the Docker container uptime report for the period 
            <strong>%s</strong> to <strong>%s</strong>.
        </p>
        <p>
            The detailed report is attached as a PDF for your reference.
        </p>
        <h3 style="color: #1565c0;">Uptime Report Summary</h3>
        <ul>
            <li><strong>Total Containers:</strong> %d</li>
            <li><strong>Running Containers:</strong> %d</li>
            <li><strong>Stopped Containers:</strong> %d</li>
            <li><strong>Total Uptime:</strong> %s</li>
        </ul>
        <p>
            If you have any questions or require further details, please let us know.
        </p>
        <p>
            Best regards,<br>
            <strong>Docker Monitoring Service</strong>
        </p>
        <img src="https://miro.medium.com/v2/resize:fit:1100/format:webp/1*QQk-kwU6qwPlIkR_rzxrYQ.gif" 
             alt="Docker Monitoring Service" style="width:320px; margin-top:16px;" />
    </div>
    `,
        startTime.Format("January 2, 2006"),
        endTime.Format("January 2, 2006"),
        numContainers,
        numRunningContainers,
        numStoppedContainers,
        totalUptime,
    )
}

func generateUptimeReportPDF(startTime, endTime time.Time, numContainers, numRunningContainers, numStoppedContainers int, totalUptime time.Duration, perContainerUptime map[string]time.Duration) (string, error) {
    pdf := fpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)	

    pdf.Cell(40, 10, "Docker Container Uptime Report")
    pdf.Ln(20) 

    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(40, 10, "Report Period:")
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")))
    pdf.Ln(15)

    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(60, 10, "Metric")
    pdf.Cell(40, 10, "Value")
    pdf.Ln(10)
    pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+100, pdf.GetY()) // Separator line
    pdf.Ln(5)

    reportData := map[string]string{
        "Total number of containers":   fmt.Sprintf("%d", numContainers),
        "Number of running containers": fmt.Sprintf("%d", numRunningContainers),
        "Number of stopped containers": fmt.Sprintf("%d", numStoppedContainers),
        "Total uptime":                 totalUptime.String(),
    }

    pdf.SetFont("Arial", "", 12)
    for key, value := range reportData {
        pdf.Cell(60, 10, key)
        pdf.Cell(40, 10, value)
        pdf.Ln(10)
    }
    
    // Add per container uptime details
    if len(perContainerUptime) > 0 {
        pdf.Ln(10)
        pdf.SetFont("Arial", "B", 14)
        pdf.Cell(40, 10, "Per Container Uptime")
        pdf.Ln(15)
        
        pdf.SetFont("Arial", "B", 12)
        pdf.Cell(60, 10, "Container ID")
        pdf.Cell(40, 10, "Uptime")
        pdf.Ln(10)
        pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+100, pdf.GetY()) // Separator line
        pdf.Ln(5)

        pdf.SetFont("Arial", "", 12)
        for containerID, duration := range perContainerUptime {
            pdf.Cell(60, 10, containerID)
            pdf.Cell(40, 10, duration.String())
            pdf.Ln(10)
        }
    }

    pdfPath := filepath.Join(os.TempDir(), fmt.Sprintf("uptime-report-%d.pdf", time.Now().Unix()))
    err := pdf.OutputFileAndClose(pdfPath)
    if err != nil {
        return "", fmt.Errorf("failed to generate PDF: %w", err)
    }

    return pdfPath, nil
}

func (s *mailService) SendUptimeReport(cfg *config.Config, startTimeInt, endTimeInt int64) error {
    resp, err := s.grpcClient.GetContainerUptimeDuration(startTimeInt, endTimeInt)
    if err != nil {
        s.logger.Error("Failed to get container uptime info", "error", err)
        return fmt.Errorf("failed to get container uptime info: %w", err)
    }
    
    startTime := time.Unix(startTimeInt, 0)
    endTime := time.Unix(endTimeInt, 0)

    numContainers := int(resp.NumContainers)
    numRunningContainers := int(resp.NumRunningContainers)
    numStoppedContainers := int(resp.NumStoppedContainers)
    
    totalUptime := time.Duration(resp.UptimeDetails.TotalUptime) * time.Millisecond 
    perContainerUptime := make(map[string]time.Duration)
    for containerID, milliseconds := range resp.UptimeDetails.PerContainerUptime {
        perContainerUptime[containerID] = time.Duration(milliseconds) * time.Millisecond
    }

    pdfPath, err := generateUptimeReportPDF(startTime, endTime, numContainers, numRunningContainers, 
        numStoppedContainers, totalUptime, perContainerUptime)
    if err != nil {
        s.logger.Error("Failed to create PDF uptime report", "error", err)
        return err 
    }

    defer func() {
        if err := os.Remove(pdfPath); err != nil {
            s.logger.Error("Failed to remove temporary PDF file", "path", pdfPath, "error", err)
        }
    }()

    subject := fmt.Sprintf("[Monitoring] Container Uptime Report: %s to %s", 
        startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
    htmlBody := generateUptimeEmailHTMLBody(startTime, endTime, numContainers, numRunningContainers, 
        numStoppedContainers, totalUptime, perContainerUptime)

    m := gomail.NewMessage()
    m.SetHeader("From", cfg.SenderEmailAddr)
    m.SetHeader("To", cfg.AdminEmailAddr)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", htmlBody)
    m.Attach(pdfPath)

    d := gomail.NewDialer("smtp.gmail.com", 587, cfg.SenderEmailAddr, cfg.SenderEmailPassword)

    if err := d.DialAndSend(m); err != nil {
        s.logger.Error("Failed to send uptime report email", "error", err)
        return fmt.Errorf("failed to send uptime report email: %w", err)
    }

    s.logger.Info("Uptime report email sent successfully", "to", cfg.AdminEmailAddr, "subject", subject)
    return nil
}
