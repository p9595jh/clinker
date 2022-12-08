package entity

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaDao struct {
	ObjectId  *primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp int64               `bson:"timestamp" transform:"unix2time"`
	Name      string              `bson:"name"`
	Address   string              `bson:"address" transform:"str2addr"`
}

type Ca struct {
	ObjectId  *primitive.ObjectID `json:"_id,omitempty"`
	Timestamp *time.Time          `json:"timestamp" transform:"time2unix"`
	Name      string              `json:"name"`
	Address   common.Address      `json:"address" transform:"addr2str"`
}
