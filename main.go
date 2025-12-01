package main

import (
	"PrestasiMhs-API/config" // Import konfigurasi database
	"PrestasiMhs-API/routes" // IMPORT BARU: Import package routes agar bisa dipanggil
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
	
	// Pastikan koneksi ditutup saat aplikasi berhenti
	defer func() {
		if config.DB != nil {
			config.DB.Close()
		}
	}()

	// 3. Setup Fiber Framework
	app := fiber.New(fiber.Config{
		// Batas upload file 10MB (Sesuai Modul 9)
		BodyLimit: 10 * 1024 * 1024, 
	})

	// 4. Middlewares Global
	app.Use(logger.New()) // Logging request
	app.Use(cors.New())   // Handle CORS untuk Frontend nanti

	// 5. Setup Routes
	// Panggil fungsi SetupRoutes dari folder routes/setup.go
	// Ini akan mendaftarkan /api/v1/auth/login dan route lainnya secara otomatis
	routes.SetupRoutes(app)
	
	// Health Check Endpoint (Manual)
	// Kita taruh manual disini untuk memudahkan pengecekan server
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "success",
			"message":   "Server is running smoothly! 🚀",
			"db_status": "connected",
		})
	})

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