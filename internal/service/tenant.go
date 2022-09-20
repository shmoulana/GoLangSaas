package service

import (
	"context"

	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"github.com/google/uuid"
	"github.com/shmoulana/Redios/pkg/database"
	"github.com/shmoulana/Redios/pkg/dto"
)

type TenantService struct {
	DatabaseRepo database.DatabaseRepo
}

func (s TenantService) CreateTenant(ctx context.Context, payload dto.TenantRequestV1) error {
	ctx = saas.NewCurrentTenant(ctx, "", "")
	db, err := s.DatabaseRepo.GetDB(ctx)
	if err != nil {
		return err
	}

	isCreateNewTenant := payload.SeparateDb

	t := &database.Tenant{
		ID:          uuid.New().String(),
		Name:        payload.Name,
		DisplayName: payload.Name,
	}

	if isCreateNewTenant {
		t3Conn, _ := s.DatabaseRepo.GetDSN().TenantDSN.Gen(ctx, saas.NewBasicTenantInfo(t.ID, t.Name))
		t.Conn = []database.TenantConn{
			{Key: data.Default, Value: t3Conn}, // use tenant3.db
		}
	}

	err = db.Model(t).Create(t).Error
	if err != nil {
		return err
	}

	return nil
}
