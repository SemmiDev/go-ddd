package persistence

import (
	"DDD/domain/entity"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		return LocalDatabase()
	}
	return CIBuild()
}

func CIBuild() (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "127.0.0.1", "5432", "sammidev", "product_app_test", "sammidev")
	conn, err := gorm.Open("postgres", DBURL)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	return conn, nil
}

func LocalDatabase() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	conn, err := gorm.Open(dbdriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", dbdriver)
	}

	err = conn.DropTableIfExists(&entity.User{}, &entity.Product{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		entity.User{},
		entity.Product{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func seedUser(db *gorm.DB) (*entity.User, error) {
	user := &entity.User{
		ID:        1,
		FirstName: "Sammi",
		LastName:  "Dev",
		Email:     "sammidev@gmail.com",
		Password:  "sammidev",
		DeletedAt: nil,
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedUsers(db *gorm.DB) ([]entity.User, error) {
	users := []entity.User{
		{
			ID:        1,
			FirstName: "Sammi",
			LastName:  "Aldhi Yanto",
			Email:     "sammidev@gmail.com",
			Password:  "sammidev",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "Rahmatul",
			LastName:  "Izzah Annisa",
			Email:     "izzaah@yahoo.com",
			Password:  "izzaah",
			DeletedAt: nil,
		},
	}
	for _, v := range users {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func seedProduct(db *gorm.DB) (*entity.Product, error) {
	product := &entity.Product{
		ID:          1,
		Title:       "product title",
		Description: "product desc",
		UserID:      1,
	}
	err := db.Create(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func seedProducts(db *gorm.DB) ([]entity.Product, error) {
	products := []entity.Product{
		{
			ID:          1,
			Title:       "first product",
			Description: "first desc",
			UserID:      1,
		},
		{
			ID:          2,
			Title:       "second product",
			Description: "second desc",
			UserID:      1,
		},
	}
	for _, v := range products {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}