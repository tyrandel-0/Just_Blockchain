package main

import (
	"PowBlockchain/blockchain"
	"PowBlockchain/network"
	"fmt"
)

const DBNAME = "blockchain.db"
const TransactionOption = 1

func main() {
	miner := blockchain.NewUser()
	err := blockchain.NewChain(DBNAME, miner.Address())
	if err != nil {
		return
	}

	chain := blockchain.LoadChain(DBNAME)

	block := blockchain.NewBlock(miner.Address(), chain.LastHash())
	tx := blockchain.NewTransaction(miner, chain.LastHash(), "MEgCQQCjalTUWYEOJ8X/9ym0JBUyAlse+AydJI6/ybs7vTq8fmP2qSponXndAM7IOrth6JVDdj+2DHuK86gZxloXRlkRAgMBAAE=", 50)
	err = block.AddTransaction(chain, tx)
	if err != nil {
		return
	}

	err = block.Accept(chain, miner, make(chan bool))
	if err != nil {
		return
	}
	chain.AddBlock(block)

	listener := network.Listen("localhost:8080", func(conn network.Conn, pack *network.Package) {
		if network.Handle(TransactionOption, conn, pack, handleTransaction) {
			fmt.Println("Transaction handled successfully")
		} else {
			fmt.Println("Failed to handle transaction")
		}
	})

	if listener == nil {
		fmt.Println("Failed to start server")
		return
	}

	select {}
}

func handleTransaction(pack *network.Package) string {
	fmt.Printf("Received transaction with option %d and data: %s\n", pack.Option, pack.Data)
	return "Transaction processed successfully"
}
