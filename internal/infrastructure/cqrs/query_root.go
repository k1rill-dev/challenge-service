package cqrs

import (
	"context"
)

type Query interface {
	GetAggregateID() int64
}

type QueryHandler[AbstractQuery Query] interface {
	Handle(ctx context.Context, query AbstractQuery) (interface{}, error)
}

type BaseQuery struct {
	AggregateID int64
}

func NewBaseQuery(aggregateID int64) BaseQuery {
	return BaseQuery{
		AggregateID: aggregateID,
	}
}

func (q BaseQuery) GetAggregateID() int64 {
	return q.AggregateID
}
