package rest

import (
	"net/http"
)

func (cr *AppChiRouter) getConfig(rw http.ResponseWriter, r *http.Request) {
	cr.renderJSON(rw, http.StatusOK, cr.config)
}
