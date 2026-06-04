package rest

import (
	"net/http"

	libhttp "github.com/ElfAstAhe/go-service-template/pkg/transport/http"
)

func (cr *AppChiRouter) getConfig(rw http.ResponseWriter, r *http.Request) {
	libhttp.RenderJSON(rw, http.StatusOK, cr.config, mapToHTTPStatus)
}
