package repositories

import (
	"PrestasiMhs-API/config"
	"database/sql"
)

type ReportRepository interface {
	CountUsersByRole(roleName string) (int, error)
	CountAchievementsByStatus() (map[string]int, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository() ReportRepository {
	return &reportRepository{
		db: config.DB,
	}
}

// Menghitung jumlah user berdasarkan role (Mahasiswa/Dosen Wali/Admin)
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

// Menghitung statistik status prestasi (Draft, Submitted, Verified, Rejected)
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