package middleware

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

const (
	HeaderXRequestID     = "X-Request-ID"
	HeaderXCorrelationID = "X-Correlation-ID"
	HeaderRequestID      = "Request-ID"
)

type AuditRequestIDExtractor struct {
	headers []string
}

func NewAuditRequestIDExtractor(headers []string) *AuditRequestIDExtractor {
	return &AuditRequestIDExtractor{
		headers: headers,
	}
}

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
			requestID = "unknown"
		}
		// wrap request
		req := r.WithContext(utils.WithRequestID(r.Context(), requestID))
		// next pipe node
		next.ServeHTTP(w, req)
	})
}
