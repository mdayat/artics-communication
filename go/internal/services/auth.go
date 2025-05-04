package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/internal/retryutil"
	"github.com/mdayat/artics-communication/go/repository"
)

type AuthServicer interface {
	CreateAccessToken(claims AccessTokenClaims) (string, error)
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
	RegisterUser(ctx context.Context, arg RegisterUserParams) (repository.User, error)
	AuthenticateUser(ctx context.Context, arg AuthenticateUserParams) (repository.User, error)
}

type auth struct {
	configs configs.Configs
}

func NewAuthService(configs configs.Configs) AuthServicer {
	return &auth{
		configs: configs,
	}
}

type AccessTokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (a auth) CreateAccessToken(claims AccessTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.configs.Env.SecretKey))
}

func (a auth) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(a.configs.Env.SecretKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(a.configs.Env.OriginURL),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, errors.New("invalid access token claims")
	}

	return claims, nil
}

type RegisterUserParams struct {
	Username string
	Email    string
	Password string
}

func (a auth) RegisterUser(ctx context.Context, arg RegisterUserParams) (repository.User, error) {
	hashedPassword, err := argon2id.CreateHash(arg.Password, argon2id.DefaultParams)
	if err != nil {
		return repository.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userUUID := uuid.New()
	return retryutil.RetryWithData(func() (repository.User, error) {
		return a.configs.Db.Queries.InsertUser(ctx, repository.InsertUserParams{
			ID:       pgtype.UUID{Bytes: userUUID, Valid: true},
			Name:     arg.Username,
			Email:    arg.Email,
			Password: hashedPassword,
			Role:     dtos.UserRole,
		})
	})
}

type AuthenticateUserParams struct {
	Email    string
	Password string
}

func (a auth) AuthenticateUser(ctx context.Context, arg AuthenticateUserParams) (repository.User, error) {
	user, err := retryutil.RetryWithData(func() (repository.User, error) {
		return a.configs.Db.Queries.SelectUserByEmail(ctx, arg.Email)
	})

	if err != nil {
		return repository.User{}, fmt.Errorf("failed to select user by email: %w", err)
	}

	match, err := argon2id.ComparePasswordAndHash(arg.Password, user.Password)
	if err != nil {
		return repository.User{}, fmt.Errorf("failed to compare password: %w", err)
	}

	if !match {
		return repository.User{}, fmt.Errorf("wrong password: %w", pgx.ErrNoRows)
	}

	return user, nil
}
