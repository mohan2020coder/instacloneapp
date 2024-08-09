package db



// Database defines the methods needed for database operations
type Database interface {
	GetUsers() ([]User, error)
	CreateUser(User) (User, error)
}


