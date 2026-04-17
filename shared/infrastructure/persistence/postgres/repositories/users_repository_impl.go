package repositories

import (
	"context"
	"time"

	"shared/domain/entities"
	"shared/domain/repositories"
	"shared/infrastructure/persistence/postgres"
	"shared/infrastructure/persistence/postgres/ent"
	"shared/infrastructure/persistence/postgres/ent/user"
	"shared/infrastructure/persistence/postgres/ent/usersonaccessgroups"

	"github.com/google/uuid"
)

type usersRepository struct {
	client *ent.Client
}

func NewUsersRepository(client *ent.Client) repositories.UsersRepository {
	return &usersRepository{client: client}
}

func (r *usersRepository) FindUserToLogin(
	ctx context.Context,
	document string,
) (*entities.User, error) {

	u, err := r.client.User.
		Query().
		Where(user.DocumentEQ(document)).
		WithAccessGroups(func(q *ent.UsersOnAccessGroupsQuery) {
			q.Select(usersonaccessgroups.FieldAccessGroupID)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return postgres.ToUserEntity(u), nil
}

func (r *usersRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entities.User, error) {

	u, err := r.client.User.
		Query().
		Where(user.IDEQ(id)).
		WithAccessGroups(func(q *ent.UsersOnAccessGroupsQuery) {
			q.Select(usersonaccessgroups.FieldAccessGroupID)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return postgres.ToUserEntity(u), nil
}

func (r *usersRepository) FindManyToReport(ctx context.Context) (<-chan repositories.UserStreamItem, error) {
	out := make(chan repositories.UserStreamItem)

	go func() {
		defer close(out)
		const pageSize = 100
		var lastCreatedAt time.Time

		for {
			users, err := r.client.User.
				Query().
				Where(user.CreatedAtGT(lastCreatedAt)).
				Order(ent.Asc(user.FieldID)).
				Limit(pageSize).
				All(ctx)
			if err != nil {
				out <- repositories.UserStreamItem{Err: err}
				return
			}

			if len(users) == 0 {
				break
			}

			for _, u := range users {
				out <- repositories.UserStreamItem{
					User: *postgres.ToUserEntity(u),
				}
				lastCreatedAt = u.CreatedAt
			}
		}
	}()

	return out, nil
}
