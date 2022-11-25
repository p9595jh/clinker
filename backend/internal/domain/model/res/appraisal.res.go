package res

import "clinker-backend/internal/infrastructure/database/entity"

type AppraisalRes struct {
	Count     int     `json:"count"`
	Appraisal float64 `json:"appraisal"`
}

func (r *AppraisalRes) FromEntity(e []entity.Appraisal) *AppraisalRes {
	if e == nil {
		return nil
	}
	sum := int64(0)
	for _, appraisal := range e {
		if appraisal.Confirmed {
			sum += appraisal.Value
		}
	}
	r.Count = len(e)
	r.Appraisal = float64(sum) / float64(r.Count)
	return r
}

type AppraisalSpecificRes struct {
	CreatedAt string      `json:"createdAt"`
	TxHash    string      `json:"txHash"`
	Value     int64       `json:"value"`
	Confirmed bool        `json:"confirmed"`
	Vestige   *VestigeRes `json:"vestige"`
	Next      *VestigeRes `json:"next"`
	User      *UserRes    `json:"user"`
}

func (r *AppraisalSpecificRes) FromEntity(e *entity.Appraisal) *AppraisalSpecificRes {
	if e == nil {
		return nil
	}
	r.TxHash = e.TxHash
	r.Value = e.Value
	r.Confirmed = e.Confirmed
	r.Vestige = new(VestigeRes).FromEntity(&e.Vestige)
	r.Next = new(VestigeRes).FromEntity(&e.Next)
	r.User = new(UserRes).FromEntity(&e.User)
	return r
}
