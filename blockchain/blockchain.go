package blockchain

import (
	"database/sql"
	"os"
	"time"
)

func NewChain(filename, receiver string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	_, err = db.Exec(CreateTable)
	chain := &BlockChain{
		DB: db,
	}

	genesis := &Block{
		PrevHash:  []byte(GenesisBlock),
		Mapping:   make(map[string]uint64),
		Miner:     receiver,
		TimeStamp: time.Now().Format(time.RFC3339),
	}

	genesis.Mapping[StorageChain] = StorageValue
	genesis.Mapping[receiver] = GenesisReward
	genesis.CurrHash = genesis.hash()
	chain.AddBlock(genesis)
	return nil
}

func LoadChain(filename string) *BlockChain {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	chain := &BlockChain{
		DB: db,
	}

	return chain
}

func NewBlock(miner string, prevHash []byte) *Block {
	return &Block{
		Difficulty: Difficulty,
		PrevHash:   prevHash,
		Mapping:    make(map[string]uint64),
		Miner:      miner,
	}
}

func NewTransactions(user *User, lastHash []byte, to string, value uint64) *Transaction {
	tx := &Transaction{
		RandBytes: GenerateRandomBytes(RandBytes),
		PrevBlock: lastHash,
		Sender:    user.Address(),
		Receiver:  to,
		Value:     value,
	}

	if value > StartPercent {
		tx.ToStorage = StorageReward
	}
	tx.CurrHash = tx.hash()
	tx.Signature = tx.sign(user.Private())
	return tx
}
