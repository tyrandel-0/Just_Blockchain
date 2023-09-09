package blockchain

const (
	CreateTable = `
	CREATE TABLE BlockChain (
		Id INTEGER PRIMARY KEY AUTOINCREMENT, 
		Hash VARCHAR(44) UNIQUE,
		Block TEXT
	);`
	GenesisBlock  = "GENESIS BLOCK"
	StorageValue  = 100
	GenesisReward = 100
	StorageChain  = "STORAGE CHAIN"
	Difficulty    = 20
	RandBytes     = 32
	StartPercent  = 10
	StorageReward = 1
	TxsLimit      = 2
)
