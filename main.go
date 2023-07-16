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
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	nclient := etherclient.ConnEth()
	dba := db.Buildconnect()

	logs := historyLogs(nclient, config.Historyfileter)
	parseLogs(logs)

	dba.AutoMigrate(&table.Transfer{})
	db.Insert(dba, logs)

	// var pendingfileter Filter = Filter{
	// 	address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
	// 	topic:   "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	// }
	// subscribeLogs(nclient, pendingfileter)

}

func subscribeLogs(eclient *ethclient.Client, filter config.Filter) {
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
			outputLog(vLog)
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
