package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"wallet-service/internal/domain"
	jwtclaims "wallet-service/internal/jwt_claims"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const userIdContext = "userId"

func (s *Server) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		partsOfHeader := strings.Split(header, "")

		if partsOfHeader[0] != "Bearer" {
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(partsOfHeader[1], &jwtclaims.Claims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwtclaims.ErrInvalidSigningMethod
			}

			return s.key, nil
		})
		if err != nil {
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		claims, ok := token.Claims.(jwtclaims.Claims)
		if !(ok && token.Valid) {
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		userInfo := domain.UserInfo{
			Id: claims.UserId,
		}

		r = r.WithContext(context.WithValue(r.Context(), userIdContext, userInfo))
		next.ServeHTTP(w, r)
	})
}
