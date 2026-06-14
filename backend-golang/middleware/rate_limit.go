package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"react-example/backend-golang/httputil"
)

type IPBucket struct {
	tokens     float64
	lastRefill time.Time
}

type RateLimiter struct {
	mu       sync.RWMutex
	clients    map[string]*IPBucket
	rate       float64
	capacity   float64
	ttl        time.Duration
}

func NewRateLimiter(rate float64, capacity float64, ttl time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*IPBucket),
		rate:     rate,
		capacity: capacity,
		ttl:      ttl,
	}
	go rl.startCleanupDaemon(5 * time.Minute)
	return rl
}

func (rl *RateLimiter) Limit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.clients[ip]
	now := time.Now()

	if !exists {
		rl.clients[ip] = &IPBucket{
			tokens:     rl.capacity - 1,
			lastRefill: now,
		}
		return false
	}

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > rl.capacity {
		bucket.tokens = rl.capacity
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return false
	}

	return true
}

func (rl *RateLimiter) startCleanupDaemon(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, bucket := range rl.clients {
			if now.Sub(bucket.lastRefill) > rl.ttl {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func GetIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func LimitMiddleware(limiter *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)
		if limiter.Limit(ip) {
			httputil.WriteErrorResponse(w, http.StatusTooManyRequests, "429", "Too Many Requests")
			return
		}
		next(w, r)
	}
}
