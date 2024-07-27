package repository

import (
	"github.com/SomchaiSPB/user-auth/internal/entity"
	"gorm.io/gorm"
)

type UserDBRepository struct {
	db *gorm.DB
}

func NewUserDBRepository(db *gorm.DB) UserDBRepository {
	return UserDBRepository{db: db}
}

func (r UserDBRepository) Create(u *entity.User) (*entity.User, error) {
	return u, r.db.Create(&u).Error
}

func (r UserDBRepository) Exists(username string) bool {
	var exists bool

	r.db.Model(&entity.User{}).Select("count(id) > 0").Where("name = ?", username).Find(&exists) // ignore errors for sake of time

	return exists
}

func (r UserDBRepository) GetByName(username string) (*entity.User, error) {
	var u *entity.User

	return u, r.db.Where("name = ?", username).First(&u).Error
}
