package db

import "log"

// User represents the user model
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SeedUsers seeds the user table with initial data
func SeedUsers(database Database) {
	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}

	for _, user := range users {
		_, err := database.CreateUser(user)
		if err != nil {
			log.Printf("Failed to seed user: %v", err)
		} else {
			log.Printf("User %s seeded successfully.", user.Name)
		}
	}
}

// SeedDatabase runs all the seeders
func SeedDatabase(database Database) {
	SeedUsers(database)
	// Add other seed functions here
}
