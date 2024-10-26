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
	Interests   string    `gorm:"type:text;not null" json:"interests"`
	EndDate     time.Time `gorm:"type:timestamptz;not null" json:"end_date"`
	Type        string    `gorm:"type:varchar(10);not null" json:"type"` // семейный, личный, общий(групповой)
	IsTeam      bool      `gorm:"not null" json:"is_team"`
	CreatorID   int64     `gorm:"not null" json:"creator_id"`
}
