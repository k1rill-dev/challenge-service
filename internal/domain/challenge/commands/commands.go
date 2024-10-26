package commands

import (
	"challenge-service/internal/infrastructure/cqrs"
	"time"
)

type CreateChallengeCommand struct {
	cqrs.BaseCommand
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon        string    `gorm:"type:varchar(255);not null" json:"icon"`
	Image       string    `gorm:"type:varchar(255);not null" json:"image"`
	Description string    `gorm:"type:text;not null" json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `gorm:"type:timestamptz;not null" json:"end_date"`
	Type        string    `gorm:"type:varchar(10);not null" json:"type"` // семейный, личный, общий(групповой)
	IsTeam      bool      `gorm:"not null" json:"is_team"`
	CreatorID   int64     `gorm:"not null" json:"creator_id"`
}

func NewCreateChallengeCommand(id int64, name *string, icon *string, description *string,
	endDate *time.Time, typeChallenge *string, isTeam *bool, creatorID *int64) *CreateChallengeCommand {
	return &CreateChallengeCommand{
		BaseCommand: cqrs.NewBaseCommand(id),
		Name:        *name,
		Icon:        *icon,
		Description: *description,
		EndDate:     *endDate,
		Type:        *typeChallenge,
		IsTeam:      *isTeam,
		CreatorID:   *creatorID,
	}
}

func NewEmptyCreateChallengeCommand() *CreateChallengeCommand {
	return &CreateChallengeCommand{}
}

type UpdateChallengeCommand struct {
	cqrs.BaseCommand
	ChallengeID int64      `json:"challenge_id"`
	Name        *string    `json:"name,omitempty"`
	Icon        *string    `json:"icon,omitempty"`
	Image       *string    `json:"image,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Type        *string    `json:"type,omitempty"` // семейный, личный, общий(групповой)
	IsTeam      *bool      `json:"is_team,omitempty"`
	CreatorID   *int64     `json:"creator_id,omitempty"`
}

func NewUpdateChallengeCommand(id int64, challengeID int64, name *string, icon *string, image *string, description *string,
	endDate *time.Time, typeChallenge *string, isTeam *bool, creatorID *int64) *UpdateChallengeCommand {
	return &UpdateChallengeCommand{
		BaseCommand: cqrs.NewBaseCommand(id),
		ChallengeID: challengeID,
		Name:        name,
		Icon:        icon,
		Image:       image,
		Description: description,
		EndDate:     endDate,
		Type:        typeChallenge,
		IsTeam:      isTeam,
		CreatorID:   creatorID,
	}
}

func NewEmptyUpdateChallengeCommand() *UpdateChallengeCommand {
	return &UpdateChallengeCommand{}
}

type DeleteChallengeCommand struct {
	cqrs.BaseCommand
	ChallengeID int64 `json:"challenge_id"`
}

func NewDeleteChallengeCommand(id int64, challengeID int64) *DeleteChallengeCommand {
	return &DeleteChallengeCommand{
		BaseCommand: cqrs.NewBaseCommand(id),
		ChallengeID: challengeID,
	}
}

func NewEmptyDeleteChallengeCommand() *DeleteChallengeCommand {
	return &DeleteChallengeCommand{}
}
