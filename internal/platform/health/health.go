package health

import (
	"context"
	"net/http"
)

type Checker interface {
	Ping(context.Context) error
}

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func Ready(checkers []Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, ch := range checkers {
			if err := ch.Ping(r.Context()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
