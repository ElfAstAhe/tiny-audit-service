package transport

// ErrorDTO описывает объект передачи данных (DTO) для унифицированного
// представления и сериализации сведений об ошибках на транспортном уровне API.
//
// @Description Error dto
type ErrorDTO struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
} // @name ErrorDTO

// NewErrorDTO создает и инициализирует новый экземпляр ErrorDTO на основе переданного HTTP/gRPC статус-кода и текстового сообщения.
func NewErrorDTO(status int, message string) *ErrorDTO {
	return &ErrorDTO{
		Code:    status,
		Message: message,
	}
}

// NewErrorDTOFromError конструирует экземпляр ErrorDTO, автоматически извлекая текстовое описание из переданного объекта ошибки error.
func NewErrorDTOFromError(code int, err error) *ErrorDTO {
	return NewErrorDTO(code, err.Error())
}
