package repository

import (
	"layer/internal/infrastructure/database/entity"
	"layer/internal/infrastructure/database/repository/reposh"

	"github.com/p9595jh/transform"
	"go.mongodb.org/mongo-driver/mongo"
)

type TxnRepository struct {
	*reposh.BaseRepository[entity.TxnDao, entity.Txn]
}

func NewTxnRepository(collection *mongo.Collection, transformer transform.Transformer) *TxnRepository {
	return &TxnRepository{
		BaseRepository: reposh.NewBaseRepository[entity.TxnDao, entity.Txn](collection, transformer),
	}
}
