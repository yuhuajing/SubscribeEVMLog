package config

var (
	EthRpc = "wss://cool-muddy-butterfly.discover.quiknode.pro/0e41f42d5a7c9611f30ef800444bfcb93d3ae9a6/"
)

type Filter struct {
	Address string
	From    int64
	To      int64
	Topic   string
}

var Historyfileter = Filter{
	Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
	From:    17704346,
	To:      17704349,
	Topic:   "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
}
var Pendingfileter Filter = Filter{
	Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
	Topic:   "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
}

type MysqlConFig struct {
	Addr            string
	Port            int
	Db              string
	Username        string
	Password        string
	MaxIdealConn    int
	MaxOpenConn     int
	ConnMaxLifetime int
}

var MysqlCon = MysqlConFig{
	"127.0.0.1",
	3306,
	"eventLog",
	"root",
	"123456",
	10,
	256,
	600,
}
