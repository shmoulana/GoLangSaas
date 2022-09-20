package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDriver struct {
	DSN        DatabaseDSN
	DBProvider *sgorm.DbProvider
}

func (d SQLiteDriver) Connect(ctx context.Context) (*sgorm.DbProvider, error) {
	conn := make(data.ConnStrings, 1)

	conn.SetDefault(d.DSN.SharedDSN)

	cache := saas.NewCache[string, *sgorm.DbWrap]()
	defer cache.Flush()

	clientProvider := sgorm.ClientProviderFunc(func(ctx context.Context, s string) (*gorm.DB, error) {
		client, _, err := cache.GetOrSet(s, func() (*sgorm.DbWrap, error) {
			var client *gorm.DB
			var err error

			db, err := sql.Open("sqlite3", s)
			if err != nil {
				return nil, err
			}

			db.SetMaxIdleConns(1)
			db.SetMaxOpenConns(1)

			client, err = gorm.Open(&sqlite.Dialector{
				DriverName: sqlite.DriverName,
				DSN:        s,
				Conn:       db,
			})

			return sgorm.NewDbWrap(client), err
		})

		if err != nil {
			return nil, err
		}

		return client.WithContext(ctx).Debug(), err
	})

	tenantStore := &TenantStore{dbProvider: sgorm.NewDbProvider(conn, clientProvider)}

	mr := saas.NewMultiTenancyConnStrResolver(tenantStore, conn)

	db := sgorm.NewDbProvider(mr, clientProvider)

	seeder := seed.NewDefaultSeeder(NewMigrationSeeder(db))
	err := seeder.Seed(context.Background(), seed.AddHost())
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (d SQLiteDriver) GetDB(ctx context.Context) (*gorm.DB, error) {
	if d.DBProvider == nil {
		db, err := d.Connect(ctx)
		if err != nil {
			return nil, err
		}

		d.DBProvider = db
	}
	provider := *d.DBProvider

	db := provider.Get(ctx, "")
	if db == nil {
		return nil, errors.New("failed to get db provider")
	}

	return db, nil
}

func (d SQLiteDriver) GetDSN() DatabaseDSN {
	return d.DSN
}

func NewSQLLiteDriver() DatabaseRepo {
	sharedDsn := "./example.db"
	connStrGen := saas.NewConnStrGenerator("./example-%s.db")

	dsn := DatabaseDSN{
		SharedDSN: sharedDsn,
		TenantDSN: connStrGen,
	}

	return SQLiteDriver{
		DSN: dsn,
	}
}
