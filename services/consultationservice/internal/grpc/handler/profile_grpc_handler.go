package handler

import (
	dto "consultationservice/internal/dto"
	pb "consultationservice/internal/grpc/gprc.grpc_generated"
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProfileGrpc struct {
	profileClient  pb.ProfileServiceClient
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
		profileClient:  pb.NewProfileServiceClient(conn),
		checkClient:    pb.NewCheckProfileServiceClient(conn),
		matchingClient: pb.NewCreateMatchingRequestClient(conn),
		conn:           conn,
	}, nil
}

func (r *ProfileGrpc) GetProfile(profileId string) (*pb.ProfileData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if profileId == "" {
		return nil, errors.New("profile id cannot be empty")
	}

	req := &pb.GetProfileRequest{
		ProfileId: profileId,
	}

	res, err := r.profileClient.GetProfile(ctx, req)
	if err != nil {
		log.Printf("gRPC GetProfile error: %v", err)
		return nil, err
	}

	if !res.Success {
		return nil, errors.New(res.ErrorMessage)
	}

	return res.Profile, nil
}

func (r *ProfileGrpc) CheckProfileExists(profileId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := r.checkClient.CheckProfileExists(ctx, &pb.ProfileRequest{ProfileId: profileId})
	if err != nil {
		log.Print(err)
		return false, err
	}
	return res.Exists, nil
}

func (r *ProfileGrpc) FriendRequestAccept(clientId, therapistId, message string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.matchingClient.CreateMatchingRequest(ctx, &pb.MatchingRequest{
		Message:            message,
		ClientProfileId:    clientId,
		TherapistProfileId: therapistId,
	})
	if err != nil {
		log.Print(err)
		return false, err
	}
	return res.GetSuccess(), nil
}

func (r *ProfileGrpc) ResponseMatchingRequest(request dto.ResponseMatchingRequest, profileId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := r.matchingClient.ResponseRequestMatching(ctx, &pb.ResponseMatchingRequest{
		ClientProfileId:    request.ClientId,
		TherapistProfileId: profileId,
		Option:             request.Option,
	})
	if err != nil {
		log.Print(err)
		return false, err
	}
	return res.GetSuccess(), nil
}

func (r *ProfileGrpc) RetriveProfileInfo(profileId string) (*pb.ProfileResponseV1, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if profileId == "" {
		return nil, errors.New("profile id cannot be empty")
	}

	profileRequest := &pb.ProfileRequestV1{
		ProfileId: profileId,
	}

	res, err := r.checkClient.CheckProfileExistsV1(ctx, profileRequest)
	if err != nil || !res.Exists {
		return nil, err
	}
	return res, nil
}

func (r *ProfileGrpc) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
