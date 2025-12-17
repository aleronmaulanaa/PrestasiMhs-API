// package repositories

// import (
// 	"PrestasiMhs-API/config"
// 	"database/sql"
// )

// type ReportRepository interface {
// 	CountUsersByRole(roleName string) (int, error)
// 	CountAchievementsByStatus() (map[string]int, error)
// }

// type reportRepository struct {
// 	db *sql.DB
// }

// func NewReportRepository() ReportRepository {
// 	return &reportRepository{
// 		db: config.DB,
// 	}
// }

// // Menghitung jumlah user berdasarkan role (Mahasiswa/Dosen Wali/Admin)
// func (r *reportRepository) CountUsersByRole(roleName string) (int, error) {
// 	var count int
// 	query := `
// 		SELECT COUNT(u.id) 
// 		FROM users u 
// 		JOIN roles r ON u.role_id = r.id 
// 		WHERE r.name = $1 AND u.is_active = true
// 	`
// 	err := r.db.QueryRow(query, roleName).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// // Menghitung statistik status prestasi (Draft, Submitted, Verified, Rejected)
// func (r *reportRepository) CountAchievementsByStatus() (map[string]int, error) {
// 	query := `
// 		SELECT status, COUNT(*) 
// 		FROM achievement_references 
// 		WHERE status != 'deleted' 
// 		GROUP BY status
// 	`
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	stats := map[string]int{
// 		"draft":     0,
// 		"submitted": 0,
// 		"verified":  0,
// 		"rejected":  0,
// 	}

// 	for rows.Next() {
// 		var status string
// 		var count int
// 		if err := rows.Scan(&status, &count); err == nil {
// 			stats[status] = count
// 		}
// 	}
// 	return stats, nil
// }


package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"database/sql"
)

type ReportRepository interface {
	// Dashboard Stats
	CountUsersByRole(roleName string) (int, error)
	CountAchievementsByStatus() (map[string]int, error)

	// [NEW] Student Report
	GetStudentHeader(studentID string) (*models.StudentHeader, error)
	GetStudentAchievements(studentID string) ([]models.AchievementReference, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository() ReportRepository {
	return &reportRepository{
		db: config.DB,
	}
}

// Menghitung jumlah user berdasarkan role
func (r *reportRepository) CountUsersByRole(roleName string) (int, error) {
	var count int
	query := `
        SELECT COUNT(u.id) 
        FROM users u 
        JOIN roles r ON u.role_id = r.id 
        WHERE r.name = $1 AND u.is_active = true
    `
	err := r.db.QueryRow(query, roleName).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Menghitung statistik status prestasi
func (r *reportRepository) CountAchievementsByStatus() (map[string]int, error) {
	query := `
        SELECT status, COUNT(*) 
        FROM achievement_references 
        WHERE status != 'deleted' 
        GROUP BY status
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]int{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err == nil {
			stats[status] = count
		}
	}
	return stats, nil
}

// [NEW] Ambil Header Biodata Mahasiswa untuk Laporan
func (r *reportRepository) GetStudentHeader(studentID string) (*models.StudentHeader, error) {
	// Query Join: Students -> Users (untuk Nama Mhs) -> Lecturers -> Users (untuk Nama Dosen)
	query := `
		SELECT u.full_name, s.student_id, s.program_study, COALESCE(u_adv.full_name, 'Belum ditentukan')
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users u_adv ON l.user_id = u_adv.id
		WHERE s.id = $1
	`
	var h models.StudentHeader
	err := r.db.QueryRow(query, studentID).Scan(&h.FullName, &h.NIM, &h.ProgramStudy, &h.AdvisorName)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// [NEW] Ambil List Referensi Prestasi Mahasiswa (SQL)
func (r *reportRepository) GetStudentAchievements(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, mongo_achievement_id, status, verified_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		var verAt sql.NullTime
		if err := rows.Scan(&ref.ID, &ref.MongoAchievementID, &ref.Status, &verAt); err != nil {
			return nil, err
		}
		if verAt.Valid {
			ref.VerifiedAt = &verAt.Time
		}
		refs = append(refs, ref)
	}
	return refs, nil
}