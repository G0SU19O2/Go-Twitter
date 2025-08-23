//go:build integration
// +build integration

package domain

import (
	"context"
	"testing"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/G0SU19O2/Go-Twitter/faker"
	"github.com/G0SU19O2/Go-Twitter/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := gotwitter.RegisterInput{
		Username:        faker.UserName(),
		Email:           faker.Email(),
		Password:        "password",
		ConfirmPassword: "password",
	}
	t.Run("can register a user", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(ctx, t, db)
		res, err := authService.Register(ctx, validInput)
		require.NoError(t, err)
		require.Equal(t, res.User.Username, validInput.Username)
		require.Equal(t, res.User.Email, validInput.Email)
		require.NotEmpty(t, res.AccessToken)
		require.NotEqual(t, res.User.Password, validInput.Password)
	})

	t.Run("existing username", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(ctx, t, db)
		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)
		_, err = authService.Register(ctx, gotwitter.RegisterInput{
			Username:   validInput.Username,
			Email:      "test2@example.com",
			Password:   validInput.Password,
			ConfirmPassword: validInput.ConfirmPassword,
		})
		require.ErrorIs(t, err, gotwitter.ErrUsernameTaken)
	})
	t.Run("existing email", func(t *testing.T) {
		ctx := context.Background()
		defer test_helpers.TeardownDB(ctx, t, db)
		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)
		_, err = authService.Register(ctx, gotwitter.RegisterInput{
			Username:   "test3",
			Email:      validInput.Email,
			Password:   validInput.Password,
			ConfirmPassword: validInput.ConfirmPassword,
		})
		require.ErrorIs(t, err, gotwitter.ErrEmailTaken)
	})
}
