package rest

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"
)

// postAPIV1AuditData godoc
// @Summary      Аудит данных
// @Description  Создание новой записи аудита
// @Tags         audit
// @Accept       json
// @Param        input  body      DataAuditDTO  true  "Аудит данных"
// @Success      201    "Всё хорошо (пустое тело)"
// @Failure      400    {object}  ErrorDTO "Кривые данные"
// @Failure      401    {object}  ErrorDTO "Не авторизован"
// @Failure      403    {object}  ErrorDTO "В доступе отказано"
// @Failure      500    "Внутренняя ошибка сервера (пустое тело)"
// @Router       /api/v1/audit/data [post]
func (cr *AppChiRouter) postAPIV1AuditData(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditData start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditData finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.DataAuditDTO{}
	err := cr.decodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditData decode income json [%v]", err)

		cr.renderError(rw, err)

		return
	}

	err = cr.dataAuditFacade.Audit(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditData post data audit error, [%v]", err)

		cr.renderError(rw, err)

		return
	}

	cr.renderEmpty(rw, http.StatusCreated)

}
