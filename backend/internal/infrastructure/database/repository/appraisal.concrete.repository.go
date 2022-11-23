package repository

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"

	"gorm.io/gorm"
)

type appraisalRepository struct {
	reposh.BaseRepository[entity.Appraisal]
}

func NewAppraisalRepository(db *gorm.DB) AppraisalRepository {
	return &appraisalRepository{
		BaseRepository: reposh.NewRepository[entity.Appraisal](db),
	}
}

func (r *appraisalRepository) FindByVestigeHead(head string) (*[]entity.Appraisal, error) {
	appraisals, err := r.Find(new(reposh.FindOption[entity.Appraisal]).Entity(
		&entity.Appraisal{
			VestigeId: head,
		},
	))
	return reposh.FilteredRecord(appraisals, err)
}

func (r *appraisalRepository) CountByUserId(userId string) (int64, error) {
	var count int64
	err := r.Model().Count(&count).Where("user_id = ?", userId).Error
	return count, err
}
