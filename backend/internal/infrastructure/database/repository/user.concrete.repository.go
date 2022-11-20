package repository

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"

	"gorm.io/gorm"
)

type userRepository struct {
	reposh.BaseRepository[entity.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: reposh.NewRepository[entity.User](db),
	}
}

func (r *userRepository) FindById(id string) (*entity.User, error) {
	user, err := r.FindOne(new(reposh.FindOption[entity.User]).Entity(&entity.User{
		Id: id,
	}))
	return reposh.FilteredRecord(user, err)
}
