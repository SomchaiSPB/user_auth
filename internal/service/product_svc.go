package service

import (
	"encoding/json"
	"errors"
	"github.com/SomchaiSPB/user-auth/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrEmptyRequestName = errors.New("product name is empty in uri query error")
	ErrProductNotFound  = errors.New("product not found error")
)

type ProductService struct {
	productRepository repository.ProductRepository
}

func NewProductSvc(pr repository.ProductRepository) *ProductService {
	return &ProductService{productRepository: pr}
}

func (s ProductService) GetProduct(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyRequestName
	}

	p, err := s.productRepository.GetByName(name)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
	}

	return json.Marshal(p)
}
