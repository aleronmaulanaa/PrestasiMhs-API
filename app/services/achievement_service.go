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
// 	GetMyAchievements(c *fiber.Ctx) error      
// 	GetAdviseeAchievements(c *fiber.Ctx) error 
// 	SubmitAchievement(c *fiber.Ctx) error      
	
// 	// [Fase 1: Update & Delete]
// 	UpdateAchievement(c *fiber.Ctx) error
// 	DeleteAchievement(c *fiber.Ctx) error

// 	// [Fase 1: Detail & History]
// 	GetAchievementByID(c *fiber.Ctx) error
// 	GetAchievementHistory(c *fiber.Ctx) error

//     // [FIX: Split Verify & Reject sesuai SRS FR-007 & FR-008]
// 	VerifyAchievement(c *fiber.Ctx) error
// 	RejectAchievement(c *fiber.Ctx) error
// }

// type achievementService struct {
// 	repo repositories.AchievementRepository
// }

// func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
// 	return &achievementService{
// 		repo: repo,
// 	}
// }

// // --- 1. FEATURE: UPLOAD & MANAGE PRESTASI (Mahasiswa) ---

// func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Hanya mahasiswa terdaftar yang boleh upload prestasi"})
// 	}

// 	var req models.CreateAchievementRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
// 	}

// 	file, err := c.FormFile("file")
// 	var attachments []models.Attachment

// 	if err == nil {
// 		ext := filepath.Ext(file.Filename)
// 		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		
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
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan file"})
// 		}

// 		attachments = append(attachments, models.Attachment{
// 			FileName:   file.Filename,
// 			FileURL:    filePath,
// 			FileType:   file.Header.Get("Content-Type"),
// 			UploadedAt: time.Now(),
// 		})
// 	}

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
// 			MedalType:        req.MedalType,
//             // Mapping field lain...
// 		},
// 	}

// 	if err := s.repo.Create(mongoData, studentID); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil disimpan sebagai draft"})
// }

// func (s *achievementService) UpdateAchievement(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"})
// 	}

// 	ref, err := s.repo.FindRefByID(id)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
// 	}

// 	if ref.StudentID != studentID {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak memiliki akses ke prestasi ini"})
// 	}

// 	if ref.Status != "draft" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Hanya prestasi berstatus draft yang bisa diubah"})
// 	}

// 	var req models.CreateAchievementRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
// 	}

// 	var attachments []models.Attachment
// 	file, err := c.FormFile("file")
	
// 	if err == nil {
// 		ext := filepath.Ext(file.Filename)
// 		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
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
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan file baru"})
// 		}
// 		attachments = append(attachments, models.Attachment{
// 			FileName:   file.Filename,
// 			FileURL:    filePath,
// 			FileType:   file.Header.Get("Content-Type"),
// 			UploadedAt: time.Now(),
// 		})
// 	}

// 	eventDate, _ := time.Parse("2006-01-02", req.EventDate)

// 	updateData := &models.AchievementMongo{
// 		AchievementType: req.AchievementType,
// 		Title:           req.Title,
// 		Description:     req.Description,
// 		Attachments:     attachments, 
// 		Details: models.AchievementDetails{
// 			CompetitionName:  req.CompetitionName,
// 			CompetitionLevel: req.CompetitionLevel,
// 			Rank:             req.Rank,
// 			OrganizationName: req.OrganizationName,
// 			Position:         req.Position,
// 			Location:         req.Location,
// 			Organizer:        req.Organizer,
// 			EventDate:        eventDate,
// 			MedalType:        req.MedalType,
// 		},
// 	}

// 	if err := s.repo.UpdateMongo(ref.MongoAchievementID, updateData); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengupdate data: " + err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil diperbarui"})
// }

// func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"})
// 	}

// 	ref, err := s.repo.FindRefByID(id)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
// 	}

// 	if ref.StudentID != studentID {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak memiliki akses menghapus prestasi ini"})
// 	}

// 	if err := s.repo.SoftDelete(id); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil dihapus (soft delete)"})
// }

// func (s *achievementService) SubmitAchievement(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "User tidak valid"})
// 	}

// 	ref, err := s.repo.FindRefByID(id)
// 	if err != nil || ref.StudentID != studentID {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
// 	}

