package startup

import (
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SslMode  string `yaml:"ssl_mode"`
}

func DatabaseConnect(config DatabaseConfig) (*ent.Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.DBName, config.Password, config.SslMode)
	driver, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, err
	}

	err = driver.DB().Ping()
	if err != nil {
		return nil, err
	}

	return ent.NewClient(ent.Driver(driver)), nil
}
