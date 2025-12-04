// package repositories

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/config"
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type AchievementRepository interface {
// 	// Create menyimpan ke Mongo dan Postgres
// 	Create(mongoData *models.AchievementMongo, studentUUID string) error
// 	// Helper untuk mencari Student ID dari User ID yang login
// 	GetStudentIDByUserID(userID string) (string, error)
// }

// type achievementRepository struct {
// 	pg    *sql.DB
// 	mongo *mongo.Database
// }

// func NewAchievementRepository() AchievementRepository {
// 	return &achievementRepository{
// 		pg:    config.DB,
// 		mongo: config.MongoDB,
// 	}
// }

// func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, error) {
// 	var studentID string
// 	query := "SELECT id FROM students WHERE user_id = $1"
// 	err := r.pg.QueryRow(query, userID).Scan(&studentID)
// 	if err != nil {
// 		return "", errors.New("data mahasiswa tidak ditemukan")
// 	}
// 	return studentID, nil
// }

// func (r *achievementRepository) Create(mongoData *models.AchievementMongo, studentUUID string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// 1. Insert ke MongoDB
// 	collection := r.mongo.Collection("achievements")
// 	result, err := collection.InsertOne(ctx, mongoData)
// 	if err != nil {
// 		return err
// 	}

// 	// Ambil ObjectID yang baru terbuat
// 	mongoID := result.InsertedID.(primitive.ObjectID).Hex()

// 	// 2. Insert ke PostgreSQL (Tabel Referensi)
// 	// Status default 'draft' sudah diatur di database (DEFAULT 'draft'), tapi kita set eksplisit biar jelas
// 	query := `
// 		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at)
// 		VALUES ($1, $2, 'draft', NOW())
// 	`
// 	_, err = r.pg.Exec(query, studentUUID, mongoID)
	
// 	if err != nil {
// 		// ROLLBACK MANUAL: Jika simpan ke Postgres gagal, hapus data sampah di Mongo
// 		_, _ = collection.DeleteOne(ctx, result.InsertedID)
// 		return errors.New("gagal menyimpan referensi prestasi")
// 	}

// 	return nil
// }


package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"context"
	"database/sql"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	// --- Commands (Write) ---
	Create(mongoData *models.AchievementMongo, studentUUID string) error
	UpdateStatus(id string, status string, notes string, verifierID string) error
	
	// --- Queries (Read) ---
	GetStudentIDByUserID(userID string) (string, error)
	GetAdvisorIDByUserID(userID string) (string, error)
	
	// Mencari referensi prestasi (Postgres)
	FindAllByStudentID(studentID string) ([]models.AchievementReference, error)
	FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error)
	FindRefByID(id string) (*models.AchievementReference, error)
	
	// Mengambil detail prestasi (Mongo)
	FindMongoDetails(mongoIDs []string) (map[string]models.AchievementMongo, error)
}

type achievementRepository struct {
	pg    *sql.DB
	mongo *mongo.Database
}

func NewAchievementRepository() AchievementRepository {
	return &achievementRepository{
		pg:    config.DB,
		mongo: config.MongoDB,
	}
}

// --- Helper User ID ---

func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, error) {
	var studentID string
	// Menggunakan QueryRow (Aturan No. 11)
	err := r.pg.QueryRow("SELECT id FROM students WHERE user_id = $1", userID).Scan(&studentID)
	if err != nil {
		return "", errors.New("data mahasiswa tidak ditemukan")
	}
	return studentID, nil
}

func (r *achievementRepository) GetAdvisorIDByUserID(userID string) (string, error) {
	var advisorID string
	err := r.pg.QueryRow("SELECT id FROM lecturers WHERE user_id = $1", userID).Scan(&advisorID)
	if err != nil {
		return "", errors.New("data dosen tidak ditemukan")
	}
	return advisorID, nil
}

// --- Write Operations ---

func (r *achievementRepository) Create(mongoData *models.AchievementMongo, studentUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.mongo.Collection("achievements")
	result, err := collection.InsertOne(ctx, mongoData)
	if err != nil {
		return err
	}

	mongoID := result.InsertedID.(primitive.ObjectID).Hex()

	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at)
		VALUES ($1, $2, 'draft', NOW())
	`
	_, err = r.pg.Exec(query, studentUUID, mongoID)
	
	if err != nil {
		_, _ = collection.DeleteOne(ctx, result.InsertedID) // Rollback Manual
		return errors.New("gagal menyimpan referensi prestasi")
	}

	return nil
}

func (r *achievementRepository) UpdateStatus(id string, status string, notes string, verifierID string) error {
	// Query update status dan waktu verifikasi
	// Menggunakan Exec (Aturan No. 11)
	query := `
		UPDATE achievement_references 
		SET status = $1, rejection_note = $2, verified_by = $3, verified_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`
	// Jika rejection_note kosong, kita kirim NULL (sql.NullString) atau string kosong tergantung setup,
	// disini kita kirim string biasa, jika kosong biarkan kosong.
	
	result, err := r.pg.Exec(query, status, notes, verifierID, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("prestasi tidak ditemukan")
	}
	return nil
}

// --- Read Operations (PostgreSQL) ---

// FindAllByStudentID: Untuk Mahasiswa melihat miliknya sendiri
func (r *achievementRepository) FindAllByStudentID(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
	`
	rows, err := r.pg.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		// Handle nullable fields appropriately if using sql.Null* types, 
		// but for simplicity assuming string/time is fine or handled by driver
		var note sql.NullString
		var verAt sql.NullTime

		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt, &verAt); err != nil {
			return nil, err
		}
		ref.RejectionNote = note.String
		if verAt.Valid {
			ref.VerifiedAt = &verAt.Time
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

// FindAllByAdvisorID: Untuk Dosen Wali melihat mahasiswa bimbingannya
func (r *achievementRepository) FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error) {
	// JOIN tabel achievement_references dengan students
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.rejection_note, ar.created_at
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`
	rows, err := r.pg.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		var note sql.NullString
		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt); err != nil {
			return nil, err
		}
		ref.RejectionNote = note.String
		refs = append(refs, ref)
	}
	return refs, nil
}

// FindRefByID: Mengambil satu data referensi (untuk detail/verifikasi)
func (r *achievementRepository) FindRefByID(id string) (*models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at, submitted_at
		FROM achievement_references
		WHERE id = $1
	`
	var ref models.AchievementReference
	var note sql.NullString
	var verAt, subAt sql.NullTime

	err := r.pg.QueryRow(query, id).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, 
		&note, &ref.CreatedAt, &verAt, &subAt,
	)
	if err != nil {
		return nil, err
	}
	
	ref.RejectionNote = note.String
	if verAt.Valid { ref.VerifiedAt = &verAt.Time }
	if subAt.Valid { ref.SubmittedAt = &subAt.Time }
	
	return &ref, nil
}

// --- Read Operations (MongoDB) ---

// FindMongoDetails mengambil data detail dari Mongo berdasarkan List ID
func (r *achievementRepository) FindMongoDetails(mongoIDs []string) (map[string]models.AchievementMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range mongoIDs {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, oid)
		}
	}

	// Query Mongo dengan operator $in
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := r.mongo.Collection("achievements").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Masukkan hasil ke Map agar mudah dicocokkan nanti
	results := make(map[string]models.AchievementMongo)
	for cursor.Next(ctx) {
		var doc models.AchievementMongo
		if err := cursor.Decode(&doc); err == nil {
			results[doc.ID.Hex()] = doc
		}
	}

	return results, nil
}