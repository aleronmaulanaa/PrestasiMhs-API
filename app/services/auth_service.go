package services

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/app/repositories"
	"PrestasiMhs-API/utils"
	"errors"
)

type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Cari user berdasarkan username di database
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah") // Pesan error umum agar aman
	}

	// 2. Cek apakah password cocok dengan hash di database
	// Menggunakan utils yang sudah kita buat di Step 1
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("username atau password salah")
	}

	// 3. Generate JWT Token
	// Menggunakan utils yang sudah kita buat di Step 1
	token, err := utils.GenerateToken(user.ID, user.RoleName)
	if err != nil {
		return nil, errors.New("gagal membuat token session")
	}

	// 4. Susun response
	response := &models.LoginResponse{
		Token: token,
	}
	// Isi data user di response (tanpa password)
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.FullName = user.FullName
	response.User.Role = user.RoleName

	return response, nil
}