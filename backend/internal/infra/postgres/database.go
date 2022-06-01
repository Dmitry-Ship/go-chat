package postgres

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection() *gorm.DB {
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

	err = autoMigrate(db)

	if err != nil {
		panic("⛔️ Could not migrate database")
	}

	log.Printf("💿 Connected to database %s", dbname)

	return db
}

func autoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&User{},
		&GroupConversation{},
		&Participant{},
		&Message{},
		&Conversation{},
	}

	// for _, model := range models {
	// 	db.Migrator().DropTable(&model)
	// }

	err := db.AutoMigrate(models...)

	if err != nil {
		return err
	}

	return nil
}
