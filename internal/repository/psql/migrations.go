package psql

import "github.com/golang-migrate/migrate/v4"

type migrations struct {
	migr *migrate.Migrate
}

func NewMigrations(src, dsn string) (*migrations, error) {
	migr, err := migrate.New(src, dsn)
	if err != nil {
		return nil, err
	}
	return &migrations{
		migr: migr,
	}, nil
}

func (m *migrations) Up() error {
	err := m.migr.Up()
	if err != nil {
		return err
	}
	return nil
}