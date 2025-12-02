package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository interface {
	CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error
	CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error
	// Kita akan tambah GetAllUsers nanti sesuai kebutuhan
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.DB,
	}
}

// CreateLecturer memasukkan data ke tabel 'users' DAN 'lecturers' dalam satu transaksi
func (r *userRepository) CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error {
	// 1. Mulai Transaksi (Wajib untuk integritas data)
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 2. Insert ke tabel USERS
	// Kita ambil Role ID untuk 'Dosen Wali' dulu
	var roleID string
	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'Dosen Wali'").Scan(&roleID)
	if err != nil {
		tx.Rollback()
		return errors.New("role 'Dosen Wali' tidak ditemukan di database")
	}

	queryUser := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	// Perhatikan: Kita pakai tx.QueryRow (bukan r.db.QueryRow) karena dalam transaksi
	err = tx.QueryRow(queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, roleID).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal insert user: %v", err)
	}

	// 3. Insert ke tabel LECTURERS
	queryLecturer := `
		INSERT INTO lecturers (user_id, lecturer_id, department)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(queryLecturer, user.ID, lecturer.LecturerID, lecturer.Department)
	if err != nil {
		tx.Rollback() // Batalkan insert user tadi jika ini gagal
		return fmt.Errorf("gagal insert lecturer profile: %v", err)
	}

	// 4. Commit (Simpan Permanen)
	return tx.Commit()
}

// CreateStudent memasukkan data ke tabel 'users' DAN 'students'
func (r *userRepository) CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Ambil Role ID 'Mahasiswa'
	var roleID string
	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'Mahasiswa'").Scan(&roleID)
	if err != nil {
		tx.Rollback()
		return errors.New("role 'Mahasiswa' tidak ditemukan")
	}

	// Insert User
	queryUser := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tx.QueryRow(queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, roleID).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal insert user: %v", err)
	}

	// Insert Student
	// advisorID bisa null/kosong string, kita perlu handle agar masuk NULL ke DB jika kosong
	var advisorUUID interface{} = nil
	if advisorID != "" {
		advisorUUID = advisorID
	}

	queryStudent := `
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(queryStudent, user.ID, student.StudentID, student.ProgramStudy, student.AcademicYear, advisorUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal insert student profile: %v", err)
	}

	return tx.Commit()
}