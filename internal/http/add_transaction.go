package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/application/service"
	"github.com/tomwright/finance-planner/internal/errs"
	"io/ioutil"
	"net/http"
)

func NewAddTransactionHandler(profileService service.Profile) *addTransactionHandler {
	return &addTransactionHandler{
		profileService: profileService,
	}
}

type addTransactionHandler struct {
	profileService service.Profile
}

func (x *addTransactionHandler) Bind(r chi.Router) {
	r.Post("/{profile}/transactions", x.handle)
}

func (x *addTransactionHandler) handle(rw http.ResponseWriter, r *http.Request) {
	var err error

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(errs.FromErr(err).PrefixMessage("could not read body: "), rw)
		return
	}

	input := struct {
		Amount int64    `json:"amount"`
		Label  string   `json:"label"`
		Tags   []string `json:"tags"`
	}{}

	if err := json.Unmarshal(bytes, &input); err != nil {
		sendError(errs.FromErr(err).PrefixMessage("could not unmarshal body: "), rw)
		return
	}

	profileName := chi.URLParam(r, "profile")

	profile, err := x.profileService.LoadOrCreateProfile(profileName)
	if err != nil {
		sendError(err, rw)
		return
	}

	t := domain.NewTransaction().
		WithLabel(input.Label).
		WithAmount(input.Amount)

	if len(input.Tags) > 0 {
		t.WithTags(input.Tags...)
	}

	profile.Transactions.Add(t)

	if err := x.profileService.SaveProfile(profile); err != nil {
		sendError(err, rw)
		return
	}

	tags := make([]string, 0)
	tags = append(tags, t.Tags...)
	transactionsData := map[string]interface{}{
		"uuid":   t.UUID,
		"label":  t.Label,
		"amount": t.Amount,
		"tags":   tags,
	}

	sendResponse(transactionsData, http.StatusCreated, rw)
}
