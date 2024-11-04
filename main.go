package main

import (
    "flag"
    "log"
)

func main() {
    protocol := flag.String("protocol", "", "Transfer protocol (azureblob, cifs, sftp, s3, local)")
    source := flag.String("source", "", "Source file path")
    destination := flag.String("destination", "", "Destination file path")
    flag.Parse()

    if *protocol == "" || *source == "" || *destination == "" {
        log.Fatal("protocol, source, and destination are required")
    }

    // Initialize configuration
    initConfig()

    // Initialize database
    db, err := initDatabase()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Transfer file
    err = transferFile(*protocol, *source, *destination)
    if err != nil {
        log.Fatalf("File transfer failed: %v", err)
    }

    // Log transfer
    err = logTransfer(db, *protocol, *source, *destination)
    if err != nil {
        log.Printf("Failed to log transfer: %v", err)
    }

    log.Println("File transfer completed successfully")
}
