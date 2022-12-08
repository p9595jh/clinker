package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClinkDao struct {
	ObjectId  *primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp int64               `bson:"timestamp" transform:"unix2time"`
	TxHash    string              `bson:"txHash"`
	Address   string              `bson:"address"`
}

type Clink struct {
	ObjectId  *primitive.ObjectID `json:"_id,omitempty"`
	Timestamp *time.Time          `json:"timestamp" transform:"time2unix"`
	TxHash    string              `json:"txHash"`
	Address   string              `json:"address"`
}
