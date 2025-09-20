package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/internal/api"
	"backend/internal/db"
	"backend/internal/repository"
	"backend/internal/service"
)

func main() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("⚠️  No .env file found, falling back to system env vars")
	}

	conn, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	//Auth setup
	userRepo := repository.NewUserRepository(conn)
	authService := service.NewAuthService(userRepo)
	authHandler := api.NewAuthHandler(authService)

	//File setup
	fileRepo := repository.NewFileRepository(conn)
	fileService := service.NewFileService(fileRepo)
	fileHandler := api.NewFileHandler(fileService , fileRepo)

	r := gin.Default()

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true, 
	}))

	
	apiRoutes := r.Group("/api")
	{
		// Public
		apiRoutes.POST("/signup", authHandler.SignUp)
		apiRoutes.POST("/login", authHandler.SignIn)
		apiRoutes.POST("/logout", authHandler.SignOut)

		// Protected routes (require auth)
		protected := apiRoutes.Group("/")
		protected.Use(api.AuthMiddleware())
		{
			protected.GET("/me", authHandler.Me)

			
			protected.POST("/upload", fileHandler.Upload)
			protected.GET("/files", fileHandler.ListFiles)
			protected.GET("/files/:id/download", fileHandler.DownloadFile)
		}
	}

	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Server running on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
