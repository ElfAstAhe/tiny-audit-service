package main

import (
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/app"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
)

// @title           Audit Service API
// @version         1.0
// @description     Сервис аудита
// @termsOfService  Free use

// @contact.name   API Support
// @contact.url    https://github.com/ElfAstAhe/tiny-audit-service
// @contact.email  elf.ast.ahe@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org

// @BasePath  /
func main() {
	var err error
	// 0. startup logger
	var startupLogger logger.Logger
	startupLogger, err = logger.NewZapLogger("info", "")
	if err != nil {
		panic(err)
	}
	defer startupLogger.Close()
	startupLogger = startupLogger.GetLogger("main")

	// 1. Загрузка конфигурации
	startupLogger.Info("init config")
	cfg, err := config.Load()
	if err != nil {
		startupLogger.Errorf("failed to load config: %v", err)

		panic(err)
	}

	// 2. Инициализация логгера
	startupLogger.Info("init logger")
	zapLogger, err := logger.NewZapLogger(cfg.Log.Level, cfg.Log.FilePath)
	if err != nil {
		startupLogger.Errorf("failed to init logger: %v", err)

		panic(err)
	}
	defer zapLogger.Close()

	// 3. Создание приложения
	startupLogger.Info("create application")
	appl, err := app.NewApplication(app.WithConfig(cfg), app.WithLogger(zapLogger))
	if err != nil {
		startupLogger.Errorf("failed to create application: %v", err)

		panic(errs.NewCommonError("failed to create application", err))
	}

	// 4. Инициализация приложения
	startupLogger.Info("init application")
	if err := appl.Init(); err != nil {
		startupLogger.Errorf("failed to init application: %v", err)
		appl.Close()
		panic(errs.NewCommonError("failed to init application", err))
	}

	// 5. Запуск приложения
	startupLogger.Info("run application")
	if err := appl.Run(); err != nil {
		startupLogger.Errorf("failed to run application [%v]", err)
	}

	// 6. Освобождение ресурсов
	startupLogger.Info("close application")
	err = appl.Close()
	if err != nil {
		startupLogger.Errorf("failed to close application [%v]", err)
	}

	startupLogger.Info("shutdown application")
}

//func mainBeforeContainers() {
//    var err error
//    // 0. startup logger
//    var startupLogger logger.Logger
//    startupLogger, err = logger.NewZapLogger("info", "")
//    if err != nil {
//        panic(err)
//    }
//    defer startupLogger.Close()
//    startupLogger = startupLogger.GetLogger("main")
//
//    // 1. Загрузка конфигурации
//    startupLogger.Info("init config")
//    cfg, err := config.Load()
//    if err != nil {
//        startupLogger.Errorf("failed to load config: %v", err)
//
//        panic(err)
//    }
//
//    // 2. Инициализация логгера
//    startupLogger.Info("init logger")
//    zapLogger, err := logger.NewZapLogger(cfg.Log.Level, cfg.Log.FilePath)
//    if err != nil {
//        startupLogger.Errorf("failed to init logger: %v", err)
//
//        panic(err)
//    }
//    defer zapLogger.Close()
//
//    // 3. Создание приложения
//    startupLogger.Info("create application")
//    application := app.NewApp(cfg, zapLogger)
//
//    // 4. Инициализация приложения
//    startupLogger.Info("init application")
//    if err := application.Init(); err != nil {
//        startupLogger.Errorf("failed application initialization [%v]", err)
//        application.Close()
//
//        panic(errs.NewCommonError("failed application initialization", err))
//    }
//
//    // 5. Запуск приложения
//    startupLogger.Info("run application")
//    if err := application.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//        application.Stop()
//        startupLogger.Errorf("run application error [%v]", err)
//    }
//
//    // 6. Ожидание завершения приложения
//    application.WaitForStop()
//
//    // 7.Остановка
//    startupLogger.Info("stop application")
//    application.Stop()
//
//    // 8. Освобождение ресурсов
//    startupLogger.Info("close application")
//    application.Close()
//
//    startupLogger.Info("shutdown application")
//}
