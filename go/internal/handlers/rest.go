package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/services"
)

func NewRestHandler(configs configs.Configs, customMiddleware MiddlewareHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chiMiddleware.CleanPath)
	router.Use(chiMiddleware.RealIP)
	router.Use(customMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	options := cors.Options{
		AllowedOrigins:   strings.Split(configs.Env.AllowedOrigins, ","),
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "Host", "Origin", "Referer", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	router.Use(cors.Handler(options))
	router.Use(chiMiddleware.Heartbeat("/ping"))

	authService := services.NewAuthService(configs)
	authHandler := NewAuthHandler(configs, authService)
	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/login", authHandler.Login)
	router.Post("/auth/logout", authHandler.Logout)

	router.Group(func(r chi.Router) {
		r.Use(customMiddleware.Authenticate)

		userHandler := NewUserHandler(configs)
		r.Get("/users/me", userHandler.GetUser)
		r.Get("/users/me/reservations", userHandler.GetUserReservations)
		r.Post("/users/me/reservations", userHandler.CreateUserReservation)
		r.Patch("/users/me/reservations/{reservationId}", userHandler.CancelUserReservation)

		meetingRoomHandler := NewMeetingRoomHandler(configs)
		r.Get("/meeting-rooms", meetingRoomHandler.GetMeetingRooms)
		r.Get("/meeting-rooms/available", meetingRoomHandler.GetAvailableMeetingRooms)

		reservationHandler := NewReservationHandler(configs)
		r.Get("/reservations", reservationHandler.GetReservations)
		r.Patch("/reservations/{reservationId}", reservationHandler.CancelReservation)
	})

	return router
}
