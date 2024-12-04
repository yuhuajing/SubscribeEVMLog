package config

var (
	EthRpc = "wss://cool-muddy-butterfly.discover.quiknode.pro/0e41f42d5a7c9611f30ef800444bfcb93d3ae9a6/"
)

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
