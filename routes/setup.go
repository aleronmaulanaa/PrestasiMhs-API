// package routes

// import (
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/app/services"
// 	"PrestasiMhs-API/middleware"
// 	"github.com/gofiber/fiber/v2"
// )

// func SetupRoutes(app *fiber.App) {
// 	api := app.Group("/api/v1")

// 	// ============================================
// 	// 1. DEPENDENCY INJECTION 
// 	// ============================================

// 	// Auth & User
// 	authRepo := repositories.NewAuthRepository()
// 	authService := services.NewAuthService(authRepo)
// 	userRepo := repositories.NewUserRepository()
// 	userService := services.NewUserService(userRepo)

// 	// Achievement
// 	achievementRepo := repositories.NewAchievementRepository()
// 	achievementService := services.NewAchievementService(achievementRepo)

// 	// [NEW FASE 3] Reporting
// 	reportRepo := repositories.NewReportRepository()
// 	reportService := services.NewReportService(reportRepo)

// 	// ============================================
// 	// 2. ROUTE DEFINITIONS
// 	// ============================================

// 	// --- Auth Routes ---
// 	auth := api.Group("/auth")
// 	auth.Post("/login", authService.Login)
	
// 	// [NEW FASE 3] Profile (User bisa melihat data dirinya sendiri)
// 	auth.Get("/profile", middleware.Protected(), authService.GetProfile)

// 	// --- A. User Management (Admin Only) ---
// 	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	users.Get("/", userService.GetAllUsers)
// 	users.Get("/:id", userService.GetUserByID)         
// 	users.Put("/:id", userService.UpdateUser)          
// 	users.Put("/:id/role", userService.ChangePassword) 
// 	users.Delete("/:id", userService.DeleteUser)
// 	users.Post("/lecturers", userService.RegisterLecturer)
// 	users.Post("/students", userService.RegisterStudent)

// 	// --- Relations Management (Admin Only) ---
// 	students := api.Group("/students", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	students.Get("/", userService.GetAllStudents)
// 	students.Put("/:id/advisor", userService.AssignAdvisor)
	
// 	lecturers := api.Group("/lecturers", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	lecturers.Get("/", userService.GetAllLecturers)

// 	// --- B. Achievement Routes ---
// 	achievements := api.Group("/achievements", middleware.Protected())
	
// 	// Admin View All
// 	achievements.Get("/", middleware.RoleMiddleware("Admin"), achievementService.GetAllAchievements)
	
// 	// Static Routes
// 	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
// 	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
// 	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)
	
// 	// Dynamic Routes (/:id)
// 	achievements.Get("/:id", achievementService.GetAchievementByID)
// 	achievements.Get("/:id/history", achievementService.GetAchievementHistory)
// 	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
// 	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)
// 	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
// 	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
// 	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)

// 	// --- [NEW FASE 3] Report Routes ---
// 	reports := api.Group("/reports", middleware.Protected())
	
// 	// Dashboard Statistics (Admin Only)
// 	// Jika Dosen Wali juga butuh akses, pastikan middleware Anda support multiple roles
// 	reports.Get("/statistics", middleware.RoleMiddleware("Admin"), reportService.GetDashboardStatistics)
// }


package routes

import (
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/app/services"
	"PrestasiMhs-API/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// ============================================
	// 1. DEPENDENCY INJECTION 
	// ============================================

	// Auth & User
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)

	// Achievement
	achievementRepo := repositories.NewAchievementRepository()
	achievementService := services.NewAchievementService(achievementRepo)

	// [NEW FASE 3] Reporting
	reportRepo := repositories.NewReportRepository()
	// [UPDATE] Inject achievementRepo juga agar Service bisa akses Mongo
	reportService := services.NewReportService(reportRepo, achievementRepo)

	// ============================================
	// 2. ROUTE DEFINITIONS
	// ============================================

	// --- Auth Routes ---
	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Get("/profile", middleware.Protected(), authService.GetProfile)

	// --- A. User Management (Admin Only) ---
	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	users.Get("/", userService.GetAllUsers)
	users.Get("/:id", userService.GetUserByID)
	users.Put("/:id", userService.UpdateUser)
	users.Put("/:id/role", userService.ChangePassword)
	users.Delete("/:id", userService.DeleteUser)
	users.Post("/lecturers", userService.RegisterLecturer)
	users.Post("/students", userService.RegisterStudent)

	// --- Relations Management ---
	students := api.Group("/students", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	students.Get("/", userService.GetAllStudents)
	students.Put("/:id/advisor", userService.AssignAdvisor)

	lecturers := api.Group("/lecturers", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	lecturers.Get("/", userService.GetAllLecturers)

	// --- B. Achievement Routes ---
	achievements := api.Group("/achievements", middleware.Protected())
	
	// Admin View All
	achievements.Get("/", middleware.RoleMiddleware("Admin"), achievementService.GetAllAchievements)

	// Static
	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)

	// Dynamic
	achievements.Get("/:id", achievementService.GetAchievementByID)
	achievements.Get("/:id/history", achievementService.GetAchievementHistory)
	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)
	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)

	// --- [NEW FASE 3] Report Routes ---
	reports := api.Group("/reports", middleware.Protected())

	// Dashboard Statistics
	reports.Get("/statistics", middleware.RoleMiddleware("Admin"), reportService.GetDashboardStatistics)
	
	// [NEW] Student Report (Transkrip)
	// Admin & Dosen Wali boleh lihat
	reports.Get("/student/:studentID", middleware.RoleMiddleware("Admin", "Dosen Wali"), reportService.GetStudentReport)
}