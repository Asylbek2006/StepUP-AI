package entity

import "time"

type University struct {
	ID             string
	Name           string
	Country        string
	AcceptanceRate float32
	Type           string
	Category       string
	CreatedAt      time.Time
}

type Grant struct {
	ID          string
	Name        string
	Description string
	Amount      float32
	Deadline    string
	Country     string
	MinGPA      float32
	CreatedAt   time.Time
}

type SavedUniversity struct {
	ID           string
	UserID       string
	UniversityID string
	SavedAt      time.Time
}

type SavedGrant struct {
	ID      string
	UserID  string
	GrantID string
	SavedAt time.Time
}
