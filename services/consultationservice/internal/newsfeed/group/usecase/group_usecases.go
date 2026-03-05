package usecase

import (
	"consultationservice/internal/newsfeed/group/domain"
	"consultationservice/internal/newsfeed/group/repository"
	"context"
	"errors"
	"time"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrAlreadyMember = errors.New("user is already a member of this group")
)


type CreateGroupUseCase struct {
	groupRepo repository.GroupRepository
}

func NewCreateGroupUseCase(groupRepo repository.GroupRepository) *CreateGroupUseCase {
	return &CreateGroupUseCase{groupRepo: groupRepo}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, group *domain.Group) error {
	if group.Name == "" {
		return errors.New("group name is required")
	}
	group.CreatedAt = time.Now().UTC()
	group.UpdatedAt = time.Now().UTC()
	return uc.groupRepo.CreateGroup(ctx, group)
}


type GetGroupsUseCase struct {
	groupRepo repository.GroupRepository
}

func NewGetGroupsUseCase(groupRepo repository.GroupRepository) *GetGroupsUseCase {
	return &GetGroupsUseCase{groupRepo: groupRepo}
}

func (uc *GetGroupsUseCase) Execute(ctx context.Context, category, query string, limit, skip int64) ([]domain.Group, error) {
	return uc.groupRepo.GetGroups(ctx, category, query, limit, skip)
}


type GetGroupByIDUseCase struct {
	groupRepo repository.GroupRepository
}

func NewGetGroupByIDUseCase(groupRepo repository.GroupRepository) *GetGroupByIDUseCase {
	return &GetGroupByIDUseCase{groupRepo: groupRepo}
}

func (uc *GetGroupByIDUseCase) Execute(ctx context.Context, id string) (*domain.Group, error) {
	group, err := uc.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group, nil
}


type JoinGroupUseCase struct {
	groupRepo repository.GroupRepository
}

func NewJoinGroupUseCase(groupRepo repository.GroupRepository) *JoinGroupUseCase {
	return &JoinGroupUseCase{groupRepo: groupRepo}
}

func (uc *JoinGroupUseCase) Execute(ctx context.Context, groupID string, userID string) error {
	
	_, err := uc.groupRepo.GetGroupByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	
	isMember, err := uc.groupRepo.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrAlreadyMember
	}

	return uc.groupRepo.AddMember(ctx, groupID, userID)
}
