package repository

import (
	"context"

	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"github.com/shmoulana/Redios/pkg/database"
)

type TenantRepository struct {
	databaseRepo database.DatabaseRepo
}

func (r TenantRepository) CreateTenant(ctx context.Context, tenant database.Tenant, isNew bool) (*string, error) {
	ctx = saas.NewCurrentTenant(ctx, "", "")
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	if isNew {
		t3Conn := r.databaseRepo.GetTenantDSN(ctx, saas.NewBasicTenantInfo(tenant.ID, tenant.Name))
		tenant.Conn = []database.TenantConn{
			{Key: data.Default, Value: t3Conn},
		}

		if r.databaseRepo.GetDriver() == "postgres" {
			err := r.databaseRepo.CreateDatabase(ctx, tenant.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	result := db.Model(&tenant).Create(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tenant.ID, nil
}

func (r TenantRepository) UpdateTenant(ctx context.Context, tenant database.Tenant) (*string, error) {
	ctx = saas.NewCurrentTenant(ctx, "", "")
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Save(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}

	// logging
	// rowsAffected := result.RowsAffected

	return &tenant.ID, nil
}

func (r TenantRepository) DeleteTenant(ctx context.Context, tenant database.Tenant) (*string, error) {
	ctx = saas.NewCurrentTenant(ctx, "", "")
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Delete(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}

	// logging
	// rowsAffected := result.RowsAffected

	return &tenant.ID, nil
}

func (r TenantRepository) Find(ctx context.Context) ([]database.Tenant, error) {
	var tenants []database.Tenant
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Find(&tenants)

	if result.Error != nil {
		return nil, result.Error
	}

	return tenants, nil
}

func (r TenantRepository) FindById(ctx context.Context, id string) (*database.Tenant, error) {
	var tenant database.Tenant
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Where("id = ?", id).Find(&tenant)
	if result.Error != nil {
		return nil, result.Error
	}

	result.Scan(&tenant)

	return &tenant, nil
}

func NewTenantRepository(db database.DatabaseRepo) TenantRepository {
	return TenantRepository{
		databaseRepo: db,
	}
}
