package middlewares

import (
	"net/http"
	"golang.org/x/time/rate"
)

func ThrottleMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(1, 1)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "The request cannot be served due to the rate limit having been exhausted for the resource", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
