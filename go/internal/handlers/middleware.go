package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mdayat/artics-communication/go/configs"
	"github.com/mdayat/artics-communication/go/internal/services"
	"github.com/rs/zerolog/log"
)

type MiddlewareHandler interface {
	Logger(next http.Handler) http.Handler
	Authenticator
}

type middleware struct {
	configs       configs.Configs
	authenticator Authenticator
}

func NewMiddlewareHandler(configs configs.Configs, authenticator Authenticator) MiddlewareHandler {
	return &middleware{
		configs:       configs,
		authenticator: authenticator,
	}
}

func (m middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		logger := log.
			With().
			Str("request_id", uuid.New().String()).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Str("client_ip", req.RemoteAddr).
			Logger()

		req = req.WithContext(logger.WithContext(req.Context()))
		next.ServeHTTP(res, req)
	})
}

func (m middleware) Authenticate(next http.Handler) http.Handler {
	return m.authenticator.Authenticate(next)
}

type Authenticator interface {
	Authenticate(next http.Handler) http.Handler
}

type userIdKey struct{}
type accountRoleKey struct{}

type prodAuthenticator struct {
	authService services.AuthServicer
}

func NewProdAuthenticator(authService services.AuthServicer) Authenticator {
	return &prodAuthenticator{
		authService: authService,
	}
}

func (p prodAuthenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger := log.Ctx(ctx).With().Logger()

		accessTokenCookie, err := req.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				logger.Error().Err(err).Caller().Int("status_code", http.StatusUnauthorized).Msg("cookie not found")
				http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				logger.Error().Err(err).Caller().Int("status_code", http.StatusInternalServerError).Msg("failed to get cookie")
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		claims, err := p.authService.ValidateAccessToken(accessTokenCookie.Value)
		if err != nil {
			logger.Error().Err(err).Caller().Int("status_code", http.StatusUnauthorized).Msg("invalid access token")
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctxWithUserId := context.WithValue(req.Context(), userIdKey{}, claims.Subject)
		ctxWithAccountRole := context.WithValue(ctxWithUserId, accountRoleKey{}, claims.Role)

		req = req.WithContext(ctxWithAccountRole)
		next.ServeHTTP(res, req)
	})
}
