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
	
	// [NEW] Admin View All
	GetAllAchievements(c *fiber.Ctx) error
	
	SubmitAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
	GetAchievementByID(c *fiber.Ctx) error
	GetAchievementHistory(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error

	GetAchievementsByStudentID(c *fiber.Ctx) error // [NEW]
}

type achievementService struct {
	repo repositories.AchievementRepository
}

func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
	return &achievementService{
		repo: repo,
	}
}

// ... [KODE FASE 1 SEBELUMNYA TETAP SAMA] ...

// CreateAchievement godoc
// @Summary      Upload Prestasi Baru
// @Description  Mahasiswa mengupload data prestasi (Draft). File upload wajib via form-data.
// @Tags         Achievements (Mahasiswa)
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        title formData string true "Judul Prestasi"
// @Param        achievement_type formData string true "Tipe: Kompetisi / Organisasi"
// @Param        description formData string true "Deskripsi"
// @Param        event_date formData string true "Tanggal (YYYY-MM-DD)"
// @Param        file formData file false "Bukti Lampiran (PDF/IMG)"
// @Success      201  {object}  map[string]interface{}
// @Router       /achievements [post]
func (s *achievementService) CreateAchievement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Hanya mahasiswa terdaftar"}) }

	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"}) }
	file, err := c.FormFile("file")
	var attachments []models.Attachment
	if err == nil {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		var subFolder string
		lowerExt := strings.ToLower(ext)
		switch lowerExt {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp": subFolder = "photos"
		default: subFolder = "documents"
		}
		filePath := fmt.Sprintf("./uploads/%s/%s", subFolder, newFileName)
		if err := c.SaveFile(file, filePath); err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal file"}) }
		attachments = append(attachments, models.Attachment{FileName: file.Filename, FileURL: filePath, FileType: file.Header.Get("Content-Type"), UploadedAt: time.Now()})
	}
	eventDate, _ := time.Parse("2006-01-02", req.EventDate)
	mongoData := &models.AchievementMongo{
		ID: primitive.NewObjectID(), StudentID: studentID, AchievementType: req.AchievementType, Title: req.Title, Description: req.Description, Attachments: attachments, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		Details: models.AchievementDetails{
			CompetitionName: req.CompetitionName, CompetitionLevel: req.CompetitionLevel, Rank: req.Rank, OrganizationName: req.OrganizationName, Position: req.Position, Location: req.Location, Organizer: req.Organizer, EventDate: eventDate, MedalType: req.MedalType,
		},
	}
	if err := s.repo.Create(mongoData, studentID); err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Draft disimpan"})
}

// UpdateAchievement godoc
// @Summary      Update Data Prestasi
// @Description  Mengubah data prestasi (Hanya jika status Draft)
// @Tags         Achievements (Mahasiswa)
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Param        title formData string true "Judul Prestasi"
// @Param        description formData string true "Deskripsi"
// @Param        file formData file false "Update File (Optional)"
// @Success      200  {object}  map[string]interface{}
// @Router       /achievements/{id} [put]
func (s *achievementService) UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Akses ditolak"}) }
	ref, err := s.repo.FindRefByID(id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Not found"}) }
	if ref.StudentID != studentID { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "Bukan milik anda"}) }
	if ref.Status != "draft" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Bukan draft"}) }
	
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error"}) }
	var attachments []models.Attachment
	file, err := c.FormFile("file")
	if err == nil {
		// (Same logic upload file)
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := fmt.Sprintf("./uploads/documents/%s", newFileName) // simplify for brevity
		c.SaveFile(file, filePath)
		attachments = append(attachments, models.Attachment{FileName: file.Filename, FileURL: filePath, FileType: file.Header.Get("Content-Type"), UploadedAt: time.Now()})
	}
	eventDate, _ := time.Parse("2006-01-02", req.EventDate)
	updateData := &models.AchievementMongo{
		AchievementType: req.AchievementType, Title: req.Title, Description: req.Description, Attachments: attachments,
		Details: models.AchievementDetails{CompetitionName: req.CompetitionName, CompetitionLevel: req.CompetitionLevel, Rank: req.Rank, OrganizationName: req.OrganizationName, Position: req.Position, Location: req.Location, Organizer: req.Organizer, EventDate: eventDate, MedalType: req.MedalType},
	}
	s.repo.UpdateMongo(ref.MongoAchievementID, updateData)
	return c.JSON(fiber.Map{"status": "success", "message": "Updated"})
}

