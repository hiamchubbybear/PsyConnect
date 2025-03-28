package handler

import (
	pb "consultationservice/grpc/gprc.grpc_generated"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"time"
)

var grpcClient pb.CheckProfileServiceClient

func InitGrpcClient() {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC profile server: %v", err)
	}
	grpcClient = pb.NewCheckProfileServiceClient(conn)
	log.Println("Connected to gRPC profile server")
}
func CheckProfileExists(profileID string) (bool, error) {
	log.Println("CheckProfileExists")
	if grpcClient == nil {
		return false, errors.New("gRPC client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := grpcClient.CheckProfileExists(ctx, &pb.ProfileRequest{ProfileId: profileID})
	if err != nil {
		log.Printf("Error calling Profile Service: %v", err)
		return false, errors.New("Failed to call gRPC profile server")
	}
	log.Printf("Profile %s exists", profileID)
	return res.Exists, nil
}
