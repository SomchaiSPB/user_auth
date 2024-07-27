package service

import "github.com/SomchaiSPB/user-auth/internal/repository"

type ProductService struct {
	productRepository repository.ProductRepository
}

func NewProductSvc(pr repository.ProductRepository) *ProductService {
	return &ProductService{productRepository: pr}
}
