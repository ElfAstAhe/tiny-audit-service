package grpc

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	pb "github.com/ElfAstAhe/tiny-audit-service/pkg/api/grpc/tiny-audit-service/v1"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

const (
	dataAuditClientNameTemplate string = "gRPC-data-audit-client-%s"
)

// DataAuditClient реализует интерфейсы client.AuditClient и client.DataAuditClient,
// предоставляя специализированную обертку над базовым gRPC-клиентом для фиксации событий аутентификации.
//
// Инкапсулирует типизированный AuditClient, генерирует строго типизированный gRPC-сервис pb.AuthAuditServiceClient,
// управляет трассировкой жизненного цикла через логгер и связывает воркеры с удаленным вызовом процедур (RPC).
type DataAuditClient struct {
	*AuditClient[*dto.DataAuditDTO]
	name string
	svc  AuditorService[*pb.DataAuditRequest]
	log  logger.Logger
}

// Контролируем реализацию контрактов SDK аудит-клиента на этапе компиляции
var _ client.AuditClient[*dto.DataAuditDTO] = (*DataAuditClient)(nil)
var _ client.DataAuditClient = (*DataAuditClient)(nil)

// NewDataAuditClient — фабричный конструктор gRPC-клиента аудита авторизационных событий.
// Настраивает опции, проверяет инварианты и инжектирует кастомный метод auditAction в базовый дженерик-клиент.
func NewDataAuditClient(options ...Option[*dto.DataAuditDTO]) (*DataAuditClient, error) {
	// options
	opts := NewOptions[*dto.DataAuditDTO]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewTlCommonError("NewDataAuditClient", "gRPC audit client options validation failed", err)
	}
	// instance
	res := &DataAuditClient{
		name: fmt.Sprintf(dataAuditClientNameTemplate, opts.Pool.Name),
		log:  opts.Pool.Logger.GetLogger(fmt.Sprintf(dataAuditClientNameTemplate, opts.Pool.Name)),
	}
	// embed
	embed, err := NewAuditClient[*dto.DataAuditDTO](append(options, withAuditAction[*dto.DataAuditDTO](res.auditAction))...)
	if err != nil {
		return nil, errs.NewTlCommonError("NewDataAuditClient", "audit client creation failed", err)
	}
	// setup
	res.AuditClient = embed

	return res, nil
}

// Start осуществляет запуск сетевых ресурсов базового клиента и инициализирует строго типизированный gRPC-сервис.
func (dac *DataAuditClient) Start(ctx context.Context) error {
	dac.GetLogger().Debugf("Start start")
	defer dac.GetLogger().Debugf("Start finish")

	if err := dac.AuditClient.Start(ctx); err != nil {
		return err
	}

	dac.svc = pb.NewDataAuditServiceClient(dac.conn)

	return nil
}

// Stop осуществляет фазу контролируемой мягкой остановки клиента, зануляет ссылку на gRPC-сервис и деаллоцирует ресурсы соединения.
func (dac *DataAuditClient) Stop(ctx context.Context) error {
	dac.GetLogger().Debugf("Stop start")
	defer dac.GetLogger().Debugf("Stop finish")

	dac.svc = nil

	return dac.AuditClient.Stop(ctx)
}

// auditAction — это внутренняя прикладная функция (экшен), обрабатывающая извлеченный из очереди лог в рамках конкретного воркера.
// Выполняет маппинг в Protobuf-структуру и осуществляет блокирующий сетевой вызов к удаленному gRPC-серверу аудита.
func (dac *DataAuditClient) auditAction(
	ctx context.Context,
	workerIndex int,
	data *dto.DataAuditDTO,
	token string,
) error {
	dac.GetLogger().Debugf("auditAction start")
	defer dac.GetLogger().Debugf("auditAction finish")

	// transform
	clientDTO := MapDataDtoSDKToGRPC(data)

	// request
	_, err := dac.svc.Audit(
		ctx,
		pb.DataAuditRequest_builder{
			Data: clientDTO,
		}.Build(),
	)
	if err != nil {
		return errs.NewCommonError("data audit action failed", err)
	}

	return nil
}
