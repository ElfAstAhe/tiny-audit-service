package middleware

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

const (
	HeaderXCloudTraceContext string = "X-Cloud-Trace-Context"
	HeaderTraceParent        string = "Traceparent"
	HeaderXTraceID           string = "X-Trace-ID"
	HeaderTraceID            string = "Trace-ID"
)

type AuditTraceIDExtractor struct {
	headers []string
}

func NewAuditTraceIDExtractor(headers []string) *AuditTraceIDExtractor {
	return &AuditTraceIDExtractor{
		headers: headers,
	}
}

func (ate *AuditTraceIDExtractor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceID string
		for _, header := range ate.headers {
			traceID = r.Header.Get(header)
			if traceID != "" {
				break
			}
		}
		// if none
		if traceID == "" {
			traceID = "unknown"
		}
		// wrap request
		req := r.WithContext(utils.WithTraceID(r.Context(), traceID))
		// next pipe node
		next.ServeHTTP(w, req)
	})
}
