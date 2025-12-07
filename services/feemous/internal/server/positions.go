package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

func (s *Server) handleCreatePosition(w http.ResponseWriter, r *http.Request) {
	var req simpleReq
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p := &store.Position{ID: uuid.NewString(), Name: req.Name, Description: req.Description}
	if err := s.store.CreatePosition(r.Context(), p); err != nil {
		http.Error(w, "create: "+err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": p.ID})
}

func (s *Server) handleDeletePosition(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.store.DeletePosition(r.Context(), id); err != nil {
		http.Error(w, "delete: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListPositions(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListPositions(r.Context())
	if err != nil {
		http.Error(w, "list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
