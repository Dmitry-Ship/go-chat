package postgres

import (
	"fmt"
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

	fmt.Println(fmt.Sprintf("💿 Connected to database %s", dbname))

	return &DatabaseConnection{
		connection: db,
	}
}

func (d *DatabaseConnection) GetConnection() *gorm.DB {
	return d.connection
}

func (d *DatabaseConnection) AutoMigrate() {
	d.connection.AutoMigrate(Message{})
	d.connection.AutoMigrate(Conversation{})
	d.connection.AutoMigrate(PublicConversation{})
	d.connection.AutoMigrate(PrivateConversation{})
	d.connection.AutoMigrate(User{})
	d.connection.AutoMigrate(Participant{})
	d.connection.AutoMigrate(TextMessage{})
	d.connection.AutoMigrate(ConversationRenamedMessage{})
	d.connection.AutoMigrate(UserNotificationTopic{})
}
