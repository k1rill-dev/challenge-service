package cqrs

import (
	"context"
)

type Command interface {
	GetAggregateID() int64
}

type CommandHandler[AbstractCommand Command] interface {
	Handle(ctx context.Context, command AbstractCommand) (interface{}, error)
}

type BaseCommand struct {
	AggregateID int64
}

func NewBaseCommand(aggregateID int64) BaseCommand {
	return BaseCommand{AggregateID: aggregateID}
}

func (c BaseCommand) GetAggregateID() int64 {
	return c.AggregateID
}
