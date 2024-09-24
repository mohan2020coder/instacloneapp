// cloudinary.go
package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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

// UploadImageToCloudinary uses the existing Cloudinary client to upload an image
func UploadImageToCloudinary(cloudinaryClient *cloudinary.Cloudinary, image multipart.File, folder string) (string, error) {
	// Read image file into a buffer
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(image); err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	// Upload image using the provided cloudinaryClient instance
	uploadParams := uploader.UploadParams{Folder: folder} // Adjust folder as needed
	resp, err := cloudinaryClient.Upload.Upload(context.TODO(), buf, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %v", err)
	}

	// Return the secure URL from Cloudinary
	return resp.SecureURL, nil
}
