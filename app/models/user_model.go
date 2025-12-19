// package models

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// // --- ENTITIES (Tabel Database) ---

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

// // --- AUTH DTOs ---

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

// // --- FASE 2: RELATIONS & ADMIN DTOs (NEW) ---
// // Struct ini digunakan untuk menampilkan List Mahasiswa & Dosen di dashboard Admin

// // StudentResponse untuk menampilkan list mahasiswa beserta Dosen Walinya
// type StudentResponse struct {
// 	ID          string `json:"id"`           // ID tabel Students (UUID)
// 	FullName    string `json:"full_name"`    // Dari tabel Users
// 	NIM         string `json:"student_id"`   // Dari tabel Students
// 	Prodi       string `json:"program_study"`
// 	AdvisorName string `json:"advisor_name"` // Nama Dosen Wali
// }

// // LecturerResponse untuk menampilkan list dosen
// type LecturerResponse struct {
// 	ID         string `json:"id"`          // ID tabel Lecturers (UUID)
// 	FullName   string `json:"full_name"`   // Dari tabel Users
// 	NIP        string `json:"lecturer_id"` // Dari tabel Lecturers
// 	Department string `json:"department"`
// }

// // AssignAdvisorRequest untuk body JSON saat assign dosen
// type AssignAdvisorRequest struct {
// 	AdvisorID string `json:"advisor_id" validate:"required"` // ID tabel Lecturers
// }


package models

import (
	"time"

	"github.com/google/uuid"
)

// ==========================================
// 1. ENTITIES (Representasi Tabel Database)
// ==========================================

// User merepresentasikan tabel 'users' di PostgreSQL
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Password tidak dikirim balik
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Field Relasi (Optional)
	RoleName string `json:"role_name,omitempty"`
}

// Role merepresentasikan tabel 'roles'
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// StudentInfo merepresentasikan tabel 'students' (Detail Mahasiswa)
type StudentInfo struct {
	StudentID    string    `json:"student_id"` // NIM
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
	AdvisorID    uuid.UUID `json:"advisor_id"`
	AdvisorName  string    `json:"advisor_name,omitempty"` // Optional untuk display
}

// LecturerInfo merepresentasikan tabel 'lecturers' (Detail Dosen)
type LecturerInfo struct {
	LecturerID string `json:"lecturer_id"` // NIP
	Department string `json:"department"`
}

// ==========================================
// 2. AUTH DTOs (Data Transfer Objects)
// ==========================================

// LoginRequest adalah format JSON input login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse adalah format JSON output login sukses
type LoginResponse struct {
	Token        string `json:"token"`
	TokenRefresh string `json:"tokenRefresh"` // [NEW] Sesuai gambar referensi
	User         struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		FullName string    `json:"full_name"`
		Role     string    `json:"role"`
	} `json:"user"`
}

// ==========================================
// 3. REGISTRATION DTOs (Register)
// ==========================================

// CreateStudentRequest untuk input Register Mahasiswa
type CreateStudentRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	StudentID    string `json:"student_id" validate:"required"` // NIM
	ProgramStudy string `json:"program_study" validate:"required"`
	AcademicYear string `json:"academic_year" validate:"required"`
	AdvisorID    string `json:"advisor_id" validate:"required"` // UUID Dosen Wali
}

// CreateLecturerRequest untuk input Register Dosen
type CreateLecturerRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	FullName   string `json:"full_name" validate:"required"`
	LecturerID string `json:"lecturer_id" validate:"required"` // NIP
	Department string `json:"department" validate:"required"`
}

// ==========================================
// 4. ADMIN & RELATIONS DTOs (Dashboard)
// ==========================================

// StudentResponse untuk list mahasiswa + dosen wali
type StudentResponse struct {
	ID          string `json:"id"`            // ID Tabel Students
	FullName    string `json:"full_name"`     // Nama User
	NIM         string `json:"student_id"`    // NIM
	Prodi       string `json:"program_study"` // Prodi
	Email       string `json:"email"`         // Email User
	AdvisorName string `json:"advisor_name"`  // Nama Dosen Wali
}

// LecturerResponse untuk list dosen
type LecturerResponse struct {
	ID         string `json:"id"`          // ID Tabel Lecturers
	FullName   string `json:"full_name"`   // Nama User
	NIP        string `json:"lecturer_id"` // NIP
	Department string `json:"department"`
}

// AssignAdvisorRequest untuk body JSON saat assign dosen
type AssignAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required"`
}

// ==========================================
// 5. LEGACY DTOs (Safety Net)
// ==========================================
// Struct ini dipindahkan dari master_user_model.go agar tidak ada kode yang error
// jika ada bagian lain yang membutuhkannya.

type UserDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`

	// Data Tambahan (tampil jika ada)
	StudentInfo  *StudentInfo  `json:"student_info,omitempty"`
	LecturerInfo *LecturerInfo `json:"lecturer_info,omitempty"`
}