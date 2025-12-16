// package repositories

// import (
// 	"PrestasiMhs-API/app/models"
// 	"PrestasiMhs-API/config"
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type AchievementRepository interface {
// 	// --- Commands (Write) ---
// 	Create(mongoData *models.AchievementMongo, studentUUID string) error
// 	UpdateStatus(id string, status string, notes string, verifierID string) error
// 	Submit(id string) error

// 	// [Fase 1: Update & Delete]
// 	FindMongoIDByRefID(refID string) (string, error)
// 	UpdateMongo(mongoID string, updateData *models.AchievementMongo) error
// 	SoftDelete(id string) error

// 	// --- Queries (Read) ---
// 	GetStudentIDByUserID(userID string) (string, error)
// 	GetAdvisorIDByUserID(userID string) (string, error)
	
// 	// [NEW] Cek apakah mahasiswa adalah bimbingan dosen tertentu (Rule 3)
// 	IsAdvisee(advisorID string, studentID string) (bool, error)

// 	// Mencari referensi prestasi (Postgres)
// 	FindAllByStudentID(studentID string) ([]models.AchievementReference, error)
// 	FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error)
// 	FindRefByID(id string) (*models.AchievementReference, error)

// 	// Mengambil detail prestasi (Mongo)
// 	FindMongoDetails(mongoIDs []string) (map[string]models.AchievementMongo, error)
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

// // --- Helper User ID ---

// func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, error) {
// 	var studentID string
// 	err := r.pg.QueryRow("SELECT id FROM students WHERE user_id = $1", userID).Scan(&studentID)
// 	if err != nil {
// 		return "", errors.New("data mahasiswa tidak ditemukan")
// 	}
// 	return studentID, nil
// }

// func (r *achievementRepository) GetAdvisorIDByUserID(userID string) (string, error) {
// 	var advisorID string
// 	err := r.pg.QueryRow("SELECT id FROM lecturers WHERE user_id = $1", userID).Scan(&advisorID)
// 	if err != nil {
// 		return "", errors.New("data dosen tidak ditemukan")
// 	}
// 	return advisorID, nil
// }

// // [NEW] Implementasi IsAdvisee
// func (r *achievementRepository) IsAdvisee(advisorID string, studentID string) (bool, error) {
// 	var exists bool
// 	// Menggunakan SELECT EXISTS agar efisien
// 	query := "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1 AND advisor_id = $2)"
// 	err := r.pg.QueryRow(query, studentID, advisorID).Scan(&exists)
// 	if err != nil {
// 		return false, err
// 	}
// 	return exists, nil
// }

// // --- Write Operations ---

// func (r *achievementRepository) Create(mongoData *models.AchievementMongo, studentUUID string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	collection := r.mongo.Collection("achievements")
// 	result, err := collection.InsertOne(ctx, mongoData)
// 	if err != nil {
// 		return err
// 	}

// 	mongoID := result.InsertedID.(primitive.ObjectID).Hex()

// 	query := `
//         INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at)
//         VALUES ($1, $2, 'draft', NOW())
//     `
// 	_, err = r.pg.Exec(query, studentUUID, mongoID)

// 	if err != nil {
// 		_, _ = collection.DeleteOne(ctx, result.InsertedID) // Rollback Manual
// 		return errors.New("gagal menyimpan referensi prestasi")
// 	}

// 	return nil
// }

// func (r *achievementRepository) UpdateStatus(id string, status string, notes string, verifierID string) error {
// 	// [PENTING] Hanya update jika status 'submitted' (untuk verifikasi) atau flow valid lainnya
// 	// Reject/Verify hanya boleh dilakukan jika status sebelumnya submitted
// 	query := `
//         UPDATE achievement_references 
//         SET status = $1, rejection_note = $2, verified_by = $3, verified_at = NOW(), updated_at = NOW()
//         WHERE id = $4 AND status = 'submitted'
//     `
// 	result, err := r.pg.Exec(query, status, notes, verifierID, id)
// 	if err != nil {
// 		return err
// 	}

// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		return errors.New("prestasi tidak ditemukan atau belum disubmit (status harus submitted)")
// 	}
// 	return nil
// }

// func (r *achievementRepository) Submit(id string) error {
// 	query := `
//         UPDATE achievement_references 
//         SET status = 'submitted', submitted_at = NOW(), updated_at = NOW()
//         WHERE id = $1 AND status = 'draft'
//     `
// 	result, err := r.pg.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		return errors.New("prestasi tidak ditemukan atau status bukan draft")
// 	}
// 	return nil
// }

// // --- [Fase 1 Implementation] ---

// func (r *achievementRepository) FindMongoIDByRefID(refID string) (string, error) {
// 	var mongoID string
// 	query := "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1"
// 	err := r.pg.QueryRow(query, refID).Scan(&mongoID)
// 	if err != nil {
// 		return "", errors.New("referensi prestasi tidak ditemukan")
// 	}
// 	return mongoID, nil
// }

// func (r *achievementRepository) UpdateMongo(mongoID string, updateData *models.AchievementMongo) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	oid, err := primitive.ObjectIDFromHex(mongoID)
// 	if err != nil {
// 		return errors.New("format ID mongo tidak valid")
// 	}

// 	filter := bson.M{"_id": oid}
// 	update := bson.M{
// 		"$set": bson.M{
// 			"achievementType": updateData.AchievementType,
// 			"title":           updateData.Title,
// 			"description":     updateData.Description,
// 			"details":         updateData.Details,
// 			"updatedAt":       time.Now(),
// 		},
// 	}

// 	if len(updateData.Attachments) > 0 {
// 		update["$set"].(bson.M)["attachments"] = updateData.Attachments
// 	}

// 	_, err = r.mongo.Collection("achievements").UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *achievementRepository) SoftDelete(id string) error {
// 	query := `
//         UPDATE achievement_references 
//         SET status = 'deleted', updated_at = NOW() 
//         WHERE id = $1 AND status = 'draft'
//     `
// 	result, err := r.pg.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rows, _ := result.RowsAffected()
// 	if rows == 0 {
// 		return errors.New("prestasi tidak ditemukan atau status bukan draft")
// 	}
// 	return nil
// }

// // --- Read Operations ---

// func (r *achievementRepository) FindAllByStudentID(studentID string) ([]models.AchievementReference, error) {
// 	query := `
//         SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at
//         FROM achievement_references
//         WHERE student_id = $1 AND status != 'deleted'
//         ORDER BY created_at DESC
//     `
// 	rows, err := r.pg.Query(query, studentID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var refs []models.AchievementReference
// 	for rows.Next() {
// 		var ref models.AchievementReference
// 		var note sql.NullString
// 		var verAt sql.NullTime

// 		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt, &verAt); err != nil {
// 			return nil, err
// 		}
// 		ref.RejectionNote = note.String
// 		if verAt.Valid {
// 			ref.VerifiedAt = &verAt.Time
// 		}
// 		refs = append(refs, ref)
// 	}
// 	return refs, nil
// }

// func (r *achievementRepository) FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error) {
// 	query := `
//         SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.rejection_note, ar.created_at
//         FROM achievement_references ar
//         JOIN students s ON ar.student_id = s.id
//         WHERE s.advisor_id = $1 AND ar.status != 'deleted' AND ar.status != 'draft'
//         ORDER BY ar.created_at DESC
//     `
// 	rows, err := r.pg.Query(query, advisorID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var refs []models.AchievementReference
// 	for rows.Next() {
// 		var ref models.AchievementReference
// 		var note sql.NullString
// 		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		ref.RejectionNote = note.String
// 		refs = append(refs, ref)
// 	}
// 	return refs, nil
// }

// func (r *achievementRepository) FindRefByID(id string) (*models.AchievementReference, error) {
// 	query := `
//         SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at, submitted_at
//         FROM achievement_references
//         WHERE id = $1
//     `
// 	var ref models.AchievementReference
// 	var note sql.NullString
// 	var verAt, subAt sql.NullTime

// 	err := r.pg.QueryRow(query, id).Scan(
// 		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
// 		&note, &ref.CreatedAt, &verAt, &subAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ref.RejectionNote = note.String
// 	if verAt.Valid {
// 		ref.VerifiedAt = &verAt.Time
// 	}
// 	if subAt.Valid {
// 		ref.SubmittedAt = &subAt.Time
// 	}

// 	return &ref, nil
// }

// func (r *achievementRepository) FindMongoDetails(mongoIDs []string) (map[string]models.AchievementMongo, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var objectIDs []primitive.ObjectID
// 	for _, id := range mongoIDs {
// 		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
// 			objectIDs = append(objectIDs, oid)
// 		}
// 	}

// 	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
// 	cursor, err := r.mongo.Collection("achievements").Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	results := make(map[string]models.AchievementMongo)
// 	for cursor.Next(ctx) {
// 		var doc models.AchievementMongo
// 		if err := cursor.Decode(&doc); err == nil {
// 			results[doc.ID.Hex()] = doc
// 		}
// 	}

// 	return results, nil
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
	Submit(id string) error

	// [Fase 1: Update & Delete]
	FindMongoIDByRefID(refID string) (string, error)
	UpdateMongo(mongoID string, updateData *models.AchievementMongo) error
	SoftDelete(id string) error

	// --- Queries (Read) ---
	GetStudentIDByUserID(userID string) (string, error)
	GetAdvisorIDByUserID(userID string) (string, error)
	IsAdvisee(advisorID string, studentID string) (bool, error)

	FindAllByStudentID(studentID string) ([]models.AchievementReference, error)
	FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error)
	
	// [NEW] Admin View All (FR-010)
	FindAllAchievements() ([]models.AchievementReference, error)
	
	FindRefByID(id string) (*models.AchievementReference, error)
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

