package repository

import (
	"layer/internal/infrastructure/database/entity"
	"layer/internal/infrastructure/database/repository/reposh"

	"github.com/p9595jh/transform"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClinkRepository struct {
	*reposh.BaseRepository[entity.ClinkDao, entity.Clink]
}

func NewClinkRepository(collection *mongo.Collection, transformer transform.Transformer) *ClinkRepository {
	return &ClinkRepository{
		BaseRepository: reposh.NewBaseRepository[entity.ClinkDao, entity.Clink](collection, transformer),
	}
}
