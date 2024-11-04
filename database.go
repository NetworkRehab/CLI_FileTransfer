package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// initDatabase initializes the SQLite database.
func initDatabase() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./transfers.db")
    if (err != nil) {
        return nil, err
    }

    sqlStmt := `
    CREATE TABLE IF NOT EXISTS transfers (
        id INTEGER NOT NULL PRIMARY KEY,
        protocol TEXT,
        source TEXT,
        destination TEXT
    );`
    _, err = db.Exec(sqlStmt)
    if err != nil {
        return nil, err
    }

    return db, nil
}

// logTransfer logs the transfer details into the database.
func logTransfer(db *sql.DB, protocol, source, destination string) error {
    stmt, err := db.Prepare("INSERT INTO transfers(protocol, source, destination) values(?,?,?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(protocol, source, destination)
    return err
}
