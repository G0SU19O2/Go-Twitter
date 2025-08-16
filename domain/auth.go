package domain

import (
	"context"
	"errors"
	"fmt"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/G0SU19O2/Go-Twitter/faker"
	"golang.org/x/crypto/bcrypt"
)

var passwordCost = bcrypt.DefaultCost

type AuthService struct {
	UserRepo gotwitter.UserRepo
}

func NewAuthService(ur gotwitter.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input gotwitter.RegisterInput) (gotwitter.AuthResponse, error) {
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return gotwitter.AuthResponse{}, err
	}
	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, gotwitter.ErrNotFound) {
		return gotwitter.AuthResponse{}, gotwitter.ErrUsernameTaken
	}
	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, gotwitter.ErrNotFound) {
		return gotwitter.AuthResponse{}, gotwitter.ErrEmailTaken
	}
	user := gotwitter.User{
		Email:    input.Email,
		Username: input.Username,
	}
	user.Password = faker.Password
	user, err := as.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return gotwitter.AuthResponse{}, fmt.Errorf("failed to create user: %w", err)
	}
	return gotwitter.AuthResponse{AccessToken: "some-jwt-token", User: user}, nil
}

func (as *AuthService) Login(ctx context.Context, input gotwitter.LoginInput) (gotwitter.AuthResponse, error) {
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return gotwitter.AuthResponse{}, err
	}
	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, gotwitter.ErrNotFound):
			return gotwitter.AuthResponse{}, gotwitter.ErrBadCredentials
		default:
			return gotwitter.AuthResponse{}, err
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return gotwitter.AuthResponse{}, gotwitter.ErrBadCredentials
	}
	return gotwitter.AuthResponse{AccessToken: "some-jwt-token", User: user}, nil
}
