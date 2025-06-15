package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)


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