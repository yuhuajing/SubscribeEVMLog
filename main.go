package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"main/common/db"
	"main/common/etherclient"
	"main/common/table"
	"main/config"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	nclient := etherclient.ConnEth()
	dba := db.Buildconnect()
	dba.AutoMigrate(&table.Transfer{})

	logs := historyLogs(nclient, config.Historyfileter)
	parseLogs(logs)

	for _, log := range *logs {
		db.Insert(dba, log)
	}

	//subscribeLogs(dba, nclient, config.Pendingfileter)

}

func subscribeLogs(dba *gorm.DB, eclient *ethclient.Client, filter config.Filter) *[]types.Log {
	contractAddress := []common.Address{
		common.HexToAddress(filter.Address),
	}
	query := ethereum.FilterQuery{
		Addresses: contractAddress,
		Topics:    [][]common.Hash{{common.HexToHash(filter.Topic)}},
	}
	logs := make(chan types.Log)

	cominglogs, err := eclient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-cominglogs.Err():
			log.Fatal(err)
		case vLog := <-logs:
			//outputLog(vLog)
			db.Insert(dba, vLog)
		}
	}
}

func historyLogs(eclient *ethclient.Client, filter config.Filter) *[]types.Log {
	contractAddress := []common.Address{
		common.HexToAddress(filter.Address),
	}
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(filter.From),
		ToBlock:   big.NewInt(filter.To),
		Addresses: contractAddress,
		Topics:    [][]common.Hash{{common.HexToHash(filter.Topic)}},
	}

	logs, err := eclient.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	return &logs

}
func parseLogs(logs *[]types.Log) {
	for _, vLog := range *logs {
		outputLog(vLog)
	}
}

func outputLog(vLog types.Log) {
	byteLog, _ := MarshalJSON(vLog)
	fmt.Println(string(byteLog))
}

func MarshalJSON(l types.Log) ([]byte, error) {
	type Log struct {
		From  string `json:"from" gencodec:"required"`
		To    string `json:"to" gencodec:"required"`
		Value uint64 `json:"value" gencodec:"required"`
	}
	var enc Log

	enc.From = "0x" + l.Topics[1].Hex()[26:]
	enc.To = "0x" + l.Topics[2].Hex()[26:]
	enc.Value = db.ParseByteToUint(l.Data)

	return json.Marshal(&enc)
}

func StartsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}
