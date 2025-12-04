package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/GoBetterAuth/go-better-auth/internal/auth"
	"github.com/GoBetterAuth/go-better-auth/internal/util"
	"github.com/GoBetterAuth/go-better-auth/pkg/domain"
)

func EndpointHooksMiddleware(config *domain.Config, authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.EndpointHooks.Before == nil && config.EndpointHooks.After == nil {
				next.ServeHTTP(w, r)
				return
			}

			hookCtx := &domain.EndpointHookContext{
				Path:    r.URL.Path,
				Method:  r.Method,
				Headers: make(map[string]string),
				Query:   make(map[string]string),
				Request: r,
			}

			for k, v := range r.Header {
				if len(v) > 0 {
					hookCtx.Headers[k] = v[0]
				}
			}

			for k, v := range r.URL.Query() {
				if len(v) > 0 {
					hookCtx.Query[k] = v[0]
				}
			}

			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
					r.GetBody = func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader(bodyBytes)), nil
					}

					if len(bodyBytes) > 0 {
						var bodyMap map[string]any
						if json.Unmarshal(bodyBytes, &bodyMap) == nil {
							hookCtx.Body = bodyMap
						}
					}
				}
			}

			cookie, err := r.Cookie(config.Session.CookieName)
			if err == nil && cookie.Value != "" {
				session, _ := authService.SessionService.GetSessionByToken(authService.TokenService.HashToken(cookie.Value))
				if session != nil {
					user, _ := authService.UserService.GetUserByID(session.UserID)
					if user != nil {
						hookCtx.User = user
					}
				}
			}

			if config.EndpointHooks.Before != nil {
				if err := config.EndpointHooks.Before(hookCtx); err != nil {
					util.JSONResponse(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
					return
				}
			}

			next.ServeHTTP(w, r)

			if config.EndpointHooks.After != nil {
				go func() {
					if err := config.EndpointHooks.After(hookCtx); err != nil {
						slog.Error("Error in After Hook for %s: %v", hookCtx.Path, err)
					}
				}()
			}
		})
	}
}
