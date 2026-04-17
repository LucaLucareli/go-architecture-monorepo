package services

import (
	"auth-api/modules/auth/dto/io"
	"context"
	"shared/application/auth"
)

type LoginService struct {
	authService *auth.AuthService
}

func NewLoginService(authService *auth.AuthService) *LoginService {
	return &LoginService{authService: authService}
}

func (s *LoginService) Execute(
	ctx context.Context,
	input io.LoginInputDTO,
) (*io.LoginOutputDTO, error) {

	authResponse, err := s.authService.Login(
		ctx,
		input.Document,
		input.Password,
	)

	if err != nil {
		return nil, err
	}

	return &io.LoginOutputDTO{
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
	}, nil
}
