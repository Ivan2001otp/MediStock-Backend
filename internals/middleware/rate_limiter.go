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

/*
type rateLimiter struct {
	ips map[string] *rate.Limiter
	mu *sync.RWMutex	// mutex to protect map
	r rate.Limit   	// requests per second
	b int			// burst size
}

func newRateLimiter(r rate.Limit, b int) *rateLimiter {
	return &rateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu : &sync.RWMutex{},
		r: r,
		b: b,
	}
}

func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock() // read lock first
	limiter, exists := rl.ips[ip]

	rl.mu.RUnlock()

	if !exists {
		// acquire a write lock
		rl.mu.Lock()
		defer rl.mu.Unlock()

		limiter , exists = rl.ips[ip]
		// double check after acquiring write lock, incase another might create it..
		if (! exists) {
			limiter = rate.NewLimiter(rl.r, rl.b)
			rl.ips[ip] = limiter;
			log.Printf("New rate limiter created for IP: %s (Rate: %v, Burst: %d)", ip, rl.r, rl.b)
		}
	}

	return limiter
}

func RateLimitMiddleWare(next http.Handler, rps float64, burst int) http.Handler {
	limiter := newRateLimiter(rate.Limit(rps), burst)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip, _ , err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			log.Printf("Error parsing IP address from RemoteAddr %s: %v", req.RemoteAddr, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}


		currentLimiter := limiter.getLimiter(ip)

		if !currentLimiter.Allow() {
			log.Printf("rate limit exceeded immanuel ")
			log.Printf("Rate limit exceeded for IP: %s", ip)
			// Optionally set rate limit headers (RFC 6585)
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", burst)) // Total requests allowed in the window (simplified to burst)
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", currentLimiter.Tokens())) // Tokens available
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(currentLimiter.Burst())*time.Second).Unix())) // Estimated reset time

			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// pass to next handler in chain
		next.ServeHTTP(w, req)
	});
}

*/
