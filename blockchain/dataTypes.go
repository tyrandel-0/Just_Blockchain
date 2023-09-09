package blockchain

import (
	"bytes"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"errors"
)

type BlockChain struct {
	DB *sql.DB
}

func (chain *BlockChain) AddBlock(block *Block) {
	_, err := chain.DB.Exec("INSERT INTO BlockChain (Hash, Block) VALUES ($1, $2)",
		base64.StdEncoding.EncodeToString(block.CurrHash),
		SerializeBlock(block),
	)
	if err != nil {
		return
	}
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
		rows.Scan(&sblock)
		block = DeserializeBlock(sblock)
		if value, ok := block.Mapping[address]; ok {
			balance = value
			break
		}
	}
	return balance
}

func (chain *BlockChain) Size() uint64 {
	var size uint64
	row := chain.DB.QueryRow("SELECT Id FROM BlockChain ORDER BY Id DESC")
	row.Scan(&size)
	return size
}

type Block struct {
	Nonce        uint64
	Difficulty   uint8
	CurrHash     []byte
	PrevHash     []byte
	Transactions []Transaction
	Mapping      map[string]uint64
	Miner        string
	Signature    []byte
	TimeStamp    string
}

func (block *Block) AddTransaction(chain *BlockChain, tx *Transaction) error {
	if tx == nil {
		return errors.New("tx is null")
	}

	if tx.Value == 0 {
		return errors.New("tx value = 0")
	}

	if tx.Sender != StorageChain && len(block.Transactions) == TxsLimit {
		return errors.New("len tx = limit")
	}

	if tx.Sender != StorageChain && tx.Value > StartPercent && tx.ToStorage != StorageReward {
		return errors.New("storage reward pass")
	}

	if !bytes.Equal(tx.PrevBlock, chain.LastHash()) {
		return errors.New("prev block in tx /= last hash in chain")
	}

	var balanceInChain uint64
	balanceInTx := tx.Value + tx.ToStorage
	if value, ok := block.Mapping[tx.Sender]; ok {
		balanceInChain = value
	} else {
		balanceInChain = chain.Balance(tx.Sender, chain.Size())
	}
	if balanceInTx > balanceInChain {
		return errors.New("insufficient funds")
	}

	block.Mapping[tx.Sender] = balanceInChain - balanceInTx
	block.addBalance(chain, tx.Receiver, tx.Value)
	block.addBalance(chain, StorageChain, tx.ToStorage)
	block.Transactions = append(block.Transactions, *tx)
	return nil
}

func (block *Block) addBalance(chain *BlockChain, receiver string, value uint64) {
	var balanceInChain uint64
	if v, ok := block.Mapping[receiver]; ok {
		balanceInChain = v
	} else {
		balanceInChain = chain.Balance(receiver, chain.Size())
	}
	block.Mapping[receiver] = balanceInChain + value
}

type Transaction struct {
	RandBytes []byte
	PrevBlock []byte
	Sender    string
	Receiver  string
	Value     uint64
	ToStorage uint64
	CurrHash  []byte
	Signature []byte
}

func (tx *Transaction) hash() []byte {
	return HashSum(bytes.Join(
		[][]byte{
			tx.RandBytes,
			tx.PrevBlock,
			[]byte(tx.Sender),
			[]byte(tx.Receiver),
			ToBytes(tx.Value),
			ToBytes(tx.ToStorage),
		},
		[]byte{},
	))
}

func (tx *Transaction) sign(priv *rsa.PrivateKey) []byte {
	return Sign(priv, tx.CurrHash)
}

type User struct {
	PrivateKey *rsa.PrivateKey
}

func (user *User) Address() string {
	return StringPublic(user.Public())
}

func (user *User) Private() *rsa.PrivateKey {
	return user.PrivateKey
}

func (user *User) Public() *rsa.PublicKey {
	return &(user.PrivateKey).PublicKey
}
