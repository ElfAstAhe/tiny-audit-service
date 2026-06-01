package rest

import (
	libhttp "github.com/ElfAstAhe/go-service-template/pkg/transport/http"
)

func mapToHTTPStatus(err error) int {
	return libhttp.MapToHTTPStatus(err)
}