// 	if err := s.repo.Submit(id); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil disubmit dan siap diverifikasi"})
// }

// // --- 2. FEATURE: READ DATA (Common) ---

// func (s *achievementService) mergeData(refs []models.AchievementReference) ([]models.AchievementReference, error) {
// 	if len(refs) == 0 {
// 		return refs, nil
// 	}
// 	var mongoIDs []string
// 	for _, ref := range refs {
// 		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
// 	}
// 	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range refs {
// 		if detail, exists := mongoDetails[refs[i].MongoAchievementID]; exists {
// 			refs[i].Detail = &detail
// 		}
// 	}
// 	return refs, nil
// }

// func (s *achievementService) GetMyAchievements(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uuid.UUID).String()
// 	studentID, err := s.repo.GetStudentIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Data mahasiswa tidak ditemukan"})
// 	}
// 	refs, err := s.repo.FindAllByStudentID(studentID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	finalData, err := s.mergeData(refs)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "data": finalData})
// }

// func (s *achievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uuid.UUID).String()
// 	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
// 	}
// 	refs, err := s.repo.FindAllByAdvisorID(advisorID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	finalData, err := s.mergeData(refs)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "data": finalData})
// }

// func (s *achievementService) GetAchievementByID(c *fiber.Ctx) error {
//     id := c.Params("id")
//     ref, err := s.repo.FindRefByID(id)
//     if err != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
//     }
//     mongoIDs := []string{ref.MongoAchievementID}
//     mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail data"})
//     }
//     if detail, exists := mongoDetails[ref.MongoAchievementID]; exists {
//         ref.Detail = &detail
//     }
//     return c.JSON(fiber.Map{"status": "success", "data": ref})
// }

// func (s *achievementService) GetAchievementHistory(c *fiber.Ctx) error {
//     id := c.Params("id")
//     ref, err := s.repo.FindRefByID(id)
//     if err != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
//     }
//     var history []fiber.Map
//     history = append(history, fiber.Map{
//         "status":    "draft",
//         "timestamp": ref.CreatedAt,
//         "note":      "Prestasi dibuat (Draft)",
//         "actor":     "Mahasiswa",
//     })
//     if ref.SubmittedAt != nil {
//         history = append(history, fiber.Map{
//             "status":    "submitted",
//             "timestamp": ref.SubmittedAt,
//             "note":      "Menunggu verifikasi Dosen Wali",
//             "actor":     "Mahasiswa",
//         })
//     }
//     if ref.VerifiedAt != nil {
//         note := "Prestasi telah diverifikasi"
//         if ref.Status == "rejected" {
//             note = "Prestasi ditolak: " + ref.RejectionNote
//         }
//         history = append(history, fiber.Map{
//             "status":    ref.Status, 
//             "timestamp": ref.VerifiedAt,
//             "note":      note,
//             "actor":     "Dosen Wali",
//         })
//     }
//     return c.JSON(fiber.Map{"status": "success", "data": history})
// }

// // --- 3. FEATURE: WORKFLOW VERIFICATION (Dosen Wali) ---

// // VerifyAchievement: Mengubah status menjadi 'verified' (FR-007)
// func (s *achievementService) VerifyAchievement(c *fiber.Ctx) error {
// 	achievementID := c.Params("id")
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// Logic FR-007: Dosen approve, status jadi verified. Tidak wajib ada notes.
// 	err := s.repo.UpdateStatus(achievementID, "verified", "", userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"message": "Prestasi berhasil diverifikasi",
// 	})
// }

// // RejectAchievement: Mengubah status menjadi 'rejected' dengan catatan (FR-008)
// func (s *achievementService) RejectAchievement(c *fiber.Ctx) error {
// 	achievementID := c.Params("id")
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// Logic FR-008: Wajib ada rejection note
// 	type RejectRequest struct {
// 		Notes string `json:"notes" validate:"required"`
// 	}

// 	var req RejectRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
// 	}

// 	if strings.TrimSpace(req.Notes) == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Catatan penolakan wajib diisi"})
// 	}

