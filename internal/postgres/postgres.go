package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/database"
)

var _ database.DB = (*db)(nil)

type db struct {
	db *sqlx.DB
}

func (db *db) WithTx(ctx context.Context, fn func(tx database.Tx) error) error {
	sqlxTx, err := db.db.Beginx()
	if err != nil {
		return err
	}

	t := &tx{db: sqlxTx}
	if err := fn(t); err != nil {
		if err := sqlxTx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return err
	}

	return sqlxTx.Commit()
}

func (db *db) License() database.LicenseStore {
	return newLicenseStore(internal.NewQueryLogger(db.db))
}

func (db *db) User() database.UserStore {
	return newUserStore(internal.NewQueryLogger(db.db))
}

var _ database.Tx = (*tx)(nil)

type tx struct {
	db *sqlx.Tx
}

func (t *tx) License() database.LicenseStore {
	return newLicenseStore(internal.NewQueryLogger(t.db))
}

func (t *tx) User() database.UserStore {
	return newUserStore(internal.NewQueryLogger(t.db))
}

func New(sqlxDB *sqlx.DB) database.DB {
	return &db{
		db: sqlxDB,
	}
}

const (
	maxIdleConns = 25
	maxOpenConns = 100
)

func Open() (*sqlx.DB, error) {
	sqlDB, err := sql.Open("postgres", dsn())
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	for {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return sqlx.NewDb(sqlDB, "postgres"), nil
}

func dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Config.Postgres.Host,
		config.Config.Postgres.Port,
		config.Config.Postgres.User,
		config.Config.Postgres.Password,
		config.Config.Postgres.DB,
	)
}

func Migrate(dir string) error {
	db, err := Open()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+dir, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
