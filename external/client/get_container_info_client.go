package client

import (
	"context"
	"fmt"
	"thanhnt208/mail-service/pkg/logger"
	"thanhnt208/mail-service/proto/pb"
)

type IGetContainerInfoClient interface {
	GetContainerInfo(startTime, endTime int64) (*pb.GetContainerInfomationResponse, error)
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

func (c *getContainerInfoClient) GetContainerInfo(startTime, endTime int64) (*pb.GetContainerInfomationResponse, error) {
	resp, err := c.client.GetContainerInfo(
		context.Background(),
		&pb.GetContainerInfomationRequest{
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
