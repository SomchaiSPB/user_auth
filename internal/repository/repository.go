package repository

import "github.com/SomchaiSPB/user-auth/internal/entity"

type UserRepository interface {
	Create(u *entity.User) (*entity.User, error)
	Exists(username string) bool
	GetByName(username string) (*entity.User, error)
}

type ProductRepository interface {
	Create(p *entity.Product) (*entity.Product, error)
	GetByID(id uint) (*entity.Product, error)
	Get(limit, offset int) ([]*entity.Product, error)
	GetByName(name string) (*entity.Product, error)
	GetWithFilters(filters ...Filter) ([]*entity.Product, error)
}
