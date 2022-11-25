package repository

import (
	"clinker-backend/common/util"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository/reposh"
	"sync"

	"gorm.io/gorm"
)

type vestigeRepository struct {
	reposh.BaseRepository[entity.Vestige]
}

func NewVestigeRepository(db *gorm.DB) VestigeRepository {
	return &vestigeRepository{
		BaseRepository: reposh.NewRepository[entity.Vestige](db),
	}
}

// return: Vestige[txHash].Children, Vestige[txHash].Children[:].Appraisals
func (r *vestigeRepository) FindAliveChildrenByTxHash(txHash string, page, take int) (*[]entity.Vestige, int64, error) {
	// find count of all children
	var (
		countWait sync.WaitGroup
		count     int64
		countErr  error
	)
	go func(count *int64) {
		defer countWait.Done()
		countWait.Add(1)
		countErr = r.Model().Count(count).Where("parent = ? AND next <> ?", txHash, util.DefaultTxHash).Error
	}(&count)

	// find all vestiges of which parent is `txHash`
	// to figure alive, 'next' must not be default (0)
	vestiges, err := r.Find(&reposh.FindOption[entity.Vestige]{
		Where: reposh.EntityParts[entity.Vestige]{
			Raw: reposh.Raw{
				"parent = ? AND next <> ?",
				txHash,
				util.DefaultTxHash,
			},
		},
		Order:  reposh.OrderBy{Column: "created_at", Desc: true},
		Limit:  take,
		Offset: take * page,
	})

	// filtering
	vestiges, err = reposh.FilteredRecord(vestiges, err)
	if err != nil {
		return nil, 0, err
	} else if vestiges == nil {
		return nil, 0, nil
	}

	// using a head of each element, find one vestige and link Appraisals
	// as async
	var wg sync.WaitGroup
	for i, vestige := range *vestiges {
		go func(i int, head string) {
			defer wg.Done()
			wg.Add(1)
			v, err := r.FindOne(&reposh.FindOption[entity.Vestige]{
				Select: []string{"tx_hash"},
				Where: reposh.EntityParts[entity.Vestige]{EntityFn: func(e *entity.Vestige) {
					e.TxHash = head
				}},
				Preload: []string{"Appraisals", "User"},
			})
			if err != nil {
				return
			} else {
				(*vestiges)[i].Appraisals = v.Appraisals
				(*vestiges)[i].User = v.User
			}
		}(i, vestige.Head)
	}
	wg.Wait()

	countWait.Wait()
	if countErr != nil {
		return nil, 0, countErr
	}

	// error already filtered so it can just be returned
	return vestiges, count, nil
}

func (r *vestigeRepository) CountAncestors() (int64, error) {
	var count int64
	err := r.Model().Count(&count).Where("parent = ?", util.DefaultTxHash).Error
	return count, err
}
