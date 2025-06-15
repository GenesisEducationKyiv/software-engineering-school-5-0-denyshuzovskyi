package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func InitDB(ctx context.Context, url string, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Error("unable to open database", "error", err)
		return nil, err
	}
	if err = db.PingContext(ctx); err != nil {
		log.Error("unable to ping database", "error", err)
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationPath string, log *slog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("unable to acquire database driver", "error", err)
		return err
	}
	d, err := iofs.New(migrations.Files, migrationPath)
	if err != nil {
		log.Error("unable to set up driver for io/fs#FS", "error", err)
		return err
	}
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Error("unable to set up migrations", "error", err)
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("all migrations have already been applied")
		} else {
			log.Error("unable to apply migrations", "error", err)
			return err
		}
	} else {
		log.Info("migration completed successfully")
	}

	return nil
}
