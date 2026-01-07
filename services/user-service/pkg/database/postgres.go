package database

import (
	"fmt"
	"log"
	"wealth-vault/user-service/configs"
	"wealth-vault/user-service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(config configs.PostgreSQL) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_ = db.AutoMigrate(
		&domain.User{},
		&domain.FriendList{},
		&domain.Group{},
		&domain.GroupItem{},
		&domain.GroupLog{},
		&domain.Notification{},
	)

	log.Println("Database connection established successfully!")
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database is not initialized")
	}

	return db
}
