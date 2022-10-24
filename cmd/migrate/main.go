package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/app/startup"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/migrate"
)

func main() {
	configFile := flag.String("config", "", "path to configuration file")
	migrationName := flag.String("name", "", "name of current migration")
	flag.Parse()

	config, err := startup.NewMigrationConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	ctx := context.Background()

	dir, err := atlas.NewLocalDir("./internal/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeInspect),
		schema.WithDialect(dialect.Postgres),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Host,
		config.Database.Port, config.Database.DBName, config.Database.SslMode)
	err = migrate.NamedDiff(ctx, url, *migrationName, opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
