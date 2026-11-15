package repository

import (
	librepo "github.com/ElfAstAhe/go-service-template/pkg/repository/v2"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
)

//goland:noinspection DuplicatedCode
func (aa *AuthAudit) pgStatementsBuilder() *librepo.Statements[*domain.AuthAudit, string] {
	return librepo.NewStatements[*domain.AuthAudit, string]().
		WithStatement(librepo.StatementFind, librepo.NewStatement(postgres.SqlAuthAuditFind, aa.mapDefault)).
		WithStatement(librepo.StatementList, librepo.NewStatement(postgres.SqlAuthAuditList, aa.mapDefault)).
		WithStatement(librepo.StatementCreate, librepo.NewStatement(postgres.SqlAuthAuditCreate, aa.mapDefault)).
		WithStatement(librepo.StatementChange, librepo.NewStatement(postgres.SqlAuthAuditChange, aa.mapDefault)).
		WithStatement(librepo.StatementDelete, librepo.NewStatement[*domain.AuthAudit, string](postgres.SqlAuthAuditDelete, nil)).
		WithStatement(statementAuthGetTail, librepo.NewStatement(postgres.SqlAuthAuditTailList, aa.mapTail)).
		WithStatement(statementAuthListByPeriod, librepo.NewStatement(postgres.SqlAuthAuditListByPeriod, aa.mapDefault)).
		WithStatement(statementAuthListByUsername, librepo.NewStatement(postgres.SqlAuthAuditListByUsername, aa.mapDefault))
}
