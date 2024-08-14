// datauri.go
package utils

import (
	"encoding/base64"
	"mime"
	"path/filepath"
)

// GetDataURI generates a Data URI from the file content
func GetDataURI(filename string, fileContent []byte) string {
	// Extract the file extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".txt" // Default to .txt if no extension is provided
	}

	// Detect MIME type based on file extension
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default MIME type
	}

	// Encode file content to Base64
	encoded := base64.StdEncoding.EncodeToString(fileContent)

	// Create the Data URI
	dataURI := "data:" + mimeType + ";base64," + encoded

	return dataURI
}
