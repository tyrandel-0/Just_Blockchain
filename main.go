package main

import (
	"PowBlockchain/blockchain"
	"PowBlockchain/network"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const DBNAME = "blockchain.db"
const NODEADDR = "localhost:8080"

func main() {
	user := blockchain.LoadUser("MIIBOwIBAAJBAKNqVNRZgQ4nxf/3KbQkFTICWx74DJ0kjr/Juzu9Orx+Y/apKmided0Azsg6u2HolUN2P7YMe4rzqBnGWhdGWRECAwEAAQJAWQMoZerDA2Ti00RccQVejjj+TWYr6MTrBMjrteSjQ9xtuoik5cA85JMp3T+UOTokD+tQG3SNFETRi42aSDky4QIhANQcqBNbhywtICGiSxT7TRyoVMK+refwJSouApcMAtqVAiEAxTpEynk1F5/I77HSquJd0VBszPalyJDahs+HxKp74Y0CIQC//ZBEtTwMyGulBflf7HdH0TWncGCI59076Jl/juemYQIgdnlyKU52HiLVyWbAbfZc9Qei09y16a1aF/FCVVkz4WECIQDCov/RUtRytQ4C0MVmmLhjz8i7YzYbbiDWV1fPo7rPww==")
	fmt.Println("Your privateKey:" + blockchain.StringPrivate(user.Private()))
	fmt.Println("Your address:" + blockchain.StringPublic(user.Public()))
	for {
		fmt.Print("Enter tx: ")
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		splitedMsg := strings.Split(msg, " ")
		toAddr := splitedMsg[0]
		amount, _ := strconv.Atoi(splitedMsg[1])
		chain := blockchain.LoadChain(DBNAME)

		tx := blockchain.NewTransaction(user, chain.LastHash(), toAddr, uint64(amount))
		packge := network.Package{
			Option: 1,
			Data:   blockchain.SerializeTX(tx),
		}
		network.Send(NODEADDR, &packge)
	}
}
