package domain

import (
	"fmt"
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// DataAudit представляет собой доменную модель (Entity) для хранения,
// валидации и логирования записей аудита изменений сущностей данных.
type DataAudit struct {
	ID               string
	Source           string
	EventDate        time.Time
	Event            string
	Status           string
	RequestID        string
	TraceID          string
	Username         string
	InternalTypeName string
	TypeName         string
	TypeDescription  string
	InstanceID       string
	InstanceName     string
	Values           []*DataAuditValue
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

var _ libdomain.Entity[string] = (*DataAudit)(nil)
var _ commonAudit = (*DataAudit)(nil)

// NewEmptyDataAudit возвращает новый экземпляр DataAudit с пустым слайсом изменений и предустановленными метками времени.
func NewEmptyDataAudit() *DataAudit {
	return &DataAudit{
		Values:    make([]*DataAuditValue, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetID возвращает уникальный строковый идентификатор текущей записи аудита данных.
func (da *DataAudit) GetID() string {
	return da.ID
}

// SetID принудительно устанавливает уникальный строковый идентификатор для текущей записи аудита данных.
func (da *DataAudit) SetID(id string) {
	da.ID = id
}

// IsExists проверяет, была ли данная сущность уже сохранена в базе данных (имеет ли она заполненный ID).
func (da *DataAudit) IsExists() bool {
	return da.ID != ""
}

// ValidateCreate осуществляет проверку бизнес-правил доменного слоя и вложенных значений полей, необходимых для создания записи.
func (da *DataAudit) ValidateCreate() error {
	if da.ID != "" {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "id must be empty", nil)
	}
	if err := validateCommon(da); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "common audit validate", err)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "event validate", err)
	}
	if len(da.Values) > 0 {
		for index, value := range da.Values {
			if err := value.ValidateCreate(); err != nil {
				return errs.NewBllValidateError("DataAudit.ValidateCreate", fmt.Sprintf("values[%d] validate failed", index), err)
			}
		}
	}

	return nil
}

// ValidateChange осуществляет проверку бизнес-правил доменного слоя и вложенных значений полей, необходимых для обновления записи.
func (da *DataAudit) ValidateChange() error {
	if da.ID == "" {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "id must not be empty", nil)
	}
	if err := validateCommon(da); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "common audit validate", err)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "event validate", err)
	}
	if len(da.Values) > 0 {
		for index, value := range da.Values {
			if err := value.ValidateChange(); err != nil {
				return errs.NewBllValidateError("DataAudit.ValidateChange", fmt.Sprintf("values[%d] validate failed", index), err)
			}
		}
	}

	return nil
}

// BeforeCreate выполняет хук жизненного цикла перед сохранением: генерирует UUIDv7 и выставляет метки времени.
func (da *DataAudit) BeforeCreate() error {
	if err := libdomain.AssignUUIDv7(da); err != nil {
		return errs.NewBllError("DataAudit.BeforeCreate", "default before create failed", err)
	}

	if da.CreatedAt.IsZero() {
		da.CreatedAt = time.Now()
	}
	da.UpdatedAt = time.Now()

	return nil
}

// BeforeChange выполняет хук жизненного цикла перед обновлением сущности, актуализируя метку времени UpdatedAt.
func (da *DataAudit) BeforeChange() error {
	da.UpdatedAt = time.Now()

	return nil
}

// GetSource возвращает наименование компонента-источника события изменения данных.
func (da *DataAudit) GetSource() string {
	return da.Source
}

// GetEventDate возвращает точную дату и время фиксации события мутации данных.
func (da *DataAudit) GetEventDate() time.Time {
	return da.EventDate
}

// GetEvent возвращает строковый маркер типа события изменения данных (create/change/remove).
func (da *DataAudit) GetEvent() string {
	return da.Event
}

// GetStatus возвращает статус завершения операции (success/fail).
func (da *DataAudit) GetStatus() string {
	return da.Status
}

// GetRequestID возвращает сквозной идентификатор HTTP/gRPC запроса.
func (da *DataAudit) GetRequestID() string {
	return da.RequestID
}

// GetTraceID возвращает идентификатор распределенной трассировки OpenTelemetry.
func (da *DataAudit) GetTraceID() string {
	return da.TraceID
}

// GetUsername возвращает имя пользователя, совершившего изменение.
func (da *DataAudit) GetUsername() string {
	return da.Username
}
