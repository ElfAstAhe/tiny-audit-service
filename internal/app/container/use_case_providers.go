package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

func (ucc *UseCaseContainer) providerTM(name string) (any, error) {
	dbCnt, err := ucc.GetOrchestrator().GetContainer(DBContainerName)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve container failed", err)
	}
	dbInst, err := container.GetInstance[db.DB](dbCnt, InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return db.NewTxManager(dbInst), nil
}
