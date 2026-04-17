package services

import (
	"context"
	"shared/domain/entities"
	"shared/domain/repositories"

	"github.com/google/uuid"
)

type FindUserByIdService struct {
	userRepo repositories.UsersRepository
}

func NewFindUserByIdService(userRepo repositories.UsersRepository) *FindUserByIdService {
	return &FindUserByIdService{userRepo: userRepo}
}

func (s *FindUserByIdService) Execute(
	ctx context.Context,
	id uuid.UUID,
) (*entities.User, error) {

	return s.userRepo.FindByID(ctx, id)
}
