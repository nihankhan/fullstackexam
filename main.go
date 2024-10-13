// Package main is the entry point of the application.
package main

import (
	"fmt"
	"log"

	"github.com/zuu-development/fullstack-examination-2024/cmd"
	_ "github.com/zuu-development/fullstack-examination-2024/docs"
	"github.com/zuu-development/fullstack-examination-2024/internal/db"
)

// @title			fullstack-examination-2024 API
// @version		0.0.1
// @description	This is a server for fullstack-examination-2024.
// @license.name	Apache 2.0
// @host			localhost:8080
// @BasePath		/api/v1
// @schemes		http
func main() {
	cmd.Execute()
	dbPath := "./tmp/gorm.db"
	fmt.Println("dbpath: ", dbPath)

	// Step 1: Create a new database connection
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Step 2: Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
