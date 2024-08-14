package db

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GORMDB implements the Database interface for both SQLite and PostgreSQL
type GORMDB struct {
	conn *gorm.DB
}

// User represents the user model in the database
type UserSql struct {
	ID             uint           `gorm:"primaryKey" json:"id,omitempty"`
	Username       string         `gorm:"unique;not null" json:"username" binding:"required"`
	Email          string         `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password       string         `gorm:"not null" json:"password" binding:"required"`
	ProfilePicture string         `gorm:"type:text" json:"profilePicture,omitempty"`
	Bio            string         `gorm:"type:text" json:"bio,omitempty"`
	Gender         string         `json:"gender,omitempty"`
	Followers      []*User        `gorm:"many2many:user_followers;association_jointable_foreignkey:follower_id" json:"followers,omitempty"`
	Following      []*User        `gorm:"many2many:user_followings;association_jointable_foreignkey:following_id" json:"following,omitempty"`
	Posts          []Post         `json:"posts,omitempty"`
	Bookmarks      []Post         `gorm:"many2many:user_bookmarks" json:"bookmarks,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Post represents a user's post (this struct should be defined based on your needs)
type PostSql struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
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
		return nil, errors.New("unsupported database type: " + dbType)
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

// GetUsers retrieves all users or applies a filter (simplified for SQL-based systems)
func (db *GORMDB) GetUsers(filter interface{}) ([]User, error) {
	var users []User
	if err := db.conn.Where(filter).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a single user by ID
func (db *GORMDB) GetUserByID(id uint) (User, error) {
	var user User
	if err := db.conn.First(&user, id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// CreateUser creates a new user
func (db *GORMDB) CreateUser(u User) (User, error) {
	if err := db.conn.Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

// GetUserByEmail retrieves a user by email
func (db *GORMDB) GetUserByEmail(email string) (User, error) {
	var user User
	if err := db.conn.Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// UpdateUser updates a user's information
func (db *GORMDB) UpdateUser(id uint, update interface{}) error {
	if err := db.conn.Model(&User{}).Where("id = ?", id).Updates(update).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID
func (db *GORMDB) DeleteUser(id uint) error {
	if err := db.conn.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// FollowOrUnfollowUser handles following or unfollowing a user
func (db *GORMDB) FollowOrUnfollowUser(followingUserID, targetUserID uint, action string) error {
	var user, target User
	if err := db.conn.First(&user, followingUserID).Error; err != nil {
		return err
	}
	if err := db.conn.First(&target, targetUserID).Error; err != nil {
		return err
	}

	if action == "follow" {
		return db.conn.Model(&user).Association("Following").Append(&target)
	} else if action == "unfollow" {
		return db.conn.Model(&user).Association("Following").Delete(&target)
	}
	return errors.New("invalid action")
}

// FollowUser handles following a user
func (db *GORMDB) FollowUser(followingUserID, targetUserID uint) error {
	return db.FollowOrUnfollowUser(followingUserID, targetUserID, "follow")
}

// UnfollowUser handles unfollowing a user
func (db *GORMDB) UnfollowUser(followingUserID, targetUserID uint) error {
	return db.FollowOrUnfollowUser(followingUserID, targetUserID, "unfollow")
}
