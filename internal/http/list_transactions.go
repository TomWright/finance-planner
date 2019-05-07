package http

import (
	"github.com/go-chi/chi"
	"github.com/tomwright/finance-planner/internal/application/service"
	"net/http"
)

func NewListTransactionsHandler(profileService service.Profile) *listTransactionsHandler {
	return &listTransactionsHandler{
		profileService: profileService,
	}
}

type listTransactionsHandler struct {
	profileService service.Profile
}

func (x *listTransactionsHandler) Bind(r chi.Router) {
	r.Get("/{profile}/transactions", x.handle)
}

func (x *listTransactionsHandler) handle(rw http.ResponseWriter, r *http.Request) {
	profileName := chi.URLParam(r, "profile")

	profile, err := x.profileService.LoadOrCreateProfile(profileName)
	if err != nil {
		sendError(err, rw)
		return
	}

	transactions := profile.Transactions.All()
	transactionsData := make([]map[string]interface{}, len(transactions))

	for i := range transactions {
		tags := make([]string, 0)
		tags = append(tags, transactions[i].Tags...)
		transactionsData[i] = map[string]interface{}{
			"uuid":   transactions[i].UUID,
			"label":  transactions[i].Label,
			"amount": transactions[i].Amount,
			"tags":   tags,
		}
	}

	resp := map[string]interface{}{
		"data": transactionsData,
	}

	sendResponse(resp, http.StatusOK, rw)
}
