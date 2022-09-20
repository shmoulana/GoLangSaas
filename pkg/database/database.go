package database

import (
	"context"

	"github.com/go-saas/saas"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	"github.com/shmoulana/Redios/configs"
	"gorm.io/gorm"
)

type Database struct {
	Config configs.Config
}

type DatabaseRepo interface {
	Connect(ctx context.Context) (*sgorm.DbProvider, error)
	GetDB(ctx context.Context) (*gorm.DB, error)
	GetDSN() DatabaseDSN
}

type DatabaseDSN struct {
	SharedDSN string
	TenantDSN *saas.DefaultConnStrGenerator
}

type MigrationSeeder struct {
	dbProvider sgorm.DbProvider
}

func NewMigrationSeeder(dbProvider sgorm.DbProvider) *MigrationSeeder {
	return &MigrationSeeder{dbProvider: dbProvider}
}

func (m *MigrationSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := m.dbProvider.Get(ctx, "")
	if sCtx.TenantId == "" {
		//host add tenant database
		err := db.AutoMigrate(&Tenant{}, &TenantConn{})
		if err != nil {
			return err
		}
	}
	// err := db.AutoMigrate(&Post{})
	// if err != nil {
	// 	return err
	// }
	return nil
}