// DeleteAchievement godoc
// @Summary      Delete Prestasi (Draft)
// @Description  Menghapus (Soft Delete) prestasi yang masih draft
// @Tags         Achievements (Mahasiswa)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /achievements/{id} [delete]
func (s *achievementService) DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error"}) }
	ref, err := s.repo.FindRefByID(id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }
	if ref.StudentID != studentID { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error"}) }
	if err := s.repo.SoftDelete(id); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()}) }
	return c.JSON(fiber.Map{"status": "success", "message": "Deleted"})
}

// SubmitAchievement godoc
// @Summary      Submit Prestasi
// @Description  Mengubah status Draft menjadi Submitted (Siap Diverifikasi)
// @Tags         Achievements (Mahasiswa)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /achievements/{id}/submit [post]
func (s *achievementService) SubmitAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error"}) }
	ref, err := s.repo.FindRefByID(id)
	if err != nil || ref.StudentID != studentID { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }
	if err := s.repo.Submit(id); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()}) }
	return c.JSON(fiber.Map{"status": "success", "message": "Submitted"})
}

// --- Read (Common) ---

func (s *achievementService) mergeData(refs []models.AchievementReference) ([]models.AchievementReference, error) {
	if len(refs) == 0 { return refs, nil }
	var mongoIDs []string
	for _, ref := range refs { mongoIDs = append(mongoIDs, ref.MongoAchievementID) }
	mongoDetails, err := s.repo.FindMongoDetails(mongoIDs)
	if err != nil { return nil, err }
	for i := range refs {
		if detail, exists := mongoDetails[refs[i].MongoAchievementID]; exists { refs[i].Detail = &detail }
	}
	return refs, nil
}

// GetMyAchievements godoc
// @Summary      List Prestasi Saya
// @Description  Melihat daftar prestasi milik mahasiswa yang sedang login
// @Tags         Achievements (Mahasiswa)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.AchievementReference
// @Router       /achievements/my [get]
func (s *achievementService) GetMyAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	studentID, err := s.repo.GetStudentIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error"}) }
	refs, err := s.repo.FindAllByStudentID(studentID)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error"}) }
	finalData, err := s.mergeData(refs)
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}

// GetAdviseeAchievements godoc
// @Summary      List Prestasi Bimbingan
// @Description  Melihat daftar prestasi mahasiswa bimbingan (Dosen Wali Only)
// @Tags         Achievements (Dosen & Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.AchievementReference
// @Router       /achievements/advisees [get]
func (s *achievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID).String()
	advisorID, err := s.repo.GetAdvisorIDByUserID(userID)
	if err != nil { return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error"}) }
	refs, err := s.repo.FindAllByAdvisorID(advisorID)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error"}) }
	finalData, err := s.mergeData(refs)
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}

// [NEW] Admin View All

