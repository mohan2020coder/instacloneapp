package main

import (
	"log"
	"os"

	"instacloneapp/server/pkg/db"
	"instacloneapp/server/routes"
	cloudinary "instacloneapp/server/utils"

	"instacloneapp/server/socket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize Cloudinary
	cloudinaryClient := cloudinary.InitCloudinary()
	router := gin.Default()

	var database db.Database
	var err error

	router.GET("/ws", socket.HandleConnection)

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
		// Retrieve MongoDB connection details from environment variables
		mongoURI := os.Getenv("MONGO_URI")
		mongoDBName := os.Getenv("MONGO_DB_NAME")
		mongoCollections := os.Getenv("MONGO_COLLECTIONS")

		// Ensure that all required environment variables are set
		if mongoURI == "" || mongoDBName == "" || mongoCollections == "" {
			log.Fatalf("MongoDB environment variables are not set")
		}

		// Initialize the MongoDB instance
		database, err := db.NewMongoDB(mongoURI, mongoDBName, mongoCollections)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Use the MongoDB instance (e.g., for further operations)
		_ = database
	} else {
		// Connect to SQLite or PostgreSQL
		dbType := os.Getenv("DB_TYPE")
		if dbType == "" {
			dbType = "sqlite" // Default to SQLite if DB_TYPE is not set
		}

		// var dsn string
		// if dbType == "sqlite" {
		// 	dsn = os.Getenv("SQLITE_DB_PATH")
		// 	if dsn == "" {
		// 		dsn = "./test.db" // Default to a local file if SQLITE_DB_PATH is not set
		// 	}
		// 	database, err = db.NewGORMDB(dsn, "sqlite")
		// } else if dbType == "postgres" {
		// 	dsn = os.Getenv("POSTGRES_DSN")
		// 	if dsn == "" {
		// 		log.Fatalf("PostgreSQL DSN is not set")
		// 	}
		// 	database, err = db.NewGORMDB(dsn, "postgres")
		// } else {
		// 	log.Fatalf("Unsupported database type: %s", dbType)
		// }

		// if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		// }
	}
	//db.SeedDatabase(database)
	// Serve static files from frontend/dist
	// Serve static files from the .next directory
	// router.Static("/_next", filepath.Join(".", "frontend", ".next"))

	// // Serve any other static assets (e.g., images, CSS, etc.)
	// router.Static("/static", filepath.Join(".", "frontend", "public", "static"))
	// router.StaticFile("/favicon.ico", filepath.Join(".", "frontend", "public", "favicon.ico"))

	// Setup routes

	// Set up routes with dependencies
	routes.SetupRoutes(router, database, cloudinaryClient)
	routes.SetupMessageRoutes(router, database, cloudinaryClient)
	routes.SetupPostRoutes(router, database, cloudinaryClient)

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
