package db

import (
	"github.com/zuu-development/fullstack-examination-2024/internal/model"
	"gorm.io/gorm"
)

// Migrate runs the auto-migration for the database
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		return err
	}

	// Define CHECK constraint for priority
	const todos = `
			CREATE TABLE IF NOT EXISTS todos (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				task TEXT NOT NULL,
				status TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				priority TEXT CHECK(priority IN ('Low', 'Medium', 'High'))
			);
		`
	if err := db.Exec(todos).Error; err != nil {
		return err
	}

	return nil
}
