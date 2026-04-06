package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// gracefulShutdown - внутренний метод "агрессивного" закрытия приложения (ctrl+c) + остальные сигналы OS на закрытие
func (app *App) gracefulShutdown() {
	defer app.wg.Done()

	log := app.logger.GetLogger("App.gracefulShutdown")
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// awaiting signal
	select {
	case <-sig:
		app.cancel()
	case <-app.ctx.Done():
		signal.Stop(sig)
	}

	var srvShutdownWg sync.WaitGroup

	srvShutdownWg.Add(1)
	go func() { // stop HTTP
		defer srvShutdownWg.Done()

		ctxTimed, cancelTimed := context.WithTimeout(context.Background(), app.config.HTTP.ShutdownTimeout)
		defer cancelTimed()

		// stop http
		log.Info("shutdown http server...")
		if err := app.httpServer.Shutdown(ctxTimed); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Error("http server shutdown timed out (force close)")
			} else {
				log.Warnf("shutdown http server with error [%v]", err)
			}
		}
		log.Info("shutdown http server complete")
	}()

	srvShutdownWg.Add(1)
	go func() { // stop gRPC
		defer srvShutdownWg.Done()

		log.Info("shutdown gRPC server...")

		doneChan := make(chan struct{})
		go func() {
			app.grpcServer.GracefulStop()
			close(doneChan)
		}()
		select {
		case <-doneChan:
			log.Info("shutdown gRPC server complete")
		case <-time.After(app.config.GRPC.ShutdownTimeout):
			log.Error("gRPC graceful shutdown timed out: forcing stop")
			app.grpcServer.Stop()
		}
	}()

	// Ожидаем завершения остановки всех серверов
	srvShutdownWg.Wait()
}
