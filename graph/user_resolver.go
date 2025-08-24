package graph

import (
	"context"
	"fmt"
)

// schema.resolvers.go (this file doesn't exist yet)
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	// Your business logic here
	// e.g., get user from JWT token, query database
	panic(fmt.Errorf("not implemented"))
}
