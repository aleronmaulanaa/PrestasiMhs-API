package services

import (
	"PrestasiMhs-API/app/repositories"
	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	GetDashboardStatistics(c *fiber.Ctx) error
}

type reportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) ReportService {
	return &reportService{
		repo: repo,
	}
}

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