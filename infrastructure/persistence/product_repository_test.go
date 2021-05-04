package persistence

import (
	"DDD/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveProduct_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var product = entity.Product{}
	product.Title = "product title"
	product.Description = "product description"
	product.UserID = 1

	repo := NewProductRepository(conn)

	f, saveErr := repo.SaveProduct(&product)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.Title, "product title")
	assert.EqualValues(t, f.Description, "product description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestSaveProduct_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the product
	_, err = seedProduct(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var product = entity.Product{}
	product.Title = "product title"
	product.Description = "product desc"
	product.UserID = 1

	repo := NewProductRepository(conn)
	f, saveErr := repo.SaveProduct(&product)

	dbMsg := map[string]string{
		"unique_title": "product title already taken",
	}
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, saveErr)
}

func TestGetProduct_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	product, err := seedProduct(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewProductRepository(conn)

	f, saveErr := repo.GetProduct(product.ID)

	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.Title, product.Title)
	assert.EqualValues(t, f.Description, product.Description)
	assert.EqualValues(t, f.UserID, product.UserID)
}

func TestGetAllProduct_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedProducts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewProductRepository(conn)
	products, getErr := repo.GetAllProduct()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(products), 2)
}

func TestUpdateProduct_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	product, err := seedProduct(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//updating
	product.Title = "product title update"
	product.Description = "product description update"

	repo := NewProductRepository(conn)
	f, updateErr := repo.UpdateProduct(product)

	assert.Nil(t, updateErr)
	assert.EqualValues(t, f.ID, 1)
	assert.EqualValues(t, f.Title, "product title update")
	assert.EqualValues(t, f.Description, "product description update")
	assert.EqualValues(t, f.UserID, 1)
}

//Duplicate title error
func TestUpdateProduct_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	products, err := seedProducts(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var secondProduct entity.Product

	//get the second product title
	for _, v := range products {
		if v.ID == 1 {
			continue
		}
		secondProduct = v
	}
	secondProduct.Title = "first product" //this title belongs to the first product already, so the second product cannot use it
	secondProduct.Description = "New description"

	repo := NewProductRepository(conn)
	f, updateErr := repo.UpdateProduct(&secondProduct)

	dbMsg := map[string]string{
		"unique_title": "title already taken",
	}
	assert.NotNil(t, updateErr)
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, updateErr)
}

func TestDeleteProduct_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	product, err := seedProduct(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewProductRepository(conn)

	deleteErr := repo.DeleteProduct(product.ID)

	assert.Nil(t, deleteErr)
}
