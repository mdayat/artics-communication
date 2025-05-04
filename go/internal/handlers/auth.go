package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/internal/httputil"
	"github.com/mdayat/artics-communication/go/internal/services"
	"github.com/rs/zerolog/log"
)

type AuthHandler interface {
	Register(res http.ResponseWriter, req *http.Request)
}

type auth struct {
	configs configs.Configs
	service services.AuthServicer
}

func NewAuthHandler(configs configs.Configs, service services.AuthServicer) AuthHandler {
	return &auth{
		configs: configs,
		service: service,
	}
}

func (a auth) Register(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.With().Ctx(ctx).Logger()

	var reqBody dtos.RegisterRequest
	if err := httputil.DecodeAndValidate(req, a.configs.Validate, &reqBody); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusBadRequest).Msg("invalid request body")
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := a.service.RegisterUser(ctx, services.RegisterUserParams{
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusConflict).Msg("user already exist")
			http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to register user")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.UserResponse{
		Id:        result.ID.String(),
		Email:     result.Email,
		Name:      result.Name,
		Role:      result.Role,
		CreatedAt: result.CreatedAt.Time.Format(time.RFC3339),
	}

	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusCreated,
		ResBody:    resBody,
	}

	res.Header().Set("Location", fmt.Sprintf("%s/users/me", a.configs.Env.OriginURL))
	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusCreated).Msg("successfully registered user")
}
