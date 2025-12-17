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

// RegisterLecturer godoc
// @Summary      Register Dosen
// @Description  Mendaftarkan akun baru untuk Dosen (Admin Only)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateLecturerRequest true "Data Dosen"
// @Success      201  {object}  map[string]interface{}
// @Router       /users/lecturers [post]
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

// RegisterStudent godoc
// @Summary      Register Mahasiswa
// @Description  Mendaftarkan akun baru untuk Mahasiswa (Admin Only)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateStudentRequest true "Data Mahasiswa"
// @Success      201  {object}  map[string]interface{}
// @Router       /users/students [post]
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

// GetAllUsers godoc
// @Summary      List All Users
// @Description  Melihat semua user di sistem (Admin Only)
// @Tags         Users (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User
// @Router       /users [get]
func (s *userService) GetAllUsers(c *fiber.Ctx) error {
	users, err := s.repo.FindAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": users})
}

// [NEW] Get Detail User

// GetUserByID godoc
// @Summary      Get User Detail
// @Description  Melihat detail user berdasarkan ID (Admin Only)
// @Tags         Users (Admin)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  models.User
// @Router       /users/{id} [get]
func (s *userService) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.repo.FindUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": user})
}

// [NEW] Update User Info

// UpdateUser godoc
// @Summary      Update User Profile
// @Description  Mengubah data dasar user (Admin Only)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Param        body body      map[string]string true "Field: full_name, username, email"
// @Success      200  {object}  map[string]interface{}
// @Router       /users/{id} [put]
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

// ChangePassword godoc
// @Summary      Reset Password User
// @Description  Mengganti password user lain (Admin Only)
// @Tags         Users (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Param        body body      map[string]string true "Field: password"
// @Success      200  {object}  map[string]interface{}
// @Router       /users/{id}/role [put]
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

// DeleteUser godoc
// @Summary      Delete User
// @Description  Menghapus user dan profil terkait (Admin Only)
// @Tags         Users (Admin)
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /users/{id} [delete]
func (s *userService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.repo.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User berhasil dihapus"})
}

// GetAllStudents godoc
// @Summary      List Students
// @Description  Melihat daftar mahasiswa beserta dosen walinya
// @Tags         Relations (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.StudentResponse
// @Router       /students [get]
func (s *userService) GetAllStudents(c *fiber.Ctx) error {
	students, err := s.repo.GetAllStudents()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": students})
}

// GetAllLecturers godoc
// @Summary      List Lecturers
// @Description  Melihat daftar dosen
// @Tags         Relations (Admin)
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.LecturerResponse
// @Router       /lecturers [get]
func (s *userService) GetAllLecturers(c *fiber.Ctx) error {
	lecturers, err := s.repo.GetAllLecturers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": lecturers})
}

// AssignAdvisor godoc
// @Summary      Assign Dosen Wali
// @Description  Menghubungkan mahasiswa dengan dosen wali
// @Tags         Relations (Admin)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Student Table ID"
// @Param        body body      models.AssignAdvisorRequest true "Advisor ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /students/{id}/advisor [put]
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