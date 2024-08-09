package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GORMDB implements the Database interface for both SQLite and PostgreSQL
type GORMDB struct {
	conn *gorm.DB
}

// NewGORMDB creates a new GORM database connection
func NewGORMDB(dsn string, dbType string) (*GORMDB, error) {
	var dialector gorm.Dialector

	switch dbType {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the User model to create the table if it doesn't exist
	err = conn.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return &GORMDB{conn: conn}, nil
}

func (db *GORMDB) GetUsers() ([]User, error) {
	var users []User
	if err := db.conn.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (db *GORMDB) CreateUser(u User) (User, error) {
	user := User{Name: u.Name}
	if err := db.conn.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}
