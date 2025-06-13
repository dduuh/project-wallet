package psql

import (
	"fmt"
	configs "wallet-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type postgresDB struct {
	db *sqlx.DB
}

func New(cfg *configs.Config) (*postgresDB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgreSQLCfg.Host,
		cfg.PostgreSQLCfg.Port,
		cfg.PostgreSQLCfg.User,
		cfg.PostgreSQLCfg.Password,
		cfg.PostgreSQLCfg.DBName,
		cfg.PostgreSQLCfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresDB{
		db: db,
	}, nil
}

func (p *postgresDB) Up() error {
	driver, err := postgres.WithInstance(p.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (p *postgresDB) Database() *sqlx.DB {
	return p.db
}
