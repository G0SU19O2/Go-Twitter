package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/G0SU19O2/Go-Twitter/config"
	"github.com/G0SU19O2/Go-Twitter/domain"
	"github.com/G0SU19O2/Go-Twitter/graph"
	"github.com/G0SU19O2/Go-Twitter/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()
	conf := config.New()
	config.LoadEnv(".env")
	db := postgres.New(ctx, conf)
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	userRepo := postgres.NewUserRepo(db)
	authService := domain.NewAuthService(userRepo)
	router.Handle("/", playground.Handler("Twitter clone", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService: authService,
				},
			},
		),
	))
	log.Fatal(http.ListenAndServe(":8080", router))
}
