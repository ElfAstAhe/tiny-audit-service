package rest

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	_ "github.com/ElfAstAhe/tiny-audit-service/internal/transport"
)

func (cr *AppChiRouter) postAPIV1AuditAuthUsername(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditAuthUsername start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditAuthUsername finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.AuditUserDTO{}
	err := cr.decodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthUsername decode income json [%v]", err)
		cr.renderError(rw, err)
		return
	}

	res, err := cr.authAuditFacade.ListByUsername(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthUsername list by username [%v]", err)
		cr.renderError(rw, err)
		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
