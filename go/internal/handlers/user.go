package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	CancelUserReservation(res http.ResponseWriter, req *http.Request)
	CreateUserReservation(res http.ResponseWriter, req *http.Request)
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
		var canceledAt *string
		if reservation.CanceledAt.Valid {
			formatted := reservation.CanceledAt.Time.Format(time.RFC3339)
			canceledAt = &formatted
		}

		resBody = append(resBody, dtos.UserReservationResponse{
			Id:         reservation.ID.String(),
			Canceled:   reservation.Canceled,
			CanceledAt: canceledAt,
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

func (u user) CancelUserReservation(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

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

		return u.configs.Db.Queries.CancelUserReservation(ctx, repository.CancelUserReservationParams{
			ID:     pgtype.UUID{Bytes: reservationUUID, Valid: true},
			UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
		})
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusNotFound).Msg("reservation not found")
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to cancel user reservation")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.ReservationResponse{
		Id:            reservation.ID.String(),
		UserId:        reservation.UserID.String(),
		MeetingRoomId: reservation.MeetingRoomID.String(),
		TimeSlotId:    reservation.TimeSlotID.String(),
		Canceled:      reservation.Canceled,
		ReservedAt:    reservation.ReservedAt.Time.Format(time.RFC3339),
	}

	if reservation.CanceledAt.Valid {
		formatted := reservation.CanceledAt.Time.Format(time.RFC3339)
		resBody.CanceledAt = &formatted
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

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully canceled user reservation")
}

func (u user) CreateUserReservation(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	var reqBody dtos.CreateReservationRequest
	if err := httputil.DecodeAndValidate(req, u.configs.Validate, &reqBody); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusBadRequest).Msg("invalid request body")
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId := ctx.Value(userIdKey{}).(string)
	reservation, err := retryutil.RetryWithData(func() (repository.Reservation, error) {
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse user Id to UUID: %w", err)
		}

		meetingRoomUUID, err := uuid.Parse(reqBody.MeetingRoomId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse meeting room Id to UUID: %w", err)
		}

		timeSlotUUID, err := uuid.Parse(reqBody.TimeSlotId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse time slot Id to UUID: %w", err)
		}

		return u.configs.Db.Queries.InsertReservation(ctx, repository.InsertReservationParams{
			ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:        pgtype.UUID{Bytes: userUUID, Valid: true},
			MeetingRoomID: pgtype.UUID{Bytes: meetingRoomUUID, Valid: true},
			TimeSlotID:    pgtype.UUID{Bytes: timeSlotUUID, Valid: true},
		})
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusConflict).Msg("time slot already reserved")
			http.Error(res, http.StatusText(http.StatusConflict), http.StatusConflict)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to create user reservation")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.ReservationResponse{
		Id:            reservation.ID.String(),
		UserId:        reservation.UserID.String(),
		MeetingRoomId: reservation.MeetingRoomID.String(),
		TimeSlotId:    reservation.TimeSlotID.String(),
		Canceled:      reservation.Canceled,
		ReservedAt:    reservation.ReservedAt.Time.Format(time.RFC3339),
	}

	if reservation.CanceledAt.Valid {
		formatted := reservation.CanceledAt.Time.Format(time.RFC3339)
		resBody.CanceledAt = &formatted
	}

	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusCreated,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusCreated).Msg("successfully created user reservation")
}
