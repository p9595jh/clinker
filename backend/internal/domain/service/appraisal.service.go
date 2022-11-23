package service

import (
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"

	"github.com/p9595jh/fpgo"
	"github.com/p9595jh/transform"
)

type AppraisalService struct {
	appraisalRepository repository.AppraisalRepository
	processService      *ProcessService
}

func NewAppraisalService(
	appraisalRepository repository.AppraisalRepository,
	processService *ProcessService,
) *AppraisalService {
	return &AppraisalService{
		appraisalRepository: appraisalRepository,
		processService:      processService,
	}
}

func (s *AppraisalService) Initializer() {
	s.processService.transformer.RegisterTransformer("appraisalE2R", transform.F2(func(a *entity.Appraisal, _ string) *res.AppraisalRes {
		appraisalRes := new(res.AppraisalRes)
		s.processService.transformer.Mapping(a, appraisalRes)
		return appraisalRes
	}))

	s.processService.transformer.RegisterTransformer("appraisalsE2R", transform.F2(func(a []entity.Appraisal, s string) []res.AppraisalSpecificRes {
		return fpgo.Pipe[[]entity.Appraisal, []res.AppraisalSpecificRes](
			a,
			fpgo.Map(func(i int, a *entity.Appraisal) *res.AppraisalSpecificRes {
				return new(res.AppraisalSpecificRes).FromEntity(a)
			}),
		)
	}))
}

func (s *AppraisalService) FindByVestigeHead(head string) (*res.AppraisalRes, *res.ErrorRes) {
	appraisals, err := s.appraisalRepository.FindByVestigeHead(head)
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if appraisals == nil {
		return new(res.AppraisalRes), nil
	} else {
		return new(res.AppraisalRes).FromEntity(*appraisals), nil
	}
}

func (s *AppraisalService) FindByUserId(skip, take int, userId string) (*res.AppraisalSpecificsRes, *res.ErrorRes) {
	appraisals, err := s.appraisalRepository.Find(&reposh.FindOption[entity.Appraisal]{
		Order:  reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:  take,
		Offset: take * skip,
		Where: reposh.EntityParts[entity.Appraisal]{
			Entity: &entity.Appraisal{UserId: userId},
		},
		Preload: []string{"Vestige", "Next"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if appraisals == nil {
		return &res.AppraisalSpecificsRes{TotalCount: 0, Data: make([]res.AppraisalSpecificRes, 0)}, nil
	} else {
		count, err := s.appraisalRepository.CountByUserId(userId)
		if err != nil {
			return nil, res.NewInternalErrorRes(err)
		}
		return &res.AppraisalSpecificsRes{
			TotalCount: count,
			Data: fpgo.Pipe[[]entity.Appraisal, []res.AppraisalSpecificRes](
				*appraisals,
				fpgo.Map(func(i int, a *entity.Appraisal) *res.AppraisalSpecificRes {
					return new(res.AppraisalSpecificRes).FromEntity(a)
				}),
			),
		}, nil
	}
}

func (s *AppraisalService) FindValidAppraisal(head string) (*res.AppraisalRes, *res.ErrorRes) {
	appraisalRes := new(res.AppraisalRes)
	err := s.appraisalRepository.Model().
		Select("SUM(value) AS Appraisal, COUNT(*) AS Count").
		Where("vestige_id = ? AND confirmed = ?", head, true).
		Take(&appraisalRes).
		Error
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if appraisalRes == nil {
		return &res.AppraisalRes{}, nil
	} else {
		return appraisalRes, nil
	}
}
