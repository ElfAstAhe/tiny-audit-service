package rest

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	_ "github.com/ElfAstAhe/tiny-audit-service/internal/transport"
)

// postAPIV1AuditAuthPeriod godoc
// @Summary      Список аудита аутентификации в разрезе периода
// @Description  Получить список аудита аутентификации в разрезе периода
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        input  body      AuditPeriodDTO  true  "Период аудита"
// @Success      200    {array}   AuthAuditDTO "список данных аудита"
// @Failure      400    {object}  ErrorDTO "Кривые данные"
// @Failure      401    {object}  ErrorDTO "Не авторизован"
// @Failure      403    {object}  ErrorDTO "В доступе отказано"
// @Failure      409    {object}  ErrorDTO "Уже существует"
// @Failure      500    "Внутренняя ошибка сервера (пустое тело)"
// @Router       /api/v1/audit/auth/period [post]
func (cr *AppChiRouter) postAPIV1AuditAuthPeriod(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditAuthPeriod start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditAuthPeriod finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.AuditPeriodDTO{}
	err := cr.decodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthPeriod decode income json [%v]", err)
		cr.renderError(rw, err)
		return
	}

	res, err := cr.authAuditFacade.ListByPeriod(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthPeriod list by period [%v]", err)
		cr.renderError(rw, err)
		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
