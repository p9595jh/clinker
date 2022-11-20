package repository

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
)

type UserRepository interface {
	reposh.BaseRepository[entity.User]
	FindById(id string) (*entity.User, error)
}
