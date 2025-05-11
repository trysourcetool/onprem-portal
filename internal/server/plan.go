package server

import (
	"net/http"
	"strconv"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type planResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	StripePriceID string `json:"stripePriceId"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type listPlansResponse struct {
	Plans []*planResponse `json:"plans"`
}

func planFromModel(p *core.Plan) *planResponse {
	if p == nil {
		return nil
	}
	return &planResponse{
		ID:            p.ID.String(),
		Name:          p.Name,
		Price:         p.Price,
		StripePriceID: p.StripePriceID,
		CreatedAt:     strconv.FormatInt(p.CreatedAt.Unix(), 10),
		UpdatedAt:     strconv.FormatInt(p.UpdatedAt.Unix(), 10),
	}
}

func (s *Server) handleListPlans(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	plans, err := s.db.Plan().List(ctx)
	if err != nil {
		return err
	}
	resp := &listPlansResponse{Plans: make([]*planResponse, 0, len(plans))}
	for _, p := range plans {
		resp.Plans = append(resp.Plans, planFromModel(p))
	}
	return s.renderJSON(w, http.StatusOK, resp)
}
