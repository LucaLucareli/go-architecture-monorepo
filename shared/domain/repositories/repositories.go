package repositories

import (
	"context"
	"shared/domain/entities"

	"github.com/google/uuid"
)

type UsersRepository interface {
	FindUserToLogin(ctx context.Context, document string) (*entities.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindManyToReport(ctx context.Context) (<-chan UserStreamItem, error)
}

type UserStreamItem struct {
	User entities.User
	Err  error
}
