package entity

import (
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TxnDao struct {
	ObjectId  *primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp int64               `bson:"timestamp,omitempty" transform:"unix2time"`
	Hash      string              `bson:"hash,omitempty"`
	GasPrice  string              `bson:"gasPrice,omitempty" transform:"str2big"`
	GasUsed   string              `bson:"gasUsed,omitempty" transform:"str2big"`
	Fee       string              `bson:"fee,omitempty" transform:"str2big"`
}

type Txn struct {
	ObjectId  *primitive.ObjectID `json:"_id,omitempty"`
	Timestamp *time.Time          `json:"timestamp,omitempty" transform:"time2unix"`
	Hash      string              `json:"hash,omitempty" transform:"lower,trimPre:0x"`
	GasPrice  *big.Int            `json:"gasPrice,omitempty" transform:"big2str"`
	GasUsed   *big.Int            `json:"gasUsed,omitempty" transform:"big2str"`
	Fee       *big.Int            `json:"fee,omitempty" transform:"big2str"`
}
