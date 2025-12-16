// package repositories

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/config"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// )

// type UserRepository interface {
// 	// Fase 1 & Basic Auth
// 	CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error
// 	CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error
	
// 	// Fase 2: User Management & Relations
// 	FindAllUsers() ([]models.User, error)
// 	DeleteUser(userID string) error
	
// 	GetAllStudents() ([]models.StudentResponse, error)
// 	GetAllLecturers() ([]models.LecturerResponse, error)
// 	AssignAdvisor(studentID string, advisorID string) error
// }

// type userRepository struct {
// 	db *sql.DB
// }

// func NewUserRepository() UserRepository {
// 	return &userRepository{
// 		db: config.DB,
// 	}
// }

// // --- Create Operations (Existing Code) ---

// func (r *userRepository) CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	var roleID string
// 	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'Dosen Wali'").Scan(&roleID)
// 	if err != nil {
// 		tx.Rollback()
// 		return errors.New("role 'Dosen Wali' tidak ditemukan di database")
// 	}

// 	queryUser := `
//         INSERT INTO users (username, email, password_hash, full_name, role_id)
//         VALUES ($1, $2, $3, $4, $5)
//         RETURNING id
//     `
// 	err = tx.QueryRow(queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, roleID).Scan(&user.ID)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("gagal insert user: %v", err)
// 	}

// 	queryLecturer := `
//         INSERT INTO lecturers (user_id, lecturer_id, department)
//         VALUES ($1, $2, $3)
//     `
// 	_, err = tx.Exec(queryLecturer, user.ID, lecturer.LecturerID, lecturer.Department)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("gagal insert lecturer profile: %v", err)
// 	}

// 	return tx.Commit()
// }

// func (r *userRepository) CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	var roleID string
// 	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'Mahasiswa'").Scan(&roleID)
// 	if err != nil {
// 		tx.Rollback()
// 		return errors.New("role 'Mahasiswa' tidak ditemukan")
// 	}

// 	queryUser := `
//         INSERT INTO users (username, email, password_hash, full_name, role_id)
//         VALUES ($1, $2, $3, $4, $5)
//         RETURNING id
//     `
// 	err = tx.QueryRow(queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, roleID).Scan(&user.ID)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("gagal insert user: %v", err)
// 	}

// 	var advisorUUID interface{} = nil
// 	if advisorID != "" {
// 		advisorUUID = advisorID
// 	}

// 	queryStudent := `
//         INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
//         VALUES ($1, $2, $3, $4, $5)
//     `
// 	_, err = tx.Exec(queryStudent, user.ID, student.StudentID, student.ProgramStudy, student.AcademicYear, advisorUUID)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("gagal insert student profile: %v", err)
// 	}

// 	return tx.Commit()
// }

// // --- FASE 2: NEW IMPLEMENTATION ---

// // FindAllUsers: Menampilkan semua user + nama role-nya
// func (r *userRepository) FindAllUsers() ([]models.User, error) {
// 	query := `
// 		SELECT u.id, u.username, u.email, u.full_name, r.name as role_name, u.is_active, u.created_at
// 		FROM users u
// 		JOIN roles r ON u.role_id = r.id
// 		ORDER BY u.created_at DESC
// 	`
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []models.User
// 	for rows.Next() {
// 		var u models.User
// 		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleName, &u.IsActive, &u.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, u)
// 	}
// 	return users, nil
// }

// // DeleteUser: Menghapus user dan profilnya (Manual Cascade)
// func (r *userRepository) DeleteUser(userID string) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// 1. Hapus profile di child tables dulu (students/lecturers)
// 	_, _ = tx.Exec("DELETE FROM students WHERE user_id = $1", userID)
// 	_, _ = tx.Exec("DELETE FROM lecturers WHERE user_id = $1", userID)
	
// 	// Catatan: Jika user adalah verifikator di achievement_references, 
// 	// set null atau biarkan (tergantung constraint on delete set null)
// 	// Untuk keamanan, kita asumsikan FK achievement_references diatur dengan baik.

// 	// 2. Hapus User
// 	query := "DELETE FROM users WHERE id = $1"
// 	result, err := tx.Exec(query, userID)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		tx.Rollback()
// 		return errors.New("user tidak ditemukan")
// 	}

// 	return tx.Commit()
// }

// // GetAllStudents: Menampilkan list mahasiswa lengkap dengan Nama Dosen Wali
// func (r *userRepository) GetAllStudents() ([]models.StudentResponse, error) {
// 	// Join Students -> Users (untuk nama Mhs)
// 	// Left Join Students -> Lecturers -> Users (untuk nama Dosen)
// 	query := `
// 		SELECT s.id, u.full_name, s.student_id, s.program_study, 
// 		       COALESCE(u_lec.full_name, 'Belum ditentukan') as advisor_name
// 		FROM students s
// 		JOIN users u ON s.user_id = u.id
// 		LEFT JOIN lecturers l ON s.advisor_id = l.id
// 		LEFT JOIN users u_lec ON l.user_id = u_lec.id
// 		ORDER BY s.student_id ASC
// 	`
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var students []models.StudentResponse
// 	for rows.Next() {
// 		var s models.StudentResponse
// 		if err := rows.Scan(&s.ID, &s.FullName, &s.NIM, &s.Prodi, &s.AdvisorName); err != nil {
// 			return nil, err
// 		}
// 		students = append(students, s)
// 	}
// 	return students, nil
// }

