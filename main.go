package main

import (
	"PowBlockchain/blockchain"
	"fmt"
)

const DBNAME = "blockchain.db"

func main() {
	miner := blockchain.NewUser()
	err := blockchain.NewChain(DBNAME, miner.Address())
	if err != nil {
		return
	}

	chain := blockchain.LoadChain(DBNAME)
	for i := 0; i < 3; i++ {
		block := blockchain.NewBlock(miner.Address(), chain.LastHash())
		tx1 := blockchain.NewTransaction(miner, chain.LastHash(), "aaa", 5)
		err := block.AddTransaction(chain, tx1)
		if err != nil {
			return
		}

		tx2 := blockchain.NewTransaction(miner, chain.LastHash(), "bbb", 3)
		err = block.AddTransaction(chain, tx2)
		if err != nil {
			return
		}

		err = block.Accept(chain, miner, make(chan bool))
		if err != nil {
			return
		}
		chain.AddBlock(block)
	}

	var sblock string
	rows, err := chain.DB.Query("SELECT Block FROM BlockChain")

	if err != nil {
		panic("error: query to db")
	}

	for rows.Next() {
		err := rows.Scan(&sblock)
		if err != nil {
			return
		}
		fmt.Println(sblock)
	}
}
