//go:build wireinject
// +build wireinject

package dependency_injection

import (
	"context"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/mocks"
	"github.com/aluferraz/go-expert-zipkin/internal/infra/web/webhandlers/get_temperature_handler"
	"github.com/aluferraz/go-expert-zipkin/internal/usecase/get_temperature"
	"github.com/google/wire"
)

func NewTemperatureUseCase(client mocks.ZipkinClientInterface) get_temperature.UseCase {
	return get_temperature.NewUseCase(client)
}
func NewTemperatureHandler(ctx *context.Context, client mocks.ZipkinClientInterface) *get_temperature_handler.WebGetTemperatureHandler {
	wire.Build(NewTemperatureUseCase, get_temperature_handler.NewGetTemperatureHandler)
	return &get_temperature_handler.WebGetTemperatureHandler{}
}

/*
var setSampleRepositoryDependency = wire.NewSet(
	database.SampleRepository,
	wire.Bind(new(entity.SampleRepositoryInterface), new(*database.SampleRepository)),
)

func NewListAllOrdersUseCase(db *sql.DB) *usecase.MyUseCase {
	wire.Build(
		setSampleRepositoryDependency,
		usecase.NewUseCase,
	)
	return &usecase.MyUseCase{}
}
*/
