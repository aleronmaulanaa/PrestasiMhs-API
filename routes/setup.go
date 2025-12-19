// package routes

// import (
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/app/services"
// 	"PrestasiMhs-API/middleware"
// 	"github.com/gofiber/fiber/v2"

// 	// [PENTING] Import ini wajib ada agar route swagger dikenali
// 	"github.com/gofiber/swagger"
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

// 	// Reporting
// 	reportRepo := repositories.NewReportRepository()
// 	// Inject achievementRepo juga agar Service bisa akses Mongo untuk detail laporan
// 	reportService := services.NewReportService(reportRepo, achievementRepo)

// 	// ============================================
// 	// 2. ROUTE DEFINITIONS
// 	// ============================================

// 	// --- Auth Routes ---
// 	auth := api.Group("/auth")
// 	auth.Post("/login", authService.Login)
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

// 	// --- Relations Management ---
	
// 	// Group Students
// 	students := api.Group("/students", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	students.Get("/", userService.GetAllStudents)
// 	students.Put("/:id/advisor", userService.AssignAdvisor)
	
// 	// [NEW] Endpoint Tambahan Sesuai SRS & Request
// 	students.Get("/:id", userService.GetStudentByID) 
// 	// Endpoint ini menggunakan AchievementService karena butuh akses DB Prestasi
// 	students.Get("/:id/achievements", achievementService.GetAchievementsByStudentID) 

// 	// Group Lecturers
// 	lecturers := api.Group("/lecturers", middleware.Protected(), middleware.RoleMiddleware("Admin"))
// 	lecturers.Get("/", userService.GetAllLecturers)
	
// 	// [NEW] Endpoint Tambahan Sesuai SRS & Request
// 	lecturers.Get("/:id/advisees", userService.GetLecturerAdviseesSRS)

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

// 	// --- Report Routes ---
// 	reports := api.Group("/reports", middleware.Protected())

// 	// Dashboard Statistics
// 	reports.Get("/statistics", middleware.RoleMiddleware("Admin"), reportService.GetDashboardStatistics)
	
// 	// Student Report (Transkrip)
// 	reports.Get("/student/:studentID", middleware.RoleMiddleware("Admin", "Dosen Wali"), reportService.GetStudentReport)

// 	// ============================================
// 	// 3. DOCUMENTATION ROUTE (SWAGGER)
// 	// ============================================
// 	app.Get("/swagger/*", swagger.HandlerDefault)
// }


package routes

import (
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/app/services"
	"PrestasiMhs-API/middleware"
	"github.com/gofiber/fiber/v2"

	// [PENTING] Import ini wajib ada agar route swagger dikenali
	"github.com/gofiber/swagger"
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

	// Reporting
	reportRepo := repositories.NewReportRepository()
	// Inject achievementRepo juga agar Service bisa akses Mongo untuk detail laporan
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
	
	// Group Students
	// [PERBAIKAN] Hapus "Admin" dari sini agar Mahasiswa bisa masuk
	students := api.Group("/students", middleware.Protected())

	// Endpoint Khusus Admin
	students.Get("/", middleware.RoleMiddleware("Admin"), userService.GetAllStudents)
	students.Put("/:id/advisor", middleware.RoleMiddleware("Admin"), userService.AssignAdvisor)
	
	// Endpoint Umum (Admin, Dosen, Mahasiswa bisa akses)
	// Mahasiswa butuh akses ini untuk melihat detail dirinya sendiri
	students.Get("/:id", middleware.RoleMiddleware("Admin", "Dosen Wali", "Mahasiswa"), userService.GetStudentByID) 
	students.Get("/:id/achievements", middleware.RoleMiddleware("Admin", "Dosen Wali", "Mahasiswa"), achievementService.GetAchievementsByStudentID) 

	// Group Lecturers
	// [PERBAIKAN] Hapus "Admin" dari sini agar Dosen bisa masuk
	lecturers := api.Group("/lecturers", middleware.Protected())
	
	// Endpoint Khusus Admin
	lecturers.Get("/", middleware.RoleMiddleware("Admin"), userService.GetAllLecturers)
	
	// Endpoint Admin & Dosen Wali
	lecturers.Get("/:id/advisees", middleware.RoleMiddleware("Admin", "Dosen Wali"), userService.GetLecturerAdviseesSRS)

	// --- B. Achievement Routes ---
	achievements := api.Group("/achievements", middleware.Protected())
	
	// Admin View All
	achievements.Get("/", middleware.RoleMiddleware("Admin"), achievementService.GetAllAchievements)

	// Static Routes
	achievements.Post("/", middleware.RoleMiddleware("Mahasiswa"), achievementService.CreateAchievement)
	achievements.Get("/my", middleware.RoleMiddleware("Mahasiswa"), achievementService.GetMyAchievements)
	achievements.Get("/advisees", middleware.RoleMiddleware("Dosen Wali"), achievementService.GetAdviseeAchievements)

	// Dynamic Routes (/:id)
	achievements.Get("/:id", achievementService.GetAchievementByID)
	achievements.Get("/:id/history", achievementService.GetAchievementHistory)
	achievements.Put("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.UpdateAchievement)
	achievements.Delete("/:id", middleware.RoleMiddleware("Mahasiswa"), achievementService.DeleteAchievement)
	achievements.Post("/:id/submit", middleware.RoleMiddleware("Mahasiswa"), achievementService.SubmitAchievement)
	achievements.Post("/:id/verify", middleware.RoleMiddleware("Dosen Wali"), achievementService.VerifyAchievement)
	achievements.Post("/:id/reject", middleware.RoleMiddleware("Dosen Wali"), achievementService.RejectAchievement)

	// --- Report Routes ---
	reports := api.Group("/reports", middleware.Protected())

	// Dashboard Statistics
	reports.Get("/statistics", middleware.RoleMiddleware("Admin"), reportService.GetDashboardStatistics)
	
	// Student Report (Transkrip)
	reports.Get("/student/:studentID", middleware.RoleMiddleware("Admin", "Dosen Wali"), reportService.GetStudentReport)

	// ============================================
	// 3. DOCUMENTATION ROUTE (SWAGGER)
	// ============================================
	app.Get("/swagger/*", swagger.HandlerDefault)
}