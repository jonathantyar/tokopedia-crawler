//go:build wireinject
// +build wireinject

package src

import (
	"jonathantyar/tokopedia-crawler/src/database"
	"jonathantyar/tokopedia-crawler/src/repository"
	"jonathantyar/tokopedia-crawler/src/service"

	"github.com/google/wire"
)

func InitializeTokopediaService() (service.ScrapperServiceInterface, error) {
	wire.Build(
		database.InitDB,
		service.NewScrapperServiceTokopedia,
		repository.NewProductRepository,
	)

	return nil, nil
}
