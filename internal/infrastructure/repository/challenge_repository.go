package repository

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/entity"
	interfaceRepo "challenge-service/internal/domain/challenge/usecases/repository_interface"
	"gorm.io/gorm"
	"log/slog"
)

type challengeRepository struct {
	interfaceRepo.ChallengeRepositoryInterface
	cfg *config.Config
	log *slog.Logger
	db  *gorm.DB
}

func NewChallengeRepository(cfg *config.Config, log *slog.Logger, db *gorm.DB) interfaceRepo.ChallengeRepositoryInterface {
	return &challengeRepository{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (c *challengeRepository) FindAll() ([]*entity.AuthenticationChallenge, error) {
	return nil, nil
}
func (c *challengeRepository) Create(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error) {
	return nil, nil
}

func (c *challengeRepository) Update(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error) {
	return nil, nil
}
func (c *challengeRepository) Delete(challengeID int64) error {
	return nil
}
func (c *challengeRepository) FindByParams(params *interfaceRepo.AuthenticationChallengeParams) ([]*entity.AuthenticationChallenge, error) {
	return nil, nil
}

func (c *challengeRepository) GetAllChallengesFromUser(userID string) ([]*entity.AuthenticationChallenge, error) {
	return nil, nil
}

func (c *challengeRepository) GetAllChallengesFromTeam(teamID string) ([]*entity.AuthenticationChallenge, error) {
	return nil, nil
}
