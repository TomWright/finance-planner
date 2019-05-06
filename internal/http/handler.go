package http

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/tomwright/finance-planner/internal/errs"
	"net/http"
)

// Handler represents a single HTTP handler.
type Handler interface {
	Bind(r chi.Router)
}

func sendError(err error, statusCode int, rw http.ResponseWriter) {
	e := errs.FromErr(err)

	resp := map[string]interface{}{
		"code":  e.Code(),
		"error": e.Message(),
	}

	sendResponse(resp, statusCode, rw)
}

func sendResponse(body interface{}, statusCode int, rw http.ResponseWriter) {
	bytes, err := json.Marshal(body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(`Internal server error: ` + err.Error()))
		return
	}

	rw.WriteHeader(statusCode)
	_, _ = rw.Write(bytes)
}
