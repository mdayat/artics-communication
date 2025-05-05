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

type ReservationHandler interface {
	CancelReservation(res http.ResponseWriter, req *http.Request)
	GetReservations(res http.ResponseWriter, req *http.Request)
}

type reservation struct {
	configs configs.Configs
}

func NewReservationHandler(configs configs.Configs) ReservationHandler {
	return &reservation{
		configs: configs,
	}
}

func (r reservation) CancelReservation(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	accountRole := ctx.Value(accountRoleKey{}).(string)
	if accountRole != dtos.AdminRole {
		logger.Error().Err(errors.New("insufficient permissions to cancel reservation")).Caller().Int("status_code", http.StatusForbidden).Send()
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	reservationId := chi.URLParam(req, "reservationId")
	reservation, err := retryutil.RetryWithData(func() (repository.Reservation, error) {
		reservationUUID, err := uuid.Parse(reservationId)
		if err != nil {
			return repository.Reservation{}, fmt.Errorf("failed to parse reservation Id to UUID: %w", err)
		}

		return r.configs.Db.Queries.CancelReservation(ctx, pgtype.UUID{Bytes: reservationUUID, Valid: true})
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusNotFound).Msg("reservation not found")
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to cancel reservation")
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	resBody := dtos.CanceledReservationResponse{
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

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully canceled reservation")
}

func (r reservation) GetReservations(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	accountRole := ctx.Value(accountRoleKey{}).(string)
	if accountRole != dtos.AdminRole {
		logger.Error().Err(errors.New("insufficient permissions to get reservations")).Caller().Int("status_code", http.StatusForbidden).Send()
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	reservations, err := retryutil.RetryWithData(func() ([]repository.SelectReservationsRow, error) {
		return r.configs.Db.Queries.SelectReservations(ctx)
	})

	if err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to select reservations")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resBody := make([]dtos.ReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		var canceledAt *string
		if reservation.CanceledAt.Valid {
			formatted := reservation.CanceledAt.Time.Format(time.RFC3339)
			canceledAt = &formatted
		}

		resBody = append(resBody, dtos.ReservationResponse{
			Id:         reservation.ID.String(),
			Canceled:   reservation.Canceled,
			CanceledAt: canceledAt,
			ReservedAt: reservation.ReservedAt.Time.Format(time.RFC3339),
			User: dtos.UserResponse{
				Id:        reservation.User.ID.String(),
				Email:     reservation.User.Email,
				Name:      reservation.User.Name,
				Role:      reservation.User.Role,
				CreatedAt: reservation.User.CreatedAt.Time.Format(time.RFC3339),
			},
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

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get reservations")
}
