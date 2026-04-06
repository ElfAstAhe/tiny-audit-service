package transport

import (
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
	transperrs "github.com/ElfAstAhe/tiny-audit-service/internal/transport/errs"
)

func IsBadRequest(err error) bool {
	var (
		errInvalidArgument *errs.InvalidArgumentError
		errBllValidate     *domerrs.BllValidateError
		errTrMapping       *transperrs.TrMappingError
	)

	return errors.As(err, &errInvalidArgument) ||
		errors.As(err, &errBllValidate) ||
		errors.As(err, &errTrMapping)
}

func IsUnauthorized(err error) bool {
	var (
		errBllUnauthorized *domerrs.BllUnauthorizedError
	)

	return errors.As(err, &errBllUnauthorized)
}

func IsForbidden(err error) bool {
	var (
		errBllForbidden *domerrs.BllForbiddenError
	)

	return errors.As(err, &errBllForbidden)
}

func IsNotFound(err error) bool {
	var (
		errBllNotFound *domerrs.BllNotFoundError
		errDalNotFound *errs.DalNotFoundError
	)

	return errors.As(err, &errBllNotFound) ||
		errors.As(err, &errDalNotFound)
}

func IsConflict(err error) bool {
	var (
		errBllUnique        *domerrs.BllUniqueError
		errDalAlreadyExists *errs.DalAlreadyExistsError
	)

	return errors.As(err, &errBllUnique) ||
		errors.As(err, &errDalAlreadyExists)
}

func IsGone(err error) bool {
	var (
		errDalSoftDeleted *errs.DalSoftDeletedError
	)

	return errors.As(err, &errDalSoftDeleted)
}
