package catalog_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/response"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
)

type storeStub struct {
	create func(ctx context.Context, venue *catalog.Venue) error
	get    func(ctx context.Context, id uuid.UUID) (*catalog.Venue, error)
}

func (s storeStub) Create(ctx context.Context, venue *catalog.Venue) error {
	return s.create(ctx, venue)
}

func (s storeStub) Get(ctx context.Context, id uuid.UUID) (*catalog.Venue, error) {
	return s.get(ctx, id)
}

func newTestRouter(t *testing.T, store storeStub) (http.Handler, *loggertest.Buffer) {
	t.Helper()

	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	handler := catalog.NewVenueHandler(store)

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Route("/v1/venues", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
	})

	return r, buf
}

func rejectingStore(t *testing.T) storeStub {
	t.Helper()

	return storeStub{
		create: func(context.Context, *catalog.Venue) error {
			t.Error("store must not be called")

			return nil
		},
		get: func(context.Context, uuid.UUID) (*catalog.Venue, error) {
			t.Error("store must not be called")

			return nil, nil
		},
	}
}

func createBody(name, city, address, timezone string) string {
	return fmt.Sprintf(
		`{"name":%q,"city":%q,"address":%q,"timezone":%q}`,
		name, city, address, timezone,
	)
}

func do(t *testing.T, router http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequestWithContext(t.Context(), method, target, strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	return w
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) (string, map[string]string) {
	t.Helper()

	var body struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	return body.Error, body.Fields
}

func errorMessage(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()

	msg, fields := decodeError(t, w)
	require.Empty(t, fields, "у этой ошибки не должно быть разбивки по полям")

	return msg
}

type wireVenue struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Timezone  string `json:"timezone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func decodeVenue(t *testing.T, w *httptest.ResponseRecorder) wireVenue {
	t.Helper()

	var v wireVenue
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &v))

	return v
}

func requireInternalError(
	t *testing.T,
	w *httptest.ResponseRecorder,
	buf *loggertest.Buffer,
	wantOp string,
) {
	t.Helper()

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Equal(t, response.InternalErrorMsg, errorMessage(t, w))
	require.NotContains(t, w.Body.String(), "boom")
	require.NotContains(t, w.Body.String(), wantOp)

	logs := loggertest.Logs(t, buf)
	require.NotEmpty(t, logs)
	require.Equal(t, "request failed", logs[0]["msg"])
	require.Contains(t, logs[0]["error"], wantOp)
	require.Contains(t, logs[0]["error"], "boom")
	require.NotEmpty(t, logs[0]["request_id"])
}

func TestVenueHandlerCreate(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		var saved *catalog.Venue
		store := storeStub{create: func(_ context.Context, venue *catalog.Venue) error {
			saved = venue

			return nil
		}}
		router, _ := newTestRouter(t, store)

		body := createBody(" Большой театр ", "Москва", "Театральная площадь, 1", "Europe/Moscow")
		w := do(t, router, http.MethodPost, "/v1/venues", body)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.NotNil(t, saved)
		require.Equal(t, "Большой театр", saved.Name)

		got := decodeVenue(t, w)
		require.Equal(t, saved.ID.String(), got.ID)
		require.Equal(t, "Большой театр", got.Name)
		require.Equal(t, "Москва", got.City)
		require.Equal(t, "Театральная площадь, 1", got.Address)
		require.Equal(t, "Europe/Moscow", got.Timezone)
		require.Equal(t, got.CreatedAt, got.UpdatedAt)
		require.Equal(t, "/v1/venues/"+got.ID, w.Header().Get("Location"))
	})

	t.Run("invalid body", func(t *testing.T) {
		tests := map[string]string{
			"invalid json":  `{"name":`,
			"unknown field": `{"name":"x","city":"y","address":"z","timezone":"UTC","o":1}`,
		}

		for name, body := range tests {
			t.Run(name, func(t *testing.T) {
				router, _ := newTestRouter(t, rejectingStore(t))

				w := do(t, router, http.MethodPost, "/v1/venues", body)

				require.Equal(t, http.StatusBadRequest, w.Code)
				require.Equal(t, "Invalid body", errorMessage(t, w))
			})
		}
	})

	t.Run("validation failed", func(t *testing.T) {
		longName := strings.Repeat("я", 201)
		longCity := strings.Repeat("я", 101)
		longAddress := strings.Repeat("я", 301)

		tests := map[string]struct {
			body       string
			wantFields map[string]string
		}{
			"empty name": {
				body:       createBody("   ", "Москва", "Тверская, 1", "Europe/Moscow"),
				wantFields: map[string]string{"name": "must not be empty"},
			},
			"empty city": {
				body:       createBody("МХТ", "", "Тверская, 1", "Europe/Moscow"),
				wantFields: map[string]string{"city": "must not be empty"},
			},
			"empty address": {
				body:       createBody("МХТ", "Москва", " ", "Europe/Moscow"),
				wantFields: map[string]string{"address": "must not be empty"},
			},
			"empty timezone": {
				body:       createBody("МХТ", "Москва", "Тверская, 1", ""),
				wantFields: map[string]string{"timezone": "must not be empty"},
			},
			"local timezone": {
				body:       createBody("МХТ", "Москва", "Тверская, 1", "Local"),
				wantFields: map[string]string{"timezone": "must not be local"},
			},
			"unknown timezone": {
				body:       createBody("МХТ", "Москва", "Тверская, 1", "Mars/Olympus"),
				wantFields: map[string]string{"timezone": "invalid timezone"},
			},
			"long name": {
				body:       createBody(longName, "Москва", "Тверская, 1", "UTC"),
				wantFields: map[string]string{"name": "must be at most 200 characters"},
			},
			"long city": {
				body:       createBody("МХТ", longCity, "Тверская, 1", "UTC"),
				wantFields: map[string]string{"city": "must be at most 100 characters"},
			},
			"long address": {
				body:       createBody("МХТ", "Москва", longAddress, "UTC"),
				wantFields: map[string]string{"address": "must be at most 300 characters"},
			},
			"several fields": {
				body: createBody("", "", "Тверская, 1", "Mars/Olympus"),
				wantFields: map[string]string{
					"name":     "must not be empty",
					"city":     "must not be empty",
					"timezone": "invalid timezone",
				},
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				router, buf := newTestRouter(t, rejectingStore(t))

				w := do(t, router, http.MethodPost, "/v1/venues", tt.body)

				require.Equal(t, http.StatusUnprocessableEntity, w.Code)

				msg, fields := decodeError(t, w)
				require.Equal(t, "Validation error", msg)
				require.Equal(t, tt.wantFields, fields)

				require.Len(t, loggertest.Logs(t, buf), 1)
			})
		}
	})

	t.Run("body over limit", func(t *testing.T) {
		router, _ := newTestRouter(t, rejectingStore(t))

		body := createBody(strings.Repeat("a", 70*1024), "Москва", "Тверская, 1", "UTC")
		w := do(t, router, http.MethodPost, "/v1/venues", body)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "Invalid body", errorMessage(t, w))
	})

	t.Run("store failed", func(t *testing.T) {
		store := storeStub{create: func(context.Context, *catalog.Venue) error {
			return errors.New("boom")
		}}
		router, buf := newTestRouter(t, store)

		body := createBody("МХТ", "Москва", "Камергерский, 3", "Europe/Moscow")
		w := do(t, router, http.MethodPost, "/v1/venues", body)

		requireInternalError(t, w, buf, "create venue")
	})
}

func TestVenueHandlerGet(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		venue, err := catalog.NewVenue(catalog.VenueCreateParams{
			Name:     "МХТ",
			City:     "Москва",
			Address:  "Камергерский, 3",
			Timezone: "Europe/Moscow",
		})
		require.NoError(t, err)

		var asked uuid.UUID
		store := storeStub{get: func(_ context.Context, id uuid.UUID) (*catalog.Venue, error) {
			asked = id

			return venue, nil
		}}
		router, _ := newTestRouter(t, store)

		w := do(t, router, http.MethodGet, "/v1/venues/"+venue.ID.String(), "")

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, venue.ID, asked)

		got := decodeVenue(t, w)
		require.Equal(t, venue.ID.String(), got.ID)
		require.Equal(t, "МХТ", got.Name)
		require.Equal(t, "Europe/Moscow", got.Timezone)
	})

	t.Run("invalid id", func(t *testing.T) {
		router, _ := newTestRouter(t, rejectingStore(t))

		w := do(t, router, http.MethodGet, "/v1/venues/not-a-uuid", "")

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "Invalid ID", errorMessage(t, w))
	})

	t.Run("not found", func(t *testing.T) {
		store := storeStub{get: func(context.Context, uuid.UUID) (*catalog.Venue, error) {
			return nil, catalog.ErrVenueNotFound
		}}
		router, _ := newTestRouter(t, store)

		w := do(t, router, http.MethodGet, "/v1/venues/"+uuid.NewV7().String(), "")

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "Venue not found", errorMessage(t, w))
	})

	t.Run("store failed", func(t *testing.T) {
		store := storeStub{get: func(context.Context, uuid.UUID) (*catalog.Venue, error) {
			return nil, errors.New("boom")
		}}
		router, buf := newTestRouter(t, store)

		w := do(t, router, http.MethodGet, "/v1/venues/"+uuid.NewV7().String(), "")

		requireInternalError(t, w, buf, "get venue")
	})
}
