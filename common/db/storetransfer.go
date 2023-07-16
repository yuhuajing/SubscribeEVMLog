package db

import (
	"fmt"
	"main/common/table"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
)

// Insertion
func Insert(dba *gorm.DB, txdata types.Log) {
	//res := dba.Model(&table.TxResult{}).Where("hash = ?", tx.Hash).First(&table.TxResult{})
	//res := dba.Table("tx_results").Where("hash = ?", tx.Hash).First(&table.TxResult{})
	//if res.RowsAffected == 0 {

	dba.Create(&table.Transfer{
		From:  "0x" + txdata.Topics[1].Hex()[26:],
		To:    "0x" + txdata.Topics[2].Hex()[26:],
		Value: ParseByteToUint(txdata.Data),
	})

	//}
}

func ParseByteToUint(b []byte) uint64 {
	str := ""
	for i := 0; i < len(b)/32; i++ {
		num := big.NewInt(0)
		num.SetBytes(b[i*32 : i*32+32])
		resInt := fmt.Sprintf("%d", num)
		str += resInt + ""
	}
	uint64Value, _ := strconv.ParseUint(str, 10, 64)

	return uint64Value
}
