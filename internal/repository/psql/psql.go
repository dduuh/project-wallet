package psql

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	configs "wallet-service/internal/config"
)

type postgresDB struct {
	db *sqlx.DB
}

func New(cfg *configs.Config) (*postgresDB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("failed to open the PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the PostgreSQL: %w", err)
	}

	return &postgresDB{
		db: db,
	}, nil
}

func (p *postgresDB) Up() error {
	driver, err := postgres.WithInstance(p.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to get Migrate instance from source URL: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to Up() migrations: %w", err)
	}

	return nil
}

func (p *postgresDB) Database() *sqlx.DB {
	return p.db
}
