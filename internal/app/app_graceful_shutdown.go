package app

import (
	"os"
	"os/signal"
	"syscall"
)

// gracefulShutdown - внутренний метод "агрессивного" закрытия приложения (ctrl+c) + остальные сигналы OS на закрытие
func (app *App) gracefulShutdown() {
	log := app.logger.GetLogger("App.gracefulShutdown")
	defer app.wg.Done()
	log.Info("graceful shutdown start")
	defer log.Info("graceful shutdown finish")

	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// awaiting signal
	select {
	case osSig := <-sig:
		log.Infof("got OS signal [%v], shutdown app", osSig)
		app.cancel()
	case <-app.ctx.Done():
		log.Info("app context done, shutdown app")
		signal.Stop(sig)
	}

	// останавливаем сервера
	log.Info("stoping servers")
	app.stopServers()
}
