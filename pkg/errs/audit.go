package errs

import (
	"fmt"
)

type AuditError struct {
	op  string
	msg string
	err error
}

// Гарантируем строгое соответствие встроенному интерфейсу error на этапе компиляции
var _ error = (*AuditError)(nil)

// NewAuditError — фабричный конструктор базовой ошибки слоя аудита.
func NewAuditError(op string, msg string, err error) *AuditError {
	return &AuditError{op: op, msg: msg, err: err}
}

// Error форматирует и возвращает развернутый строковый паспорт ошибки.
// Динамически выстраивает текстовое сообщение, изолируя операцию квадратными скобками для читаемости в логах.
func (e *AuditError) Error() string {
	msg := "AUDIT: error"
	if e.op != "" {
		msg = fmt.Sprintf("%s: [%s]", msg, e.op)
	}
	if e.msg != "" {
		msg = fmt.Sprintf("%s %s", msg, e.msg)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.err)
	}

	return msg
}

// Unwrap возвращает исходную ошибку драйвера или пула соединений, инициировавшую сбой.
// Требуется для корректной работы функций errors.Is и errors.As во внешних классификаторах транспорта.
func (e *AuditError) Unwrap() error {
	return e.err
}
