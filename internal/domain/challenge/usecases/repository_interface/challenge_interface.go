package repository_interface

import (
	"challenge-service/internal/domain/challenge/entity"
)

type AuthenticationChallengeParams struct {
	Name   *string
	Type   *string // семейный, личный, общий(групповой)
	IsTeam *bool
}

type ChallengeRepositoryInterface interface {
	Create(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error)
	Delete(challengeID int64) error
	Update(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error)
	FindAll() ([]*entity.AuthenticationChallenge, error)
	FindByParams(params *AuthenticationChallengeParams) ([]*entity.AuthenticationChallenge, error)
	GetAllChallengesFromUser(userID string) ([]*entity.AuthenticationChallenge, error)
	GetAllChallengesFromTeam(teamID string) ([]*entity.AuthenticationChallenge, error)
}
