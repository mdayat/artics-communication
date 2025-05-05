package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/internal/httputil"
	"github.com/mdayat/artics-communication/go/internal/retryutil"
	"github.com/mdayat/artics-communication/go/repository"
	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	GetUser(res http.ResponseWriter, req *http.Request)
	GetUserReservations(res http.ResponseWriter, req *http.Request)
	UpdateUserReservation(res http.ResponseWriter, req *http.Request)
}

type user struct {
	configs configs.Configs
}

func NewUserHandler(configs configs.Configs) UserHandler {
	return &user{
		configs: configs,
	}
}

func (u user) GetUser(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	userId := ctx.Value(userIdKey{}).(string)
	user, err := retryutil.RetryWithData(func() (repository.User, error) {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return repository.User{}, fmt.Errorf("failed to parse user Id to UUID: %w", err)
		}

		return u.configs.Db.Queries.SelectUser(ctx, pgtype.UUID{Bytes: userUUID, Valid: true})
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusNotFound).Msg("user not found")
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to select user")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.UserResponse{
		Id:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}

	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusOK,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get user")
}

func (u user) GetUserReservations(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	userId := ctx.Value(userIdKey{}).(string)
	reservations, err := retryutil.RetryWithData(func() ([]repository.SelectUserReservationsRow, error) {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user Id to UUID: %w", err)
		}

		return u.configs.Db.Queries.SelectUserReservations(ctx, pgtype.UUID{Bytes: userUUID, Valid: true})
	})

	if err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to select user reservations")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resBody := make([]dtos.UserReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		resBody = append(resBody, dtos.UserReservationResponse{
			Id:         reservation.ID.String(),
			Status:     reservation.Status,
			ReservedAt: reservation.ReservedAt.Time.Format(time.RFC3339),
			MeetingRoom: dtos.MeetingRoom{
				Id:        reservation.MeetingRoom.ID.String(),
				Name:      reservation.MeetingRoom.Name,
				CreatedAt: reservation.MeetingRoom.CreatedAt.Time.Format(time.RFC3339),
			},
			TimeSlot: dtos.TimeSlot{
				Id:        reservation.TimeSlot.ID.String(),
				StartDate: reservation.TimeSlot.StartDate.Time.Format(time.RFC3339),
				EndDate:   reservation.TimeSlot.EndDate.Time.Format(time.RFC3339),
				CreatedAt: reservation.TimeSlot.CreatedAt.Time.Format(time.RFC3339),
			},
		})
	}

	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusOK,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get user reservations")
}

func (u user) UpdateUserReservation(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	var reqBody dtos.UpdateUserReservationRequest
	if err := httputil.DecodeAndValidate(req, u.configs.Validate, &reqBody); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusBadRequest).Msg("invalid request body")
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	reservationId := chi.URLParam(req, "reservationId")
	userId := ctx.Value(userIdKey{}).(string)

	reservation, err := retryutil.RetryWithData(func() (repository.Reservation, error) {
		reservationUUID, err := uuid.Parse(reservationId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse reservation Id to UUID: %w", err)
		}

		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse user Id to UUID: %w", err)
		}

		return u.configs.Db.Queries.UpdateReservation(ctx, repository.UpdateReservationParams{
			ID:     pgtype.UUID{Bytes: reservationUUID, Valid: true},
			UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
			Status: reqBody.Status,
		})
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusNotFound).Msg("reservation not found")
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to update user reservation")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.UserReservationResponse{
		Id:         reservation.ID.String(),
		Status:     reservation.Status,
		ReservedAt: reservation.ReservedAt.Time.Format(time.RFC3339),
	}

	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusOK,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully updated user reservation")
}
