package catalog

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"path"
	"uuid"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/response"
)

type venueStore interface {
	Create(ctx context.Context, v *Venue) error
	Get(ctx context.Context, id uuid.UUID) (*Venue, error)
}

type VenueHandler struct {
	store venueStore
}

func NewVenueHandler(store venueStore) *VenueHandler {
	return &VenueHandler{
		store: store,
	}
}

func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "create venue"

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

	var req createVenueRequest
	if err := json.UnmarshalRead(r.Body, &req, json.RejectUnknownMembers(true)); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid body")
		return
	}

	venue, err := NewVenue(VenueCreateParams(req))
	if err != nil {
		httpx.WriteError(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	if err := h.store.Create(r.Context(), venue); err != nil {
		httpx.WriteError(w, r, fmt.Errorf("%s: %w", op, err))
		return
	}

	resp := venueResponseFromModel(venue)
	w.Header().Set("Location", path.Join(r.URL.Path, resp.ID.String()))
	response.JSON(w, http.StatusCreated, resp)
}

func (h *VenueHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "get venue"

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	venue, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrVenueNotFound) {
			response.Error(w, http.StatusNotFound, "Venue not found")
			return
		}

		httpx.WriteError(w, r, fmt.Errorf("%s %q: %w", op, id, err))
		return
	}

	resp := venueResponseFromModel(venue)
	response.JSON(w, http.StatusOK, resp)
}
