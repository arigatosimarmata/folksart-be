package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

// IPBucket defines a thread-safe token bucket for individual client IPs
type IPBucket struct {
	tokens     float64
	lastRefill time.Time
}

// RateLimiter manages the registry of individual IP rate limits
type RateLimiter struct {
	mu         sync.RWMutex
	clients    map[string]*IPBucket
	rate       float64       // Tokens refilled per second
	capacity   float64       // Maximum burst capacity
	ttl        time.Duration // Time-to-live for inactive IPs to prevent memory leaks
}

// NewRateLimiter instantiates a new RateLimiter controller
func NewRateLimiter(rate float64, capacity float64, ttl time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*IPBucket),
		rate:     rate,
		capacity: capacity,
		ttl:      ttl,
	}

	// Active background cleanup operation to prune stale client buckets (Single Responsibility)
	go rl.startCleanupDaemon(5 * time.Minute)

	return rl
}

// Limit checks whether the client IP has exceeded the specified threshold parameters
func (rl *RateLimiter) Limit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.clients[ip]
	now := time.Now()

	if !exists {
		rl.clients[ip] = &IPBucket{
			tokens:     rl.capacity - 1, // Consume first token on creation
			lastRefill: now,
		}
		return false // Allowed
	}

	// Calculate refilled tokens based on time elapsed
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * rl.rate
	if bucket.tokens > rl.capacity {
		bucket.tokens = rl.capacity
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return false // Allowed
	}

	return true // Rate Limit Exceeded
}

// startCleanupDaemon runs a ticker loop in the background checking for stale entries
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

// GetIP extracts the client IP address, checking standard proxy headers if behind load balancers
func GetIP(r *http.Request) string {
	// 1. Check Cloudflare / Reverse Proxy headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := stringsSplitCommon(xff, ",")
		if len(parts) > 0 {
			return stringsTrimSpaceCommon(parts[0])
		}
	}
	
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 2. Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Helper methods to keep package zero-dependency
func stringsSplitCommon(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func stringsTrimSpaceCommon(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// LimitMiddleware wraps the HTTP requests with the rate limiting algorithms
func LimitMiddleware(limiter *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)

		if limiter.Limit(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":       "Too Many Requests",
				"message":     "Identity operations quota exceeded. Please reduce request frequency.",
				"retry_after": "1s",
			})
			return
		}

		next(w, r)
	}
}
