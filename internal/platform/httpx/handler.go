package httpx

import "net/http"

type Handler func(http.ResponseWriter, *http.Request) error

func Handle(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		if err := fn(w, r); err != nil {
			WriteError(w, r, err)
		}
	}
}
