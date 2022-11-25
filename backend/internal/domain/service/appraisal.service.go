package service

import (
	"clinker-backend/common/util"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"

	"github.com/gofiber/fiber/v2"
	"github.com/p9595jh/fpgo"
	"github.com/p9595jh/transform"
)

type AppraisalService struct {
	appraisalRepository repository.AppraisalRepository
	vestigeRepository   repository.VestigeRepository
	processService      *ProcessService
}

func NewAppraisalService(
	appraisalRepository repository.AppraisalRepository,
	vestigeRepository repository.VestigeRepository,
	processService *ProcessService,
) *AppraisalService {
	return &AppraisalService{
		appraisalRepository: appraisalRepository,
		vestigeRepository:   vestigeRepository,
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

func (s *AppraisalService) FindByUserId(page, take int, userId string) (*res.ProfuseRes[res.AppraisalSpecificRes], *res.ErrorRes) {
	appraisals, err := s.appraisalRepository.Find(&reposh.FindOption[entity.Appraisal]{
		Order:  reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:  take,
		Offset: take * page,
		Where: reposh.EntityParts[entity.Appraisal]{
			Entity: &entity.Appraisal{UserId: userId},
		},
		Preload: []string{"Vestige", "Next"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if appraisals == nil {
		return &res.ProfuseRes[res.AppraisalSpecificRes]{TotalCount: 0, Data: make([]res.AppraisalSpecificRes, 0)}, nil
	} else {
		count, err := s.appraisalRepository.CountByUserId(userId)
		if err != nil {
			return nil, res.NewInternalErrorRes(err)
		}
		return &res.ProfuseRes[res.AppraisalSpecificRes]{
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

func (s *AppraisalService) Save(userId string, appraisal *dto.AppraisalDto) (*res.SaveTxHashRes, *res.ErrorRes) {
	appraisalEntity := &entity.Appraisal{
		Value:  appraisal.Value,
		NextId: appraisal.NextId,
		UserId: userId,
	}

	vestige, err := s.vestigeRepository.FindOne(&reposh.FindOption[entity.Vestige]{
		Select: []string{"tx_hash", "head"},
		Where: reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) {
			e.TxHash = appraisal.NextId
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else {
		appraisalEntity.VestigeId = vestige.Head
	}

	prev, err := s.appraisalRepository.FindOne(&reposh.FindOption[entity.Appraisal]{
		Select: []string{"vestige_id", "user_id", "tx_hash"},
		Where: reposh.EntityParts[entity.Appraisal]{Entity: &entity.Appraisal{
			VestigeId: vestige.Head,
			UserId:    userId,
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if prev != nil {
		return nil, res.NewErrorfRes(fiber.StatusConflict, "already did")
	}

	// temp
	appraisalEntity.TxHash = util.RandHex(64)
	appraisalEntity.Confirmed = true

	if newAppraisal, err := s.appraisalRepository.Save(appraisalEntity); err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else {
		return &res.SaveTxHashRes{TxHash: newAppraisal.TxHash}, nil
	}
}
