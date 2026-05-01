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

var refreshTestPasswordHash string

func init() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	refreshTestPasswordHash = string(hash)
}

func setupRefreshTest(t *testing.T) (*echo.Echo, *ent.Client) {
	client := testutils.GetTestClient(t)
	testutils.CleanupDatabase(t, client)

	ctx := context.Background()

	client.UserStatus.Create().SetID(1).SetExternalID("ACT").SetName("Active").ExecX(ctx)
	client.AccessGroup.Create().SetID(int(enums.AccessGroupAdmin)).SetName("Admin").ExecX(ctx)
	biz := client.Business.Create().SetID(1).SetName("Test Business").SaveX(ctx)
	
	u := client.User.Create().
		SetID(uuid.New()).
		SetDocument("12345678901").
		SetPassword(refreshTestPasswordHash).
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

	refreshTokenService := services.NewRefreshTokenService(authService)
	refreshTokenController := controllers.NewRefreshTokenController(refreshTokenService)
	e.POST("/refresh-token", refreshTokenController.Handle)

	return e, client
}

func getTokensForRefresh(t *testing.T, e *echo.Echo, document, password string) (string, string) {
	reqBody := request.LoginRequestDTO{
		Document: document,
		Password: password,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Fatalf("Failed to login for refresh test: expected 200, got %d", rec.Code)
	}
	
	var resp struct {
		Data io.LoginOutputDTO `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	
	return resp.Data.AccessToken, resp.Data.RefreshToken
}

func TestRefreshTokenE2E(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		e, client := setupRefreshTest(t)
		defer client.Close()
		
		_, refreshToken := getTokensForRefresh(t, e, "12345678901", "password123")

		reqBody := request.RefreshTokenRequestDTO{
			RefreshToken: refreshToken,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var resp struct {
			Message string                  `json:"message"`
			Data    io.RefreshTokenOutputDTO `json:"data"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Token atualizado com sucesso", resp.Message)
		assert.NotEmpty(t, resp.Data.AccessToken)
		assert.NotEmpty(t, resp.Data.RefreshToken)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		e, client := setupRefreshTest(t)
		defer client.Close()

		reqBody := request.RefreshTokenRequestDTO{
			RefreshToken: "invalid-token",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Empty Token", func(t *testing.T) {
		e, client := setupRefreshTest(t)
		defer client.Close()

		reqBody := request.RefreshTokenRequestDTO{
			RefreshToken: "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