// GetAllAchievements godoc
// @Summary      List ALL Prestasi (Admin)
// @Description  Melihat seluruh prestasi yang masuk (Admin Only)
// @Tags         Achievements (Dosen & Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.AchievementReference
// @Router       /achievements [get]
func (s *achievementService) GetAllAchievements(c *fiber.Ctx) error {
	refs, err := s.repo.FindAllAchievements()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	finalData, err := s.mergeData(refs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal merge data"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}

// GetAchievementByID godoc
// @Summary      Detail Prestasi
// @Description  Melihat detail lengkap prestasi (termasuk data MongoDB)
// @Tags         Achievements (Common)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Success      200  {object}  models.AchievementReference
// @Router       /achievements/{id} [get]
func (s *achievementService) GetAchievementByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()
	ref, err := s.repo.FindRefByID(id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }

	isAllowed := false
	studentID, errMhs := s.repo.GetStudentIDByUserID(userID)
	if errMhs == nil && ref.StudentID == studentID { isAllowed = true }
	if !isAllowed {
		advisorID, errDos := s.repo.GetAdvisorIDByUserID(userID)
		if errDos == nil {
			isAdvisee, _ := s.repo.IsAdvisee(advisorID, ref.StudentID)
			if isAdvisee && (ref.Status != "draft" && ref.Status != "deleted") { isAllowed = true }
		}
	}
	// [NEW] Admin Allow Bypass (Opsional jika ingin Admin lihat detail juga)
	// if role == "Admin" { isAllowed = true }

	if !isAllowed { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }
	
	mongoDetails, _ := s.repo.FindMongoDetails([]string{ref.MongoAchievementID})
	if detail, exists := mongoDetails[ref.MongoAchievementID]; exists { ref.Detail = &detail }
	return c.JSON(fiber.Map{"status": "success", "data": ref})
}

// GetAchievementHistory godoc
// @Summary      History Prestasi
// @Description  Melihat jejak status (Draft -> Submitted -> Verified)
// @Tags         Achievements (Common)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Success      200  {array}   map[string]interface{}
// @Router       /achievements/{id}/history [get]
func (s *achievementService) GetAchievementHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uuid.UUID).String()
	ref, err := s.repo.FindRefByID(id)
	if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }

	isAllowed := false
	studentID, errMhs := s.repo.GetStudentIDByUserID(userID)
	if errMhs == nil && ref.StudentID == studentID { isAllowed = true }
	if !isAllowed {
		advisorID, errDos := s.repo.GetAdvisorIDByUserID(userID)
		if errDos == nil {
			isAdvisee, _ := s.repo.IsAdvisee(advisorID, ref.StudentID)
			if isAdvisee && (ref.Status != "draft" && ref.Status != "deleted") { isAllowed = true }
		}
	}
	if !isAllowed { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error"}) }

	// Construct History
	var history []fiber.Map
	history = append(history, fiber.Map{"status": "draft", "timestamp": ref.CreatedAt, "note": "Prestasi dibuat (Draft)", "actor": "Mahasiswa"})
	if ref.SubmittedAt != nil { history = append(history, fiber.Map{"status": "submitted", "timestamp": ref.SubmittedAt, "note": "Menunggu verifikasi", "actor": "Mahasiswa"}) }
	if ref.VerifiedAt != nil { 
		note := "Prestasi verified"
		if ref.Status == "rejected" { note = "Prestasi ditolak: " + ref.RejectionNote }
		history = append(history, fiber.Map{"status": ref.Status, "timestamp": ref.VerifiedAt, "note": note, "actor": "Dosen Wali"}) 
	}
	return c.JSON(fiber.Map{"status": "success", "data": history})
}

// VerifyAchievement godoc
// @Summary      Verify Prestasi
// @Description  Menyetujui prestasi (Status -> Verified)
// @Tags         Achievements (Dosen & Admin)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /achievements/{id}/verify [post]
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

// RejectAchievement godoc
// @Summary      Reject Prestasi
// @Description  Menolak prestasi dengan catatan (Status -> Rejected)
// @Tags         Achievements (Dosen & Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement Ref ID"
// @Param        body body      map[string]string true "Field: notes"
// @Success      200  {object}  map[string]interface{}
// @Router       /achievements/{id}/reject [post]
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

// GET /api/v1/students/:id/achievements
// @Summary      Get Student Achievements (Public/Admin)
// @Description  Melihat list prestasi berdasarkan ID Mahasiswa (Tabel Students)
// @Tags         Students & Lecturers
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Student Table ID"
// @Success      200  {array}   models.AchievementReference
// @Router       /students/{id}/achievements [get]
func (s *achievementService) GetAchievementsByStudentID(c *fiber.Ctx) error {
	studentID := c.Params("id")
	// Kita reuse fungsi FindAllByStudentID yang sudah ada di repo
	refs, err := s.repo.FindAllByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	finalData, _ := s.mergeData(refs)
	return c.JSON(fiber.Map{"status": "success", "data": finalData})
}