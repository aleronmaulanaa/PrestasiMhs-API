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
	
// 	authRepo := repositories.NewAuthRepository()
// 	authService := services.NewAuthService(authRepo)

// 	userRepo := repositories.NewUserRepository()
// 	userService := services.NewUserService(userRepo)

// 	achievementRepo := repositories.NewAchievementRepository()
// 	achievementService := services.NewAchievementService(achievementRepo)

// 	// ============================================
// 	// 2. ROUTE DEFINITIONS
// 	// ============================================

// 	// --- Public ---
// 	auth := api.Group("/auth")
// 	auth.Post("/login", authService.Login)

// 	// --- Protected ---
	
// 	// A. User Management (Admin Only - SRS 5.2 & 5.5)
// 	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	users.Get("/", userService.GetAllUsers)         // List All Users
// 	users.Delete("/:id", userService.DeleteUser)    // Delete User
// 	users.Post("/lecturers", userService.RegisterLecturer)
// 	users.Post("/students", userService.RegisterStudent)

// 	// Relations Management (Admin Only)
// 	students := api.Group("/students", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	students.Get("/", userService.GetAllStudents)           // List Students
// 	students.Put("/:id/advisor", userService.AssignAdvisor) // Assign Dosen Wali

// 	lecturers := api.Group("/lecturers", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	lecturers.Get("/", userService.GetAllLecturers)         // List Lecturers

// 	// B. Achievement Routes
// 	achievements := api.Group("/achievements", middleware.Protected())
	
// 	// --- 1. STATIC ROUTES (Harus ditaruh DI ATAS route /:id) ---
	
// 	// Upload (Create)
// 	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
	
// 	// List Milik Mahasiswa Sendiri
// 	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
	
// 	// [FIX] List Bimbingan Dosen Wali (Dipindahkan ke sini agar tidak tertutup oleh /:id)
// 	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)
	
// 	// --- 2. DYNAMIC ROUTES (/:id) (Ditaruh DI BAWAH) ---
// 	// Karena /:id bersifat wildcard, dia akan menangkap apa saja jika ditaruh paling atas.
	
// 	// Detail & History
// 	achievements.Get("/:id", achievementService.GetAchievementByID)
// 	achievements.Get("/:id/history", achievementService.GetAchievementHistory)

// 	// Actions (Update, Delete, Submit)
// 	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
// 	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)
// 	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
	
// 	// Verify & Reject Actions
// 	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
// 	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)
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

	// Injection
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	achievementRepo := repositories.NewAchievementRepository()
	achievementService := services.NewAchievementService(achievementRepo)

	// --- Routes ---
	
	api.Post("/auth/login", authService.Login)

	// A. User Management (Admin)
	users := api.Group("/users", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	users.Get("/", userService.GetAllUsers)
	users.Get("/:id", userService.GetUserByID)       // [NEW]
	users.Put("/:id", userService.UpdateUser)        // [NEW]
	users.Put("/:id/role", userService.ChangePassword) // [NEW] (Mapping Update Role/Password)
	users.Delete("/:id", userService.DeleteUser)
	users.Post("/lecturers", userService.RegisterLecturer)
	users.Post("/students", userService.RegisterStudent)

	// Relations
	students := api.Group("/students", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	students.Get("/", userService.GetAllStudents)
	students.Put("/:id/advisor", userService.AssignAdvisor)
	lecturers := api.Group("/lecturers", middleware.Protected(), middleware.RoleMiddleware("Admin"))
	lecturers.Get("/", userService.GetAllLecturers)

	// B. Achievement Routes
	achievements := api.Group("/achievements", middleware.Protected())
	
	// [NEW] Admin View All
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
}