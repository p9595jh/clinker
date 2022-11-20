package repository

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
)

type VestigeRepository interface {
	reposh.BaseRepository[entity.Vestige]
	FindAliveChildrenByTxHash(txHash string, skip, take int) (*[]entity.Vestige, int64, error)
	CountOrphans() (int64, error)
}
