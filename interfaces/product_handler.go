package interfaces

import (
	"DDD/application"
	"DDD/domain/entity"
	"DDD/infrastructure/auth"
	"DDD/interfaces/fileupload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Product struct {
	productApp application.ProductAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

//Product constructor
func NewProduct(fApp application.ProductAppInterface, uApp application.UserAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Product {
	return &Product{
		productApp:    fApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
	}
}

func (fo *Product) SaveProduct(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var saveProductError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new product for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyProduct := entity.Product{}
	emptyProduct.Title = title
	emptyProduct.Description = description
	saveProductError = emptyProduct.Validate("")
	if len(saveProductError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveProductError)
		return
	}
	file, err := c.FormFile("product_image")
	if err != nil {
		saveProductError["invalid_file"] = "a valid file is required"
		c.JSON(http.StatusUnprocessableEntity, saveProductError)
		return
	}
	//check if the user exist
	_, err = fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	uploadedFile, err := fo.fileUpload.UploadFile(file)
	if err != nil {
		saveProductError["upload_err"] = err.Error() //this error can be any we defined in the UploadFile method
		c.JSON(http.StatusUnprocessableEntity, saveProductError)
		return
	}
	var product = entity.Product{}
	product.UserID = userId
	product.Title = title
	product.Description = description
	product.ProductImage = uploadedFile
	savedProduct, saveErr := fo.productApp.SaveProduct(&product)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedProduct)
}

func (fo *Product) UpdateProduct(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updateProductError = make(map[string]string)

	productId, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new product for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyProduct := entity.Product{}
	emptyProduct.Title = title
	emptyProduct.Description = description
	updateProductError = emptyProduct.Validate("update")
	if len(updateProductError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateProductError)
		return
	}
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the product exist:
	product, err := fo.productApp.GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.ID != product.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this product")
		return
	}
	file, _ := c.FormFile("product_image")
	if file != nil {
		product.ProductImage, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		product.ProductImage = os.Getenv("DO_SPACES_URL") + product.ProductImage
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	//we dont need to update user's id
	product.Title = title
	product.Description = description
	product.UpdatedAt = time.Now()
	updatedProduct, dbUpdateErr := fo.productApp.UpdateProduct(product)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}

func (fo *Product) GetAllProduct(c *gin.Context) {
	allproduct, err := fo.productApp.GetAllProduct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allproduct)
}

func (fo *Product) GetProductAndCreator(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	product, err := fo.productApp.GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(product.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	productAndUser := map[string]interface{}{
		"product":    product,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, productAndUser)
}

func (fo *Product) DeleteProduct(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	productId, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.productApp.DeleteProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "product deleted")
}
