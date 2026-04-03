package rest

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	_ "github.com/ElfAstAhe/tiny-audit-service/internal/transport"
)

func (cr *AppChiRouter) postAPIV1AuditAuth(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditAuth start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditAuth finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.AuthAuditDTO{}
	err := cr.decodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuth decode income json [%v]", err)

		cr.renderError(rw, err)

		return
	}

	err = cr.authAuditFacade.Audit(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuth post auth audit error, [%v]", err)

		cr.renderError(rw, err)

		return
	}

	cr.renderEmpty(rw, http.StatusCreated)
}
