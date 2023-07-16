package etherclient

import (
	"log"
	"main/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnEth() *ethclient.Client {
	nclient, err := ethclient.Dial(config.EthRpc)
	if err != nil {
		log.Fatal(err)
	}
	return nclient
}
