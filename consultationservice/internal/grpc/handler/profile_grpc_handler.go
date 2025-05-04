package handler

import (
	pb "consultationservice/internal/grpc/gprc.grpc_generated"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type ProfileGrpc struct {
	checkClient    pb.CheckProfileServiceClient
	matchingClient pb.CreateMatchingRequestClient
	conn           *grpc.ClientConn
}

func NewProfileGrpc(addr string) (*ProfileGrpc, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProfileGrpc{
		checkClient:    pb.NewCheckProfileServiceClient(conn),
		matchingClient: pb.NewCreateMatchingRequestClient(conn),
		conn:           conn,
	}, nil
}

func (r *ProfileGrpc) CheckProfileExists(profileId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.checkClient.CheckProfileExists(ctx, &pb.ProfileRequest{ProfileId: profileId})
	if err != nil {
		log.Printf("gRPC CheckProfileExists error: %v", err)
		return false, errors.New("failed to call CheckProfileExists")
	}
	return res.Exists, nil
}

func (r *ProfileGrpc) FriendRequestAccept(clientId, therapistId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.matchingClient.CreateMatchingRequest(ctx, &pb.MatchingRequest{
		ClientProfileId:    clientId,
		TherapistProfileId: therapistId,
	})
	if err != nil {
		log.Printf("gRPC FriendRequestAccept error: %v", err)
		return false, errors.New("failed to call CreateMatchingRequest")
	}
	return res.GetSuccess(), nil
}

func (r *ProfileGrpc) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
