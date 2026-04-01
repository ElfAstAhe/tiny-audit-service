package app

// Init инициализирует тяжелые ресурсы: БД, Кеш, Репозитории
func (app *App) Init() error {
	log := app.logger.GetLogger("App.Init")

	log.Info("initializing application resources...")

	log.Info("init helpers")
	if err := app.initHelpers(); err != nil {
		return err
	}

	log.Info("init database")
	if err := app.initDB(); err != nil {
		return err
	}

	log.Info("launch migrations")
	if err := app.migrateDB(); err != nil {
		return err
	}

	log.Info("init telemetry")
	if err := app.initTelemetry(); err != nil {
		return err
	}

	log.Info("init metrics")
	if err := app.initMetrics(); err != nil {
		return err
	}

	log.Info("init dependencies")
	if err := app.initDependencies(); err != nil {
		return err
	}

	log.Info("init startup services")
	if err := app.initStartupServices(); err != nil {
		return err
	}

	log.Info("init health")
	if err := app.initHealth(); err != nil {
		return err
	}

	log.Info("init http router")
	if err := app.initHTTPRouter(); err != nil {
		return err
	}

	log.Info("init http server")
	if err := app.initHTTPServer(); err != nil {
		return err
	}

	log.Info("init gRPC service")
	if err := app.initGRPCService(); err != nil {
		return err
	}

	log.Info("init gRPC server")
	if err := app.initGRPCServer(); err != nil {
		return err
	}

	return nil
}
