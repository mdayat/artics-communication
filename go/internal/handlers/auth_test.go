package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mdayat/artics-communication/go/internal/dtos"
)

func TestAuthHandlers(t *testing.T) {
	var createdUser dtos.UserResponse
	registerTable := []struct {
		name           string
		reqBody        string
		expectedStatus int
		expectedResult dtos.UserResponse
	}{
		{
			name:           "Register/Success",
			reqBody:        `{"username": "Marie", "email": "marie@gmail.com", "password": "password"}`,
			expectedStatus: http.StatusCreated,
			expectedResult: dtos.UserResponse{
				Name:  "Marie",
				Email: "marie@gmail.com",
			},
		},
		{
			name:           "Register/Bad_Username",
			reqBody:        `{"username": "X", "email": "marie@gmail.com", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register/Bad_Email",
			reqBody:        `{"username": "Marie", "email": "marie@marie", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register/Bad_Password",
			reqBody:        `{"username": "Marie", "email": "marie@gmail.com", "password": "Marie"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register/Conflict",
			reqBody:        `{"username": "Marie", "email": "marie@gmail.com", "password": "password"}`,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, v := range registerTable {
		t.Run(v.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/auth/register", testServer.URL)
			res, err := testClient.Post(url, "application/json", bytes.NewBuffer([]byte(v.reqBody)))
			if err != nil {
				t.Fatalf("wasn't expecting error, got: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != v.expectedStatus {
				t.Fatalf("expected status %d, got %d", v.expectedStatus, res.StatusCode)
			}

			if v.expectedStatus == http.StatusCreated {
				if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
					t.Fatalf("unexpected response body: %v", res)
				}

				diff := cmp.Diff(v.expectedResult, createdUser, cmpopts.IgnoreFields(dtos.UserResponse{}, "Id", "Role", "CreatedAt"))
				if diff != "" {
					t.Error(diff)
				}
			}
		})
	}

	loginTable := []struct {
		name           string
		reqBody        string
		expectedStatus int
		expectedResult dtos.UserResponse
	}{
		{
			name:           "Login/Success",
			reqBody:        `{"email": "marie@gmail.com", "password": "password"}`,
			expectedStatus: http.StatusOK,
			expectedResult: dtos.UserResponse{
				Name:  "Marie",
				Email: "marie@gmail.com",
			},
		},
		{
			name:           "Login/Bad_Email",
			reqBody:        `{"email": "marie@marie", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Login/Bad_Password",
			reqBody:        `{"email": "marie@gmail.com", "password": "Marie"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Login/Not_Found",
			reqBody:        `{"email": "marie@gmail.com", "password": "wrong password"}`,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, v := range loginTable {
		t.Run(v.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/auth/login", testServer.URL)
			res, err := testClient.Post(url, "application/json", bytes.NewBuffer([]byte(v.reqBody)))
			if err != nil {
				t.Fatalf("wasn't expecting error, got: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != v.expectedStatus {
				t.Fatalf("expected status %d, got %d", v.expectedStatus, res.StatusCode)
			}

			if v.expectedStatus == http.StatusCreated {
				if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
					t.Fatalf("unexpected response body: %v", res)
				}

				diff := cmp.Diff(v.expectedResult, createdUser, cmpopts.IgnoreFields(dtos.UserResponse{}, "Id", "Role", "CreatedAt"))
				if diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
