// package main

// import (
// 	"PrestasiMhs-API/config" // Import konfigurasi database
// 	"PrestasiMhs-API/routes" // Import package routes agar bisa dipanggil
// 	"log"
// 	"os"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/joho/godotenv"
// )

// // @title Sistem Pelaporan Prestasi Mahasiswa API
// // @version 1.0
// // @description API untuk sistem prestasi mahasiswa (PostgreSQL + MongoDB).
// // @contact.name Tim Backend
// // @host localhost:3000
// // @BasePath /api/v1
// func main() {
// 	// 1. Load Environment Variables
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan environment sistem.")
// 	}

// 	// 2. Inisialisasi Database (Postgres & Mongo)
// 	config.ConnectDB()
	
// 	// Pastikan koneksi ditutup saat aplikasi berhenti
// 	defer func() {
// 		if config.DB != nil {
// 			config.DB.Close()
// 		}
// 	}()

// 	// --- [NEW] Auto-Create Folder Uploads ---
// 	// Membuat folder penyimpanan file jika belum ada
// 	uploadDirs := []string{"./uploads/documents", "./uploads/photos"}
// 	for _, dir := range uploadDirs {
// 		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
// 			log.Fatalf("❌ Gagal membuat folder upload %s: %v", dir, err)
// 		}
// 	}
// 	log.Println("✅ Folder penyimpanan file siap.")
// 	// ----------------------------------------

// 	// 3. Setup Fiber Framework
// 	app := fiber.New(fiber.Config{
// 		// Batas upload file 10MB (Sesuai Modul 9)
// 		BodyLimit: 10 * 1024 * 1024, 
// 	})

// 	// 4. Middlewares Global
// 	app.Use(logger.New()) // Logging request
// 	app.Use(cors.New())   // Handle CORS untuk Frontend nanti

// 	// 5. Setup Routes
// 	// Panggil fungsi SetupRoutes dari folder routes/setup.go
// 	routes.SetupRoutes(app)
	
// 	// Health Check Endpoint (Manual)
// 	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"status":    "success",
// 			"message":   "Server is running smoothly! 🚀",
// 			"db_status": "connected",
// 		})
// 	})

// 	// 6. Start Server
// 	port := os.Getenv("APP_PORT")
// 	if port == "" {
// 		port = "3000"
// 	}

// 	log.Printf("Server starting on port :%s", port)
// 	if err := app.Listen(":" + port); err != nil {
// 		log.Fatalf("❌ Server failed to start: %v", err)
// 	}
// }


package main

import (
	"PrestasiMhs-API/config" // Import konfigurasi database
	"PrestasiMhs-API/routes" // Import package routes
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	// [PENTING] Import docs yang digenerate oleh swag init
	// Tanda underscore (_) berarti kita hanya menjalankan fungsi init() dari package tersebut
	_ "PrestasiMhs-API/docs"
)

// @title           Sistem Pelaporan Prestasi Mahasiswa API
// @version         1.0
// @description     API untuk sistem prestasi mahasiswa (PostgreSQL + MongoDB).
// @termsOfService  http://swagger.io/terms/

// @contact.name    Tim Backend
// @contact.email   support@prestasi.id

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <token_jwt_anda>
func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// 2. Inisialisasi Database
	// PostgreSQL
	config.ConnectDB()
	
	// MongoDB (WAJIB ADA untuk Fase 3)
	config.ConnectMongo()

	// Pastikan koneksi SQL ditutup saat aplikasi berhenti
	defer func() {
		if config.DB != nil {
			config.DB.Close()
		}
	}()

	// --- Auto-Create Folder Uploads ---
	uploadDirs := []string{"./uploads/documents", "./uploads/photos"}
	for _, dir := range uploadDirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("❌ Gagal membuat folder upload %s: %v", dir, err)
		}
	}
	log.Println("✅ Folder penyimpanan file siap.")
	// ----------------------------------------

	// 3. Setup Fiber Framework
	app := fiber.New(fiber.Config{
		// Batas upload file 10MB (Sesuai Modul)
		BodyLimit: 10 * 1024 * 1024,
	})

	// 4. Middlewares Global
	app.Use(logger.New()) // Logging request
	app.Use(cors.New())   // Handle CORS

	// 5. Setup Routes
	routes.SetupRoutes(app)

	// Health Check Endpoint (Manual)
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

	log.Printf("Swagger UI tersedia di: http://localhost:%s/swagger/index.html", port)
	log.Printf("Server starting on port :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}