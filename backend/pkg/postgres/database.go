package postgres

import (
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

	db.AutoMigrate(Message{})
	db.AutoMigrate(Conversation{})
	db.AutoMigrate(PublicConversation{})
	db.AutoMigrate(User{})
	db.AutoMigrate(Participant{})
	db.AutoMigrate(TextMessage{})
	db.AutoMigrate(ConversationRenamedMessage{})

	return db
}
