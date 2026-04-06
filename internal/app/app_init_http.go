package app

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/rest"
)

func (app *App) initHTTPRouter() error {
	app.httpRouter = rest.NewAppChiRouter(
		app.config,
		app.logger,
		app.jwtHelper,
		app.jwtHTTPHelper,
		app.authHelper,
		app.health,
		nil,
		nil,
		app.authFacade,
		app.dataFacade,
	)

	return nil
}

func (app *App) initHTTPServer() error {
	app.httpServer = &http.Server{
		Addr:         app.config.HTTP.Address,
		Handler:      app.httpRouter.GetRouter(),
		ReadTimeout:  app.config.HTTP.ReadTimeout,
		WriteTimeout: app.config.HTTP.WriteTimeout,
		IdleTimeout:  app.config.HTTP.IdleTimeout,
	}

	return nil
}
