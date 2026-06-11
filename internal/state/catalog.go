package state

import (
	"fmt"
	"log"
	"sync"

	"github.com/joyboy1210/tex/internal/models"
)

var (
	catalogMut         sync.RWMutex
	categoryList       []models.Category
	productsByCategory map[uint][]models.Product
	productsByID       map[uint]models.Product
)

func Innitcatalog() error {
	catalogMut.Lock()
	defer catalogMut.Unlock()

	cats, err := models.GetAllCategories()
	if err != nil {
		return fmt.Errorf("failed to load categories: %w", err)
	}

	categoryList = cats
	prods, err := models.GetAllProducts()
	if err != nil {
		return fmt.Errorf("failed to load products: %w", err)
	}
	productsByCategory = make(map[uint][]models.Product)
	productsByID = make(map[uint]models.Product)
	for _, prod := range prods {
		productsByCategory[prod.CategoryID] = append(productsByCategory[prod.CategoryID], prod)
		productsByID[prod.ID] = prod
	}
	log.Printf("prodcuts loaded into cache")
	return nil
}

func GetCategories() []models.Category {
	catalogMut.RLock()
	defer catalogMut.RUnlock()
	return categoryList
}

func GetProductsByCategoryID(categoryID uint) []models.Product {
	catalogMut.RLock()
	defer catalogMut.RUnlock()
	return productsByCategory[categoryID]
}

func GetProductByID(productID uint) (models.Product, bool) {
	catalogMut.RLock()
	defer catalogMut.RUnlock()
	prod, exists := productsByID[productID]
	return prod, exists
}
