package res

import (
	"clinker-backend/common/util"
	"clinker-backend/internal/infrastructure/database/entity"

	"github.com/p9595jh/fpgo"
)

type UserRes struct {
	CreatedAt  string                 `json:"createdAt"`
	Id         string                 `json:"id"`
	Nickname   string                 `json:"nickname"`
	Address    string                 `json:"address"`
	StopUntil  string                 `json:"stopUntil"`
	Vestiges   []VestigeRes           `json:"vestiges"`
	Appraisals []AppraisalSpecificRes `json:"appraisals"`
}

func (r *UserRes) FromEntity(e *entity.User) *UserRes {
	if e == nil {
		return nil
	}
	r.CreatedAt = e.CreatedAt.Format(util.DateFormat)
	r.Id = e.Id
	r.Nickname = e.Nickname
	r.Address = e.Address
	r.Vestiges = fpgo.Pipe[[]entity.Vestige, []VestigeRes](
		e.Vestiges,
		fpgo.Map(func(i int, v *entity.Vestige) *VestigeRes {
			return new(VestigeRes).FromEntity(v)
		}),
	)
	r.Appraisals = fpgo.Pipe[[]entity.Appraisal, []AppraisalSpecificRes](
		e.Appraisals,
		fpgo.Map(func(i int, v *entity.Vestige) *VestigeRes {
			return new(VestigeRes).FromEntity(v)
		}),
	)
	return r
}

type UserIdRes struct {
	Id string `json:"id"`
}
