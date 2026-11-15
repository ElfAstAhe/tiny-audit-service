package repository

import (
	librepo "github.com/ElfAstAhe/go-service-template/pkg/repository/v2"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
)

//goland:noinspection DuplicatedCode
func (da *DataAudit) pgStatementsBuilder() *librepo.Statements[*domain.DataAudit, string] {
	return librepo.NewStatements[*domain.DataAudit, string]().
		WithStatement(librepo.StatementFind, librepo.NewStatement(postgres.SqlDataAuditFind, da.mapDefault)).
		WithStatement(librepo.StatementList, librepo.NewStatement(postgres.SqlDataAuditList, da.mapDefault)).
		WithStatement(librepo.StatementCreate, librepo.NewStatement(postgres.SqlDataAuditCreate, da.mapDefault)).
		WithStatement(librepo.StatementChange, librepo.NewStatement(postgres.SqlDataAuditChange, da.mapDefault)).
		WithStatement(librepo.StatementDelete, librepo.NewStatement[*domain.DataAudit, string](postgres.SqlDataAuditDelete, nil)).
		WithStatement(statementDataGetTail, librepo.NewStatement(postgres.SqlDataAuditTailList, da.mapTail)).
		WithStatement(statementDataListByPeriod, librepo.NewStatement(postgres.SqlDataAuditListByPeriod, da.mapDefault)).
		WithStatement(statementDataListByInstance, librepo.NewStatement(postgres.SqlDataAuditListByInstance, da.mapDefault))
}
