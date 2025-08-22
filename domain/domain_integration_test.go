//go:build integration
// +build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/G0SU19O2/Go-Twitter/config"
	"github.com/G0SU19O2/Go-Twitter/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	conf        *config.Config
	db          *postgres.DB
	authService gotwitter.AuthService
	userRepo    gotwitter.UserRepo
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")
	passwordCost = bcrypt.MinCost
	ctx := context.Background()
	conf = config.New()
	db = postgres.New(ctx, conf)
	defer db.Close()
	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}
	userRepo = postgres.NewUserRepo(db)
	authService = NewAuthService(userRepo)

	os.Exit(m.Run())
}
