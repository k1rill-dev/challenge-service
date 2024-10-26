package repository

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/entity"
	interfaceRepo "challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/lib/log"
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

// Получение всех вызовов
func (c *challengeRepository) FindAll() ([]*entity.AuthenticationChallenge, error) {
	var challenges []*entity.AuthenticationChallenge
	if err := c.db.Find(&challenges).Error; err != nil {
		c.log.Error("failed to fetch challenges", log.Err(err))
		return nil, err
	}
	return challenges, nil
}

// Создание нового вызова
func (c *challengeRepository) Create(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error) {
	if err := c.db.Create(&challenge).Error; err != nil {
		c.log.Error("failed to create challenge", log.Err(err))
		return nil, err
	}
	return &challenge, nil
}

// Обновление существующего вызова
func (c *challengeRepository) Update(challenge entity.AuthenticationChallenge) (*entity.AuthenticationChallenge, error) {
	if err := c.db.Save(&challenge).Error; err != nil {
		c.log.Error("failed to update challenge", log.Err(err))
		return nil, err
	}
	return &challenge, nil
}

// Удаление вызова по ID
func (c *challengeRepository) Delete(challengeID int64) error {
	if err := c.db.Delete(&entity.AuthenticationChallenge{}, challengeID).Error; err != nil {
		c.log.Error("failed to delete challenge", log.Err(err))
		return err
	}
	return nil
}

// Поиск вызовов по параметрам
func (c *challengeRepository) FindByParams(params *interfaceRepo.AuthenticationChallengeParams) ([]*entity.AuthenticationChallenge, error) {
	var challenges []*entity.AuthenticationChallenge
	query := c.db.Model(&entity.AuthenticationChallenge{})

	if *params.Name != "" {
		query = query.Where("name = ?", params.Name)
	}
	if *params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.IsTeam != nil {
		query = query.Where("end_date = ?", params.IsTeam)
	}

	if err := query.Find(&challenges).Error; err != nil {
		c.log.Error("failed to find challenges by params", log.Err(err))
		return nil, err
	}
	return challenges, nil
}

// Получение всех вызовов пользователя по userID
func (c *challengeRepository) GetAllChallengesFromUser(userID string) ([]*entity.AuthenticationChallenge, error) {
	var challenges []*entity.AuthenticationChallenge
	if err := c.db.Where("creator_id = ?", userID).Find(&challenges).Error; err != nil {
		c.log.Error("failed to fetch user challenges", log.Err(err))
		return nil, err
	}
	return challenges, nil
}

// Получение всех вызовов для команды по teamID
func (c *challengeRepository) GetAllChallengesFromTeam(teamID string) ([]*entity.AuthenticationChallenge, error) {
	var challenges []*entity.AuthenticationChallenge
	if err := c.db.Joins("JOIN authentication_participants ON authentication_participants.challenge_id = authentication_challenges.id").
		Where("authentication_participants.team_id = ?", teamID).
		Find(&challenges).Error; err != nil {
		c.log.Error("failed to fetch team challenges", log.Err(err))
		return nil, err
	}
	return challenges, nil
}
