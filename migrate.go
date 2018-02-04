package main

import (
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"fmt"
)

func main() {
	connectString := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@db:5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Printf("Cannot connect to database %v %v", connectString, err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "go_schema_migrations"})
	if err != nil {
		fmt.Printf("Cannot init driver %v", err)
		os.Exit(1)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///go/src/app/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Printf("Cannot init migrate %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments, possible arguments: up, down")
		os.Exit(1)
	}

	action := os.Args[1]
	switch action {
	case "up":
		fmt.Println("Migrating up")
		if err := m.Up(); err != nil {
			fmt.Printf("Cannot migrate %v", err)
		}
	case "down":
		fmt.Println("Migrating down")
		if err := m.Down(); err != nil {
			fmt.Printf("Cannot migrate %v", err)
		}
	default:
		fmt.Println("Not supported argument")
		os.Exit(1)
	}

}
