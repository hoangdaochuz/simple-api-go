package repository

import "example.com/go-api/internal/database"

type Product struct {
	ID 				uint   `json:"id" gorm:"primaryKey"`
	Name 			string `json:"name"`
	Price 		float64 `json:"price"`
}

func GetAllProducts() ([]Product, error) {
	var products []Product
	err := database.DB.Find(&products).Error
	return products, err
}

func GetProductById(id uint)(Product,error){
	var product Product
	err :=database.DB.First(&product,id).Error
	return product,err
}

func CreateProduct(product *Product) error {
	return database.DB.Create(product).Error
}

func UpdateProductHandler(product *Product) error{
	return database.DB.Save(product).Error
}

func DeleteProductHandler(id uint) error{
	return database.DB.Delete(&Product{},id).Error
}
