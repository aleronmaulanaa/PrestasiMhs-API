package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/utils"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	RegisterLecturer(c *fiber.Ctx) error
	RegisterStudent(c *fiber.Ctx) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// RegisterLecturer menangani pendaftaran Dosen Baru
func (s *userService) RegisterLecturer(c *fiber.Ctx) error {
	var req models.CreateLecturerRequest

	// 1. Parsing Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Hash Password (SECURITY)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi password",
		})
	}

	// 3. Mapping ke Model
	userModel := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
	}

	lecturerModel := &models.LecturerInfo{
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	// 4. Simpan ke Database via Repository
	if err := s.repo.CreateLecturer(userModel, lecturerModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mendaftarkan dosen: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Dosen berhasil didaftarkan",
	})
}

// RegisterStudent menangani pendaftaran Mahasiswa Baru
func (s *userService) RegisterStudent(c *fiber.Ctx) error {
	var req models.CreateStudentRequest

	// 1. Parsing Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengenkripsi password",
		})
	}

	// 3. Mapping ke Model
	userModel := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
	}

	studentModel := &models.StudentInfo{
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
	}

	// 4. Simpan ke Database
	// req.AdvisorID dikirim terpisah karena bisa jadi string kosong (jika belum ada dosen wali)
	if err := s.repo.CreateStudent(userModel, studentModel, req.AdvisorID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mendaftarkan mahasiswa: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Mahasiswa berhasil didaftarkan",
	})
}