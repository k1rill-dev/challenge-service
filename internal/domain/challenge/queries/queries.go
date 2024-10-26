package queries

import (
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
)

type FindAllQuery struct {
	cqrs.BaseQuery
}

func NewFindAllQuery() *FindAllQuery {
	return &FindAllQuery{}
}

type FindByParamsQuery struct {
	cqrs.BaseQuery
	Params *repository_interface.AuthenticationChallengeParams `json:"params"`
}

func NewFindByParamsQuery(id int64, params *repository_interface.AuthenticationChallengeParams) *FindByParamsQuery {
	return &FindByParamsQuery{
		BaseQuery: cqrs.NewBaseQuery(id),
		Params:    params,
	}
}

func NewEmptyFindByParamsQuery() *FindByParamsQuery {
	return &FindByParamsQuery{}
}

type GetAllChallengesFromUserQuery struct {
	cqrs.BaseQuery
	UserID string `json:"user_id"`
}

func NewGetAllChallengesFromUserQuery(id int64, userID string) *GetAllChallengesFromUserQuery {
	return &GetAllChallengesFromUserQuery{
		BaseQuery: cqrs.NewBaseQuery(id),
		UserID:    userID,
	}
}

func NewEmptyGetAllChallengesFromUserQuery() *GetAllChallengesFromUserQuery {
	return &GetAllChallengesFromUserQuery{}
}

type GetAllChallengesFromTeamQuery struct {
	cqrs.BaseQuery
	TeamID string `json:"team_id"`
}

func NewGetAllChallengesFromTeamQuery(id int64, teamID string) *GetAllChallengesFromTeamQuery {
	return &GetAllChallengesFromTeamQuery{
		BaseQuery: cqrs.NewBaseQuery(id),
		TeamID:    teamID,
	}
}

func NewEmptyGetAllChallengesFromTeamQuery() *GetAllChallengesFromTeamQuery {
	return &GetAllChallengesFromTeamQuery{}
}
