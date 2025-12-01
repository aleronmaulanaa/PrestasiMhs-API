package routes

import (
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/app/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// --- Dependency Injection ---
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo) // Service langsung jadi handler

	// --- Routes Definition ---
	auth := api.Group("/auth")
	// Panggil langsung fungsi di Service
	auth.Post("/login", authService.Login)
}