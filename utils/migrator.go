package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Start migrations")

	var (
		path       string
		dbUser     string
		dbPassword string
		dbAddres   string
		dbName     string
		sslMode    string
	)
	flag.StringVar(&path, "migrations_path", "", "Path for migrations")
	flag.StringVar(&dbUser, "db_user", "", "User for migrations")
	flag.StringVar(&dbPassword, "db_password", "", "Password for migrations")
	flag.StringVar(&dbAddres, "db_addres", "", "Addres for migrations")
	flag.StringVar(&dbName, "db_name", "", "Db name for migrations")
	flag.StringVar(&sslMode, "ssl_mode", "", "Ssl mode for migrations")

	flag.Parse()

	m, err := migrate.New(
		"file://"+path,
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			dbUser,
			dbPassword,
			dbAddres,
			dbName,
			sslMode,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
