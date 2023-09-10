package blockchain

const (
	CreateTable = `CREATE TABLE BlockChain (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Hash VARCHAR(44) UNIQUE,
    Block TEXT
);`
	DEBUG         = true
	KeySize       = 512
	StorageChain  = "STORAGE-CHAIN"
	StorageValue  = 100
	StorageReward = 1
	GenesisBlock  = "GENESIS-BLOCK"
	GenesisReward = 100
	DIFFICULTY    = 20
	TxsLimit      = 2
	StartPercent  = 10
	RandBytes     = 32
)
