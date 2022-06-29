package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	brs "goblockchain/blockchain_server/pkg/dto/blockchain_requests"
	"goblockchain/domain/wallet"
	wrs "goblockchain/wallet_server/pkg/dto/wallet_requests"
	"goblockchain/wallet_server/utils"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const (
	TEMP_DIR   = "wallet_server/templates"
	INDEX_HTML = "index.html"
)

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
	}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(TEMP_DIR, INDEX_HTML))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	failMessage, _ := utils.JsonStatus("fail")

	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		t := wrs.TransactionRequest{}
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(failMessage))

			return
		}

		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(failMessage))

			return
		}

		publicKey := wallet.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := wallet.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)

		if err != nil {
			log.Println("ERROR: parse error")
			io.WriteString(w, string(failMessage))

			return
		}

		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(
			privateKey,
			publicKey,
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			value32,
		)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &brs.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}

		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)
		url := fmt.Sprintf("%s/transactions", ws.Gateway())

		res, _ := http.Post(
			url,
			"application/json",
			buf,
		)
		if res.StatusCode == 201 {
			m, _ := utils.JsonStatus("success")
			io.WriteString(w, string(m))
			return
		}

		io.WriteString(w, string(failMessage))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	port := strconv.Itoa(int(ws.Port()))
	host := fmt.Sprintf("0.0.0.0:%s", port)

	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe(host, nil))
}
