// package routes

// import (
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/app/services"
// 	"github.com/gofiber/fiber/v2"
// )

// func SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api/v1")

// 	// --- Dependency Injection ---
// 	authRepo := repositories.NewAuthRepository()
// 	authService := services.NewAuthService(authRepo) // Service langsung jadi handler

// 	// --- Routes Definition ---
// 	auth := api.Group("/auth")
// 	// Panggil langsung fungsi di Service
// 	auth.Post("/login", authService.Login)
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

	// ============================================
	// 2. ROUTE DEFINITIONS
	// ============================================

	// --- Public Routes (Tidak butuh token) ---
	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)

	// --- Protected Routes (Butuh Token) ---
	// Middleware.Protected() wajib ada untuk mengecek Token valid
	// Middleware.RoleMiddleware("Admin") wajib ada karena hanya Admin yang boleh akses
	
	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	
	// Endpoint: POST /api/v1/users/lecturers
	users.Post("/lecturers", userService.RegisterLecturer)
	
	// Endpoint: POST /api/v1/users/students
	users.Post("/students", userService.RegisterStudent)
}