package repository

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
)

type AppraisalRepository interface {
	reposh.BaseRepository[entity.Appraisal]
	FindByVestigeHead(head string) (*[]entity.Appraisal, error)
	CountByUserId(userId string) (int64, error)
}
