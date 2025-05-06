package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/dtos"
	"github.com/mdayat/artics-communication/go/internal/httputil"
	"github.com/mdayat/artics-communication/go/internal/retryutil"
	"github.com/mdayat/artics-communication/go/repository"
	"github.com/rs/zerolog/log"
)

type MeetingRoomHandler interface {
	GetMeetingRooms(res http.ResponseWriter, req *http.Request)
	GetAvailableMeetingRooms(res http.ResponseWriter, req *http.Request)
}

type meetingRoom struct {
	configs configs.Configs
}

func NewMeetingRoomHandler(configs configs.Configs) MeetingRoomHandler {
	return &meetingRoom{
		configs: configs,
	}
}

func (mr meetingRoom) GetMeetingRooms(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	accountRole := ctx.Value(accountRoleKey{}).(string)
	if accountRole != dtos.AdminRole {
		logger.Error().Err(errors.New("insufficient permissions to get meeting rooms")).Caller().Int("status_code", http.StatusForbidden).Send()
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	meetingRooms, err := retryutil.RetryWithData(func() ([]repository.SelectMeetingRoomsRow, error) {
		return mr.configs.Db.Queries.SelectMeetingRooms(ctx)
	})

	if err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to get meeting rooms")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resBody := make([]dtos.MeetingRoomWithTimeSlotsResponse, 0)
	currentMeetingRoom := dtos.MeetingRoomWithTimeSlotsResponse{}

	for _, meetingRoom := range meetingRooms {
		if currentMeetingRoom.Id == "" {
			currentMeetingRoom = dtos.MeetingRoomWithTimeSlotsResponse{
				Id:        meetingRoom.ID.String(),
				Name:      meetingRoom.Name,
				CreatedAt: meetingRoom.CreatedAt.Time.Format(time.RFC3339),
				TimeSlots: make([]dtos.TimeSlot, 0),
			}
		} else if meetingRoom.Name != currentMeetingRoom.Name {
			resBody = append(resBody, currentMeetingRoom)
			currentMeetingRoom = dtos.MeetingRoomWithTimeSlotsResponse{
				Id:        meetingRoom.ID.String(),
				Name:      meetingRoom.Name,
				CreatedAt: meetingRoom.CreatedAt.Time.Format(time.RFC3339),
				TimeSlots: make([]dtos.TimeSlot, 0),
			}
		}

		currentMeetingRoom.TimeSlots = append(currentMeetingRoom.TimeSlots, dtos.TimeSlot{
			Id:        meetingRoom.TimeSlot.ID.String(),
			StartDate: meetingRoom.TimeSlot.StartDate.Time.Format(time.RFC3339),
			EndDate:   meetingRoom.TimeSlot.EndDate.Time.Format(time.RFC3339),
			CreatedAt: meetingRoom.TimeSlot.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	resBody = append(resBody, currentMeetingRoom)
	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusOK,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get meeting rooms")
}

func (mr meetingRoom) GetAvailableMeetingRooms(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	availableMeetingRooms, err := retryutil.RetryWithData(func() ([]repository.SelectAvailableMeetingRoomsRow, error) {
		return mr.configs.Db.Queries.SelectAvailableMeetingRooms(ctx)
	})

	if err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to get available meeting rooms")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resBody := make([]dtos.MeetingRoomWithTimeSlotsResponse, 0)
	currentMeetingRoom := dtos.MeetingRoomWithTimeSlotsResponse{}

	for _, meetingRoom := range availableMeetingRooms {
		if currentMeetingRoom.Id == "" {
			currentMeetingRoom = dtos.MeetingRoomWithTimeSlotsResponse{
				Id:        meetingRoom.ID.String(),
				Name:      meetingRoom.Name,
				CreatedAt: meetingRoom.CreatedAt.Time.Format(time.RFC3339),
				TimeSlots: make([]dtos.TimeSlot, 0),
			}
		} else if meetingRoom.Name != currentMeetingRoom.Name {
			resBody = append(resBody, currentMeetingRoom)
			currentMeetingRoom = dtos.MeetingRoomWithTimeSlotsResponse{
				Id:        meetingRoom.ID.String(),
				Name:      meetingRoom.Name,
				CreatedAt: meetingRoom.CreatedAt.Time.Format(time.RFC3339),
				TimeSlots: make([]dtos.TimeSlot, 0),
			}
		}

		currentMeetingRoom.TimeSlots = append(currentMeetingRoom.TimeSlots, dtos.TimeSlot{
			Id:        meetingRoom.TimeSlot.ID.String(),
			StartDate: meetingRoom.TimeSlot.StartDate.Time.Format(time.RFC3339),
			EndDate:   meetingRoom.TimeSlot.EndDate.Time.Format(time.RFC3339),
			CreatedAt: meetingRoom.TimeSlot.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	resBody = append(resBody, currentMeetingRoom)
	params := httputil.SendSuccessResponseParams{
		StatusCode: http.StatusOK,
		ResBody:    resBody,
	}

	if err := httputil.SendSuccessResponse(res, params); err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to send success response")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get available meeting rooms")
}
