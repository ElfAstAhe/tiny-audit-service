package repository

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
)

type AuditOwnedRepository[E AuditableEntity[ID], ID comparable, OwnerID comparable] struct {
	next        domain.OwnedRepository[E, ID, OwnerID]
	source      string
	auditClient client.DataAuditClient
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) Find(ctx context.Context, ownerID OwnerID, id ID) (E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) List(ctx context.Context, ownerID OwnerID, limit, offset int) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) ListAll(ctx context.Context, ownerID OwnerID) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) ListAllByOwners(ctx context.Context, ownerIDs ...OwnerID) (map[OwnerID][]E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) Save(ctx context.Context, ownerID OwnerID, owned []E) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) Create(ctx context.Context, ownerID OwnerID, entity E) (E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) Change(ctx context.Context, ownerID OwnerID, entity E) (E, error) {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) DeleteAll(ctx context.Context, ownerID OwnerID) error {
	//TODO implement me
	panic("implement me")
}

func (aor *AuditOwnedRepository[E, ID, OwnerID]) Delete(ctx context.Context, ownerID OwnerID, id ID) error {
	//TODO implement me
	panic("implement me")
}