// ... [KODE HELPER & WRITE OPS SAMA SEPERTI SEBELUMNYA] ...
// (Untuk menghemat tempat, bagian Helper, Create, UpdateStatus, Submit, SoftDelete, UpdateMongo
//  tidak saya tulis ulang karena TIDAK ADA PERUBAHAN. Pastikan bagian itu tetap ada!)
// ---------------------------------------------------------
// Copy Paste bagian Helper & Write dari kode sebelumnya disini
// ---------------------------------------------------------

// --- Helper User ID ---
func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, error) {
	var studentID string
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

func (r *achievementRepository) IsAdvisee(advisorID string, studentID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1 AND advisor_id = $2)"
	err := r.pg.QueryRow(query, studentID, advisorID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
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
	query := `INSERT INTO achievement_references (student_id, mongo_achievement_id, status, created_at) VALUES ($1, $2, 'draft', NOW())`
	_, err = r.pg.Exec(query, studentUUID, mongoID)
	if err != nil {
		_, _ = collection.DeleteOne(ctx, result.InsertedID)
		return errors.New("gagal menyimpan referensi prestasi")
	}
	return nil
}

func (r *achievementRepository) UpdateStatus(id string, status string, notes string, verifierID string) error {
	query := `
        UPDATE achievement_references 
        SET status = $1, rejection_note = $2, verified_by = $3, verified_at = NOW(), updated_at = NOW()
        WHERE id = $4 AND status = 'submitted'
    `
	result, err := r.pg.Exec(query, status, notes, verifierID, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return errors.New("prestasi tidak ditemukan atau belum disubmit") }
	return nil
}

func (r *achievementRepository) Submit(id string) error {
	query := `UPDATE achievement_references SET status = 'submitted', submitted_at = NOW(), updated_at = NOW() WHERE id = $1 AND status = 'draft'`
	result, err := r.pg.Exec(query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return errors.New("prestasi tidak ditemukan atau status bukan draft") }
	return nil
}

func (r *achievementRepository) FindMongoIDByRefID(refID string) (string, error) {
	var mongoID string
	query := "SELECT mongo_achievement_id FROM achievement_references WHERE id = $1"
	err := r.pg.QueryRow(query, refID).Scan(&mongoID)
	if err != nil { return "", errors.New("referensi prestasi tidak ditemukan") }
	return mongoID, nil
}

func (r *achievementRepository) UpdateMongo(mongoID string, updateData *models.AchievementMongo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil { return errors.New("format ID mongo tidak valid") }
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"achievementType": updateData.AchievementType,
			"title":           updateData.Title,
			"description":     updateData.Description,
			"details":         updateData.Details,
			"updatedAt":       time.Now(),
		},
	}
	if len(updateData.Attachments) > 0 {
		update["$set"].(bson.M)["attachments"] = updateData.Attachments
	}
	_, err = r.mongo.Collection("achievements").UpdateOne(ctx, filter, update)
	return err
}

