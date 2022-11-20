package res

import (
	"clinker-backend/internal/infrastructure/database/entity"
)

type VestigeRes struct {
	CreatedAt string       `json:"createdAt"`
	TxHash    string       `json:"txHash"`
	Parent    string       `json:"parent"`
	Head      string       `json:"head"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Hit       int64        `json:"hit"`
	Confirmed bool         `json:"confirmed"`
	User      *UserRes     `json:"user"`
	Appraisal AppraisalRes `json:"appraisal"`
	Children  []VestigeRes `json:"children"`
	Friends   []string     `json:"friends"`
}

func (r *VestigeRes) FromEntity(e *entity.Vestige) *VestigeRes {
	if e == nil {
		return nil
	}
	r.TxHash = e.TxHash
	r.Parent = e.Parent
	r.Head = e.Head
	r.Title = e.Title
	r.Content = e.Content
	r.Hit = e.Hit
	r.Confirmed = e.Confirmed
	r.User = new(UserRes).FromEntity(&e.User)
	r.Appraisal = *new(AppraisalRes).FromEntity(e.Appraisals)
	return r
}

type VestigesRes struct {
	TotalCount int64        `json:"totalCount"`
	Data       []VestigeRes `json:"data"`
}
