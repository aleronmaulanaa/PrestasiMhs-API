// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// // --- MongoDB Models (Sesuai SRS 3.2.1) ---

// type AchievementMongo struct {
// 	ID              primitive.ObjectID `bson:"_id,omitempty"`
// 	StudentID       string             `bson:"studentId"` // Disimpan sebagai string UUID
// 	AchievementType string             `bson:"achievementType"`
// 	Title           string             `bson:"title"`
// 	Description     string             `bson:"description"`
// 	Details         AchievementDetails `bson:"details"`
// 	Attachments     []Attachment       `bson:"attachments"`
// 	Tags            []string           `bson:"tags"`
// 	Points          int                `bson:"points"`
// 	CreatedAt       time.Time          `bson:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt"`
// }

// // AchievementDetails menangani field dinamis
// type AchievementDetails struct {
// 	// Competition
// 	CompetitionName  string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
// 	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
// 	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
// 	MedalType        string `bson:"medalType,omitempty" json:"medal_type,omitempty"`
	
// 	// Organization
// 	OrganizationName string    `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
// 	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
// 	StartDate        time.Time `bson:"startDate,omitempty" json:"start_date,omitempty"`
// 	EndDate          time.Time `bson:"endDate,omitempty" json:"end_date,omitempty"`

// 	// General
// 	EventDate time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
// 	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
// 	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
// }

// type Attachment struct {
// 	FileName   string    `bson:"fileName"`
// 	FileURL    string    `bson:"fileUrl"`
// 	FileType   string    `bson:"fileType"`
// 	UploadedAt time.Time `bson:"uploadedAt"`
// }

// // --- Request DTO (Input dari Postman/Form-Data) ---

// // CreateAchievementRequest menangkap inputan form-data
// type CreateAchievementRequest struct {
// 	AchievementType string `form:"achievement_type" validate:"required"`
// 	Title           string `form:"title" validate:"required"`
// 	Description     string `form:"description" validate:"required"`
	
// 	// Details fields (Flat input from form-data)
// 	CompetitionName  string `form:"competition_name"`
// 	CompetitionLevel string `form:"competition_level"`
// 	Rank             int    `form:"rank"`
// 	OrganizationName string `form:"organization_name"`
// 	Position         string `form:"position"`
// 	Location         string `form:"location"`
// 	Organizer        string `form:"organizer"`
// 	EventDate        string `form:"event_date"` // String YYYY-MM-DD
// }


package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MongoDB Models (Sesuai SRS 3.2.1) ---

type AchievementMongo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID       string             `bson:"studentId" json:"student_id"` // Disimpan sebagai string UUID
	AchievementType string             `bson:"achievementType" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	CreatedAt       time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updated_at"`
}

// AchievementDetails menangani field dinamis
// type AchievementDetails struct {
// 	// Competition
// 	CompetitionName  string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
// 	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
// 	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
// 	MedalType        string `bson:"medalType,omitempty" json:"medal_type,omitempty"`
	
// 	// Organization
// 	OrganizationName string    `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
// 	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
// 	StartDate        time.Time `bson:"startDate,omitempty" json:"start_date,omitempty"`
// 	EndDate          time.Time `bson:"endDate,omitempty" json:"end_date,omitempty"`

// 	// General
// 	EventDate time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
// 	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
// 	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
// }
// app/models/achievement_model.go

// AchievementDetails menangani field dinamis
type AchievementDetails struct {
	// --- Competition ---
	CompetitionName  string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string `bson:"medalType,omitempty" json:"medal_type,omitempty"`

	// --- Organization ---
	OrganizationName string    `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
	StartDate        time.Time `bson:"startDate,omitempty" json:"start_date,omitempty"`
	EndDate          time.Time `bson:"endDate,omitempty" json:"end_date,omitempty"`

	// --- Publication (NEW - Sesuai SRS) ---
	PublicationType  string   `bson:"publicationType,omitempty" json:"publication_type,omitempty"`
	PublicationTitle string   `bson:"publicationTitle,omitempty" json:"publication_title,omitempty"`
	Authors          []string `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher        string   `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN             string   `bson:"issn,omitempty" json:"issn,omitempty"`

	// --- Certification (NEW - Sesuai SRS) ---
	CertificationName   string    `bson:"certificationName,omitempty" json:"certification_name,omitempty"`
	IssuedBy            string    `bson:"issuedBy,omitempty" json:"issued_by,omitempty"`
	CertificationNumber string    `bson:"certificationNumber,omitempty" json:"certification_number,omitempty"`
	ValidUntil          time.Time `bson:"validUntil,omitempty" json:"valid_until,omitempty"`

	// --- General ---
	EventDate time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score     float64   `bson:"score,omitempty" json:"score,omitempty"` // NEW

	// --- Dynamic Extension (NEW - Sesuai SRS) ---
	// Menampung field lain yang tidak terdefinisi (SRS: customFields?: Object)
	CustomFields map[string]interface{} `bson:"customFields,omitempty" json:"custom_fields,omitempty"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}

// --- Request DTO (Input dari Postman/Form-Data) ---

// CreateAchievementRequest menangkap inputan form-data
type CreateAchievementRequest struct {
	AchievementType string `form:"achievement_type" validate:"required"`
	Title           string `form:"title" validate:"required"`
	Description     string `form:"description" validate:"required"`
	
	// Details fields (Flat input from form-data)
	CompetitionName  string `form:"competition_name"`
	CompetitionLevel string `form:"competition_level"`
	Rank             int    `form:"rank"`
	OrganizationName string `form:"organization_name"`
	Position         string `form:"position"`
	Location         string `form:"location"`
	Organizer        string `form:"organizer"`
	EventDate        string `form:"event_date"` // String YYYY-MM-DD
}

// --- PostgreSQL Reference Model (TAMBAHAN BARU) ---
// Struct ini digunakan oleh Repository untuk menampung hasil query tabel 'achievement_references'
type AchievementReference struct {
	ID                 string     `json:"id"` // UUID dari Postgres
	StudentID          string     `json:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Status             string     `json:"status"`
	RejectionNote      string     `json:"rejection_note,omitempty"`
	SubmittedAt        *time.Time `json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	
	// Field ini diisi setelah menggabungkan data dengan MongoDB di Service layer
	Detail *AchievementMongo `json:"detail,omitempty"`
}