package postgres

import (
	"shared/domain/entities"
	"shared/infrastructure/persistence/postgres/ent"
)

func ToUserEntity(u *ent.User) *entities.User {
	if u == nil {
		return nil
	}

	entity := &entities.User{
		ID:            u.ID,
		Name:          u.Name,
		Password:      u.Password,
		IsManager:     u.IsManager,
		PhotoURL:      u.PhotoURL,
		Document:      u.Document,
		Email:         u.Email,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}

	if !u.DeactivatedAt.IsZero() {
		entity.DeactivatedAt = &u.DeactivatedAt
	}

	if u.ManagerID != [16]byte{} {
		entity.ManagerID = &u.ManagerID
	}

	if u.UserStatusID != 0 {
		entity.UserStatusID = &u.UserStatusID
	}

	// Map access groups if loaded
	if u.Edges.AccessGroups != nil {
		ids := make([]int16, 0, len(u.Edges.AccessGroups))
		for _, ag := range u.Edges.AccessGroups {
			ids = append(ids, int16(ag.AccessGroupID))
		}
		entity.AccessGroupIds = ids
	}

	return entity
}

func ToBusinessEntity(b *ent.Business) *entities.Business {
	if b == nil {
		return nil
	}
	entity := &entities.Business{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
	if !b.DeactivatedAt.IsZero() {
		entity.DeactivatedAt = &b.DeactivatedAt
	}
	return entity
}
