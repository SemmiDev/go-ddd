package persistence

import (
	"DDD/domain/entity"
	"DDD/domain/repository"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repositories struct {
	User 	repository.UserRepository
	Product repository.ProductRepository
	db *gorm.DB
}

func NewRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		Product: NewProductRepository(db),
		db:   db,
	}, nil
}

func (s *Repositories) Close() error {
	return s.db.Close()
}

func (s *Repositories) Automigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.Product{}).Error
}