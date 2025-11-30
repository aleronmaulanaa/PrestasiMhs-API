package models

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan tabel 'users' di PostgreSQL
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Password tidak boleh dikirim balik di JSON response
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relasi (Optional, diisi jika join)
	RoleName     string    `json:"role_name,omitempty"`
}

// Role merepresentasikan tabel 'roles'
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// LoginRequest adalah format JSON yang dikirim user saat login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse adalah apa yang kita kirim balik jika login sukses
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		FullName string    `json:"full_name"`
		Role     string    `json:"role"`
	} `json:"user"`
}