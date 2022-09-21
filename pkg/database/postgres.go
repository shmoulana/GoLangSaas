package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	_ "github.com/lib/pq"
	"github.com/shmoulana/Redios/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreDriver struct {
	DSN        DatabaseDSN
	DBProvider *sgorm.DbProvider
	conf       *configs.Config
}

func (d PostgreDriver) Connect(ctx context.Context) (*sgorm.DbProvider, error) {
	conn := make(data.ConnStrings, 1)

	conn.SetDefault(d.DSN.SharedDSN)

	cache := saas.NewCache[string, *sgorm.DbWrap]()
	defer cache.Flush()

	clientProvider := sgorm.ClientProviderFunc(func(ctx context.Context, s string) (*gorm.DB, error) {
		client, _, err := cache.GetOrSet(s, func() (*sgorm.DbWrap, error) {
			var client *gorm.DB
			var err error

			psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
				"password=%s dbname=%s sslmode=disable", d.conf.DBHost, d.conf.DBPort, d.conf.DBUser, d.conf.DBPassword, d.conf.DBName)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				return nil, err
			}

			client, err = gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}))

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

func (d PostgreDriver) GetDB(ctx context.Context) (*gorm.DB, error) {
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

func (d PostgreDriver) GetDSN() DatabaseDSN {
	return d.DSN
}

func NewPostgreDriver(conf configs.Config) DatabaseRepo {
	hostDbName := conf.DBNameTenant

	psqlInfoTenant := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, hostDbName)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)

	connStrGen := saas.NewConnStrGenerator(psqlInfoTenant)

	dsn := DatabaseDSN{
		SharedDSN: psqlInfo,
		TenantDSN: connStrGen,
	}

	//open without db name
	psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	query := fmt.Sprintf("SELECT 'CREATE DATABASE %s' as column WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')", conf.DBName, conf.DBName)

	row := db.QueryRowContext(context.Background(), query)

	var column *string

	err = row.Scan(
		&column,
	)

	if err != nil {
		panic(err)
	}

	if column != nil {
		_, err = db.ExecContext(context.Background(), *column)
		if err != nil {
			panic(err)
		}
	}

	db.Close()

	return PostgreDriver{
		DSN:  dsn,
		conf: &conf,
	}
}
