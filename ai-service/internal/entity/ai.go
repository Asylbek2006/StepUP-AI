package entity

import "time"

type AdmissionAnalysis struct {
	ID                        string
	UserID                    string
	UniversityID              string
	AdmissionChancePercentage float32
	WeakAreas                 []string
	ImprovementSuggestions    []string
	GapAnalysis               string
	CreatedAt                 time.Time
}

type Roadmap struct {
	ID        string
	UserID    string
	CreatedAt time.Time
}

type RoadmapStep struct {
	ID          string
	RoadmapID   string
	Title       string
	Description string
	MonthNumber int32
	Category    string
}

type EssayReview struct {
	ID                     string
	UserID                 string
	EssayText              string
	UniversityName         string
	ProgramName            string
	GrammarScore           float32
	CoherenceScore         float32
	UniquenessScore        float32
	RelevanceScore         float32
	ImprovementSuggestions []string
	CreatedAt              time.Time
}
