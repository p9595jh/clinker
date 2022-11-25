package service

import (
	"clinker-backend/common/asyncer"
	"clinker-backend/common/enum/Authority"
	"clinker-backend/internal/domain/model/dto"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/p9595jh/fpgo"
	"github.com/p9595jh/transform"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
	processService *ProcessService
	authService    *AuthService
}

func NewUserService(
	userRepository repository.UserRepository,
	processService *ProcessService,
	authService *AuthService,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		processService: processService,
		authService:    authService,
	}
}

func (s *UserService) Initializer() {
	s.processService.transformer.RegisterTransformer("userE2R", transform.F2(func(u *entity.User, _ string) *res.UserRes {
		userRes := new(res.UserRes)
		s.processService.transformer.Mapping(u, userRes)
		return userRes
	}))
}

func (s *UserService) FindUsers(authority Authority.Type, page, take int) (*res.ProfuseRes[res.UserRes], *res.ErrorRes) {
	if authority != Authority.ADMIN {
		return nil, res.NewErrorfRes(fiber.StatusForbidden, "Forbidden")
	}

	ress, errs := asyncer.Multiple(
		func(good, bad *any) {
			*good, *bad = s.userRepository.Find(&reposh.FindOption[entity.User]{
				Order:  reposh.OrderBy{Column: "created_at", Desc: true},
				Limit:  take,
				Offset: take * page,
			})
		},
		func(good, bad *any) {
			var i int64
			*bad = s.userRepository.Model().Count(&i).Error
			*good = i
		},
	)

	for _, err := range errs {
		if err != nil {
			return nil, res.NewInternalErrorRes(err.(error))
		}
	}

	var (
		users = ress[0].(*[]entity.User)
		count = ress[1].(int64)
	)

	users, err := s.userRepository.Find(&reposh.FindOption[entity.User]{
		Order:   reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:   take,
		Offset:  take * page,
		Preload: []string{"Vestiges", "Appraisals"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if users == nil {
		return &res.ProfuseRes[res.UserRes]{TotalCount: 0, Data: make([]res.UserRes, 0)}, nil
	} else {
		return &res.ProfuseRes[res.UserRes]{
			TotalCount: count,
			Data: fpgo.Pipe[[]entity.User, []res.UserRes](
				*users,
				fpgo.Map(func(i int, u *entity.User) *res.UserRes {
					return new(res.UserRes).FromEntity(u)
				}),
			),
		}, nil
	}
}

func (s *UserService) FindOneUser(userId string) (*res.UserRes, *res.ErrorRes) {
	user, err := s.userRepository.FindOne(&reposh.FindOption[entity.User]{
		Where:   reposh.EntityParts[entity.User]{Entity: &entity.User{Id: userId}},
		Preload: []string{"Vestiges", "Appraisals"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if user == nil {
		return nil, res.NewErrorfRes(fiber.StatusNotFound, "user '%s' not found", userId)
	} else {
		return new(res.UserRes).FromEntity(user), nil
	}
}

func (s *UserService) FindByAddress(address string) (*res.UserRes, *res.ErrorRes) {
	user, err := s.userRepository.FindOne(&reposh.FindOption[entity.User]{
		Where:   reposh.EntityParts[entity.User]{Entity: &entity.User{Address: address}},
		Preload: []string{"Vestiges", "Appraisals"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if user == nil {
		return nil, res.NewErrorfRes(fiber.StatusNotFound, "address '%s' not found", address)
	} else {
		return new(res.UserRes).FromEntity(user), nil
	}
}

func (s *UserService) Stop(authority Authority.Type, userId, reason string, date time.Time) (*res.UserIdRes, *res.ErrorRes) {
	if authority != Authority.ADMIN {
		return nil, res.NewErrorfRes(fiber.StatusForbidden, "Forbidden")
	}

	err := s.userRepository.Update(
		&reposh.EntityParts[entity.User]{Entity: &entity.User{Id: userId}},
		&reposh.EntityParts[entity.User]{Entity: &entity.User{
			StopUntil:  date,
			StopReason: reason,
		}},
	).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.NewErrorfRes(fiber.StatusNotFound, "user '%s' not found", userId)
		} else {
			return nil, res.NewInternalErrorRes(err)
		}
	} else {
		return &res.UserIdRes{Id: userId}, nil
	}
}

func (s *UserService) checkDuplicate(field, value string) *res.ErrorRes {
	if prev, err := s.userRepository.FindOne(&reposh.FindOption[entity.User]{
		Select: []string{field},
		Where:  reposh.EntityParts[entity.User]{Entity: &entity.User{Id: value}},
	}); err != nil {
		return res.NewInternalErrorRes(err)
	} else if prev != nil {
		return res.NewErrorfRes(fiber.StatusConflict, "%s '%s' already exists", field, value)
	}
	return nil
}

func (s *UserService) Register(userDto *dto.UserDto) (*res.UserIdRes, *res.ErrorRes) {
	_, errs := asyncer.Multiple(
		func(good, bad *any) {
			*bad = s.checkDuplicate("id", userDto.Id)
		},
		func(good, bad *any) {
			*bad = s.checkDuplicate("address", userDto.Address)
		},
		func(good, bad *any) {
			*bad = s.checkDuplicate("nickname", userDto.Nickname)
		},
	)
	for _, err := range errs {
		if err != nil {
			return nil, err.(*res.ErrorRes)
		}
	}

	userEntity := &entity.User{
		Id:        userDto.Id,
		Password:  userDto.Password,
		Nickname:  userDto.Nickname,
		Address:   userDto.Address,
		StopUntil: time.Now(),
	}

	// temp
	userEntity.Confirmed = true

	newUser, err := s.authService.Insert(userEntity)
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else {
		return &res.UserIdRes{Id: newUser.Id}, nil
	}
}
