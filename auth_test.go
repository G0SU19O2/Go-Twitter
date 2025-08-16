package gotwitter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	input := RegisterInput{
		Username:        " validuser ",
		Email:           " Test@example.com ",
		Password:        "password",
		ConfirmPassword: "password",
	}
	want := RegisterInput{
		Username:        "validuser",
		Email:           "test@example.com",
		Password:        "password",
		ConfirmPassword: "password",
	}
	input.Sanitize()
	require.Equal(t, want, input)
}

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid input",
			input: RegisterInput{
				Username:        "validuser",
				Email:           "test@example.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				Username:        "validuser",
				Email:           "invalid-email",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "short password",
			input: RegisterInput{
				Username:        "validuser",
				Email:           "test@example.com",
				Password:        "short",
				ConfirmPassword: "short",
			},
			err: ErrValidation,
		},
		{
			name: "username too short",
			input: RegisterInput{
				Email:           "test@example.com",
				Password:        "password",
				Username:        "a",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		}, {
			name: "passwords do not match",
			input: RegisterInput{
				Email:           "test@example.com",
				Password:        "password",
				Username:        "username",
				ConfirmPassword: "differentpassword",
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoginInput_Sanitize(t *testing.T) {
	input := LoginInput{
		Email:           " Test@example.com ",
		Password:        "password",
	}
	want := LoginInput{
		Email:           "test@example.com",
		Password:        "password",
	}
	input.Sanitize()
	require.Equal(t, want, input)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input LoginInput
		err   error
	}{
		{
			name: "valid input",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: LoginInput{
				Email:    "invalid-email",
				Password: "password",
			},
			err: ErrValidation,
		},
		{
			name: "short password",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "short",
			},
			err: ErrValidation,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}