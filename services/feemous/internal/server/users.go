package server

import "net/http"

type userResponse struct {
	ID         string  `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Patronymic string  `json:"patronymic"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone,omitempty"`
	Telegram   *string `json:"telegram,omitempty"`
	TeamID     *string `json:"team_id"`
	UnitID     *string `json:"unit_id"`
	PositionID *string `json:"position_id"`
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]userResponse, 0, len(items))
	for _, u := range items {
		resp = append(resp, userResponse{
			ID:         u.ID,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Patronymic: u.Patronymic,
			Email:      u.Email,
			Phone:      u.Phone,
			Telegram:   u.Telegram,
			TeamID:     u.TeamID,
			UnitID:     u.UnitID,
			PositionID: u.PositionID,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}
