package etherclient

import (
	"fmt"
	"log"
	"main/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnEth() *ethclient.Client {
	nclient, err := ethclient.Dial(config.EthRpc)
	if err != nil {
		fmt.Printf("Eth connect error:%s\n", err)
		log.Fatal(err)
	}
	return nclient
}
