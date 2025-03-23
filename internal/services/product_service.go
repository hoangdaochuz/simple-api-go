package services

import (
	"encoding/json"
	"fmt"

	"example.com/go-api/internal/caching"
	"example.com/go-api/internal/repository"
)

func GetAllProductsService()([]repository.Product,error){
	// cacheListProject
	// cacheKey := fmt.Sprintf("products_all")
	// productsRes,err := caching.GetOrSetCache(cacheKey,func() (map[string]interface{},error){
	// 	products,err := repository.GetAllProducts()
	// 	if err != nil {
	// 		return nil,err
	// 	}
	// 	return map[string]interface{}{
	// 		"data": products,
	// 	},nil
	// })

	// if err != nil {
	// 	return nil,err
	// }

	// // products,ok := productsRes.([]repository.Product)

	// switch value := productsRes.(type) {
	// 	case []repository.Product:
	// 		return value,nil
	// 	case string:
	// 		var products []repository.Product
	// 		err:=json.Unmarshal([]byte(value),&products);
	// 		if err != nil {
	// 			return nil,err
	// 		}
	// 		return products,nil
	// 	default:
	// 		return nil,errors.New("Invalid type")
	// }
	return repository.GetAllProducts()
}

func GetProductByIdService(id uint)(repository.Product,error){

	cacheProductKey := fmt.Sprintf("product_%d",id)
	productRes,err := caching.GetOrSetCache(cacheProductKey,func() (map[string]interface{}, error) {
		product,err := repository.GetProductById(id)
		if err != nil {
			return nil,err
		}
		return map[string]interface{}{
			"data":product,
		},nil
	})
	if err != nil {
		return repository.Product{},err
	}
	product,ok := productRes.(repository.Product)
	if !ok {
		json.Unmarshal([]byte(productRes.(string)),&product)
		return product,nil
	}

	return productRes.(repository.Product),nil
}

func CreateProductService(product *repository.Product) error{
	err:= repository.CreateProduct(product)
	if err != nil {
		return err
	}
	err = caching.SetCache(fmt.Sprintf("product_%d",product.ID),*product)
	return err
}
func UpdateProductService(product *repository.Product)error{
	err :=repository.UpdateProductHandler(product)
	if err != nil {
		return err
	}
	err = caching.SetCache(fmt.Sprintf("product_%d",product.ID),*product)
	return err
}

func DeleteProductService(id uint)error{
	
	err := caching.InvalidCacheKey(fmt.Sprintf("product_%d",id))
	if err != nil {
		return err
	}
	return repository.DeleteProductHandler(id)
}