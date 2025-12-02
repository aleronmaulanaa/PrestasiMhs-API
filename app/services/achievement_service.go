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

// 	// 2. Cari Student ID berdasarkan User ID
// 	// Karena di tabel achievements butuh student_id, bukan user_id
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
// 	file, err := c.FormFile("file") // Nama key di Postman harus "file"
// 	var attachments []models.Attachment

// 	if err == nil { // Jika ada file yang diupload
// 		// Generate nama unik agar tidak bentrok
// 		ext := filepath.Ext(file.Filename)
// 		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
// 		filePath := fmt.Sprintf("./uploads/documents/%s", newFileName)

// 		// Simpan file ke folder
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
// 	// Kita konversi input string tanggal ke time.Time
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

// 	// 6. Simpan ke Database (Repo Hybrid)
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
	"os" // Import OS untuk cek/buat folder
	"path/filepath"
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

	// 2. Cari Student ID berdasarkan User ID
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

	if err == nil { // Jika ada file yang diupload
		// --- PERBAIKAN DI SINI ---
		// Tentukan folder tujuan
		uploadDir := "./uploads/documents"

		// Cek apakah folder ada, jika tidak BUATKAN
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			fmt.Println("📂 Folder belum ada, sedang membuat folder:", uploadDir)
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				fmt.Println("❌ Gagal membuat folder:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Gagal menyiapkan folder penyimpanan",
				})
			}
		}

		// Generate nama unik
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := fmt.Sprintf("%s/%s", uploadDir, newFileName)

		// Simpan file
		if err := c.SaveFile(file, filePath); err != nil {
			fmt.Println("❌ Error SaveFile:", err) // Debugging Log
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal menyimpan file: " + err.Error(),
			})
		}
		// -------------------------

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