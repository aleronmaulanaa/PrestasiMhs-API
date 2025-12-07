// package services

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/app/repositories"
// 	"fmt"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type AchievementService interface {
// 	CreateAchievement(c *fiber.Ctx) error
// 	GetMyAchievements(c *fiber.Ctx) error      // Untuk Mahasiswa
// 	GetAdviseeAchievements(c *fiber.Ctx) error // Untuk Dosen Wali
// 	VerifyAchievement(c *fiber.Ctx) error      // Untuk Dosen Wali (Approve/Reject)
// }

// type achievementService struct {
// 	repo repositories.AchievementRepository
// }

// func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
// 	return &achievementService{
// 		repo: repo,
// 	}
// }

// // --- 1. FEATURE: UPLOAD PRESTASI ---

// func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
// 	// 1. Ambil User ID dari Token
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// 2. Validasi & Ambil Student ID
// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Hanya mahasiswa terdaftar yang boleh upload prestasi",
// 		})
// 	}

// 	// 3. Parsing Form Data
// 	var req models.CreateAchievementRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Format input salah",
// 		})
// 	}

// 	// 4. Handle File Upload
// 	file, err := c.FormFile("file")
// 	var attachments []models.Attachment

// 	if err == nil {
// 		ext := filepath.Ext(file.Filename)
// 		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		
// 		// Logika pemisahan folder
// 		var subFolder string
// 		lowerExt := strings.ToLower(ext)
// 		switch lowerExt {
// 		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
// 			subFolder = "photos"
// 		default:
// 			subFolder = "documents"
// 		}

// 		filePath := fmt.Sprintf("./uploads/%s/%s", subFolder, newFileName)

// 		if err := c.SaveFile(file, filePath); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  "error",
// 				"message": "Gagal menyimpan file",
// 			})
// 		}

// 		attachments = append(attachments, models.Attachment{
// 			FileName:   file.Filename,
// 			FileURL:    filePath,
// 			FileType:   file.Header.Get("Content-Type"),
// 			UploadedAt: time.Now(),
// 		})
// 	}

// 	// 5. Mapping ke MongoDB Model
// 	eventDate, _ := time.Parse("2006-01-02", req.EventDate)

// 	mongoData := &models.AchievementMongo{
// 		ID:              primitive.NewObjectID(),
// 		StudentID:       studentID,
// 		AchievementType: req.AchievementType,
// 		Title:           req.Title,
// 		Description:     req.Description,
// 		Attachments:     attachments,
// 		CreatedAt:       time.Now(),
// 		UpdatedAt:       time.Now(),
// 		Details: models.AchievementDetails{
// 			CompetitionName:  req.CompetitionName,
// 			CompetitionLevel: req.CompetitionLevel,
// 			Rank:             req.Rank,
// 			OrganizationName: req.OrganizationName,
// 			Position:         req.Position,
// 			Location:         req.Location,
// 			Organizer:        req.Organizer,
// 			EventDate:        eventDate,
// 		},
// 	}

// 	// 6. Simpan ke Database
// 	if err := s.repo.Create(mongoData, studentID); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "Prestasi berhasil disimpan sebagai draft",
// 	})
// }

// // --- 2. FEATURE: READ DATA (HYBRID MERGE) ---

// // helper untuk menggabungkan data Postgres dan Mongo
// func (s *achievementService) mergeData(refs []models.AchievementReference) ([]models.AchievementReference, error) {
// 	if len(refs) == 0 {
// 		return refs, nil
// 	}

// 	// Kumpulkan semua Mongo ID dari hasil query Postgres
// 	var mongoIDs []string
// 	for _, ref := range refs {
// 		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
// 	}

// 	// Ambil detail dari MongoDB dalam satu query (Bulk Read)
// 	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Gabungkan data (Mapping)
// 	for i := range refs {
// 		if detail, exists := mongoDetails[refs[i].MongoAchievementID]; exists {
// 			refs[i].Detail = &detail
// 		}
// 	}

// 	return refs, nil
// }

// // GetMyAchievements: Mahasiswa melihat prestasi sendiri
// func (s *achievementService) GetMyAchievements(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uuid.UUID).String()
	
// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Data mahasiswa tidak ditemukan"})
// 	}

// 	// 1. Ambil Referensi dari Postgres
// 	refs, err := s.repo.FindAllByStudentID(studentID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	// 2. Gabungkan dengan detail Mongo
// 	finalData, err := s.mergeData(refs)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"data":   finalData,
// 	})
// }

// // GetAdviseeAchievements: Dosen Wali melihat prestasi mahasiswa bimbingan
// func (s *achievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
// 	}

// 	// 1. Ambil Referensi dari Postgres (Filter by Advisor ID)
// 	refs, err := s.repo.FindAllByAdvisorID(advisorID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	// 2. Gabungkan dengan detail Mongo
// 	finalData, err := s.mergeData(refs)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"data":   finalData,
// 	})
// }

// // --- 3. FEATURE: VERIFICATION (DOSEN WALI) ---

// type VerifyRequest struct {
// 	Status string `json:"status" validate:"required,oneof=verified rejected"`
// 	Notes  string `json:"notes"`
// }

// func (s *achievementService) VerifyAchievement(c *fiber.Ctx) error {
// 	// Ambil ID Prestasi dari URL parameter
// 	achievementID := c.Params("id")
	
// 	// Ambil ID Dosen dari Token
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// Parse Body
// 	var req VerifyRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format status salah"})
// 	}

