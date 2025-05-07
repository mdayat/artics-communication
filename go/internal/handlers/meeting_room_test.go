package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/mdayat/artics-communication/go/internal/dtos"
)

func TestMeetingRoomHandlers(t *testing.T) {
	t.Run("GetAvailableMeetingRooms/Success", func(t *testing.T) {
		res, err := testClient.Get(fmt.Sprintf("%s/meeting-rooms/available", testServer.URL))
		if err != nil {
			t.Fatalf("wasn't expecting error, got: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var tasks []dtos.MeetingRoomWithTimeSlotsResponse
		if err = json.NewDecoder(res.Body).Decode(&tasks); err != nil {
			t.Fatalf("unexpected response body: %v", res)
		}
	})
}
