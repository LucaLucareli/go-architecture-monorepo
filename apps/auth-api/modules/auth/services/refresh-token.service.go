package services

import (
	"auth-api/modules/auth/dto/io"
	"context"
	"shared/application/auth"
)

type RefreshTokenService struct {
	authService *auth.AuthService
}

func NewRefreshTokenService(authService *auth.AuthService) *RefreshTokenService {
	return &RefreshTokenService{authService: authService}
}

func (s *RefreshTokenService) Execute(
	ctx context.Context,
	input io.RefreshTokenInputDTO,
) (*io.RefreshTokenOutputDTO, error) {

	authResponse, err := s.authService.RefreshToken(
		ctx,
		input.RefreshToken,
	)

	if err != nil {
		return nil, err
	}

	return &io.RefreshTokenOutputDTO{
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
	}, nil
}
