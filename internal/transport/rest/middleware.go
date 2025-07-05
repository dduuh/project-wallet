package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"wallet-service/internal/domain"
	claims_jwt "wallet-service/internal/jwt_claims"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const userIdContext = "userId"

func (s *Server) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			logrus.Warn("header is empty")
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		partsOfHeader := strings.Split(header, " ")
		if len(partsOfHeader) != 2 {
			logrus.Warn("header length is not equal 2")
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		if partsOfHeader[0] != "Bearer" {
			logrus.Warn("first part of header is not equal 'Bearer'")
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(partsOfHeader[1], &claims_jwt.Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, claims_jwt.ErrInvalidSigningMethod
			}

			return s.key, nil
		})
		if err != nil {
			logrus.Warn("failed to parse the token")
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		claims, ok := token.Claims.(*claims_jwt.Claims)
		if !ok {
			logrus.Warn("no such claims or token is not valid")
			response(w, http.StatusUnauthorized, ErrUnauthorized)

			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			logrus.Warn("token expired")
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
