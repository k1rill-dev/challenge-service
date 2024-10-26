package fabric

import (
	"challenge-service/internal/infrastructure/cqrs"
	"fmt"
	"reflect"
)

type HandlerFabric struct {
	commandHandlers map[reflect.Type]cqrs.CommandHandler[cqrs.Command]
	queryHandlers   map[reflect.Type]cqrs.QueryHandler[cqrs.Query]
}

func NewHandlerFabric() *HandlerFabric {
	return &HandlerFabric{
		commandHandlers: make(map[reflect.Type]cqrs.CommandHandler[cqrs.Command]),
		queryHandlers:   make(map[reflect.Type]cqrs.QueryHandler[cqrs.Query]),
	}
}

func (handlerFabric *HandlerFabric) RegisterCommandHandler(command cqrs.Command, handler cqrs.CommandHandler[cqrs.Command]) {
	handlerFabric.commandHandlers[reflect.TypeOf(command)] = handler
}
func (handlerFabric *HandlerFabric) RegisterQueryHandler(query cqrs.Query, handler cqrs.QueryHandler[cqrs.Query]) {
	handlerFabric.queryHandlers[reflect.TypeOf(query)] = handler
}
func (handlerFabric *HandlerFabric) GetCommandHandler(command cqrs.Command) (cqrs.CommandHandler[cqrs.Command], error) {
	handler, ok := handlerFabric.commandHandlers[reflect.TypeOf(command)]
	if !ok {
		return nil, fmt.Errorf("command handler not registered")
	}
	return handler, nil
}
func (handlerFabric *HandlerFabric) GetQueryHandler(query cqrs.Query) (cqrs.QueryHandler[cqrs.Query], error) {
	handler, ok := handlerFabric.queryHandlers[reflect.TypeOf(query)]
	if !ok {
		return nil, fmt.Errorf("query handler not registered")
	}
	return handler, nil
}
