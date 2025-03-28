package handler

import (
	pb "consultationservice/grpc/gprc.grpc_generated"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

// CheckProfileExists kiểm tra xem profile có tồn tại không thông qua gRPC
func CheckProfileExists(profileId string) (bool, error) {
	// Kết nối với gRPC server
	conn, err := grpc.Dial(
		"127.0.0.1:9091",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Thay thế WithInsecure()
	)
	if err != nil {
		return false, errors.New("failed to connect to gRPC profile server")
	}
	log.Println("Connected to gRPC profile server")
	defer conn.Close()

	// Tạo client từ connection
	client := pb.NewCheckProfileServiceClient(conn)

	// Thiết lập timeout cho request
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Gửi request đến gRPC server
	res, err := client.CheckProfileExists(ctx, &pb.ProfileRequest{ProfileId: profileId})
	if err != nil {
		log.Printf("Error calling Profile Service: %v", err)
		return false, errors.New("failed to call gRPC profile server")
	}

	return res.Exists, nil
}
