package table

import (
	"github.com/jinzhu/gorm"
)

type Transfer struct {
	gorm.Model
	From  string `json:"from" gencodec:"required"`
	To    string `json:"yp" gencodec:"required"`
	Value uint64 `json:"blockNumber"`
}
