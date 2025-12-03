// package services

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/app/repositories"
// 	"fmt"
// 	"path/filepath"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type AchievementService interface {
// 	CreateAchievement(c *fiber.Ctx) error
// }

// type achievementService struct {
// 	repo repositories.AchievementRepository
// }

// func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
// 	return &achievementService{
// 		repo: repo,
// 	}
// }

// func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
// 	// 1. Ambil User ID dari Token (Middleware)
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// 2. Validasi & Ambil Student ID
// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Hanya mahasiswa terdaftar yang boleh upload prestasi",
// 		})
// 	}

// 	// 3. Parsing Form Data (Text Fields)
// 	var req models.CreateAchievementRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Format input salah",
// 		})
// 	}

// 	// 4. Handle File Upload
// 	// Folder ./uploads/documents SUDAH DIJAMIN ADA oleh main.go
// 	file, err := c.FormFile("file") 
// 	var attachments []models.Attachment

// 	if err == nil { // Jika ada file yang diupload
// 		// Generate nama unik
// 		ext := filepath.Ext(file.Filename)
// 		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
// 		filePath := fmt.Sprintf("./uploads/documents/%s", newFileName)

// 		// Simpan file
// 		if err := c.SaveFile(file, filePath); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  "error",
// 				"message": "Gagal menyimpan file",
// 			})
// 		}

// 		// Tambahkan ke struct attachment
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


package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"fmt"
	"path/filepath"
	"strings" // Import strings untuk cek ekstensi file
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
}

type achievementService struct {
	repo repositories.AchievementRepository
}

func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	// 1. Ambil User ID dari Token (Middleware)
	userID := c.Locals("user_id").(uuid.UUID).String()

	// 2. Validasi & Ambil Student ID
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Hanya mahasiswa terdaftar yang boleh upload prestasi",
		})
	}

	// 3. Parsing Form Data (Text Fields)
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

	if err == nil { // Jika ada file yang diupload
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		// --- LOGIKA PEMISAHAN FOLDER ---
		var subFolder string
		
		// Ubah ekstensi ke huruf kecil agar .JPG dan .jpg dianggap sama
		lowerExt := strings.ToLower(ext)

		switch lowerExt {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			subFolder = "photos"
		default:
			subFolder = "documents" // Default untuk PDF, DOCX, ZIP, dll
		}

		// Tentukan path penyimpanan (Folder documents/photos dijamin ada oleh main.go)
		filePath := fmt.Sprintf("./uploads/%s/%s", subFolder, newFileName)
		// --------------------------------

		// Simpan file
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal menyimpan file",
			})
		}

		// Tambahkan ke struct attachment
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