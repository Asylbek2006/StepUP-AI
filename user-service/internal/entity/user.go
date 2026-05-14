package entity

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserProfile struct {
	UserID     string
	GPA        float32
	SATScore   int32
	IELTSScore float32
	TOEFLScore int32
	Country    string
	UpdatedAt  time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	Used      bool
}
