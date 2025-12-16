package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"database/sql"
	"errors"
)

type AuthRepository interface {
	FindByUsername(username string) (*models.User, error)
	// Kita siapkan fungsi CreateUser untuk nanti membuat Admin pertama kali
	CreateUser(user *models.User) error
	GetUserDetail(userID string) (*models.User, error) // [NEW]
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		db: config.DB,
	}
}

// FindByUsername mencari user berdasarkan username dan mengambil nama role-nya sekalian
func (r *authRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	
	// Query JOIN antara users dan roles untuk mendapatkan nama role
	// Sesuai SRS: users punya role_id, roles punya name
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 AND u.is_active = true
	`

	// Eksekusi QueryRow (Sesuai Aturan No. 11)
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser untuk menyimpan user baru (akan dipakai saat seeding admin nanti)
func (r *authRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(query, 
		user.Username, 
		user.Email, 
		user.PasswordHash, 
		user.FullName, 
		user.RoleID,
	).Scan(&user.ID)

	return err
}

// Tambahkan Implementasi
func (r *authRepository) GetUserDetail(userID string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, r.name as role_name, u.is_active, u.created_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	var u models.User
	err := r.db.QueryRow(query, userID).Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleName, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}