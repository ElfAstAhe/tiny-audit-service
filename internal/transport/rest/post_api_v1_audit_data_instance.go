package rest

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	_ "github.com/ElfAstAhe/tiny-audit-service/internal/transport"
)

// postAPIV1AuditAuthUsername godoc
// @Summary      Список аудита данных в разрезе инстанса данных
// @Description  Получить список аудита данных в разрезе инстанса данных
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        input  body      AuditInstanceDTO  true  "инстанс данных"
// @Success      200    {array}   DataAuditDTO "список данных аудита"
// @Failure      400    {object}  ErrorDTO "Кривые данные"
// @Failure      401    {object}  ErrorDTO "Не авторизован"
// @Failure      403    {object}  ErrorDTO "В доступе отказано"
// @Failure      409    {object}  ErrorDTO "Уже существует"
// @Failure      500    "Внутренняя ошибка сервера (пустое тело)"
// @Router       /api/v1/audit/data/instance [post]
func (cr *AppChiRouter) postAPIV1AuditDataInstance(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditDataInstance start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditDataInstance finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.AuditInstanceDTO{}
	err := cr.decodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditDataInstance decode income json [%v]", err)
		cr.renderError(rw, err)
		return
	}

	res, err := cr.dataAuditFacade.ListByInstance(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditDataInstance list by instance [%v]", err)
		cr.renderError(rw, err)
		return
	}

	cr.renderJSON(rw, http.StatusOK, res)
}
