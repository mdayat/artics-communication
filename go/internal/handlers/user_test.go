package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/mdayat/artics-communication/go/internal/dtos"
)

func TestUserHandlers(t *testing.T) {
	t.Run("GetUser/Success", func(t *testing.T) {
		res, err := testClient.Get(fmt.Sprintf("%s/users/me", testServer.URL))
		if err != nil {
			t.Fatalf("wasn't expecting error, got: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var user dtos.UserResponse
		if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
			t.Fatalf("unexpected response body: %v", res)
		}
	})

	var reservation dtos.UserReservationResponse
	t.Run("GetUserReservations/Success", func(t *testing.T) {
		res, err := testClient.Get(fmt.Sprintf("%s/users/me/reservations", testServer.URL))
		if err != nil {
			t.Fatalf("wasn't expecting error, got: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		var reservations []dtos.UserReservationResponse
		if err = json.NewDecoder(res.Body).Decode(&reservations); err != nil {
			t.Fatalf("unexpected response body: %v", res)
		}

		if len(reservations) == 0 {
			t.Log("No reservations found, skipping the rest of test")
			return
		}

		for _, item := range reservations {
			if item.Canceled {
				continue
			}

			reservation = item
		}

		if reservation.Id == "" {
			t.Log("No active reservation found, skipping the rest of test")
			return
		}
	})

	ctx := context.TODO()
	cancelReservationTable := []struct {
		name           string
		reservationId  string
		expectedStatus int
		expectedResult dtos.ReservationResponse
	}{
		{
			name:           "CancelUserReservation/Success",
			reservationId:  reservation.Id,
			expectedStatus: http.StatusOK,
			expectedResult: dtos.ReservationResponse{
				Id:       reservation.Id,
				Canceled: true,
			},
		},
		{
			name:           "CancelUserReservation/Not_Found",
			reservationId:  uuid.NewString(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, v := range cancelReservationTable {
		t.Run(v.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/users/me/reservations//%s", testServer.URL, v.reservationId)
			req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
			if err != nil {
				t.Fatalf("wasn't expecting error, got: %v", err)
			}

			res, err := testClient.Do(req)
			if err != nil {
				t.Fatalf("wasn't expecting error, got: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != v.expectedStatus {
				t.Fatalf("expected status %d, got %d", v.expectedStatus, res.StatusCode)
			}

			if v.expectedStatus == http.StatusOK {
				var canceledReservation dtos.ReservationResponse
				if err := json.NewDecoder(res.Body).Decode(&canceledReservation); err != nil {
					t.Fatalf("unexpected response body: %v", res)
				}

				diff := cmp.Diff(v.expectedResult, canceledReservation, cmpopts.IgnoreFields(
					dtos.ReservationResponse{},
					"UserId", "MeetingRoomId", "TimeSlotId", "CanceledAt", "ReservedAt",
				))

				if diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
