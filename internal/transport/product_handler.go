package transport

import (
	"net/http"
	"strconv"

	"example.com/go-api/internal/repository"
	"example.com/go-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetAllProductsHandler handles the HTTP GET request to retrieve all products.
// @Summary Get all products
// @Description Retrieve a list of all products from the service
// @Tags Products
// @Produce json
// @Success 200 {array} repository.Product "List of products"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /products [get]
func GetAllProductsHandler(ctx *gin.Context){
	products,err :=services.GetAllProductsService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK,products)
}

// GetProductById godoc
// @Summary Get product by ID
// @Description Retrieve a product using its unique ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} repository.Product
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Product not found"
// @Router /products/{id} [get]
func GetProductById(ctx *gin.Context){
	idPrams := ctx.Params.ByName("id")
	id,err := strconv.Atoi(idPrams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid ID"})
		return
	}
	products,err := services.GetProductByIdService(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound,gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK,products)

}


// CreateProductHandler handles the creation of a new product.
// @Summary Create a new product
// @Description Create a new product by providing product details in the request body
// @Tags Products
// @Accept json
// @Produce json
// @Param product body repository.Product true "Product details"
// @Success 201 {object} repository.Product "Created product"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [post]
func CreateProductHandler(ctx *gin.Context){
	var product repository.Product
	err:=ctx.ShouldBindBodyWithJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	err = services.CreateProductService(&product);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated,product)
}

// generate a swagger documentation for the UpdateProductHandler
// @Summary Update an existing product
// @Description Update an existing product by providing the product details in the request body
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body repository.Product true "Product details"
// @Success 200 {object} repository.Product "Updated product"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [put]
func UpdateProductHandler(ctx *gin.Context){
	idParams := ctx.Params.ByName("id")
	id,err := strconv.Atoi(idParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid ID"})
		return
	}
	var product repository.Product
	err = ctx.ShouldBindBodyWithJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err})
		return
	}
	product.ID = uint(id)
	err = services.UpdateProductService(&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK,product)
}
// generate a swagger documentation for the DeleteProductHandler
// @Summary Delete a product
// @Description Delete a product using its unique ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "Product deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [delete]
func DeleteProductHandler(ctx *gin.Context){
	idParams :=ctx.Params.ByName("id")
	id,err :=strconv.Atoi(idParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid ID"})
		return
	}
	err = services.DeleteProductService(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"message":"Product deleted successfully"})
}
func GetRecentlyProductsHandler(ctx *gin.Context){}