// // GetAllLecturers: Menampilkan list dosen
// func (r *userRepository) GetAllLecturers() ([]models.LecturerResponse, error) {
// 	query := `
// 		SELECT l.id, u.full_name, l.lecturer_id, l.department
// 		FROM lecturers l
// 		JOIN users u ON l.user_id = u.id
// 		ORDER BY l.lecturer_id ASC
// 	`
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var lecturers []models.LecturerResponse
// 	for rows.Next() {
// 		var l models.LecturerResponse
// 		if err := rows.Scan(&l.ID, &l.FullName, &l.NIP, &l.Department); err != nil {
// 			return nil, err
// 		}
// 		lecturers = append(lecturers, l)
// 	}
// 	return lecturers, nil
// }

// // AssignAdvisor: Menghubungkan Student dengan Lecturer (Rule 3 Enable)
// func (r *userRepository) AssignAdvisor(studentID string, advisorID string) error {
// 	// Cek apakah lecturer ID valid
// 	var exists bool
// 	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM lecturers WHERE id = $1)", advisorID).Scan(&exists)
// 	if err != nil || !exists {
// 		return errors.New("data dosen wali tidak ditemukan / ID salah")
// 	}

// 	query := `UPDATE students SET advisor_id = $1 WHERE id = $2`
// 	result, err := r.db.Exec(query, advisorID, studentID)
// 	if err != nil {
// 		return err
// 	}
	
// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		return errors.New("data mahasiswa tidak ditemukan")
// 	}
// 	return nil
// }


package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository interface {
	// Fase 1 & Basic Auth
	CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error
	CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error

	// Fase 2: User Management & Relations
	FindAllUsers() ([]models.User, error)
	FindUserByID(id string) (*models.User, error) // [NEW] Detail
	UpdateUser(user *models.User) error           // [NEW] Update Profil
	UpdatePassword(userID string, newHash string) error // [NEW] Ganti Password
	DeleteUser(userID string) error

	GetAllStudents() ([]models.StudentResponse, error)
	GetAllLecturers() ([]models.LecturerResponse, error)
	AssignAdvisor(studentID string, advisorID string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.DB,
	}
}

// --- Create Operations ---

func (r *userRepository) CreateLecturer(user *models.User, lecturer *models.LecturerInfo) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

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
	err = tx.QueryRow(queryUser, user.Username, user.Email, user.PasswordHash, user.FullName, roleID).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal insert user: %v", err)
	}

	queryLecturer := `
        INSERT INTO lecturers (user_id, lecturer_id, department)
        VALUES ($1, $2, $3)
    `
	_, err = tx.Exec(queryLecturer, user.ID, lecturer.LecturerID, lecturer.Department)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal insert lecturer profile: %v", err)
	}

	return tx.Commit()
}

func (r *userRepository) CreateStudent(user *models.User, student *models.StudentInfo, advisorID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var roleID string
	err = tx.QueryRow("SELECT id FROM roles WHERE name = 'Mahasiswa'").Scan(&roleID)
	if err != nil {
		tx.Rollback()
		return errors.New("role 'Mahasiswa' tidak ditemukan")
	}

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

// --- FASE 2 IMPLEMENTATION ---

func (r *userRepository) FindAllUsers() ([]models.User, error) {
	query := `
        SELECT u.id, u.username, u.email, u.full_name, r.name as role_name, u.is_active, u.created_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        ORDER BY u.created_at DESC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleName, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// [NEW] FindUserByID untuk Detail User
func (r *userRepository) FindUserByID(id string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, r.name as role_name, u.is_active, u.created_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	var u models.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleName, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &u, nil
}

// [NEW] UpdateUser untuk edit profil dasar
func (r *userRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET full_name = $1, username = $2, email = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(query, user.FullName, user.Username, user.Email, user.ID)
	return err
}

// [NEW] UpdatePassword untuk ganti password
func (r *userRepository) UpdatePassword(userID string, newHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, newHash, userID)
	return err
}

func (r *userRepository) DeleteUser(userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, _ = tx.Exec("DELETE FROM students WHERE user_id = $1", userID)
	_, _ = tx.Exec("DELETE FROM lecturers WHERE user_id = $1", userID)

	query := "DELETE FROM users WHERE id = $1"
	result, err := tx.Exec(query, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		tx.Rollback()
		return errors.New("user tidak ditemukan")
	}

	return tx.Commit()
}

func (r *userRepository) GetAllStudents() ([]models.StudentResponse, error) {
	query := `
        SELECT s.id, u.full_name, s.student_id, s.program_study, 
               COALESCE(u_lec.full_name, 'Belum ditentukan') as advisor_name
        FROM students s
        JOIN users u ON s.user_id = u.id
        LEFT JOIN lecturers l ON s.advisor_id = l.id
        LEFT JOIN users u_lec ON l.user_id = u_lec.id
        ORDER BY s.student_id ASC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.StudentResponse
	for rows.Next() {
		var s models.StudentResponse
		if err := rows.Scan(&s.ID, &s.FullName, &s.NIM, &s.Prodi, &s.AdvisorName); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *userRepository) GetAllLecturers() ([]models.LecturerResponse, error) {
	query := `
        SELECT l.id, u.full_name, l.lecturer_id, l.department
        FROM lecturers l
        JOIN users u ON l.user_id = u.id
        ORDER BY l.lecturer_id ASC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []models.LecturerResponse
	for rows.Next() {
		var l models.LecturerResponse
		if err := rows.Scan(&l.ID, &l.FullName, &l.NIP, &l.Department); err != nil {
			return nil, err
		}
		lecturers = append(lecturers, l)
	}
	return lecturers, nil
}

func (r *userRepository) AssignAdvisor(studentID string, advisorID string) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM lecturers WHERE id = $1)", advisorID).Scan(&exists)
	if err != nil || !exists {
		return errors.New("data dosen wali tidak ditemukan / ID salah")
	}

	query := `UPDATE students SET advisor_id = $1 WHERE id = $2`
	result, err := r.db.Exec(query, advisorID, studentID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("data mahasiswa tidak ditemukan")
	}
	return nil
}