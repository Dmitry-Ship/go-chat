package postgres

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewDatabaseConnection(conf DbConfig) *gorm.DB {
	options := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", conf.Host, conf.Port, conf.User, conf.Name, conf.Password)

	db, err := gorm.Open(postgres.Open(options), &gorm.Config{})
	if err != nil {
		panic("⛔️ Could not connect to database")
	}

	err = autoMigrate(db)

	if err != nil {
		panic("⛔️ Could not migrate database")
	}

	log.Printf("💿 Connected to database %s", conf.Name)

	return db
}

var models = []interface{}{
	&User{},
	&GroupConversation{},
	&Participant{},
	&Message{},
	&Conversation{},
}

func Drop(db *gorm.DB) error {
	for _, model := range models {
		err := db.Migrator().DropTable(&model)
		if err != nil {
			return fmt.Errorf("drop table error: %w", err)
		}
	}

	log.Printf("Dropped database %s", db.Name())

	return autoMigrate(db)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(models...)
}
