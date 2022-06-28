package server

import (
	"fmt"
	"goblockchain/domain/blockchain"
	"goblockchain/domain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{
		port: port,
	}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]

	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(
			minersWallet.BlockchainAddress(),
			bcs.Port(),
		)
		cache["blockchain"] = bc
	}

	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		bc := bcs.GetBlockchain()

		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Run() {
	port := strconv.Itoa(int(bcs.Port()))
	host := fmt.Sprintf("0.0.0.0:%s", port)

	http.HandleFunc("/", bcs.GetChain)

	log.Fatal(http.ListenAndServe(host, nil))
}