// 	// Validasi Status
// 	if req.Status != "verified" && req.Status != "rejected" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Status hanya boleh 'verified' atau 'rejected'"})
// 	}

// 	// Update Status di Database
// 	err := s.repo.UpdateStatus(achievementID, req.Status, req.Notes, userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"message": "Status prestasi berhasil diperbarui",
// 	})
// }


package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error      // Untuk Mahasiswa
	GetAdviseeAchievements(c *fiber.Ctx) error // Untuk Dosen Wali
	VerifyAchievement(c *fiber.Ctx) error      // Untuk Dosen Wali (Approve/Reject)
	SubmitAchievement(c *fiber.Ctx) error      // [NEW] Untuk Mahasiswa (Draft -> Submitted)
}

type achievementService struct {
	repo repositories.AchievementRepository
}

func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

// --- 1. FEATURE: UPLOAD PRESTASI ---

func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	// 1. Ambil User ID dari Token
	userID := c.Locals("user_id").(uuid.UUID).String()

	// 2. Validasi & Ambil Student ID
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Hanya mahasiswa terdaftar yang boleh upload prestasi",
		})
	}

	// 3. Parsing Form Data
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input salah",
		})
	}

	// 4. Handle File Upload
	file, err := c.FormFile("file")
	var attachments []models.Attachment

	if err == nil {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		
		// Logika pemisahan folder
		var subFolder string
		lowerExt := strings.ToLower(ext)
		switch lowerExt {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			subFolder = "photos"
		default:
			subFolder = "documents"
		}

		filePath := fmt.Sprintf("./uploads/%s/%s", subFolder, newFileName)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal menyimpan file",
			})
		}

		attachments = append(attachments, models.Attachment{
			FileName:   file.Filename,
			FileURL:    filePath,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		})
	}

	// 5. Mapping ke MongoDB Model
	eventDate, _ := time.Parse("2006-01-02", req.EventDate)

	mongoData := &models.AchievementMongo{
		ID:              primitive.NewObjectID(),
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Attachments:     attachments,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Details: models.AchievementDetails{
			CompetitionName:  req.CompetitionName,
			CompetitionLevel: req.CompetitionLevel,
			Rank:             req.Rank,
			OrganizationName: req.OrganizationName,
			Position:         req.Position,
			Location:         req.Location,
			Organizer:        req.Organizer,
			EventDate:        eventDate,
		},
	}

	// 6. Simpan ke Database
	if err := s.repo.Create(mongoData, studentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil disimpan sebagai draft",
	})
}

// [NEW] SubmitAchievement: Mahasiswa mengirim draft untuk diverifikasi
func (s *achievementService) SubmitAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	// 1. Validasi Kepemilikan (Cek apakah prestasi ini milik user yang login)
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "User tidak valid"})
	}

	// Cek di database apakah ID prestasi ini milik studentID tersebut
	ref, err := s.repo.FindRefByID(id)
	if err != nil || ref.StudentID != studentID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	// 2. Lakukan Submit
	if err := s.repo.Submit(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(), // Pesan error jika status bukan draft
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil disubmit dan siap diverifikasi",
	})
}

// --- 2. FEATURE: READ DATA (HYBRID MERGE) ---

// helper untuk menggabungkan data Postgres dan Mongo
func (s *achievementService) mergeData(refs []models.AchievementReference) ([]models.AchievementReference, error) {
	if len(refs) == 0 {
		return refs, nil
	}

	// Kumpulkan semua Mongo ID dari hasil query Postgres
	var mongoIDs []string
	for _, ref := range refs {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}

	// Ambil detail dari MongoDB dalam satu query (Bulk Read)
	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
	if err != nil {
		return nil, err
	}

	// Gabungkan data (Mapping)
	for i := range refs {
		if detail, exists := mongoDetails[refs[i].MongoAchievementID]; exists {
			refs[i].Detail = &detail
		}
	}

	return refs, nil
}

// GetMyAchievements: Mahasiswa melihat prestasi sendiri
func (s *achievementService) GetMyAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Data mahasiswa tidak ditemukan"})
	}

	// 1. Ambil Referensi dari Postgres
	refs, err := s.repo.FindAllByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// 2. Gabungkan dengan detail Mongo
	finalData, err := s.mergeData(refs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   finalData,
	})
}

// GetAdviseeAchievements: Dosen Wali melihat prestasi mahasiswa bimbingan
func (s *achievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()

	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
	}

	// 1. Ambil Referensi dari Postgres (Filter by Advisor ID)
	refs, err := s.repo.FindAllByAdvisorID(advisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// 2. Gabungkan dengan detail Mongo
	finalData, err := s.mergeData(refs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   finalData,
	})
}

// --- 3. FEATURE: VERIFICATION (DOSEN WALI) ---

type VerifyRequest struct {
	Status string `json:"status" validate:"required,oneof=verified rejected"`
	Notes  string `json:"notes"`
}

func (s *achievementService) VerifyAchievement(c *fiber.Ctx) error {
	// Ambil ID Prestasi dari URL parameter
	achievementID := c.Params("id")
	
	// Ambil ID Dosen dari Token
	userID := c.Locals("user_id").(uuid.UUID).String()

	// Parse Body
	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format status salah"})
	}

	// Validasi Status
	if req.Status != "verified" && req.Status != "rejected" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Status hanya boleh 'verified' atau 'rejected'"})
	}

	// Update Status di Database
	err := s.repo.UpdateStatus(achievementID, req.Status, req.Notes, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Status prestasi berhasil diperbarui",
	})
}