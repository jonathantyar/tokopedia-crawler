package service

import "context"

type ScrapperServiceInterface interface {
	HandleScraping(ctx context.Context, cfg Config)
	ScrapingData(page int, total uint64)
	GetData(jsonFile string, page int, total uint64)
	GetProducts(detail DetailProduct)
}
