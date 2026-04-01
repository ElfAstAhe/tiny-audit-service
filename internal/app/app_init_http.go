package app

import "net/http"

func (app *App) initHTTPRouter() error {
	//app.httpRouter = rest.NewAppChiRouter(
	//    app.config,
	//    app.logger,
	//    app.authHelper,
	//    app.health,
	//    nil,
	//    nil,
	//    app.authFacade,
	//    app.userFacade,
	//    app.userAdminFacade,
	//    app.roleAdminFacade,
	//)

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
