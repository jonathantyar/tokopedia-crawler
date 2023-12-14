package repository

import (
	"context"
	"jonathantyar/tokopedia-crawler/src/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(ctx context.Context) []model.Product
	Create(ctx context.Context, product []model.Product) error
}

type ProductRepositoryImpl struct {
	*gorm.DB
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context) (res []model.Product) {
	tx := r.WithContext(ctx)
	tx.Find(&res)
	return
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product []model.Product) error {
	tx := r.WithContext(ctx)
	return tx.Save(&product).Error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}
