package graph

import (
	"context"
	"errors"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
)

func mapAuthResponse(a gotwitter.AuthResponse) *AuthResponse {
	return &AuthResponse{
		AccessToken: a.AccessToken,
		User:        mapUser(a.User),
	}
}

func (r *mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	res, err := r.AuthService.Register(ctx, gotwitter.RegisterInput{
		Username:        input.Username,
		Email:           input.Email,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, gotwitter.ErrValidation) || errors.Is(err, gotwitter.ErrEmailTaken) || errors.Is(err, gotwitter.ErrUsernameTaken):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}

	}
	return mapAuthResponse(res), nil
}

func (r *mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	res, err := r.AuthService.Login(ctx, gotwitter.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, gotwitter.ErrValidation) || errors.Is(err, gotwitter.ErrBadCredentials):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}

	}
	return mapAuthResponse(res), nil
}
