package dao

import "time"

type Common struct {
	Create int64 `bson:"create"` // millisecond
	Update int64 `bson:"update"` // millisecond
	Valid  bool  `bson:"valid"`
}

func newCommon() *Common {
	t := time.Now().UnixNano() / 1e6
	return &Common{
		Create: t,
		Update: t,
		Valid:  true,
	}
}
