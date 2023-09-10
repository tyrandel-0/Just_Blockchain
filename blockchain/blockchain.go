package blockchain

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type BlockChain struct {
	DB *sql.DB
}

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
	chain := &BlockChain{
		DB: db,
	}
	return chain
}

func NewBlock(miner string, prevHash []byte) *Block {
	return &Block{
		Difficulty: DIFFICULTY,
		PrevHash:   prevHash,
		Miner:      miner,
		Mapping:    make(map[string]uint64),
	}
}

func (chain *BlockChain) Size() uint64 {
	var size uint64
	row := chain.DB.QueryRow("SELECT Id FROM BlockChain ORDER BY Id DESC")
	err := row.Scan(&size)
	if err != nil {
		return 0
	}
	return size
}

func (chain *BlockChain) Balance(address string, size uint64) uint64 {
	var (
		sblock  string
		block   *Block
		balance uint64
	)
	rows, err := chain.DB.Query("SELECT Block FROM BlockChain WHERE Id <= $1 ORDER BY Id DESC", size)
	if err != nil {
		return balance
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		err := rows.Scan(&sblock)
		if err != nil {
			return 0
		}
		block = DeserializeBlock(sblock)
		if value, ok := block.Mapping[address]; ok {
			balance = value
			break
		}
	}
	return balance
}

func (chain *BlockChain) LastHash() []byte {
	var hash string
	row := chain.DB.QueryRow("SELECT Hash FROM BlockChain ORDER BY Id DESC")
	err := row.Scan(&hash)
	if err != nil {
		return nil
	}
	return Base64Decode(hash)
}

func (chain *BlockChain) AddBlock(block *Block) {
	_, err := chain.DB.Exec("INSERT INTO BlockChain (Hash, Block) VALUES ($1, $2)",
		Base64Encode(block.CurrHash),
		SerializeBlock(block),
	)
	if err != nil {
		return
	}
}
