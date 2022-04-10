package database

import (
	"GitHub/go-chat/backend/domain"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseConnection() *gorm.DB {
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	dbpassword := os.Getenv("DB_PASSWORD")

	options := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, dbpassword)

	db, err := gorm.Open(postgres.Open(options), &gorm.Config{})
	if err != nil {
		panic("⛔️ Could not connect to database")
	}
	fmt.Println(fmt.Sprintf("💿 Connected to database %s", dbname))

	// Migrate the schema

	db.AutoMigrate(domain.MessageDAO{})
	db.AutoMigrate(domain.ConversationDAO{})
	db.AutoMigrate(domain.UserDAO{})
	db.AutoMigrate(domain.ParticipantDAO{})

	return db
}
