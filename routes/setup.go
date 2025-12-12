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
// 	// 1. DEPENDENCY INJECTION 
// 	// ============================================
	
// 	// -- Auth Feature --
// 	authRepo := repositories.NewAuthRepository()
// 	authService := services.NewAuthService(authRepo)

// 	// -- User Management Feature --
// 	userRepo := repositories.NewUserRepository()
// 	userService := services.NewUserService(userRepo)

// 	// -- Achievement Feature --
// 	achievementRepo := repositories.NewAchievementRepository()
// 	achievementService := services.NewAchievementService(achievementRepo)

// 	// ============================================
// 	// 2. ROUTE DEFINITIONS
// 	// ============================================

// 	// --- Public Routes ---
// 	auth := api.Group("/auth")
// 	auth.Post("/login", authService.Login)

// 	// --- Protected Routes ---
	
// 	// A. User Management (Admin)
// 	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	users.Post("/lecturers", userService.RegisterLecturer)
// 	users.Post("/students", userService.RegisterStudent)

// 	// B. Achievement Routes
// 	achievements := api.Group("/achievements", middleware.Protected())
	
// 	// 1. Fitur Mahasiswa
// 	// Upload & List
// 	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
// 	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
	
// 	// Detail & History (Fase 1)
// 	achievements.Get("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetAchievementByID)
// 	achievements.Get("/:id/history", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetAchievementHistory)

// 	// Update & Delete (Draft Only)
// 	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
// 	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)

// 	// Submit (Finalisasi) - POST sesuai SRS
// 	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
	
// 	// 2. Fitur Dosen Wali
// 	// List Bimbingan
// 	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)
	
// 	// Verify & Reject (Sesuai SRS Endpoint 5.4 dan FR-007, FR-008)
// 	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
// 	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)
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
	// 1. DEPENDENCY INJECTION 
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

	// --- Public Routes ---
	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)

	// --- Protected Routes ---
	
	// A. User Management (Admin)
	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	users.Post("/lecturers", userService.RegisterLecturer)
	users.Post("/students", userService.RegisterStudent)

	// B. Achievement Routes
	achievements := api.Group("/achievements", middleware.Protected())
	
	// 1. Fitur Mahasiswa
	// Upload & List
	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
	
	// Detail & History (Fase 1 - DIPERBAIKI)
	// Tidak pakai RoleMiddleware khusus disini karena logic cek akses ada di Service
	// (Agar Dosen Wali juga bisa lihat detail/history untuk verifikasi)
	achievements.Get("/:id", achievementService.GetAchievementByID)
	achievements.Get("/:id/history", achievementService.GetAchievementHistory)

	// Update & Delete (Draft Only)
	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)

	// Submit (Finalisasi) - POST sesuai SRS
	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
	
	// 2. Fitur Dosen Wali
	// List Bimbingan
	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)
	
	// Verify & Reject (Sesuai SRS Endpoint 5.4 dan FR-007, FR-008)
	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)
}