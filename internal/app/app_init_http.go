package app

import (
	"net/http"

	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/rest"
)

func (app *App) initHTTPRouter() error {
	var err error
	app.httpRouter, err = rest.NewAppRouter(
		rest.WithConfig(app.config),
		rest.WithLogger(app.logger),
		rest.WithJwtHelper(app.jwtHelper),
		rest.WithJwtHTTPHelper(app.jwtHTTPHelper),
		rest.WithAuthHelper(app.authHelper),
		rest.WithHealth(app.health),
		rest.WithAuthAuditFacade(app.authFacade),
		rest.WithDataAuditFacade(app.dataFacade),
	)

	return err
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
