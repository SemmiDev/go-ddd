package persistence

import (
	"DDD/domain/entity"
	"DDD/domain/repository"
	"errors"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db}
}

//ProductRepo implements the repository.ProductRepository interface
var _ repository.ProductRepository = &ProductRepo{}

func (r *ProductRepo) SaveProduct(product *entity.Product) (*entity.Product, map[string]string) {
	dbErr := map[string]string{}
	product.ProductImage = os.Getenv("DO_SPACES_URL") + product.ProductImage

	err := r.db.Debug().Create(&product).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "product title already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return product, nil
}

func (r *ProductRepo) GetProduct(id uint64) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Debug().Where("id = ?", id).Take(&product).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func (r *ProductRepo) GetAllProduct() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&products).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return products, nil
}

func (r *ProductRepo) UpdateProduct(product *entity.Product) (*entity.Product, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&product).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return product, nil
}

func (r *ProductRepo) DeleteProduct(id uint64) error {
	var product entity.Product
	err := r.db.Debug().Where("id = ?", id).Delete(&product).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
