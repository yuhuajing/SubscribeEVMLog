package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/common/table"
	"sync"

	"main/common/db"
	"main/common/etherclient"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	wg       sync.WaitGroup
	logsChan = make(chan types.Log, 0)
	client   *ethclient.Client
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error = %v", err)
	}
}

func main() {
	client = etherclient.ConnEth()

	contracts := []string{"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"}
	topics := []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}
	wg.Add(2)
	go FilterLogs(17704346, 17704349, contracts, topics)
	go subscribeLogs(contracts, topics)
	wg.Wait()

}

func subscribeLogs(addresses, topics []string) *[]types.Log {
	defer wg.Done()
	query := ethereum.FilterQuery{}

	for _, address := range addresses {
		query.Addresses = append(query.Addresses, common.HexToAddress(address))
	}
	top := make([]common.Hash, 0)
	for _, topic := range topics {
		top = append(top, common.HexToHash(topic))
	}
	query.Topics = append(query.Topics, top)
	events, err := client.SubscribeFilterLogs(context.Background(), query, logsChan)
	if err != nil {
		log.Fatalf("Subscribe Event error: %v", err)
	}

	dba := db.BuildConnect()
	dba.AutoMigrate(&table.Transfer{})
	for {
		select {
		case err := <-events.Err():
			fmt.Println(fmt.Errorf("parse Event error: %v", err))
			subscribeLogs(addresses, topics)
		case vLog := <-logsChan:
			//outputLog(vLog)
			db.Insert(dba, vLog)
		}
	}
}

func FilterLogs(startBlockHeight, latestBlockNum int64, addresses, topics []string) {
	defer wg.Done()
	i := startBlockHeight
	for i <= latestBlockNum {
		from := &big.Int{}
		from = from.SetInt64(startBlockHeight)
		i += 5000
		to := &big.Int{}
		if i > latestBlockNum {
			to = to.SetInt64(latestBlockNum)
		} else {
			to = to.SetInt64(i)
		}
		query := ethereum.FilterQuery{
			FromBlock: from,
			ToBlock:   to,
		}
		for _, address := range addresses {
			query.Addresses = append(query.Addresses, common.HexToAddress(address))
		}
		top := make([]common.Hash, 0)
		for _, topic := range topics {
			top = append(top, common.HexToHash(topic))
		}
		query.Topics = append(query.Topics, top)

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			checkError(errors.New(fmt.Sprintf("Error in filter logs :%v", err)))
		}

		for _, OpLog := range logs {
			logsChan <- OpLog
		}
	}
}

func MarshalJSON(l types.Log) (string, error) {
	var enc table.Transfer
	enc.From = "0x" + l.Topics[1].Hex()[26:]
	enc.To = "0x" + l.Topics[2].Hex()[26:]
	enc.BlockHash = l.BlockHash.Hex()
	enc.BlockIndex = l.Index

	byteLog, err := json.Marshal(&enc)
	if err != nil {
		return string(byteLog), err
	}
	return string(byteLog), nil
}
