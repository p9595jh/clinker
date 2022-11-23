package service

import (
	"clinker-backend/common/util"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"

	"github.com/gofiber/fiber/v2"
	"github.com/p9595jh/fpgo"
	"github.com/p9595jh/transform"
)

type VestigeService struct {
	vestigeRepository   repository.VestigeRepository
	apprasialRepository repository.AppraisalRepository
	userRepository      repository.UserRepository
	processService      *ProcessService
}

func NewVestigeService(
	vestigeRepository repository.VestigeRepository,
	apprasialRepository repository.AppraisalRepository,
	userRepository repository.UserRepository,
) *VestigeService {
	return &VestigeService{
		vestigeRepository:   vestigeRepository,
		apprasialRepository: apprasialRepository,
		userRepository:      userRepository,
	}
}

func (s *VestigeService) Initializer() {
	s.processService.transformer.RegisterTransformer("vestigesE2R", transform.F2(func(v []entity.Vestige, s string) []res.VestigeRes {
		return fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
			v,
			fpgo.Map(func(i int, v *entity.Vestige) *res.VestigeRes {
				return new(res.VestigeRes).FromEntity(v)
			}),
		)
	}))
}

// will be shown in the main page
func (s *VestigeService) FindOrphans(skip, take int) (*res.VestigesRes, *res.ErrorRes) {
	vestiges, err := s.vestigeRepository.Find(&reposh.FindOption[entity.Vestige]{
		Preload: []string{"Appraisals", "User"},
		Order:   reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:   take,
		Offset:  take * skip,
		Where: reposh.EntityParts[entity.Vestige]{Entity: &entity.Vestige{
			Parent: util.DefaultTxHash,
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if vestiges == nil {
		return &res.VestigesRes{TotalCount: 0, Data: make([]res.VestigeRes, 0)}, nil
	} else {
		count, err := s.vestigeRepository.CountOrphans()
		if err != nil {
			return nil, res.NewInternalErrorRes(err)
		}
		return &res.VestigesRes{
			TotalCount: count,
			Data: fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
				*vestiges,
				fpgo.Map(func(i int, v *entity.Vestige) *res.VestigeRes {
					return new(res.VestigeRes).FromEntity(v)
				}),
			),
		}, nil
	}
}

// find one with its appraisals and user
func (s *VestigeService) FindOneByTxHash(txHash string) (*res.VestigeRes, *res.ErrorRes) {
	vestige, err := s.vestigeRepository.FindOne(&reposh.FindOption[entity.Vestige]{
		Preload: []string{"Appraisals", "User"},
		Where: reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) {
			e.TxHash = txHash
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if vestige == nil {
		return nil, res.NewErrorfRes(fiber.StatusNotFound, "vestige '%s' not found", txHash)
	} else {
		return new(res.VestigeRes).FromEntity(vestige), nil
	}
}

// parameter is `head`
func (s *VestigeService) FindFriendsByHead(head string) (*res.VestigesRes, *res.ErrorRes) {
	vestiges, err := s.vestigeRepository.Find(&reposh.FindOption[entity.Vestige]{
		Preload: []string{"Appraisals", "User"},
		Where: reposh.EntityParts[entity.Vestige]{Entity: &entity.Vestige{
			Head: head,
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if vestiges == nil {
		return &res.VestigesRes{TotalCount: 0, Data: []res.VestigeRes{}}, nil
	} else {
		return &res.VestigesRes{
			TotalCount: int64(len(*vestiges)),
			Data: fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
				*vestiges,
				fpgo.Map(func(i int, v *entity.Vestige) *res.VestigeRes {
					return new(res.VestigeRes).FromEntity(v)
				}),
			),
		}, nil
	}
}

func (s *VestigeService) FindChildren(txHash string, skip, take int) (*res.VestigesRes, *res.ErrorRes) {
	children, count, err := s.vestigeRepository.FindAliveChildrenByTxHash(txHash, skip, take)
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if children == nil {
		return &res.VestigesRes{TotalCount: 0, Data: make([]res.VestigeRes, 0)}, nil
	} else {
		return &res.VestigesRes{
			TotalCount: count,
			Data: fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
				*children,
				fpgo.Map(func(i int, v *entity.Vestige) *res.VestigeRes {
					return new(res.VestigeRes).FromEntity(v)
				}),
			),
		}, nil
	}
}

func (s *VestigeService) Save(vestige *entity.Vestige) {
	// s.re
}
