package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	t "goblockchain/internal/domain/transaction"
	"time"
)

type Block struct {
	Timestamp    int64            `json:"timestamp"`
	Nonce        int              `json:"nonce"`
	PreviousHash [32]byte         `json:"previous_hash"`
	Transactions []*t.Transaction `json:"transactions"`
}

func NewBlock(none int, previousHash [32]byte, transactions []*t.Transaction) *Block {
	return &Block{
		Timestamp:    time.Now().UnixNano(),
		Nonce:        none,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) Print() {
	fmt.Printf("timestamp %d\n", b.Timestamp)
	fmt.Printf("nonce %d\n", b.Nonce)
	fmt.Printf("previous_hash %x\n", b.PreviousHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}
