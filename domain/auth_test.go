package domain

import (
	"context"
	"errors"
	"testing"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/G0SU19O2/Go-Twitter/faker"
	"github.com/G0SU19O2/Go-Twitter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	validInput := gotwitter.RegisterInput{
		Username:        "testuser",
		Email:           "test@example.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(gotwitter.User{ID: "1902", Username: validInput.Username, Email: validInput.Email}, nil)
		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)
		userRepo.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(gotwitter.User{}, nil)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, gotwitter.ErrUsernameTaken)
		userRepo.AssertNotCalled(t, "CreateUser")
		userRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrEmailTaken)
		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, gotwitter.ErrEmailTaken)
		userRepo.AssertNotCalled(t, "CreateUser")
		userRepo.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		userRepo.On("GetByUsername", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(gotwitter.User{}, errors.New("create user error"))
		res, err := service.Register(ctx, validInput)
		require.Error(t, err)
		require.Empty(t, res.AccessToken)
		require.Empty(t, res.User.Email)
		require.Empty(t, res.User.Username)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		_, err := service.Register(ctx, gotwitter.RegisterInput{})
		require.ErrorIs(t, err, gotwitter.ErrValidation)
		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "CreateUser")
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	validInput := gotwitter.LoginInput{
		Email:    "test@example.com",
		Password: "password",
	}
	t.Run("can login", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{ID: "1902", Email: validInput.Email, Password: faker.Password}, nil)
		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, validInput)
		require.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("wrongpassword"), passwordCost)
		require.NoError(t, err)
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{ID: "1902", Email: validInput.Email, Password: string(hashedPassword)}, nil)
		service := NewAuthService(userRepo)
		_, err = service.Login(ctx, validInput)
		require.ErrorIs(t, err, gotwitter.ErrBadCredentials)
		userRepo.AssertExpectations(t)
	})

	t.Run("email not found", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{}, gotwitter.ErrNotFound)
		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, validInput)
		require.ErrorIs(t, err, gotwitter.ErrBadCredentials)
		userRepo.AssertExpectations(t)
	})

	t.Run("get user by email error", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(gotwitter.User{}, errors.New("some error"))
		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, validInput)
		require.Error(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()
		userRepo := &mocks.MockUserRepo{}
		service := NewAuthService(userRepo)
		_, err := service.Login(ctx, gotwitter.LoginInput{Email: "bob", Password: ""})
		require.ErrorIs(t, err, gotwitter.ErrValidation)
		userRepo.AssertExpectations(t)
	})
}
