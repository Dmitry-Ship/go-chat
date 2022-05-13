package postgres

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	connection *gorm.DB
}

func NewDatabaseConnection() *DatabaseConnection {
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

	log.Printf("💿 Connected to database %s", dbname)

	return &DatabaseConnection{
		connection: db,
	}
}

func (d *DatabaseConnection) GetConnection() *gorm.DB {
	return d.connection
}

func (d *DatabaseConnection) AutoMigrate() {
	// d.connection.Migrator().DropTable(&Message{})
	// d.connection.Migrator().DropTable(&User{})
	// d.connection.Migrator().DropTable(&Participant{})
	// d.connection.Migrator().DropTable(&Conversation{})
	// d.connection.Migrator().DropTable(&TextMessage{})
	// d.connection.Migrator().DropTable(&ConversationRenamedMessage{})
	// d.connection.Migrator().DropTable(&UserNotificationTopic{})
	// d.connection.Migrator().DropTable(&PrivateConversation{})
	// d.connection.Migrator().DropTable(&PublicConversation{})

	err := d.connection.AutoMigrate(Message{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(Conversation{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(PublicConversation{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(PrivateConversation{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(User{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(Participant{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(TextMessage{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(ConversationRenamedMessage{})
	if err != nil {
		panic(err)
	}
	err = d.connection.AutoMigrate(UserNotificationTopic{})
	if err != nil {
		panic(err)
	}
}
