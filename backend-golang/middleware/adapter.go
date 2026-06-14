package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"react-example/backend-golang/httputil"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Adapt(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Printf("[ERROR] %v", err)
			httputil.WriteErrorResponse(w, err)
		}
	}
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				log.Printf("[PANIC] %v\n%s", rcv, debug.Stack())
				httputil.WriteErrorResponse(w, fmt.Errorf("Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
