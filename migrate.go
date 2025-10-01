package main

import (
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5/stdlib"
)

func RunMigrations(pool *pgxpool.Pool) error {
    db := stdlib.OpenDBFromPool(pool)
    
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create database driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://internal/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    // Run all up migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("could not run migrations: %w", err)
    }

    version, dirty, err := m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("could not get migration version: %w", err)
    }

    if err == migrate.ErrNilVersion {
        log.Println("No migrations to run")
    } else {
        log.Printf("Migrations completed successfully! Current version: %d, Dirty: %v\n", version, dirty)
    }

    return nil
}

func RollbackMigrations(pool *pgxpool.Pool, steps int) error {
    db := stdlib.OpenDBFromPool(pool)
    
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create database driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://internal/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    if err := m.Steps(-steps); err != nil {
        return fmt.Errorf("could not rollback migrations: %w", err)
    }

    log.Printf("Rolled back %d migration(s) successfully!\n", steps)
    return nil
}