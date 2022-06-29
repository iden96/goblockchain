package server

import (
	"encoding/json"
	"fmt"
	brs "goblockchain/blockchain_server/pkg/dto/blockchain_requests"
	"goblockchain/domain/blockchain"
	"goblockchain/domain/transaction"
	"goblockchain/domain/wallet"
	"goblockchain/wallet_server/utils"
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

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	failMessage, _ := utils.JsonStatus("fail")

	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()

		m, _ := json.Marshal(struct {
			Transactions []*transaction.Transaction `json:"transactions"`
			Length       int                        `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})

		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		t := brs.TransactionRequest{}
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(failMessage))
			return
		}

		if !t.Validate() {
			log.Print("ERROR: missing filed(s)")
			io.WriteString(w, string(failMessage))
			return
		}

		publicKey := wallet.PublicKeyFromString(*t.SenderPublicKey)
		signature := wallet.SignatureFromString(*t.Signature)
		bc := bcs.GetBlockchain()

		isCreated := bc.CreateTransaction(
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			*t.Value,
			publicKey,
			signature,
		)

		w.Header().Add("Content-Type", "application/json")

		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m, _ = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m, _ = utils.JsonStatus("success")
		}

		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	port := strconv.Itoa(int(bcs.Port()))
	host := fmt.Sprintf("0.0.0.0:%s", port)

	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)

	log.Fatal(http.ListenAndServe(host, nil))
}
