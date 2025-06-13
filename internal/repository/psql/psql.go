package psql

import (
	"fmt"
	configs "wallet-service/internal/config"

	"github.com/jmoiron/sqlx"
)

func New(cfg *configs.Config) (*sqlx.DB, error) {
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

	return db, nil
}
