package table

import (
	"github.com/jinzhu/gorm"
)

type Transfer struct {
	gorm.Model
	From       string `json:"from" gencodec:"required"`
	To         string `json:"to" gencodec:"required"`
	BlockHash  string `json:"blockHash"`
	BlockIndex uint   `json:"blockIndex"`
}
