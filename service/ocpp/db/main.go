package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/squishmeist/ocpp-go/internal/core/utils"
	"github.com/squishmeist/ocpp-go/service/ocpp/db/schemas"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed schema.sql
var ddl string

var GlobalSQLCtx = context.Background()

func Init(dbInfo utils.DatabaseConfiguration) (string, error) {
	address := dbInfo.Address
	// if the driver is sqlite3 and .db exists in /tmp/ then delete it.
	if dbInfo.Driver == "sqlite3" {
		os.Remove(dbInfo.Address)
		slog.Debug("Sqlite3", "Removing existing database file", dbInfo.Address)
		address = fmt.Sprintf("%s:%s", dbInfo.Protocol, dbInfo.Address)
	}

	db, err := sql.Open(dbInfo.Driver, address)
	if err != nil {
		return "", err
	}

	defer db.Close()

	if _, err := db.ExecContext(GlobalSQLCtx, ddl); err != nil {
		fmt.Println("failed to exec content", "error", err)
		slog.Error("failed to exec content", "error", err)
	}

	slog.Info("Database initialised successfully.")

	return address, nil
}

func Connect(dbInfo utils.DatabaseConfiguration) (*schemas.Queries, *sql.DB, error) {
	addr, err := Init(dbInfo)
	if err != nil {
		return nil, nil, err
	}

	db, err := sql.Open(dbInfo.Driver, addr)
	if err != nil {
		return nil, nil, err
	}

	queries := schemas.New(db)

	return queries, db, nil
}
