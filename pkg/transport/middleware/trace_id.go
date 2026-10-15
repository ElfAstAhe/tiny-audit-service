package middleware

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

const (
	// HeaderXCloudTraceContext maps the canonical HTTP header key used by Google Cloud infrastructure for trace propagation.
	HeaderXCloudTraceContext string = "X-Cloud-Trace-Context"
	// HeaderTraceParent maps the W3C standard distributed tracing context header specification key.
	HeaderTraceParent string = "Traceparent"
	// HeaderXTraceID maps a common customized legacy corporate identification tracing span header key.
	HeaderXTraceID string = "X-Trace-ID"
	// HeaderTraceID maps a flat simplified diagnostic sequence tracing key.
	HeaderTraceID string = "Trace-ID"
)

// AuditTraceIDExtractor structures configuration parameters tracking targeted HTTP headers list arrays for execution contexts correlation.
type AuditTraceIDExtractor struct {
	headers []string // Sequential priority array of target header keys to evaluate during inbound traffic screening
}

// NewAuditTraceIDExtractor acts as a factory constructor allocating custom header tracking sequences for the extractor.
func NewAuditTraceIDExtractor(headers ...string) *AuditTraceIDExtractor {
	return &AuditTraceIDExtractor{
		headers: headers,
	}
}

// NewDefaultAuditTraceIDExtractor returns an initialization blueprint pre-populating standard cloud and corporate trace context header keys.
func NewDefaultAuditTraceIDExtractor() *AuditTraceIDExtractor {
	return NewAuditTraceIDExtractor(
		HeaderXCloudTraceContext,
		HeaderTraceParent,
		HeaderXTraceID,
		HeaderTraceID,
	)
}

// Handle intercepts inbound standard net/http server traffic, executing a fail-fast priority lookups loop to extract and inject active correlation trace identifiers.
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
