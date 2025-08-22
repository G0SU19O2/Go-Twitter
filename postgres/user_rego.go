package postgres

import (
	"context"
	"fmt"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	DB *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user gotwitter.User) (gotwitter.User, error) {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return gotwitter.User{}, fmt.Errorf("error begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)
	user, err = createUser(ctx, tx, user)
	if err != nil {
		return gotwitter.User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return gotwitter.User{}, fmt.Errorf("error commit transaction: %v", err)
	}
	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user gotwitter.User) (gotwitter.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *`
	u := gotwitter.User{}
	if err := pgxscan.Get(ctx, tx, &u, query, user.Username, user.Email, user.Password); err != nil {
		return gotwitter.User{}, fmt.Errorf("error insert: %v", err)
	}
	return u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (gotwitter.User, error) {
	query := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	u := gotwitter.User{}
	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return gotwitter.User{}, gotwitter.ErrNotFound
		}
		return gotwitter.User{}, fmt.Errorf("error select: %v", err)
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (gotwitter.User, error) {
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	u := gotwitter.User{}
	if err := pgxscan.Get(ctx, r.DB.Pool, &u, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return gotwitter.User{}, gotwitter.ErrNotFound
		}
		return gotwitter.User{}, fmt.Errorf("error select: %v", err)
	}
	return u, nil
}

