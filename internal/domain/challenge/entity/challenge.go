package entity

import (
	"time"
)

type AuthenticationChallenge struct {
	ID          int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon        string    `gorm:"type:varchar(255);not null" json:"icon"`
	Image       string    `gorm:"type:varchar(255);not null" json:"image"`
	Description string    `gorm:"type:text;not null" json:"description"`
	StartDate   time.Time `gorm:"type:timestamptz;not null" json:"start_date"`
	EndDate     time.Time `gorm:"type:timestamptz;not null" json:"end_date"`
	Type        string    `gorm:"type:varchar(10);not null" json:"type"` // семейный, личный, общий(групповой)
	IsTeam      bool      `gorm:"not null" json:"is_team"`
	IsFinished  bool      `gorm:"not null" json:"is_finished"`
	CreatorID   int64     `gorm:"not null" json:"creator_id"`
}

func (AuthenticationChallenge) TableName() string {
	return "authentication_challenge"
}

type AuthenticationParticipant struct {
	ID          int64                   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Status      string                  `gorm:"type:varchar(10);not null" json:"status"`
	Progress    string                  `gorm:"type:jsonb;not null" json:"progress"`
	Achievement string                  `gorm:"type:text;not null" json:"achievement"`
	ChallengeID int64                   `gorm:"not null" json:"challenge_id"`
	Challenge   AuthenticationChallenge `gorm:"foreignKey:ChallengeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID      int64                   `gorm:"not null" json:"creator_id"`
	TeamID      int64                   `gorm:"not null" json:"team_id"`
}
