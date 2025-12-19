// package services

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/app/repositories"
// 	"PrestasiMhs-API/utils"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// )

// type AuthService interface {
// 	Login(c *fiber.Ctx) error
// 	GetProfile(c *fiber.Ctx) error // [NEW]
// }

// type authService struct {
// 	repo repositories.AuthRepository
// }

// func NewAuthService(repo repositories.AuthRepository) AuthService {
// 	return &authService{
// 		repo: repo,
// 	}
// }

// // Login godoc
// // @Summary      Login User
// // @Description  Masuk ke sistem untuk mendapatkan Token JWT
// // @Tags         Auth
// // @Accept       json
// // @Produce      json
// // @Param        request body models.LoginRequest true "Username & Password"
// // @Success      200  {object}  map[string]interface{}
// // @Failure      400  {object}  map[string]interface{}
// // @Failure      401  {object}  map[string]interface{}
// // @Router       /auth/login [post]
// func (s *authService) Login(c *fiber.Ctx) error {
// 	var req models.LoginRequest

// 	// 1. Parsing Data
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Format input tidak valid",
// 		})
// 	}

// 	// 2. Logic Layer: Cari User
// 	user, err := s.repo.FindByUsername(req.Username)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Username atau password salah",
// 		})
// 	}

// 	// 3. Logic Layer: Cek Password
// 	if !utils.CheckPassword(req.Password, user.PasswordHash) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Username atau password salah",
// 		})
// 	}

// 	// 4. Logic Layer: Generate Token
// 	token, err := utils.GenerateToken(user.ID, user.RoleName)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Gagal membuat token session",
// 		})
// 	}

// 	// 5. Susun Response menggunakan Struct
// 	response := models.LoginResponse{
// 		Token: token,
// 	}
// 	response.User.ID = user.ID
// 	response.User.Username = user.Username
// 	response.User.FullName = user.FullName
// 	response.User.Role = user.RoleName

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "Login berhasil",
// 		"data":    response,
// 	})
// }

// // GetProfile godoc
// // @Summary      Get User Profile
// // @Description  Melihat data diri user yang sedang login (Butuh Token)
// // @Tags         Auth
// // @Accept       json
// // @Produce      json
// // @Security     BearerAuth
// // @Success      200  {object}  models.User
// // @Failure      404  {object}  map[string]interface{}
// // @Router       /auth/profile [get]
// func (s *authService) GetProfile(c *fiber.Ctx) error {
// 	// KOREKSI: Ambil sebagai uuid.UUID dulu, baru convert ke String
// 	userID := c.Locals("user_id").(uuid.UUID).String()

// 	// Panggil Repository (pastikan Anda sudah update auth_repository.go di langkah sebelumnya)
// 	user, err := s.repo.GetUserDetail(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan"})
// 	}

// 	return c.JSON(fiber.Map{
// 		"status": "success",
// 		"data":   user,
// 	})
// }


package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

// Login godoc
// @Summary      Login User
// @Description  Masuk ke sistem untuk mendapatkan Token JWT dan Refresh Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Username & Password"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
func (s *authService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// 1. Parsing Data
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format input tidak valid",
		})
	}

	// 2. Logic Layer: Cari User
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Username atau password salah",
		})
	}

	// 3. Logic Layer: Cek Password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Username atau password salah",
		})
	}

	// 4. Logic Layer: Generate Access Token
	token, err := utils.GenerateToken(user.ID, user.RoleName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat token session",
		})
	}

	// [NEW] 5. Generate Refresh Token (Simulasi)
	// Kita membuat token kedua sebagai refresh token. 
	// Dalam implementasi nyata, ini bisa punya expiry time yang lebih lama.
	refreshToken, _ := utils.GenerateToken(user.ID, "refresh_token")

	// 6. Susun Response menggunakan Struct
	response := models.LoginResponse{
		Token:        token,
		TokenRefresh: refreshToken, // [NEW] Field ini sekarang terisi
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.FullName = user.FullName
	response.User.Role = user.RoleName

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login berhasil",
		"data":    response,
	})
}

// GetProfile godoc
// @Summary      Get User Profile
// @Description  Melihat data diri user yang sedang login (Butuh Token)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]interface{}
// @Router       /auth/profile [get]
func (s *authService) GetProfile(c *fiber.Ctx) error {
	// KOREKSI: Ambil sebagai uuid.UUID dulu, baru convert ke String
	userID := c.Locals("user_id").(uuid.UUID).String()

	// Panggil Repository
	user, err := s.repo.GetUserDetail(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}