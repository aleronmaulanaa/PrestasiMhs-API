// package routes

// import (
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/app/services"
// 	"PrestasiMhs-API/middleware" // Import Middleware
// 	"github.com/gofiber/fiber/v2"
// )

// func SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api/v1")

// 	// ============================================
// 	// 1. DEPENDENCY INJECTION (Menyiapkan Layer)
// 	// ============================================
	
// 	// -- Auth Feature --
// 	authRepo := repositories.NewAuthRepository()
// 	authService := services.NewAuthService(authRepo)

// 	// -- User Management Feature --
// 	userRepo := repositories.NewUserRepository()
// 	userService := services.NewUserService(userRepo)

// 	// -- Achievement Feature (BARU) --
// 	achievementRepo := repositories.NewAchievementRepository()
// 	achievementService := services.NewAchievementService(achievementRepo)

// 	// ============================================
// 	// 2. ROUTE DEFINITIONS
// 	// ============================================

// 	// --- Public Routes (Tidak butuh token) ---
// 	auth := api.Group("/auth")
// 	auth.Post("/login", authService.Login)

// 	// --- Protected Routes (Butuh Token) ---
// 	// Middleware.Protected() wajib ada untuk mengecek Token valid
	
// 	// A. User Management Routes (Admin Only)
// 	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	users.Post("/lecturers", userService.RegisterLecturer)
// 	users.Post("/students", userService.RegisterStudent)

// 	// B. Achievement Routes (Mahasiswa Only)
// 	// Kita buat group '/achievements' yang dilindungi middleware Protected
// 	achievements := api.Group("/achievements", middleware.Protected())
	
// 	// Endpoint: POST /api/v1/achievements
// 	// Hanya boleh diakses oleh role "Mahasiswa"
// 	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
// }


package routes

import (
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/app/services"
	"PrestasiMhs-API/middleware" // Import Middleware
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// ============================================
	// 1. DEPENDENCY INJECTION (Menyiapkan Layer)
	// ============================================
	
	// -- Auth Feature --
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)

	// -- User Management Feature --
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)

	// -- Achievement Feature --
	achievementRepo := repositories.NewAchievementRepository()
	achievementService := services.NewAchievementService(achievementRepo)

	// ============================================
	// 2. ROUTE DEFINITIONS
	// ============================================

	// --- Public Routes (Tidak butuh token) ---
	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)

	// --- Protected Routes (Butuh Token) ---
	
	// A. User Management Routes (Admin Only)
	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	users.Post("/lecturers", userService.RegisterLecturer)
	users.Post("/students", userService.RegisterStudent)

	// B. Achievement Routes
	// Group umum yang diproteksi token (login wajib)
	achievements := api.Group("/achievements", middleware.Protected())
	
	// 1. Fitur Mahasiswa
	// Upload Prestasi
	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
	// Lihat Prestasi Sendiri
	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)

	// 2. Fitur Dosen Wali
	// Lihat Prestasi Mahasiswa Bimbingan
	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)
	// Verifikasi Prestasi (Approve/Reject)
	achievements.Put("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
}