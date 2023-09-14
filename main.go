package main

import (
	"PowBlockchain/blockchain"
	"PowBlockchain/network"
	"fmt"
)

const DBNAME = "blockchain.db"
const (
	AddBlock = iota + 1
	AddTx
	GetBlock
	GetLastHash
	GetBalance
	GetChainSize
)

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

	listener := network.Listen("localhost:8080", handleServer)

	if listener == nil {
		fmt.Println("Failed to start server")
		return
	}

	select {}
}

func handleServer(conn network.Conn, pack *network.Package) {
	network.Handle(AddBlock, conn, pack, addBlock)
	network.Handle(AddTx, conn, pack, addTransaction)
	network.Handle(GetBlock, conn, pack, getBlock)
	network.Handle(GetLastHash, conn, pack, getLastHash)
	network.Handle(GetBalance, conn, pack, getBalance)
	network.Handle(GetChainSize, conn, pack, getChainSize)
}

func addBlock(pack *network.Package) string {
	return ""
}

func addTransaction(pack *network.Package) string {
	sTx := pack.Data
	tx := blockchain.DeserializeTX(sTx)
	chain := blockchain.LoadChain(DBNAME)
	blockchain.NewBlock()
	fmt.Println()
	return ""
}

func getBlock(pack *network.Package) string {
	return ""
}

func getLastHash(pack *network.Package) string {
	return ""
}

func getBalance(pack *network.Package) string {
	return ""
}

func getChainSize(pack *network.Package) string {
	return ""
}
