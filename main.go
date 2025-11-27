package main

import (
	"PrestasiMhs-API/config" // Import dari module name Anda
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description API untuk sistem prestasi mahasiswa (PostgreSQL + MongoDB).
// @contact.name Tim Backend
// @host localhost:3000
// @BasePath /api/v1
func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// 2. Inisialisasi Database (Postgres & Mongo)
	config.ConnectDB()
	
	// Pastikan koneksi ditutup saat aplikasi berhenti (Optional, tp best practice di main flow)
	defer func() {
		if config.DB != nil {
			config.DB.Close()
		}
		// MongoDB client disconnect bisa ditambahkan disini jika client di-export
	}()

	// 3. Setup Fiber Framework
	app := fiber.New(fiber.Config{
		// Batas upload file 10MB (Sesuai Modul 9)
		BodyLimit: 10 * 1024 * 1024, 
	})

	// 4. Middlewares Global
	app.Use(logger.New()) // Logging request
	app.Use(cors.New())   // Handle CORS untuk Frontend nanti

	// 5. Setup Routes Dasar
	api := app.Group("/api/v1")
	
	// Health Check Endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running smoothly! 🚀",
			"db_status": "connected",
		})
	})

	// TODO: Register route lain disini (User, Auth, Achievements) nanti

	// 6. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}