package handlers

import (
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

func (mr meetingRoom) GetAvailableMeetingRooms(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := log.Ctx(ctx).With().Logger()

	availableMeetingRooms, err := retryutil.RetryWithData(func() ([]repository.MeetingRoom, error) {
		return mr.configs.Db.Queries.SelectAvailableMeetingRooms(ctx)
	})

	if err != nil {
		logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to get available meeting rooms")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resBody := make([]dtos.MeetingRoom, 0, len(availableMeetingRooms))
	for _, meetingRoom := range availableMeetingRooms {
		resBody = append(resBody, dtos.MeetingRoom{
			Id:        meetingRoom.ID.String(),
			Name:      meetingRoom.Name,
			CreatedAt: meetingRoom.CreatedAt.Time.Format(time.RFC3339),
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

	logger.Info().Int("status_code", http.StatusOK).Msg("successfully get available meeting rooms")
}
