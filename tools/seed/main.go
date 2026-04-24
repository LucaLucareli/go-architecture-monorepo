package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"shared/domain/enums"
	"shared/infrastructure/persistence/postgres/database"
	"shared/infrastructure/persistence/postgres/ent"
	"shared/infrastructure/persistence/postgres/ent/accessgroup"
	"shared/infrastructure/persistence/postgres/ent/business"
	"shared/infrastructure/persistence/postgres/ent/user"
	"shared/infrastructure/persistence/postgres/ent/userstatus"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load("../../.env") // Load from root if running from tools/seed

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	client, err := database.NewEntClient(dbURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Seed UserStatus
	seedUserStatus(ctx, client)

	// Seed AccessGroups
	seedAccessGroups(ctx, client)

	// Seed Business
	biz := seedBusiness(ctx, client)

	// Seed Users
	seedUsers(ctx, client, biz)

	fmt.Println("Seeding completed successfully!")
}

func seedUserStatus(ctx context.Context, client *ent.Client) {
	statuses := []struct {
		ID         int
		ExternalID string
		Name       string
	}{
		{ID: 1, ExternalID: "ACT", Name: "Active"},
		{ID: 2, ExternalID: "INA", Name: "Inactive"},
	}

	for _, s := range statuses {
		exists, _ := client.UserStatus.Query().Where(userstatus.ID(s.ID)).Exist(ctx)
		if exists {
			continue
		}
		err := client.UserStatus.Create().
			SetID(s.ID).
			SetExternalID(s.ExternalID).
			SetName(s.Name).
			Exec(ctx)
		if err != nil {
			log.Printf("could not seed status %s: %v", s.Name, err)
		}
	}
}

func seedAccessGroups(ctx context.Context, client *ent.Client) {
	groups := []struct {
		ID   int
		Name string
	}{
		{ID: int(enums.AccessGroupEmployee), Name: "Employee"},
		{ID: int(enums.AccessGroupAdmin), Name: "Admin"},
		{ID: int(enums.AccessGroupSuperAdmin), Name: "SuperAdmin"},
	}

	for _, g := range groups {
		exists, _ := client.AccessGroup.Query().Where(accessgroup.ID(g.ID)).Exist(ctx)
		if exists {
			continue
		}
		err := client.AccessGroup.Create().
			SetID(g.ID).
			SetName(g.Name).
			Exec(ctx)
		if err != nil {
			log.Printf("could not seed access group %s: %v", g.Name, err)
		}
	}
}

func seedBusiness(ctx context.Context, client *ent.Client) *ent.Business {
	biz, err := client.Business.Query().Where(business.ID(1)).Only(ctx)
	if err == nil {
		return biz
	}

	biz, err = client.Business.Create().
		SetID(1).
		SetName("Goodwe").
		Save(ctx)
	
	if err != nil {
		log.Printf("could not seed business: %v", err)
	}

	return biz
}

func seedUsers(ctx context.Context, client *ent.Client, biz *ent.Business) {
	if biz == nil {
		log.Println("Skipping users seeding because business is nil")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []struct {
		ID           uuid.UUID
		Name         string
		Document     string
		Email        string
		IsManager    bool
		AccessGroups []int
	}{
		{
			ID:           uuid.New(),
			Name:         "Admin User",
			Document:     "00000000001",
			Email:        "admin@goodwe.com",
			IsManager:    true,
			AccessGroups: []int{int(enums.AccessGroupAdmin), int(enums.AccessGroupSuperAdmin)},
		},
		{
			ID:           uuid.New(),
			Name:         "Regular Employee",
			Document:     "00000000002",
			Email:        "employee@goodwe.com",
			IsManager:    false,
			AccessGroups: []int{int(enums.AccessGroupEmployee)},
		},
	}

	for _, u := range users {
		// Check if user with document exists
		exists, _ := client.User.Query().Where(user.Document(u.Document)).Exist(ctx)
		if exists {
			continue
		}

		userRecord, err := client.User.Create().
			SetID(u.ID).
			SetName(u.Name).
			SetDocument(u.Document).
			SetEmail(u.Email).
			SetPassword(string(password)).
			SetIsManager(u.IsManager).
			SetUserStatusID(1). // Active
			SetBusiness(biz).
			Save(ctx)
		
		if err != nil {
			log.Printf("could not seed user %s: %v", u.Name, err)
			continue
		}

		// Add access groups
		for _, groupID := range u.AccessGroups {
			err = client.UsersOnAccessGroups.Create().
				SetUserID(userRecord.ID).
				SetAccessGroupID(groupID).
				Exec(ctx)
			if err != nil {
				log.Printf("could not link user %s to group %d: %v", u.Name, groupID, err)
			}
		}
	}
}