func (r *achievementRepository) SoftDelete(id string) error {
	query := `UPDATE achievement_references SET status = 'deleted', updated_at = NOW() WHERE id = $1 AND status = 'draft'`
	result, err := r.pg.Exec(query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return errors.New("prestasi tidak ditemukan atau status bukan draft") }
	return nil
}

// --- Read Operations ---

func (r *achievementRepository) FindAllByStudentID(studentID string) ([]models.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at
        FROM achievement_references
        WHERE student_id = $1 AND status != 'deleted'
        ORDER BY created_at DESC
    `
	rows, err := r.pg.Query(query, studentID)
	if err != nil { return nil, err }
	defer rows.Close()
	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		var note sql.NullString
		var verAt sql.NullTime
		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt, &verAt); err != nil { return nil, err }
		ref.RejectionNote = note.String
		if verAt.Valid { ref.VerifiedAt = &verAt.Time }
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *achievementRepository) FindAllByAdvisorID(advisorID string) ([]models.AchievementReference, error) {
	query := `
        SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.rejection_note, ar.created_at
        FROM achievement_references ar
        JOIN students s ON ar.student_id = s.id
        WHERE s.advisor_id = $1 AND ar.status != 'deleted' AND ar.status != 'draft'
        ORDER BY ar.created_at DESC
    `
	rows, err := r.pg.Query(query, advisorID)
	if err != nil { return nil, err }
	defer rows.Close()
	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		var note sql.NullString
		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt); err != nil { return nil, err }
		ref.RejectionNote = note.String
		refs = append(refs, ref)
	}
	return refs, nil
}

// [NEW] Admin View All
func (r *achievementRepository) FindAllAchievements() ([]models.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at, submitted_at
        FROM achievement_references
        WHERE status != 'deleted'
        ORDER BY created_at DESC
    `
	rows, err := r.pg.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()
	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		var note sql.NullString
		var verAt, subAt sql.NullTime
		if err := rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt, &verAt, &subAt); err != nil { return nil, err }
		ref.RejectionNote = note.String
		if verAt.Valid { ref.VerifiedAt = &verAt.Time }
		if subAt.Valid { ref.SubmittedAt = &subAt.Time }
		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *achievementRepository) FindRefByID(id string) (*models.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, rejection_note, created_at, verified_at, submitted_at
        FROM achievement_references
        WHERE id = $1
    `
	var ref models.AchievementReference
	var note sql.NullString
	var verAt, subAt sql.NullTime
	err := r.pg.QueryRow(query, id).Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &note, &ref.CreatedAt, &verAt, &subAt)
	if err != nil { return nil, err }
	ref.RejectionNote = note.String
	if verAt.Valid { ref.VerifiedAt = &verAt.Time }
	if subAt.Valid { ref.SubmittedAt = &subAt.Time }
	return &ref, nil
}

func (r *achievementRepository) FindMongoDetails(mongoIDs []string) (map[string]models.AchievementMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var objectIDs []primitive.ObjectID
	for _, id := range mongoIDs {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil { objectIDs = append(objectIDs, oid) }
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := r.mongo.Collection("achievements").Find(ctx, filter)
	if err != nil { return nil, err }
	defer cursor.Close(ctx)
	results := make(map[string]models.AchievementMongo)
	for cursor.Next(ctx) {
		var doc models.AchievementMongo
		if err := cursor.Decode(&doc); err == nil { results[doc.ID.Hex()] = doc }
	}
	return results, nil
}