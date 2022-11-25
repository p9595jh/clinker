package service

import (
	"clinker-backend/common/asyncer"
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
	processService *ProcessService,
) *VestigeService {
	return &VestigeService{
		vestigeRepository:   vestigeRepository,
		apprasialRepository: apprasialRepository,
		userRepository:      userRepository,
		processService:      processService,
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
// func (s *VestigeService) FindAncestors(page, take int) (*res.VestigesRes, *res.ErrorRes) {
func (s *VestigeService) FindAncestors(page, take int) (*res.ProfuseRes[res.VestigeRes], *res.ErrorRes) {
	ress, errs := asyncer.Multiple(
		func(good, bad *any) {
			vestiges, err := s.vestigeRepository.Find(&reposh.FindOption[entity.Vestige]{
				Preload: []string{"Appraisals", "User"},
				Order:   reposh.OrderBy{Column: "created_at", Desc: true},
				Limit:   take,
				Offset:  take * page,
				Where: reposh.EntityParts[entity.Vestige]{Entity: &entity.Vestige{
					Parent: util.DefaultTxHash,
				}},
			})
			if err != nil {
				*bad = res.NewInternalErrorRes(err)
			} else if vestiges == nil {
				*good = &res.ProfuseRes[res.VestigeRes]{TotalCount: 0, Data: make([]res.VestigeRes, 0)}
			} else {
				*good = *vestiges
			}
		},
		func(good, bad *any) {
			if count, err := s.vestigeRepository.CountAncestors(); err != nil {
				*bad = res.NewInternalErrorRes(err)
			} else {
				*good = count
			}
		},
	)

	for _, err := range errs {
		if err != nil {
			return nil, err.(*res.ErrorRes)
		}
	}

	var (
		vestiges = ress[0]
		count    = ress[1].(int64)
	)

	if vestigesRes, ok := vestiges.(*res.ProfuseRes[res.VestigeRes]); ok {
		return vestigesRes, nil
	} else {
		return &res.ProfuseRes[res.VestigeRes]{
			TotalCount: count,
			Data: fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
				vestiges.([]entity.Vestige),
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
func (s *VestigeService) FindFriendsByHead(head string) ([]res.VestigeRes, *res.ErrorRes) {
	vestiges, err := s.vestigeRepository.Find(&reposh.FindOption[entity.Vestige]{
		Preload: []string{"Appraisals", "User"},
		Where: reposh.EntityParts[entity.Vestige]{Entity: &entity.Vestige{
			Head: head,
		}},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if vestiges == nil {
		return make([]res.VestigeRes, 0), nil
	} else {
		return fpgo.Pipe[[]entity.Vestige, []res.VestigeRes](
			*vestiges,
			fpgo.Map(func(i int, v *entity.Vestige) *res.VestigeRes {
				return new(res.VestigeRes).FromEntity(v)
			}),
		), nil
	}
}

func (s *VestigeService) FindChildren(txHash string, page, take int) (*res.ProfuseRes[res.VestigeRes], *res.ErrorRes) {
	children, count, err := s.vestigeRepository.FindAliveChildrenByTxHash(txHash, page, take)
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if children == nil {
		return &res.ProfuseRes[res.VestigeRes]{TotalCount: 0, Data: make([]res.VestigeRes, 0)}, nil
	} else {
		return &res.ProfuseRes[res.VestigeRes]{
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

func (s *VestigeService) Save(userId string, vestige *dto.VestigeDto) (*res.SaveTxHashRes, *res.ErrorRes) {
	vestigeEntity := &entity.Vestige{
		Title:   vestige.Title,
		Content: vestige.Content,
		Parent:  vestige.Parent,
		Head:    vestige.Head,
	}

	for _, hash := range [][]string{{"parent", vestige.Parent}, {"head", vestige.Head}} {
		if hash[1] != "" {
			if v, err := s.vestigeRepository.FindOne(&reposh.FindOption[entity.Vestige]{
				Select: []string{"tx_hash"},
				Where: reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) {
					e.TxHash = hash[1]
				}},
			}); err != nil {
				return nil, res.NewInternalErrorRes(err)
			} else if v == nil {
				return nil, res.NewErrorfRes(fiber.StatusNotFound, "%s '%s' not found", hash[0], hash[1])
			}
		}
	}

	// temp
	vestigeEntity.TxHash = util.RandHex(64)
	vestigeEntity.Confirmed = true

	if newVestige, err := s.vestigeRepository.Save(vestigeEntity); err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else {
		return &res.SaveTxHashRes{TxHash: newVestige.TxHash}, nil
	}
}
