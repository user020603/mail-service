package client

import (
	"context"
	"fmt"
	"thanhnt208/mail-service/pkg/logger"
	"thanhnt208/mail-service/proto/pb"
)

type IGetContainerInfoClient interface {
	GetContainerInformation(startTime, endTime int64) (*pb.GetContainerInformationResponse, error)
	GetContainerUptimeDuration(startTime, endTime int64) (*pb.GetContainerUptimeDurationResponse, error)
}

type getContainerInfoClient struct {
	client pb.ContainerAdmServiceClient
	logger logger.ILogger
}

func NewGetContainerInfoClient(client pb.ContainerAdmServiceClient, logger logger.ILogger) IGetContainerInfoClient {
	return &getContainerInfoClient{
		client: client,
		logger: logger,
	}
}

func (c *getContainerInfoClient) GetContainerInformation(startTime, endTime int64) (*pb.GetContainerInformationResponse, error) {
	resp, err := c.client.GetContainerInformation(
		context.Background(),
		&pb.GetContainerInformationRequest{
			StartTime: startTime,
			EndTime:   endTime,
		},
	)

	if err != nil {
		c.logger.Error("Failed to get container info", "error", err)
		return nil, fmt.Errorf("failed to get container info: %w", err)
	}

	c.logger.Info("Successfully retrieved container info", "numContainers", resp.NumContainers, "numRunningContainers", resp.NumRunningContainers, "numStoppedContainers", resp.NumStoppedContainers, "meanUptimeRatio", resp.MeanUptimeRatio)
	return resp, nil
}

func (c *getContainerInfoClient) GetContainerUptimeDuration(startTime, endTime int64) (*pb.GetContainerUptimeDurationResponse, error) {
	resp, err := c.client.GetContainerUptimeDuration(
		context.Background(),
		&pb.GetContainerInformationRequest{
			StartTime: startTime,
			EndTime:   endTime,
		},
	)

	if err != nil {
		c.logger.Error("Failed to get container uptime", "error", err)
		return nil, fmt.Errorf("failed to get container uptime: %w", err)
	}

	c.logger.Info("Successfully retrieved container uptime duration", "uptimeDuration", resp.NumContainers, "numRunningContainers", resp.NumRunningContainers, "numStoppedContainers", resp.NumStoppedContainers, "UptimeDetails", resp.UptimeDetails)
	return resp, nil
}