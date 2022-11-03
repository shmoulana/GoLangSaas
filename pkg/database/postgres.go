package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
	_ "github.com/jackc/pgx/v5"
	"github.com/shmoulana/Redios/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreDriver struct {
	DSN          DatabaseDSN
	DBProvider   *sgorm.DbProvider
	conf         *configs.Config
	tenantPrefix string
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

			client, err = gorm.Open(postgres.New(postgres.Config{
				DSN: psqlInfo,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})

			return sgorm.NewDbWrap(client), err
		})

		if err != nil {
			return nil, err
		}

		return client.WithContext(ctx), err
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

func (d PostgreDriver) GetDriver() string {
	return "postgres"
}

func (d PostgreDriver) GetTenantDSN(ctx context.Context, tenantInfo saas.TenantInfo) string {
	t3Conn, _ := d.DSN.TenantDSN.Gen(ctx, tenantInfo)

	return strings.ReplaceAll(t3Conn, "-", "_")
}

func (d PostgreDriver) CreateDatabase(ctx context.Context, id string) error {
	//open without db name
	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
	// 	"password=%s sslmode=disable", d.conf.DBHost, d.conf.DBPort, d.conf.DBUser, d.conf.DBPassword)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: d.DSN.SharedDSN,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	dbName := fmt.Sprintf(d.tenantPrefix, id)
	dbName = strings.ReplaceAll(dbName, "-", "_")

	var column *string

	// Find database if not exist return create database query
	result := db.Raw(fmt.Sprintf("SELECT 'CREATE DATABASE %s' as column WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')", dbName, dbName))
	err = result.Error
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	result.Scan(&column)

	if column != nil {
		result := db.Exec(*column)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func NewPostgreDriver(conf configs.Config) DatabaseRepo {
	hostDbName := conf.DBNameTenant + "_%s"

	psqlInfoTenant := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, hostDbName)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)

	connStrGen := saas.NewConnStrGenerator(psqlInfoTenant)

	dsn := DatabaseDSN{
		SharedDSN: psqlInfo,
		TenantDSN: connStrGen,
	}

	// //open without db name
	// psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable dbname=%s", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, "")

	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DSN: psqlInfo,
	// }), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// var column *string

	// // Find database if not exist return create database query
	// result := db.Raw(fmt.Sprintf("SELECT 'CREATE DATABASE %s' as column WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')", conf.DBName, conf.DBName))
	// err = result.Error
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		panic(err)
	// 	}
	// }

	// result.Scan(&column)

	// if column != nil {
	// 	result := db.Exec(*column)
	// 	if result.Error != nil {
	// 		panic(result.Error)
	// 	}
	// }

	return PostgreDriver{
		DSN:          dsn,
		conf:         &conf,
		tenantPrefix: hostDbName,
	}
}
