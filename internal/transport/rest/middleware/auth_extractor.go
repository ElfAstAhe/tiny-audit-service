package middleware

import (
	"net/http"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport"
)

type AuthExtractor struct {
	authHelper     auth.Helper
	log            logger.Logger
	ignorancePaths *transport.HTTPPathMatchers
}

func NewAuthExtractor(
	ignorancePaths *transport.HTTPPathMatchers,
	authHelper auth.Helper,
	logger logger.Logger,
) *AuthExtractor {
	return &AuthExtractor{
		ignorancePaths: ignorancePaths,
		authHelper:     authHelper,
		log:            logger.GetLogger("HTTP-JWT-Extractor"),
	}
}

func (aem *AuthExtractor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		aem.log.Debug("AuthExtractor.Handle start")
		defer aem.log.Debug("AuthExtractor.Handle finish")

		// check for ignorance
		if aem.ignorancePaths.Match(r.Method, r.RequestURI) {
			next.ServeHTTP(rw, r)

			return
		}

		// here and so on we extract subject and gen 401 or continue pipeline
		// extract subject
		subj, err := aem.authHelper.SubjectFromHTTPRequest(r)
		if err != nil {
			aem.log.Errorf("AuthExtractor.Handle error [%v]", err)

			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(rw, r.WithContext(auth.WithSubject(r.Context(), subj)))
	})
}
