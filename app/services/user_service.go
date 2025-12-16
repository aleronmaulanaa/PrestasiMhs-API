// package services

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/utils"
// 	"github.com/gofiber/fiber/v2"
// )

// type UserService interface {
// 	// Register Routes (Existing)
// 	RegisterLecturer(c *fiber.Ctx) error
// 	RegisterStudent(c *fiber.Ctx) error
	
// 	// Admin Management (Fase 2)
// 	GetAllUsers(c *fiber.Ctx) error
// 	DeleteUser(c *fiber.Ctx) error
	
// 	// Relations (Fase 2)
// 	GetAllStudents(c *fiber.Ctx) error
// 	GetAllLecturers(c *fiber.Ctx) error
// 	AssignAdvisor(c *fiber.Ctx) error
// }

// type userService struct {
// 	repo repositories.UserRepository
// }

// func NewUserService(repo repositories.UserRepository) UserService {
// 	return &userService{
// 		repo: repo,
// 	}
// }

// // --- Register Logic (Existing) ---

// func (s *userService) RegisterLecturer(c *fiber.Ctx) error {
// 	var req models.CreateLecturerRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
// 	}

// 	hashedPassword, err := utils.HashPassword(req.Password)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal hash password"})
// 	}

// 	userModel := &models.User{
// 		Username:     req.Username,
// 		Email:        req.Email,
// 		PasswordHash: hashedPassword,
// 		FullName:     req.FullName,
// 	}
// 	lecturerModel := &models.LecturerInfo{
// 		LecturerID: req.LecturerID,
// 		Department: req.Department,
// 	}

// 	if err := s.repo.CreateLecturer(userModel, lecturerModel); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Dosen berhasil didaftarkan"})
// }

// func (s *userService) RegisterStudent(c *fiber.Ctx) error {
// 	var req models.CreateStudentRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
// 	}

// 	hashedPassword, err := utils.HashPassword(req.Password)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal hash password"})
// 	}

// 	userModel := &models.User{
// 		Username:     req.Username,
// 		Email:        req.Email,
// 		PasswordHash: hashedPassword,
// 		FullName:     req.FullName,
// 	}
// 	studentModel := &models.StudentInfo{
// 		StudentID:    req.StudentID,
// 		ProgramStudy: req.ProgramStudy,
// 		AcademicYear: req.AcademicYear,
// 	}

// 	if err := s.repo.CreateStudent(userModel, studentModel, req.AdvisorID); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Mahasiswa berhasil didaftarkan"})
// }

// // --- FASE 2: NEW SERVICE LOGIC ---

// func (s *userService) GetAllUsers(c *fiber.Ctx) error {
// 	users, err := s.repo.FindAllUsers()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "data": users})
// }

// func (s *userService) DeleteUser(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if err := s.repo.DeleteUser(id); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
// }

// func (s *userService) GetAllStudents(c *fiber.Ctx) error {
// 	students, err := s.repo.GetAllStudents()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "data": students})
// }

// func (s *userService) GetAllLecturers(c *fiber.Ctx) error {
// 	lecturers, err := s.repo.GetAllLecturers()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}
// 	return c.JSON(fiber.Map{"status": "success", "data": lecturers})
// }

// func (s *userService) AssignAdvisor(c *fiber.Ctx) error {
// 	studentID := c.Params("id") // ID dari tabel students (bukan user_id)
	
// 	var req models.AssignAdvisorRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
// 	}

// 	if err := s.repo.AssignAdvisor(studentID, req.AdvisorID); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"status": "success", "message": "Dosen Wali berhasil ditugaskan"})
// }


package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserService interface {
	RegisterLecturer(c *fiber.Ctx) error
	RegisterStudent(c *fiber.Ctx) error
	
	// Admin Management (Fase 2 Completed)
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error    // [NEW]
	UpdateUser(c *fiber.Ctx) error     // [NEW]
	ChangePassword(c *fiber.Ctx) error // [NEW]
	DeleteUser(c *fiber.Ctx) error
	
	// Relations
	GetAllStudents(c *fiber.Ctx) error
	GetAllLecturers(c *fiber.Ctx) error
	AssignAdvisor(c *fiber.Ctx) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// --- Register Logic (Existing) ---

func (s *userService) RegisterLecturer(c *fiber.Ctx) error {
	var req models.CreateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal hash password"})
	}

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

	if err := s.repo.CreateLecturer(userModel, lecturerModel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Dosen berhasil didaftarkan"})
}

func (s *userService) RegisterStudent(c *fiber.Ctx) error {
	var req models.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input tidak valid"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal hash password"})
	}

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

	if err := s.repo.CreateStudent(userModel, studentModel, req.AdvisorID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Mahasiswa berhasil didaftarkan"})
}

// --- FASE 2: ADMIN MANAGEMENT ---

func (s *userService) GetAllUsers(c *fiber.Ctx) error {
	users, err := s.repo.FindAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": users})
}

// [NEW] Get Detail User
func (s *userService) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.repo.FindUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": user})
}

// [NEW] Update User Info
func (s *userService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	type UpdateReq struct {
		FullName string `json:"full_name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	var req UpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input"})
	}

	uid, _ := uuid.Parse(id)
	user := &models.User{
		ID:       uid,
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User updated"})
}

// [NEW] Change Password
func (s *userService) ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	type PwdReq struct {
		Password string `json:"password"`
	}
	var req PwdReq
	c.BodyParser(&req)
	
	newHash, _ := utils.HashPassword(req.Password)
	if err := s.repo.UpdatePassword(id, newHash); err != nil {
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Password updated"})
}

func (s *userService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.repo.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}

func (s *userService) GetAllStudents(c *fiber.Ctx) error {
	students, err := s.repo.GetAllStudents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": students})
}

func (s *userService) GetAllLecturers(c *fiber.Ctx) error {
	lecturers, err := s.repo.GetAllLecturers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": lecturers})
}

func (s *userService) AssignAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")
	var req models.AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
	}

	if err := s.repo.AssignAdvisor(studentID, req.AdvisorID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Dosen Wali berhasil ditugaskan"})
}