package services

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/internal/retryutil"
	"github.com/mdayat/artics-communication/go/repository"
)

type AuthServicer interface {
	RegisterUser(ctx context.Context, arg RegisterUserParams) (repository.User, error)
}

type auth struct {
	configs configs.Configs
}

func NewAuthService(configs configs.Configs) AuthServicer {
	return &auth{
		configs: configs,
	}
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
