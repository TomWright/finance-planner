package http

import (
	"github.com/go-chi/chi"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
	"net/http"
)

func NewStatsTransactionsHandler(profileService service.Profile) *statsTransactionsHandler {
	return &statsTransactionsHandler{
		profileService: profileService,
	}
}

type statsTransactionsHandler struct {
	profileService service.Profile
}

func (x *statsTransactionsHandler) Bind(r chi.Router) {
	r.Get("/{profile}/transactions/stats", x.handle)
}

func (x *statsTransactionsHandler) handle(rw http.ResponseWriter, r *http.Request) {
	var err error
	profileName := chi.URLParam(r, "profile")

	profile, err := x.profileService.LoadOrCreateProfile(profileName)
	if err != nil {
		sendError(err, rw)
		return
	}

	res := struct {
		Sum int64 `json:"sum"`
	}{}

	err = profile.Transactions.Range(nil, func(t domain.Transaction) error {
		res.Sum += t.Amount
		return nil
	})
	if err != nil {
		sendError(err, rw)
		return
	}

	sendResponse(res, http.StatusOK, rw)
}
