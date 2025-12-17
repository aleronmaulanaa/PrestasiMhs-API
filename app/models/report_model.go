package models

import "time"

// Struktur Response Transkrip Prestasi
type StudentReportResponse struct {
	StudentInfo StudentHeader         `json:"student_info"`
	Summary     ReportSummary         `json:"summary"`
	Achievements []AchievementReportItem `json:"achievements"`
}

type StudentHeader struct {
	FullName     string `json:"full_name"`
	NIM          string `json:"student_id"` // NIM
	ProgramStudy string `json:"program_study"`
	AdvisorName  string `json:"advisor_name"`
}

type ReportSummary struct {
	TotalEntries int `json:"total_entries"`
	TotalVerified int `json:"total_verified"`
	TotalPoints   int `json:"total_points"` // Opsional: Jika nanti ada sistem poin SKP
}

type AchievementReportItem struct {
	Title           string     `json:"title"`
	Type            string     `json:"type"` // Kompetisi/Organisasi
	EventDate       *time.Time `json:"event_date,omitempty"` // Dari Mongo Details
	Status          string     `json:"status"`
	VerificationDate *time.Time `json:"verified_at,omitempty"`
}