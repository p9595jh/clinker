package service

import (
	"clinker-backend/common/asyncer"
	"clinker-backend/common/enum/Authority"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
	"database/sql"
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
}

func NewUserService(
	userRepository repository.UserRepository,
	processService *ProcessService,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		processService: processService,
	}
}

func (s *UserService) Initializer() {
	s.processService.transformer.RegisterTransformer("userE2R", transform.F2(func(u *entity.User, _ string) *res.UserRes {
		userRes := new(res.UserRes)
		s.processService.transformer.Mapping(u, userRes)
		return userRes
	}))
}

func (s *UserService) FindUsers(authority Authority.Type, skip, take int) (*res.UsersRes, *res.ErrorRes) {
	if authority != Authority.ADMIN {
		return nil, res.NewErrorfRes(fiber.StatusForbidden, "Forbidden")
	}

	ress, errs := asyncer.Multiple(
		func(a *any, e *error) {
			*a, *e = s.userRepository.Find(&reposh.FindOption[entity.User]{
				Order:  reposh.OrderBy{Column: "created_at", Desc: true},
				Limit:  take,
				Offset: take * skip,
			})
		},
		func(a *any, e *error) {
			var i int64
			*e = s.userRepository.Model().Count(&i).Error
			*a = i
		},
	)

	for _, err := range errs {
		if err != nil {
			return nil, res.NewInternalErrorRes(err)
		}
	}

	var (
		users = ress[0].(*[]entity.User)
		count = ress[1].(int64)
	)

	users, err := s.userRepository.Find(&reposh.FindOption[entity.User]{
		Order:   reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:   take,
		Offset:  take * skip,
		Preload: []string{"Vestiges", "Appraisals"},
	})
	if err != nil {
		return nil, res.NewInternalErrorRes(err)
	} else if users == nil {
		return &res.UsersRes{TotalCount: 0, Data: make([]res.UserRes, 0)}, nil
	} else {
		return &res.UsersRes{
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

func (s *UserService) Stop(authority Authority.Type, userId, reason string, date time.Time) (*res.UserStopRes, *res.ErrorRes) {
	if authority != Authority.ADMIN {
		return nil, res.NewErrorfRes(fiber.StatusForbidden, "Forbidden")
	}

	err := s.userRepository.Update(
		&reposh.EntityParts[entity.User]{Entity: &entity.User{Id: userId}},
		&reposh.EntityParts[entity.User]{Entity: &entity.User{
			StopUntil:  sql.NullTime{Valid: true, Time: date},
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
		return &res.UserStopRes{Id: userId}, nil
	}
}
