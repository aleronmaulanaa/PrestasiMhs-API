// package models

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // User merepresentasikan tabel 'users' di PostgreSQL
// type User struct {
// 	ID           uuid.UUID `json:"id"`
// 	Username     string    `json:"username"`
// 	Email        string    `json:"email"`
// 	PasswordHash string    `json:"-"` // Password tidak boleh dikirim balik di JSON response
// 	FullName     string    `json:"full_name"`
// 	RoleID       uuid.UUID `json:"role_id"`
// 	IsActive     bool      `json:"is_active"`
// 	CreatedAt    time.Time `json:"created_at"`
// 	UpdatedAt    time.Time `json:"updated_at"`
	
// 	// Relasi (Optional, diisi jika join)
// 	RoleName     string    `json:"role_name,omitempty"`
// }

// // Role merepresentasikan tabel 'roles'
// type Role struct {
// 	ID          uuid.UUID `json:"id"`
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// }

// // LoginRequest adalah format JSON yang dikirim user saat login
// type LoginRequest struct {
// 	Username string `json:"username" validate:"required"`
// 	Password string `json:"password" validate:"required"`
// }

// // LoginResponse adalah apa yang kita kirim balik jika login sukses
// type LoginResponse struct {
// 	Token string `json:"token"`
// 	User  struct {
// 		ID       uuid.UUID `json:"id"`
// 		Username string    `json:"username"`
// 		FullName string    `json:"full_name"`
// 		Role     string    `json:"role"`
// 	} `json:"user"`
// }


package models

import (
	"time"

	"github.com/google/uuid"
)

// --- ENTITIES (Tabel Database) ---

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

// --- AUTH DTOs ---

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

// --- FASE 2: RELATIONS & ADMIN DTOs (NEW) ---
// Struct ini digunakan untuk menampilkan List Mahasiswa & Dosen di dashboard Admin

// StudentResponse untuk menampilkan list mahasiswa beserta Dosen Walinya
type StudentResponse struct {
	ID          string `json:"id"`           // ID tabel Students (UUID)
	FullName    string `json:"full_name"`    // Dari tabel Users
	NIM         string `json:"student_id"`   // Dari tabel Students
	Prodi       string `json:"program_study"`
	AdvisorName string `json:"advisor_name"` // Nama Dosen Wali
}

// LecturerResponse untuk menampilkan list dosen
type LecturerResponse struct {
	ID         string `json:"id"`          // ID tabel Lecturers (UUID)
	FullName   string `json:"full_name"`   // Dari tabel Users
	NIP        string `json:"lecturer_id"` // Dari tabel Lecturers
	Department string `json:"department"`
}

// AssignAdvisorRequest untuk body JSON saat assign dosen
type AssignAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required"` // ID tabel Lecturers
}