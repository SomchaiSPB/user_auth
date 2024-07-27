package repository

import (
	"github.com/SomchaiSPB/user-auth/internal/entity"
	"gorm.io/gorm"
)

type ProductDBRepository struct {
	db *gorm.DB
}

func NewProductDBRepository(db *gorm.DB) ProductDBRepository {
	return ProductDBRepository{db: db}
}

func (r ProductDBRepository) Create(p *entity.Product) (*entity.Product, error) {
	return p, r.db.Create(&p).Error
}

func (r ProductDBRepository) GetByID(id uint) (*entity.Product, error) {
	var p *entity.Product

	return p, r.db.Find(&p, id).Error
}

func (r ProductDBRepository) Get(limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product

	return products, r.db.Limit(limit).Offset(offset).Find(&products).Error
}

func (r ProductDBRepository) GetByName(name string) (*entity.Product, error) {
	var p *entity.Product

	return p, r.db.Where("LOWER(name) = LOWER(?)", name).First(&p).Error
}

func (r ProductDBRepository) GetWithFilters(filters ...Filter) ([]*entity.Product, error) {
	var products []*entity.Product

	db := r.db

	for _, filter := range filters {
		db = db.Where(filter.Query(), filter.Args())
	}

	return products, db.Find(products).Error
}