// 	err := s.repo.UpdateStatus(achievementID, "rejected", req.Notes, userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"message": "Prestasi berhasil ditolak",
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
	GetMyAchievements(c *fiber.Ctx) error
	GetAdviseeAchievements(c *fiber.Ctx) error
	SubmitAchievement(c *fiber.Ctx) error

	// [Fase 1: Update & Delete]
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error

	// [Fase 1: Detail & History (FIXED)]
	GetAchievementByID(c *fiber.Ctx) error
	GetAchievementHistory(c *fiber.Ctx) error

	// [FIX: Verify & Reject]
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
}

type achievementService struct {
	repo repositories.AchievementRepository
}

func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

// --- 1. FEATURE: UPLOAD & MANAGE PRESTASI (Mahasiswa) ---

func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()

	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Hanya mahasiswa terdaftar yang boleh upload prestasi"})
	}

	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
	}

	file, err := c.FormFile("file")
	var attachments []models.Attachment

	if err == nil {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan file"})
		}

		attachments = append(attachments, models.Attachment{
			FileName:   file.Filename,
			FileURL:    filePath,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		})
	}

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
			MedalType:        req.MedalType,
			// Mapping field lain bisa ditambahkan disini
		},
	}

	if err := s.repo.Create(mongoData, studentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil disimpan sebagai draft"})
}

func (s *achievementService) UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"})
	}

	ref, err := s.repo.FindRefByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak memiliki akses ke prestasi ini"})
	}

	if ref.Status != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Hanya prestasi berstatus draft yang bisa diubah"})
	}

	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
	}

	var attachments []models.Attachment
	file, err := c.FormFile("file")

	if err == nil {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan file baru"})
		}
		attachments = append(attachments, models.Attachment{
			FileName:   file.Filename,
			FileURL:    filePath,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		})
	}

	eventDate, _ := time.Parse("2006-01-02", req.EventDate)

	updateData := &models.AchievementMongo{
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Attachments:     attachments,
		Details: models.AchievementDetails{
			CompetitionName:  req.CompetitionName,
			CompetitionLevel: req.CompetitionLevel,
			Rank:             req.Rank,
			OrganizationName: req.OrganizationName,
			Position:         req.Position,
			Location:         req.Location,
			Organizer:        req.Organizer,
			EventDate:        eventDate,
			MedalType:        req.MedalType,
		},
	}

	if err := s.repo.UpdateMongo(ref.MongoAchievementID, updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengupdate data: " + err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil diperbarui"})
}

func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"})
	}

	ref, err := s.repo.FindRefByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	if ref.StudentID != studentID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak memiliki akses menghapus prestasi ini"})
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil dihapus (soft delete)"})
}

func (s *achievementService) SubmitAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "User tidak valid"})
	}

	ref, err := s.repo.FindRefByID(id)
	if err != nil || ref.StudentID != studentID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	if err := s.repo.Submit(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Prestasi berhasil disubmit dan siap diverifikasi"})
}

// --- 2. FEATURE: READ DATA (Common) ---

func (s *achievementService) mergeData(refs []models.AchievementReference) ([]models.AchievementReference, error) {
	if len(refs) == 0 {
		return refs, nil
	}
	var mongoIDs []string
	for _, ref := range refs {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
	}
	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
	if err != nil {
		return nil, err
	}
	for i := range refs {
		if detail, exists := mongoDetails[refs[i].MongoAchievementID]; exists {
			refs[i].Detail = &detail
		}
	}
	return refs, nil
}

func (s *achievementService) GetMyAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Data mahasiswa tidak ditemukan"})
	}
	refs, err := s.repo.FindAllByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	finalData, err := s.mergeData(refs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}

