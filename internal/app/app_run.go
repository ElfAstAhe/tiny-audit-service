package app

import (
	"errors"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Run запускает серверы (HTTP/gRPC) и блокирует поток до сигнала завершения
func (app *App) Run() error {
	log := app.logger.GetLogger("App.Run")

	log.Info("start graceful shutdown")
	app.wg.Add(1)
	go app.gracefulShutdown()

	var eg errgroup.Group
	log.Info("start servers...")
	// http
	eg.Go(func() error {
		if err := app.launchHTTPServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})
	// gRPC
	eg.Go(func() error {
		if err := app.launchGRPCServer(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Errorf("Error starting gRPC server with error [%v]", err)

			return err
		}

		return nil
	})

	return eg.Wait()
}

func (app *App) launchHTTPServer() error {
	log := app.logger.GetLogger("App.launchHTTPServer")
	if app.config.HTTP.Secure {
		log.Info("enable https")

		return app.httpServer.ListenAndServeTLS(app.config.HTTP.CertificatePath, app.config.HTTP.PrivateKeyPath)
	}

	log.Info("enable http")

	return app.httpServer.ListenAndServe()
}

func (app *App) launchGRPCServer() error {
	log := app.logger.GetLogger("App.launchGRPCServer")
	ls, err := net.Listen("tcp", app.config.GRPC.Address)
	if err != nil {
		return err
	}
	log.Info("start gRPC server")

	return app.grpcServer.Serve(ls)
}
