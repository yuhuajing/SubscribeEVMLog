package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type node struct {
	val  [1024 * 1024]bool
	next *node
}

func finalizedNode() {
	fmt.Println("start")
	printAlloc()
	a := &node{val: [1024 * 1024]bool{true}}
	b := &node{val: [1024 * 1024]bool{false}}
	a.next = b
	b.next = a
	runtime.SetFinalizer(a, func(obj *node) {
		fmt.Println("a is finalized")
	})
	runtime.SetFinalizer(b, func(obj *node) {
		fmt.Println("b is finalized")
	})
	fmt.Println("process")
	printAlloc()
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}

func main() {
	// fmt.Println("start")
	// printAlloc()
	//finalizedNode()
	// time.Sleep(1 * time.Second)
	// runtime.GC()
	// time.Sleep(1 * time.Second)
	// runtime.GC()
	// time.Sleep(1 * time.Second)
	// runtime.GC()
	// time.Sleep(1 * time.Second)
	// runtime.GC()
	// time.Sleep(1 * time.Second)
	// runtime.GC()
	// fmt.Println("end")
	// printAlloc()
	ethServer := "https://cloudflare-eth.com"
	nclient, err := ethclient.Dial(ethServer)
	//https://cool-muddy-butterfly.discover.quiknode.pro/0e41f42d5a7c9611f30ef800444bfcb93d3ae9a6/
	if err != nil {
		log.Fatal(err)
	}
	receipt, _ := nclient.TransactionReceipt(context.Background(), common.HexToHash("0x300ad2476edb1657926cee7744aecc5c375e70fac7506af0612dcc66d8f8c342"))

	for _, log := range receipt.Logs {
		// //fmt.Println(log.Topics[0].Hex())
		// if StartsWith(log.Topics[0].Hex(), "0x7fcf532c") {
		// 	fmt.Println("123")
		// }
		byteLog, _ := MarshalJSON(*log)
		fmt.Println(string(byteLog))

	}
	// contractAddress := common.HexToAddress("0xa6be9422d48DE5C75aD4f48f22cAc4f5F60f0B98")
	// query := ethereum.FilterQuery{
	// 	FromBlock: big.NewInt(17412598),
	// 	ToBlock:   big.NewInt(17412600),
	// 	Addresses: []common.Address{
	// 		contractAddress,
	// 	},
	// }
	// logs, err := nclient.FilterLogs(context.Background(), query)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, vLog := range logs {

	// 	byteLog, _ := MarshalJSON(vLog)
	// 	fmt.Println(string(byteLog))
	// }
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

func MarshalJSON(l types.Log) ([]byte, error) {
	type Log struct {
		Address     common.Address `json:"address" gencodec:"required"`
		Topics      []common.Hash  `json:"topics" gencodec:"required"`
		Data        string         `json:"data" gencodec:"required"` //hexutil.Bytes  `json:"data" gencodec:"required"`
		BlockNumber uint64         `json:"blockNumber"`              //hexutil.Uint64 `json:"blockNumber"`
		TxHash      common.Hash    `json:"transactionHash" gencodec:"required"`
		TxIndex     uint           `json:"transactionIndex"` // hexutil.Uint   `json:"transactionIndex"`
		BlockHash   common.Hash    `json:"blockHash"`
		Index       uint           `json:"logIndex"` //  hexutil.Uint   `json:"logIndex"`
		Removed     bool           `json:"removed"`
	}
	var enc Log
	enc.Address = l.Address
	enc.Topics = l.Topics
	enc.Data = parseByteToString(l.Data) //l.Data
	enc.BlockNumber = l.BlockNumber      //hexutil.Uint64(l.BlockNumber)
	enc.TxHash = l.TxHash
	enc.TxIndex = l.TxIndex //hexutil.Uint(l.TxIndex)
	enc.BlockHash = l.BlockHash
	enc.Index = l.Index //hexutil.Uint(l.Index)
	enc.Removed = l.Removed
	return json.Marshal(&enc)
}

func parseByteToString(b []byte) string {
	str := ""
	//fmt.Println(len(b))
	for i := 0; i < len(b)/32; i++ {
		num := big.NewInt(0)
		num.SetBytes(b[i*32 : i*32+32])
		resInt := fmt.Sprintf("%d", num)
		str += "data" + strconv.Itoa(i) + ": " + resInt + " "
		//num, err := strconv.ParseInt(_b[i:i+64], 16, 32)
	}
	return str
}
