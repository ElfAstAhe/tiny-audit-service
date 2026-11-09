package transport

import (
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// IsBadRequest проверяет, относится ли переданная ошибка к категории некорректных входных данных
// (ошибки валидации бизнес-логики, невалидные аргументы или ошибки маппинга транспортного слоя).
func IsBadRequest(err error) bool {
	var (
		errInvalidArgument *errs.InvalidArgumentError
		errBllValidate     *errs.BllValidateError
		errTlMapping       *errs.TlMappingError
	)

	return errors.As(err, &errInvalidArgument) ||
		errors.As(err, &errBllValidate) ||
		errors.As(err, &errTlMapping)
}

// IsUnauthorized проверяет, представляет ли переданная ошибка нарушение аутентификации субъекта
// (отсутствие, истечение срока действия или некорректность предоставленного токена).
func IsUnauthorized(err error) bool {
	var (
		errBllUnauthorized *errs.BllUnauthorizedError
	)

	return errors.As(err, &errBllUnauthorized)
}

// IsForbidden проверяет, представляет ли переданная ошибка нарушение прав доступа (RBAC),
// когда субъект успешно аутентифицирован, но не обладает необходимыми ролями для выполнения операции.
func IsForbidden(err error) bool {
	var (
		errBllForbidden *errs.BllForbiddenError
	)

	return errors.As(err, &errBllForbidden)
}

// IsNotFound проверяет, указывает ли ошибка на отсутствие запрашиваемого ресурса
// в бизнес-логике или непосредственно в самом хранилище данных (DAL).
func IsNotFound(err error) bool {
	var (
		errBllNotFound *errs.BllNotFoundError
		errDalNotFound *errs.DalNotFoundError
	)

	return errors.As(err, &errBllNotFound) ||
		errors.As(err, &errDalNotFound)
}

// IsConflict проверяет, возникла ли при выполнении операции коллизия данных
// (нарушение ограничений уникальности бизнес-слоя или попытка повторного создания уже существующей записи в БД).
func IsConflict(err error) bool {
	var (
		errBllUnique        *errs.BllUniqueError
		errDalAlreadyExists *errs.DalAlreadyExistsError
	)

	return errors.As(err, &errBllUnique) ||
		errors.As(err, &errDalAlreadyExists)
}

// IsGone проверяет, указывает ли ошибка на то, что запрашиваемый ресурс был ранее удален
// (например, при использовании механизмов мягкого удаления Soft Delete).
func IsGone(err error) bool {
	var (
		errDalSoftDeleted *errs.DalSoftDeletedError
	)

	return errors.As(err, &errDalSoftDeleted)
}
