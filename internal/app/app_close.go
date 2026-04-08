package app

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Close - метод освобождения ресурсов приложения
//
//	if err := app.Close(); err != nil {
//		log.Errorf("app close error [%v]", err)
//
//		panic(errs.NewAppCommonError("app close failed", err))
//	}
func (app *App) Close() {
	log := app.logger.GetLogger("App.Close")

	log.Info("stop all")
	app.stopAll()

	log.Info("close db connection")
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			log.Errorf("failed to close db [%v]", err)
		}
	}
}

func (app *App) stopAll() {
	log := app.logger.GetLogger("App.stopAll")
	var stopWG sync.WaitGroup

	// stop startup
	stopWG.Add(1)
	go func() {
		defer stopWG.Done()
		log.Info("stoping startup services...")
		defer log.Info("done stop startup services")

		app.stopStartup()
	}()

	// stop servers
	stopWG.Add(1)
	go func() {
		defer stopWG.Done()
		log.Info("stoping servers (http and/or gRPC)...")
		defer log.Info("done stop servers (http and/or gRPC)")

		app.stopServers()
	}()

	stopWG.Wait()
}

func (app *App) stopStartup() {
	log := app.logger.GetLogger("App.stopStartup")
	var stopWG sync.WaitGroup

	if app.authAuditTailCutter != nil {
		stopWG.Add(1)
		go func() {
			defer stopWG.Done()
			log.Info("stoping auth tail cutter...")
			defer log.Info("done stop auth tail cutter")

			if err := app.authAuditTailCutter.Stop(app.config.HTTP.ShutdownTimeout); err != nil {
				log.Errorf("failed to stop auth tail cutter [%v]", err)
			}
			app.authAuditTailCutter = nil
		}()
	}

	if app.dataAuditTailCutter != nil {
		stopWG.Add(1)
		go func() {
			defer stopWG.Done()
			log.Info("stoping data tail cutter...")
			defer log.Info("done stop data tail cutter")

			if err := app.dataAuditTailCutter.Stop(app.config.HTTP.ShutdownTimeout); err != nil {
				log.Errorf("failed to stop data tail cutter [%v]", err)
			}
			app.dataAuditTailCutter = nil
		}()
	}

	log.Info("close telemetry service")
	if app.telemetryShutdown != nil {
		stopWG.Add(1)
		go func() {
			defer stopWG.Done()
			log.Info("shutting down telemetry batcher...")
			defer log.Info("done shutting down telemetry batcher")

			// Используем свежий контекст, так как app.ctx может быть уже отменен
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := app.telemetryShutdown(ctx); err != nil {
				log.Errorf("telemetry shutdown error: [%v]", err)
			} else {
				log.Info("telemetry flushed and closed")
			}

			app.telemetryShutdown = nil
		}()
	}

	stopWG.Wait()
}

func (app *App) stopServers() {
	log := app.logger.GetLogger("App.stopServers")
	var stopWG sync.WaitGroup

	if app.httpServer != nil {
		stopWG.Add(1)
		go func() { // stop HTTP
			defer stopWG.Done()
			log.Info("shutting down http server...")
			defer log.Info("shutdown http server complete")

			ctxTimed, cancelTimed := context.WithTimeout(context.Background(), app.config.HTTP.ShutdownTimeout)
			defer cancelTimed()

			// stop http
			if err := app.httpServer.Shutdown(ctxTimed); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.Error("http server shutdown timed out (force close)")
				} else {
					log.Warnf("shutdown http server with error [%v]", err)
				}
			}

			app.httpServer = nil
		}()
	}

	if app.grpcServer != nil {
		stopWG.Add(1)
		go func() { // stop gRPC
			defer stopWG.Done()
			log.Info("shutting down gRPC server...")
			defer log.Info("shutdown gRPC server complete")

			doneChan := make(chan struct{})
			go func() {
				app.grpcServer.GracefulStop()
				close(doneChan)
			}()
			select {
			case <-doneChan:
				log.Info("shutdown gRPC server gracefully complete")
			case <-time.After(app.config.GRPC.ShutdownTimeout):
				log.Error("shutdown gRPC server gracefully timeout, force stop")
				app.grpcServer.Stop()
			}

			app.grpcServer = nil
		}()
	}

	// Ожидаем завершения остановки всех серверов
	stopWG.Wait()
}
