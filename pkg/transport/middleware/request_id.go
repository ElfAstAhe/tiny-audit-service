package middleware

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

// Константы, определяющие наименования стандартных HTTP-заголовков для передачи идентификаторов запросов.
const (
	HeaderXRequestID     = "X-Request-ID"
	HeaderXCorrelationID = "X-Correlation-ID"
	HeaderRequestID      = "Request-ID"
)

// AuditRequestIDExtractor инкапсулирует в себе перечень отслеживаемых HTTP-заголовков
// для извлечения или автоматической фиксации сквозного идентификатора запроса для аудита.
type AuditRequestIDExtractor struct {
	headers []string
}

// NewAuditRequestIDExtractor представляет собой фабричный конструктор экстрактора с произвольным набором отслеживаемых заголовков.
func NewAuditRequestIDExtractor(headers ...string) *AuditRequestIDExtractor {
	return &AuditRequestIDExtractor{
		headers: headers,
	}
}

// NewDefaultAuditRequestIDExtractor собирает экстрактор, наполненный базовыми дефолтными заголовками трассировки фреймворка.
func NewDefaultAuditRequestIDExtractor() *AuditRequestIDExtractor {
	return NewAuditRequestIDExtractor(
		HeaderXRequestID,
		HeaderXCorrelationID,
		HeaderRequestID,
	)
}

// Handle встраивает экстрактор в каскадную цепочку обработки HTTP-запросов (http.Handler Middleware Pattern).
// Последовательно сканирует заголовки, осуществляет генерацию фолбека при их отсутствии, обогащает
// context.Context идентификатором и передает управление дальше по стеку декораторов.
func (are *AuditRequestIDExtractor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get request id from header (first not empty)
		var requestID string
		for _, header := range are.headers {
			requestID = r.Header.Get(header)
			if requestID != "" {
				break
			}
		}
		// if none
		if requestID == "" {
			requestID = utils.UnknownValue
		}
		// wrap request
		req := r.WithContext(utils.WithRequestID(r.Context(), requestID))
		// next pipe node
		next.ServeHTTP(w, req)
	})
}
