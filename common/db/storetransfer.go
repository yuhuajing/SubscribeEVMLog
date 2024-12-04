package db

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	"main/common/table"
	"strings"
)

// Insertion
func Insert(dba *gorm.DB, data types.Log) {
	res := dba.Model(&table.Transfer{}).Where("block_hash = ? AND block_index = ?", strings.ToLower(data.BlockHash.Hex()), data.Index).Find(&table.Transfer{})
	if res.RowsAffected == 0 {
		dba.Create(&table.Transfer{
			From:       data.Topics[1].Hex(),
			To:         data.Topics[2].Hex(),
			BlockHash:  strings.ToLower(data.BlockHash.Hex()),
			BlockIndex: data.Index,
		})
	}
}
