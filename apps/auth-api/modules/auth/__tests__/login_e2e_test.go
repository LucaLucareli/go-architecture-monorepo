package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-api/modules/auth/controllers"
	"auth-api/modules/auth/dto/io"
	"auth-api/modules/auth/dto/request"
	"auth-api/modules/auth/services"
	"shared/application/auth"
	"shared/domain/enums"
	"shared/infrastructure/persistence/postgres/ent"
	"shared/infrastructure/persistence/postgres/repositories"
	"shared/pkg/interceptors"
	"shared/pkg/testutils"
	"shared/pkg/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var testPasswordHash string

func init() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testPasswordHash = string(hash)
}

func setupLoginTest(t *testing.T) (*echo.Echo, *ent.Client) {
	client := testutils.GetTestClient(t)
	testutils.CleanupDatabase(t, client)
  
	ctx := context.Background()

	client.UserStatus.Create().SetID(1).SetExternalID("ACT").SetName("Active").ExecX(ctx)
	client.AccessGroup.Create().SetID(int(enums.AccessGroupAdmin)).SetName("Admin").ExecX(ctx)
	biz := client.Business.Create().SetID(1).SetName("Test Business").SaveX(ctx)
	
	u := client.User.Create().
		SetID(uuid.New()).
		SetDocument("12345678901").
		SetPassword(testPasswordHash).
		SetUserStatusID(1).
		SetBusiness(biz).
		SaveX(ctx)
	
	client.UsersOnAccessGroups.Create().
		SetUserID(u.ID).
		SetAccessGroupID(int(enums.AccessGroupAdmin)).
		ExecX(ctx)

	userRepo := repositories.NewUsersRepository(client)
	authService := &auth.AuthService{
		AccessSecret:      "access-secret",
		RefreshSecret:     "refresh-secret",
		AccessExpiryHours: 1,
		RefreshExpiryDays: 7,
		UserRepo:          userRepo,
	}

	e := echo.New()
	e.Validator = validation.NewValidator()
	e.Use(interceptors.TransformInterceptor)

	loginService := services.NewLoginService(authService)
	loginController := controllers.NewLoginController(loginService)
	e.POST("/login", loginController.Handle)

	return e, client
}

func TestLoginE2E(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		e, client := setupLoginTest(t)
		defer client.Close()

		reqBody := request.LoginRequestDTO{
			Document: "12345678901",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var resp struct {
			Message string           `json:"message"`
			Data    io.LoginOutputDTO `json:"data"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Login realizado com sucesso", resp.Message)
		assert.NotEmpty(t, resp.Data.AccessToken)
		assert.NotEmpty(t, resp.Data.RefreshToken)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		e, client := setupLoginTest(t)
		defer client.Close()

		reqBody := request.LoginRequestDTO{
			Document: "12345678901",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Validation Error - Document Too Short", func(t *testing.T) {
		e, client := setupLoginTest(t)
		defer client.Close()

		reqBody := request.LoginRequestDTO{
			Document: "123",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
