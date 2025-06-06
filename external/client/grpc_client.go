package client

import (
	"thanhnt208/mail-service/config"
	"thanhnt208/mail-service/proto/pb"

	"google.golang.org/grpc"
)

func StartGrpcClient() (client pb.ContainerAdmServiceClient, err error) {
	cfg := config.LoadConfig()
	conn, err := grpc.Dial(cfg.GrpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client = pb.NewContainerAdmServiceClient(conn)

	return client, nil
}
