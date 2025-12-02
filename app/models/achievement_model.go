package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MongoDB Models (Sesuai SRS 3.2.1) ---

type AchievementMongo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	StudentID       string             `bson:"studentId"` // Disimpan sebagai string UUID
	AchievementType string             `bson:"achievementType"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`
	Details         AchievementDetails `bson:"details"`
	Attachments     []Attachment       `bson:"attachments"`
	Tags            []string           `bson:"tags"`
	Points          int                `bson:"points"`
	CreatedAt       time.Time          `bson:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
}

// AchievementDetails menangani field dinamis
type AchievementDetails struct {
	// Competition
	CompetitionName  string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string `bson:"medalType,omitempty" json:"medal_type,omitempty"`
	
	// Organization
	OrganizationName string    `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
	StartDate        time.Time `bson:"startDate,omitempty" json:"start_date,omitempty"`
	EndDate          time.Time `bson:"endDate,omitempty" json:"end_date,omitempty"`

	// General
	EventDate time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
}

type Attachment struct {
	FileName   string    `bson:"fileName"`
	FileURL    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
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