func (s *achievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
	}
	refs, err := s.repo.FindAllByAdvisorID(advisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	finalData, err := s.mergeData(refs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail prestasi"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}

// [FIX] GetAchievementByID dengan Security Check
func (s *achievementService) GetAchievementByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	// Ambil data prestasi
	ref, err := s.repo.FindRefByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	// --- SECURITY CHECK (Aturan 3 & 4) ---
	// Cek role user berdasarkan keberadaan datanya di tabel mahasiswa/dosen
	isAllowed := false

	// 1. Cek apakah Mahasiswa & Pemilik
	studentID, errMhs := s.repo.GetStudentIDByUserID(userID)
	if errMhs == nil {
		if ref.StudentID == studentID {
			isAllowed = true
		}
	}

	// 2. Cek apakah Dosen Wali & Bimbingannya
	if !isAllowed {
		advisorID, errDos := s.repo.GetAdvisorIDByUserID(userID)
		if errDos == nil {
			isAdvisee, _ := s.repo.IsAdvisee(advisorID, ref.StudentID)
			if isAdvisee {
				isAllowed = true
			}
		}
	}

	// 3. (Opsional) Jika Admin, bisa set isAllowed = true disini (Logic nanti Fase 2)
	// Namun untuk Fase 1, kita fokus Mhs & Dosen

	if !isAllowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak memiliki hak akses untuk melihat prestasi ini"})
	}
	// -------------------------------------

	mongoIDs := []string{ref.MongoAchievementID}
	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil detail data"})
	}
	if detail, exists := mongoDetails[ref.MongoAchievementID]; exists {
		ref.Detail = &detail
	}
	return c.JSON(fiber.Map{"status": "success", "data": ref})
}

// [FIX] GetAchievementHistory dengan Security Check
func (s *achievementService) GetAchievementHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	ref, err := s.repo.FindRefByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	// --- SECURITY CHECK (Aturan 3 & 4) - Logic Sama dengan GetByID ---
	isAllowed := false
	studentID, errMhs := s.repo.GetStudentIDByUserID(userID)
	if errMhs == nil && ref.StudentID == studentID {
		isAllowed = true
	}
	if !isAllowed {
		advisorID, errDos := s.repo.GetAdvisorIDByUserID(userID)
		if errDos == nil {
			isAdvisee, _ := s.repo.IsAdvisee(advisorID, ref.StudentID)
			if isAdvisee {
				isAllowed = true
			}
		}
	}

	if !isAllowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"})
	}
	// -------------------------------------------------------------

	var history []fiber.Map
	history = append(history, fiber.Map{
		"status":    "draft",
		"timestamp": ref.CreatedAt,
		"note":      "Prestasi dibuat (Draft)",
		"actor":     "Mahasiswa",
	})
	if ref.SubmittedAt != nil {
		history = append(history, fiber.Map{
			"status":    "submitted",
			"timestamp": ref.SubmittedAt,
			"note":      "Menunggu verifikasi Dosen Wali",
			"actor":     "Mahasiswa",
		})
	}
	if ref.VerifiedAt != nil {
		note := "Prestasi telah diverifikasi"
		if ref.Status == "rejected" {
			note = "Prestasi ditolak: " + ref.RejectionNote
		}
		history = append(history, fiber.Map{
			"status":    ref.Status,
			"timestamp": ref.VerifiedAt,
			"note":      note,
			"actor":     "Dosen Wali",
		})
	}
	return c.JSON(fiber.Map{"status": "success", "data": history})
}

// --- 3. FEATURE: WORKFLOW VERIFICATION (Dosen Wali) ---

func (s *achievementService) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	// 1. Ambil data prestasi
	ref, err := s.repo.FindRefByID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	// 2. [SECURITY CHECK] Pastikan Dosen Memverifikasi Bimbingannya (Rule 3)
	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
	}

	isAdvisee, err := s.repo.IsAdvisee(advisorID, ref.StudentID)
	if err != nil || !isAdvisee {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak berhak memverifikasi mahasiswa ini"})
	}

	// 3. Proses Update
	err = s.repo.UpdateStatus(achievementID, "verified", "", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil diverifikasi",
	})
}

func (s *achievementService) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()

	// 1. Ambil data prestasi
	ref, err := s.repo.FindRefByID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Prestasi tidak ditemukan"})
	}

	// 2. [SECURITY CHECK] Rule 3
	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda bukan dosen wali"})
	}

	isAdvisee, err := s.repo.IsAdvisee(advisorID, ref.StudentID)
	if err != nil || !isAdvisee {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Anda tidak berhak menolak prestasi mahasiswa ini"})
	}

	// 3. Validasi Notes
	type RejectRequest struct {
		Notes string `json:"notes" validate:"required"`
	}
	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
	}
	if strings.TrimSpace(req.Notes) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Catatan penolakan wajib diisi"})
	}

	// 4. Proses Update
	err = s.repo.UpdateStatus(achievementID, "rejected", req.Notes, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil ditolak",
	})
}