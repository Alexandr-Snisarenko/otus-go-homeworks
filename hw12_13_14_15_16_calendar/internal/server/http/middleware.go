package http

import (
	"net/http"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/ctxmeta"
	"github.com/google/uuid"
)

const headerRequestID = "X-Request-ID"

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused, revive
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //nolint:revive
		// TODO
	})
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(headerRequestID)
		if _, err := uuid.Parse(reqID); err != nil || reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := ctxmeta.WithRequestID(r.Context(), reqID)
		w.Header().Set(headerRequestID, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
