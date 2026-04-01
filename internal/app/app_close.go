package app

import (
	"context"
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

	log.Info("close db connection")
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			log.Errorf("failed to close db [%v]", err)
		}
	}

	log.Info("close telemetry service")
	if app.telemetryShutdown != nil {
		log.Info("shutting down telemetry batcher...")
		// Используем свежий контекст, так как app.ctx может быть уже отменен
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.telemetryShutdown(ctx); err != nil {
			log.Errorf("telemetry shutdown error: [%v]", err)
		} else {
			log.Info("telemetry flushed and closed")
		}
	}
}
