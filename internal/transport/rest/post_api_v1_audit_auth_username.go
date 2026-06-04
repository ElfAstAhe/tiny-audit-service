package rest

import (
	"net/http"

	libhttp "github.com/ElfAstAhe/go-service-template/pkg/transport/http"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/ElfAstAhe/tiny-audit-service/internal/transport"
)

// postAPIV1AuditAuthUsername godoc
// @Summary      Список аудита данных в разрезе пользователя
// @Description  Получить список аудита данных в разрезе пользователя
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        input  body      AuditUserDTO  true  "Пользователь"
// @Success      200    {array}   AuthAuditDTO "список данных аудита"
// @Failure      400    {object}  ErrorDTO "Кривые данные"
// @Failure      401    {object}  ErrorDTO "Не авторизован"
// @Failure      403    {object}  ErrorDTO "В доступе отказано"
// @Failure      409    {object}  ErrorDTO "Уже существует"
// @Failure      500    "Внутренняя ошибка сервера (пустое тело)"
// @Router       /api/v1/audit/data/username [post]
func (cr *AppChiRouter) postAPIV1AuditAuthUsername(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debugf("postAPIV1AuditAuthUsername start, requestID [%s]", middleware.GetReqID(r.Context()))
	defer cr.log.Debugf("postAPIV1AuditAuthUsername finish, requestID [%s]", middleware.GetReqID(r.Context()))

	var income = &dto.AuditUserDTO{}
	err := libhttp.DecodeJSON(r, income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthUsername decode income json [%v]", err)
		libhttp.RenderError(rw, err, mapToHTTPStatus)
		return
	}

	res, err := cr.authAuditFacade.ListByUsername(r.Context(), income)
	if err != nil {
		cr.log.Errorf("postAPIV1AuditAuthUsername list by username [%v]", err)
		libhttp.RenderError(rw, err, mapToHTTPStatus)
		return
	}

	libhttp.RenderJSON(rw, http.StatusOK, res, mapToHTTPStatus)
}
