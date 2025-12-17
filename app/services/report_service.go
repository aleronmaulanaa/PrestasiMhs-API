// package services

// import (
// 	"PrestasiMhs-API/app/repositories"
// 	"github.com/gofiber/fiber/v2"
// )

// type ReportService interface {
// 	GetDashboardStatistics(c *fiber.Ctx) error
// }

// type reportService struct {
// 	repo repositories.ReportRepository
// }

// func NewReportService(repo repositories.ReportRepository) ReportService {
// 	return &reportService{
// 		repo: repo,
// 	}
// }

// func (s *reportService) GetDashboardStatistics(c *fiber.Ctx) error {
// 	// 1. Hitung User
// 	totalStudents, _ := s.repo.CountUsersByRole("Mahasiswa")
// 	totalLecturers, _ := s.repo.CountUsersByRole("Dosen Wali")

// 	// 2. Hitung Prestasi
// 	achievementStats, err := s.repo.CountAchievementsByStatus()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	// 3. Format Data
// 	data := fiber.Map{
// 		"users": fiber.Map{
// 			"total_students":  totalStudents,
// 			"total_lecturers": totalLecturers,
// 		},
// 		"achievements": achievementStats,
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"data":   data,
// 	})
// }


package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportService interface {
	GetDashboardStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error // [NEW]
}

type reportService struct {
	repo            repositories.ReportRepository
	achievementRepo repositories.AchievementRepository // [NEW] Butuh ini untuk ambil detail Mongo
}

// [UPDATE] Constructor menerima 2 Repository
func NewReportService(repo repositories.ReportRepository, aRepo repositories.AchievementRepository) ReportService {
	return &reportService{
		repo:            repo,
		achievementRepo: aRepo,
	}
}

// GetDashboardStatistics godoc
// @Summary      Dashboard Statistics
// @Description  Melihat total user dan statistik prestasi (Admin Only)
// @Tags         Reports (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /reports/statistics [get]
func (s *reportService) GetDashboardStatistics(c *fiber.Ctx) error {
	// 1. Hitung User
	totalStudents, _ := s.repo.CountUsersByRole("Mahasiswa")
	totalLecturers, _ := s.repo.CountUsersByRole("Dosen Wali")

	// 2. Hitung Prestasi
	achievementStats, err := s.repo.CountAchievementsByStatus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// 3. Format Data
	data := fiber.Map{
		"users": fiber.Map{
			"total_students":  totalStudents,
			"total_lecturers": totalLecturers,
		},
		"achievements": achievementStats,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// [NEW] Logic Transkrip Prestasi Mahasiswa

// GetStudentReport godoc
// @Summary      Student Transcript
// @Description  Melihat laporan/transkrip prestasi lengkap mahasiswa
// @Tags         Reports (Admin)
// @Produce      json
// @Security     BearerAuth
// @Param        studentID  path  string  true  "ID Tabel Students (UUID)"
// @Success      200  {object}  models.StudentReportResponse
// @Router       /reports/student/{studentID} [get]
func (s *reportService) GetStudentReport(c *fiber.Ctx) error {
	targetStudentID := c.Params("studentID") // ID dari tabel students
	_ = c.Locals("user_id").(uuid.UUID).String() // User yang request (bisa dipakai untuk validasi role)

	// 1. Ambil Header Mahasiswa
	header, err := s.repo.GetStudentHeader(targetStudentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Mahasiswa tidak ditemukan"})
	}

	// 2. Ambil List Prestasi (SQL)
	refs, err := s.repo.GetStudentAchievements(targetStudentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// 3. Ambil Detail Mongo (untuk Judul & Event Date)
	var mongoIDs []string
	for _, ref := range refs {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}
	// Reuse method dari AchievementRepository
	mongoDetails, _ := s.achievementRepo.FindMongoDetails(mongoIDs)

	// 4. Merge Data
	var reportItems []models.AchievementReportItem
	totalVerified := 0

	for _, ref := range refs {
		var title, pType string
		var eventDate *time.Time

		if detail, ok := mongoDetails[ref.MongoAchievementID]; ok {
			title = detail.Title
			pType = detail.AchievementType
			eventDate = &detail.Details.EventDate
		}

		if ref.Status == "verified" {
			totalVerified++
		}

		reportItems = append(reportItems, models.AchievementReportItem{
			Title:            title,
			Type:             pType,
			EventDate:        eventDate,
			Status:           ref.Status,
			VerificationDate: ref.VerifiedAt,
		})
	}

	// 5. Final Response
	response := models.StudentReportResponse{
		StudentInfo: *header,
		Summary: models.ReportSummary{
			TotalEntries:  len(refs),
			TotalVerified: totalVerified,
			TotalPoints:   totalVerified * 10, // Simulasi poin
		},
		Achievements: reportItems,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}