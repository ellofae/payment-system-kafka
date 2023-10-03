package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ellofae/payment-system-kafka/config"
	"github.com/ellofae/payment-system-kafka/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func OpenPoolConnection(ctx context.Context, cfg *config.Config) (conn *pgxpool.Pool) {
	log := logger.GetLogger()

	err := ConnectionAttemps(func() error {
		var err error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		conn, err = pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.PostgresDB.User,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName,
			cfg.PostgresDB.SSLmode,
		))

		return err
	}, 3, time.Duration(2)*time.Second)

	if err != nil {
		log.Error("Didn't manage to make connection with database", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Database connection is established successfully.")

	return conn
}

func RunMigrationsUp(ctx context.Context, cfg *config.Config) {
	log := logger.GetLogger()

	db_conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PostgresDB.Host,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.DBName,
		cfg.PostgresDB.SSLmode,
	)

	migration, err := migrate.New("file://client/migrations", db_conn)
	if err != nil {
		log.Error("Unable to get a migrate instance", "error", err.Error())
		os.Exit(1)
	}

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Warn("No changes while migrating")
			return
		} else {
			log.Error("Unable to migrate up", "error", err.Error())
			os.Exit(1)
		}
	}

	log.Info("Migrations are up successfully")
}
