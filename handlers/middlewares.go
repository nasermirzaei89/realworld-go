package handlers

import (
	"context"
	"encoding/json"
	"github.com/nasermirzaei89/realworld-go/libs/jwt"
	"net/http"
	"strconv"
	"strings"
)

func (h *handler) middlewareAuthentication(next http.HandlerFunc, force bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "missing authorization header",
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		if !strings.HasPrefix(authHeader, "Token ") {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		tokenStr := authHeader[6:]
		err := jwt.Verify(tokenStr, h.secret)
		if err != nil {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
						"error":   err.Error(),
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		token, err := jwt.Parse(tokenStr)
		if err != nil {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
						"error":   err.Error(),
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		sub, err := token.GetSubject()
		if err != nil {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
						"error":   err.Error(),
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		userID, err := strconv.Atoi(sub)
		if err != nil {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
						"error":   err.Error(),
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		user, err := h.userRepo.GetByID(userID)
		if err != nil {
			if force {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "invalid authorization header",
						"error":   err.Error(),
					},
				})
			} else {
				next(w, r)
			}

			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "current_user", user))
		next(w, r)
	}
}
