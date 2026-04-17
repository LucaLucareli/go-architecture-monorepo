package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Name          string
	Password      string
	IsManager     bool
	PhotoURL      string
	Document      string
	Email         string
	ManagerID     *uuid.UUID
	UserStatusID  *int
	DeactivatedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relations
	AccessGroupIds []int16
	BusinessID     *int
}

type Business struct {
	ID            int
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeactivatedAt *time.Time
}

type AccessGroup struct {
	ID            int
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeactivatedAt *time.Time
}
