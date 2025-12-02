package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Create Request DTOs (Input dari Postman) ---

// CreateLecturerRequest: Input untuk mendaftarkan Dosen
type CreateLecturerRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	FullName   string `json:"full_name" validate:"required"`
	// Data khusus tabel lecturers
	LecturerID string `json:"lecturer_id" validate:"required"` // NIP/NIDN
	Department string `json:"department" validate:"required"`
}

// CreateStudentRequest: Input untuk mendaftarkan Mahasiswa
type CreateStudentRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	FullName     string `json:"full_name" validate:"required"`
	// Data khusus tabel students
	StudentID    string `json:"student_id" validate:"required"` // NIM
	ProgramStudy string `json:"program_study" validate:"required"`
	AcademicYear string `json:"academic_year" validate:"required"`
	// AdvisorID adalah ID dari Dosen Wali (UUID dari tabel lecturers)
	AdvisorID    string `json:"advisor_id"` 
}

// --- Response DTOs (Output ke Client) ---

// UserDetailResponse: Format data lengkap saat melihat detail user
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

type StudentInfo struct {
	StudentID    string    `json:"student_id"`
	ProgramStudy string    `json:"program_study"`
	AcademicYear string    `json:"academic_year"`
	AdvisorName  string    `json:"advisor_name,omitempty"` // Nama Dosen Wali
}

type LecturerInfo struct {
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}