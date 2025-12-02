package repositories

import (
	"PrestasiMhs-API/app/models"
	"PrestasiMhs-API/config"
	"context"
	"database/sql"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	// Create menyimpan ke Mongo dan Postgres
	Create(mongoData *models.AchievementMongo, studentUUID string) error
	// Helper untuk mencari Student ID dari User ID yang login
	GetStudentIDByUserID(userID string) (string, error)
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

func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, error) {
	var studentID string
	query := "SELECT id FROM students WHERE user_id = $1"
	err := r.pg.QueryRow(query, userID).Scan(&studentID)
	if err != nil {
		return "", errors.New("data mahasiswa tidak ditemukan")
	}
	return studentID, nil
}

func (r *achievementRepository) Create(mongoData *models.AchievementMongo, studentUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Insert ke MongoDB
	collection := r.mongo.Collection("achievements")
	result, err := collection.InsertOne(ctx, mongoData)
	if err != nil {
		return err
	}

	// Ambil ObjectID yang baru terbuat
	mongoID := result.InsertedID.(primitive.ObjectID).Hex()

	// 2. Insert ke PostgreSQL (Tabel Referensi)
	// Status default 'draft' sudah diatur di database (DEFAULT 'draft'), tapi kita set eksplisit biar jelas
	query := `
		INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at)
		VALUES ($1, $2, 'draft', NOW())
	`
	_, err = r.pg.Exec(query, studentUUID, mongoID)
	
	if err != nil {
		// ROLLBACK MANUAL: Jika simpan ke Postgres gagal, hapus data sampah di Mongo
		_, _ = collection.DeleteOne(ctx, result.InsertedID)
		return errors.New("gagal menyimpan referensi prestasi")
	}

	return nil
}