package main

import (
	"log"
	"os"

	"instacloneapp/server/pkg/db"
	"instacloneapp/server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	router := gin.Default()

	var database db.Database
	var err error

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Load environment variables
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // Default to development if ENV is not set
	}

	if env == "production" {
		// Connect to MongoDB
		mongoURI := os.Getenv("MONGO_URI")
		mongoDBName := os.Getenv("MONGO_DB_NAME")
		mongoCollection := os.Getenv("MONGO_COLLECTION")
		if mongoURI == "" || mongoDBName == "" || mongoCollection == "" {
			log.Fatalf("MongoDB environment variables are not set")
		}
		database, err = db.NewMongoDB(mongoURI, mongoDBName, mongoCollection)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
	} else {
		// Connect to SQLite or PostgreSQL
		dbType := os.Getenv("DB_TYPE")
		if dbType == "" {
			dbType = "sqlite" // Default to SQLite if DB_TYPE is not set
		}

		var dsn string
		if dbType == "sqlite" {
			dsn = os.Getenv("SQLITE_DB_PATH")
			if dsn == "" {
				dsn = "./test.db" // Default to a local file if SQLITE_DB_PATH is not set
			}
			database, err = db.NewGORMDB(dsn, "sqlite")
		} else if dbType == "postgres" {
			dsn = os.Getenv("POSTGRES_DSN")
			if dsn == "" {
				log.Fatalf("PostgreSQL DSN is not set")
			}
			database, err = db.NewGORMDB(dsn, "postgres")
		} else {
			log.Fatalf("Unsupported database type: %s", dbType)
		}

		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	//db.SeedDatabase(database)
	// Serve static files from frontend/dist
	// Serve static files from the .next directory
	// router.Static("/_next", filepath.Join(".", "frontend", ".next"))

	// // Serve any other static assets (e.g., images, CSS, etc.)
	// router.Static("/static", filepath.Join(".", "frontend", "public", "static"))
	// router.StaticFile("/favicon.ico", filepath.Join(".", "frontend", "public", "favicon.ico"))

	// Setup routes

	routes.SetupRoutes(router, database)

	// Catch-all route to serve index.html for SPA
	// router.NoRoute(func(c *gin.Context) {
	// 	c.File(filepath.Join(".", "frontend", ".next", "server", "pages", "index.html"))
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
