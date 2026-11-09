package container

// Константы, определяющие уникальные имена (идентификаторы) для каждого DI-контейнера в составе оркестратора.
const (
	// OrchestratorName определяет имя корневого контейнера-оркестратора.
	OrchestratorName string = "orchestrator-container"

	// AppContainerName определяет имя контейнера ядра приложения (логирование, базовые сущности).
	AppContainerName string = "app-container"

	// ToolsContainerName определяет имя контейнера вспомогательных инструментов и утилит.
	ToolsContainerName string = "tools-container"

	// ClientContainerName определяет имя контейнера внешних клиентов и брокеров сообщений (AMQP/Kafka).
	ClientContainerName string = "client-container"

	// InfraContainerName определяет имя контейнера базовой инфраструктуры и телеметрии.
	InfraContainerName string = "infra-container"

	// DBContainerName определяет имя контейнера для управления пулами и коннектами баз данных.
	DBContainerName string = "db-container"

	// RepositoryContainerName определяет имя контейнера хранения данных (слой репозиториев/DAL).
	RepositoryContainerName string = "repository-container"

	// UseCaseContainerName определяет имя контейнера прикладных сценариев и интеракторов бизнес-логики (BLL).
	UseCaseContainerName string = "useCase-container"

	// FacadeContainerName определяет имя контейнера фасадов приложения слоя представления.
	FacadeContainerName string = "facade-container"

	// ServiceContainerName определяет имя контейнера внутренних вспомогательных сервисов.
	ServiceContainerName string = "service-container"

	// HTTPContainerName определяет имя контейнера серверов и роутеров транспортного HTTP-слоя.
	HTTPContainerName string = "http-container"

	// GRPCContainerName определяет имя контейнера серверов транспортного gRPC-слоя.
	GRPCContainerName string = "gRPC-container"

	// WorkerContainerName определяет имя контейнера фоновых обработчиков, слушателей очередей и планировщиков.
	WorkerContainerName string = "worker-container"
)
