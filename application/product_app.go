package application

import (
	"DDD/domain/entity"
	"DDD/domain/repository"
)

type productApp struct {
	fr repository.ProductRepository
}

var _ ProductAppInterface = &productApp{}

type ProductAppInterface interface {
	SaveProduct(*entity.Product) (*entity.Product, map[string]string)
	GetAllProduct() ([]entity.Product, error)
	GetProduct(uint64) (*entity.Product, error)
	UpdateProduct(*entity.Product) (*entity.Product, map[string]string)
	DeleteProduct(uint64) error
}

func (f *productApp) SaveProduct(product *entity.Product) (*entity.Product, map[string]string) {
	return f.fr.SaveProduct(product)
}

func (f *productApp) GetAllProduct() ([]entity.Product, error) {
	return f.fr.GetAllProduct()
}

func (f *productApp) GetProduct(productId uint64) (*entity.Product, error) {
	return f.fr.GetProduct(productId)
}

func (f *productApp) UpdateProduct(product *entity.Product) (*entity.Product, map[string]string) {
	return f.fr.UpdateProduct(product)
}

func (f *productApp) DeleteProduct(productId uint64) error {
	return f.fr.DeleteProduct(productId)
}