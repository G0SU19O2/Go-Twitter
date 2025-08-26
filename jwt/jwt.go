package jwt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"github.com/G0SU19O2/Go-Twitter/config"
	"github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
)

var signatureType = jwa.HS256

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (gotwitter.AuthToken, error) {
	token, err := jwtGo.ParseRequest(
		r,
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)
	if err != nil {
		return gotwitter.AuthToken{}, gotwitter.ErrInvalidAccessToken
	}
	return buildToken(token), nil
}

func buildToken(token jwtGo.Token) gotwitter.AuthToken {
	return gotwitter.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (gotwitter.AuthToken, error) {
	token, err := jwtGo.Parse(
		[]byte(payload),
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)
	if err != nil {
		return gotwitter.AuthToken{}, gotwitter.ErrInvalidAccessToken
	}
	return buildToken(token), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user gotwitter.User) (string, error) {
	t := jwtGo.New()
	if err := setDefaultToken(t, user, gotwitter.AccessTokenLifeTime, ts.Conf); err != nil {
		return "", err
	}
	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user gotwitter.User, tokenID string) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, gotwitter.RefreshTokenLifeTime, ts.Conf); err != nil {
		return "", err
	}
	if err := t.Set(jwtGo.JwtIDKey, tokenID); err != nil {
		return "", fmt.Errorf("error set jwt jti: %v", err)
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func setDefaultToken(t jwtGo.Token, user gotwitter.User, lifetime time.Duration, conf *config.Config) error {
	if err := t.Set(jwtGo.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("error set jwt sub: %v", err)
	}
	if err := t.Set(jwtGo.IssuerKey, conf.JWT.Issuer); err != nil {
		return fmt.Errorf("error set jwt iss: %v", err)
	}
	if err := t.Set(jwtGo.IssuedAtKey, time.Now().Unix()); err != nil {
		return fmt.Errorf("error set jwt iat: %v", err)
	}
	if err := t.Set(jwtGo.ExpirationKey, time.Now().Add(lifetime).Unix()); err != nil {
		return fmt.Errorf("error set jwt exp: %v", err)
	}
	return nil
}
