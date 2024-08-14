// cloudinary.go
package utils

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/joho/godotenv"
)

var cloudinaryClient *cloudinary.Cloudinary

func InitCloudinary() *cloudinary.Cloudinary {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get environment variables
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	// Initialize Cloudinary
	cloudinaryClient, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Error initializing Cloudinary: %v", err)
	}

	return cloudinaryClient
}
