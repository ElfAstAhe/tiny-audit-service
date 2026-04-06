package middleware

import (
	"errors"
	"net/http"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport"
	"github.com/golang-jwt/jwt/v5"
)

type AuthExtractor struct {
	jwtHelper     *helper.JWTHelper
	jwtHTTPHelper *helper.JWTHTTPHelper
	authHelper    auth.Helper
	log           logger.Logger
	nonSecure     *transport.HTTPPathMatchers
	acceptIssuers map[string]struct{}
}

func NewAuthExtractor(
	jwtHelper *helper.JWTHelper,
	jwtHTTPHelper *helper.JWTHTTPHelper,
	authHelper auth.Helper,
	logger logger.Logger,
	nonSecure *transport.HTTPPathMatchers,
	acceptIssuers []string,
) *AuthExtractor {
	acceptIssuersMap := make(map[string]struct{}, len(acceptIssuers))
	for _, issuer := range acceptIssuers {
		acceptIssuersMap[issuer] = struct{}{}
	}

	return &AuthExtractor{
		jwtHelper:     jwtHelper,
		jwtHTTPHelper: jwtHTTPHelper,
		authHelper:    authHelper,
		log:           logger.GetLogger("HTTP-JWT-Extractor"),
		nonSecure:     nonSecure,
		acceptIssuers: acceptIssuersMap,
	}
}

func (aem *AuthExtractor) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		aem.log.Debug("AuthExtractor.Handle start")
		defer aem.log.Debug("AuthExtractor.Handle finish")

		// check for nonsecure path
		if aem.nonSecure.Match(r.Method, r.RequestURI) {
			next.ServeHTTP(rw, r)

			return
		}

		// token
		token, err := aem.extractTokenString(r)
		if err != nil {
			aem.log.Errorf("AuthExtractor.Handle extract token [%v]", err)

			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}
		// claims
		claims, err := aem.jwtHelper.ExtractClaims(token)
		if err != nil {
			aem.log.Errorf("AuthExtractor.Handle extract claims [%v]", err)

			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}
		// check issuer via white list
		if !aem.isAcceptIssuer(claims.Issuer) {
			aem.log.Errorf("AuthExtractor.Handle invalid auth token: invalid issuer [%v]", claims.Issuer)

			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		// extract subject
		subj, err := aem.authHelper.SubjectFromToken(token)
		if err != nil {
			aem.log.Errorf("AuthExtractor.Handle error [%v]", err)

			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(rw, r.WithContext(auth.WithSubject(r.Context(), subj)))
	})
}

func (aem *AuthExtractor) isAcceptIssuer(issuer string) bool {
	_, ok := aem.acceptIssuers[issuer]

	return ok
}

func (aem *AuthExtractor) extractTokenString(r *http.Request) (*jwt.Token, error) {
	tokenHeaderStr, errHeader := aem.jwtHTTPHelper.ExtractTokenStringFromRequestHeader(auth.DefaultHeaderName, r)
	tokenCookieStr, errCookie := aem.jwtHTTPHelper.ExtractTokenStringFromRequestCookie(auth.DefaultCookieName, r)
	if errHeader != nil && errCookie != nil {
		return nil, errors.Join(errHeader, errCookie)
	}
	if errHeader == nil && tokenHeaderStr != "" {
		return aem.jwtHelper.ExtractTokenFromString(tokenHeaderStr)
	} else if errCookie == nil && tokenCookieStr != "" {
		return aem.jwtHelper.ExtractTokenFromString(tokenCookieStr)
	}

	return nil, errs.NewCommonError("invalid token (empty token)", nil)
}
