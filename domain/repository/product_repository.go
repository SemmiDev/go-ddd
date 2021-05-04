package repository

import "DDD/domain/entity"

type ProductRepository interface {
	SaveProduct(*entity.Product) (*entity.Product, map[string]string)
	GetProduct(uint64) (*entity.Product, error)
	GetAllProduct() ([]entity.Product, error)
	UpdateProduct(*entity.Product) (*entity.Product, map[string]string)
	DeleteProduct(uint64) error
}