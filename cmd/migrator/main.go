package main

import (
	"chat-bots-api/internal/config"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	cfg := config.MustLoad()

	if migrationsPath == "" {
		panic("flag --migrations-path is required")
	}
	const dbName = "bot_factory"
	mySqlPath := fmt.Sprintf("mysql://%s:%s@tcp(%s:%v)/%s",
		cfg.MySQLConfig.Username,
		cfg.MySQLConfig.Password,
		cfg.MySQLConfig.Host,
		cfg.MySQLConfig.Port,
		dbName,
	)
	m, err := migrate.New("file://"+migrationsPath, mySqlPath)
	if err != nil {
		panic("migration init error: " + err.Error())
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic("migrations error: " + err.Error())
	}
	fmt.Println("migrations applied successfully!")
}
