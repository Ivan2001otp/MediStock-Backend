package middleware

import (
	"Medistock_Backend/internals/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/time/rate"
)

type Message map[string]interface{}

func CheckRoleMiddleware(expectedRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("actor").(string)

			if !ok || role != expectedRole {
				http.Error(w, "Forbidden:Access Denied", http.StatusForbidden)
				return
			}
			log.Println("Role values are passed as SAME ... ")
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or Malformed token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			log.Println("token sent in header : ", tokenStr)

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(services.GetSecretKey()), nil
			})

			if err != nil {
				log.Println("Something wrong happened in Auth-Middlewaure : ", err.Error())
				if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
					http.Error(w, "Access Token Expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Invalid Token", http.StatusUnauthorized)
				}
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return

			}

			claims := token.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), "email", claims["email"])
			ctx = context.WithValue(ctx, "actor", claims["actor"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func RateLimitMiddleWare(next func(w http.ResponseWriter, r *http.Request)) http.Handler {

	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				log.Println("The rate-limiter is at max...")

				message := Message{
					"status":  http.StatusTooManyRequests,
					"error":   "Too many requests to handle",
					"message": "The API is at max-capacity. Try again later !",
				}

				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(&message)
			} else {
				log.Println("The rate-limiter just passed it on !")
				next(w, r)
			}
		})
}

