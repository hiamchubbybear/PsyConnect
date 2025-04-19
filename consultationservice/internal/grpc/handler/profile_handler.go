package handler

import (
	"consultationservice/internal/grpc/gprc.grpc_generated"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func CheckProfileExists(profileId string) (bool, error) {
	conn, err := grpc.Dial(
		"127.0.0.1:9091",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return false, errors.New("Failed to connect to gRPC profile server")
	}
	log.Println("Connected to gRPC profile server")
	defer conn.Close()
	client := gprc_grpc_generated.NewCheckProfileServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := client.CheckProfileExists(ctx, &gprc_grpc_generated.ProfileRequest{ProfileId: profileId})
	if err != nil {
		log.Printf("Error calling Profile Service: %v", err)
		return false, errors.New("Failed to call gRPC profile server")
	}

	return res.Exists, nil
}
