package graph

import (
	"context"
	"fmt"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
)

func mapUser(u gotwitter.User) *User {
	return &User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		CreateAt: u.CreatedAt,
	}
}

// schema.resolvers.go (this file doesn't exist yet)
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	// Your business logic here
	// e.g., get user from JWT token, query database
	panic(fmt.Errorf("not implemented"))
}
