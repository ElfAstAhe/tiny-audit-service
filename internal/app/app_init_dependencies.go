package app

import (
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/trace"
)

//goland:noinspection DuplicatedCode
func (app *App) initDependencies() error {
	var err error
	// transaction manager
	app.tm = db.NewTxManager(app.db)
	var (
		authAudit domain.AuthAuditRepository
		dataAudit domain.DataAuditRepository
	)
	// repositories
	{
		authAudit, err = postgres.NewAuthAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		authAudit = trace.NewAuthAuditTraceRepository(metrics.NewAuthAuditMetricsRepository(authAudit))

		dataAudit, err = postgres.NewDataAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		dataAudit = trace.NewDataAuditTraceRepository(metrics.NewDataAuditMetricsRepository(dataAudit))
	}

	/*
	   var (
	       roleAdminRepo      domain.RoleAdminRepository
	       userRolesAdminRepo domain.UserRolesAdminRepository
	       userAdminRepo      domain.UserAdminRepository

	       roleRepo      domain.RoleRepository
	       userRolesRepo domain.UserRolesRepository
	       userRepo      domain.UserRepository

	       loginUC       usecase.LoginUseCase
	       loginSimpleUC usecase.LoginSimpleUseCase

	       registerUC       usecase.RegisterUseCase
	       profileUC        usecase.ProfileUseCase
	       changeKeysUC     usecase.ChangeKeysUseCase
	       changePasswordUC usecase.ChangePasswordUseCase

	       userAdminGetUC       usecase.UserAdminGetUseCase
	       userAdminGetByNameUC usecase.UserAdminGetNameUseCase
	       userAdminListUC      usecase.UserAdminListUseCase
	       userAdminSaveUC      usecase.UserAdminSaveUseCase
	       userAdminDeleteUC    usecase.UserAdminDeleteUseCase

	       roleAdminGetUC       usecase.RoleAdminGetUseCase
	       roleAdminGetByNameUC usecase.RoleAdminGetNameUseCase
	       roleAdminListUC      usecase.RoleAdminListUseCase
	       roleAdminSaveUC      usecase.RoleAdminSaveUseCase
	       roleAdminDeleteUC    usecase.RoleAdminDeleteUseCase
	   )
	   // repositories
	   {
	       // role repo
	       roleRepo, err = postgres.NewRolePgRepository(app.db, app.db)
	       if err != nil {
	           return err
	       }
	       roleRepo = repository.NewRoleTraceRepository(repository.NewRoleMetricsRepository(roleRepo))

	       // user roles repo
	       userRolesRepo, err = postgres.NewUserRolesPgRepository(app.db, app.db)
	       if err != nil {
	           return err
	       }
	       // user roles metrics repo
	       userRolesRepo = repository.NewUserRolesTraceRepository(repository.NewUserRolesMetricsRepository(userRolesRepo))
	       // user repo
	       userRepo, err = postgres.NewUserPgRepository(app.db, app.db, app.hashCipher, app.cipherHelper, userRolesRepo)
	       if err != nil {
	           return err
	       }
	       userRepo = repository.NewUserTraceRepository(repository.NewUserMetricsRepository(userRepo))
	       // role admin repo
	       roleAdminRepo, err = postgres.NewRoleAdminPgRepository(app.db, app.db)
	       if err != nil {
	           return err
	       }
	       roleAdminRepo = repository.NewRoleAdminTraceRepository(repository.NewRoleAdminMetricsRepository(roleAdminRepo))
	       // user roles admin repo
	       userRolesAdminRepo, err = postgres.NewUserRolesAdminPgRepository(app.db, app.db)
	       if err != nil {
	           return err
	       }
	       userRolesAdminRepo = repository.NewUserRolesAdminTraceRepository(repository.NewUserRolesAdminMetricsRepository(userRolesAdminRepo))
	       // user admin repo
	       userAdminRepo, err = postgres.NewUserAdminPgRepository(app.db, app.db, app.cipherHelper, app.hashCipher, userRolesAdminRepo)
	       if err != nil {
	           return err
	       }
	       userAdminRepo = repository.NewUserAdminTraceRepository(repository.NewUserAdminMetricsRepository(userAdminRepo))
	   }
	   // use cases
	   {
	       // auth
	       loginUC = telemetry.NewLoginTraceUseCase("LoginUseCase", usecase.NewLoginUseCase(app.hashCipher, app.keysHelper, app.authHelper, userRepo))
	       loginSimpleUC = telemetry.NewLoginSimpleTraceUseCase("LoginSimpleUseCase", usecase.NewLoginSimpleUseCase(app.hashCipher, app.authHelper, userRepo))
	       // users
	       registerUC = telemetry.NewRegisterTraceUseCase("RegisterUseCase", usecase.NewRegisterUseCase(app.tm, app.hashCipher, app.keysHelper, userRepo))
	       profileUC = telemetry.NewProfileTraceUseCase("ProfileUseCase", usecase.NewProfileUseCase(userRepo))
	       changeKeysUC = telemetry.NewChangeKeysTraceUseCase("ChangeKeysUseCase", usecase.NewChangeKeysUseCase(app.keysHelper, app.tm, userRepo))
	       changePasswordUC = telemetry.NewChangePasswordTraceUseCase("ChangePasswordUseCase", usecase.NewChangePasswordUseCase(app.hashCipher, app.tm, userRepo))
	       // role admin
	       roleAdminGetUC = telemetry.NewRoleAdminGetTraceUseCase("RoleAdminGetUseCase", usecase.NewRoleAdminGetUseCase(roleAdminRepo))
	       roleAdminGetByNameUC = telemetry.NewROleAdminGetNameTraceUseCase("RoleAdminGetNameUseCase", usecase.NewRoleAdminGetNameUseCase(roleAdminRepo))
	       roleAdminListUC = telemetry.NewRoleAdminListTraceUseCase("RoleAdminListUseCase", usecase.NewRoleAdminListUseCase(roleAdminRepo, app.config.App.MaxListLimit))
	       roleAdminSaveUC = telemetry.NewRoleAdminSaveTraceUseCase("RoleAdminSaveUseCase", usecase.NewRoleAdminSaveUseCase(app.tm, roleAdminRepo))
	       roleAdminDeleteUC = telemetry.NewRoleAdminDeleteTraceUseCase("RoleAdminDeleteUseCase", usecase.NewRoleAdminDeleteUseCase(app.tm, roleAdminRepo))
	       // user admin
	       userAdminGetUC = telemetry.NewUserAdminGetTraceUseCase("UserAdminGetUseCase", usecase.NewUserAdminGetUseCase(userAdminRepo))
	       userAdminGetByNameUC = telemetry.NewUserAdminGetNameTraceUseCase("UserAdminGetNameUseCase", usecase.NewUserAdminGetNameUseCase(userAdminRepo))
	       userAdminListUC = telemetry.NewUserAdminListTraceUseCase("UserAdminListUseCase", usecase.NewUserAdminListUseCase(userAdminRepo, app.config.App.MaxListLimit))
	       userAdminSaveUC = telemetry.NewUserAdminSaveTraceUseCase("UserAdminSaveUseCase", usecase.NewUserAdminSaveUseCase(app.tm, app.hashCipher, app.keysHelper, userAdminRepo))
	       userAdminDeleteUC = telemetry.NewUserAdminDeleteTraceUseCase("UserAdminDeleteUseCase", usecase.NewUserAdminDeleteUseCase(app.tm, userAdminRepo))
	   }
	   // facades
	   {
	       // auth
	       app.authFacade = facade.NewAuthFacade(app.jwtHelper, loginUC, loginSimpleUC)
	       app.userFacade = facade.NewUserFacade(
	           app.authHelper,
	           registerUC,
	           profileUC,
	           changePasswordUC,
	           changeKeysUC,
	       )
	       app.roleAdminFacade = facade.NewRoleAdminFacade(
	           app.authHelper,
	           roleAdminGetUC,
	           roleAdminGetByNameUC,
	           roleAdminListUC,
	           roleAdminSaveUC,
	           roleAdminDeleteUC,
	           app.config.App.MaxListLimit,
	       )
	       app.userAdminFacade = facade.NewUserAdminFacade(
	           app.authHelper,
	           userAdminGetUC,
	           userAdminGetByNameUC,
	           userAdminListUC,
	           userAdminSaveUC,
	           userAdminDeleteUC,
	           app.config.App.MaxListLimit,
	       )
	   }
	*/
	return nil
}
