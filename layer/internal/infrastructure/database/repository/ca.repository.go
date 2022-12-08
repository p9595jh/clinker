package repository

import (
	"layer/internal/infrastructure/database/entity"
	"layer/internal/infrastructure/database/repository/reposh"

	"github.com/p9595jh/transform"
	"go.mongodb.org/mongo-driver/mongo"
)

type CaRepository struct {
	*reposh.BaseRepository[entity.CaDao, entity.Ca]
}

func NewCaRepository(collection *mongo.Collection, transformer transform.Transformer) *CaRepository {
	return &CaRepository{
		BaseRepository: reposh.NewBaseRepository[entity.CaDao, entity.Ca](collection, transformer),
	}
}
