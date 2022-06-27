package blockchain

import (
	"fmt"
	"goblockchain/internal/domain/block"
	"goblockchain/internal/domain/transaction"
	"log"
	"strings"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	transactionPool   []*transaction.Transaction
	chain             []*block.Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &block.Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())

	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *block.Block {
	b := block.NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*transaction.Transaction{}

	return b
}

func (bc *Blockchain) LastBlock() *block.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	t := transaction.NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*transaction.Transaction {
	transactions := make([]*transaction.Transaction, 0)

	for _, t := range bc.transactionPool {
		transactions = append(transactions, transaction.NewTransaction(
			t.SenderBlockchainAddress,
			t.RecipientBlockchainAddress,
			t.Value,
		))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*transaction.Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)

	guessBlock := block.Block{
		Timestamp:    0,
		Nonce:        nonce,
		PreviousHash: previousHash,
		Transactions: transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())

	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0

	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)

	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()

	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")

	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var total float32 = 0.0

	for _, b := range bc.chain {
		for _, t := range b.Transactions {
			if t.SenderBlockchainAddress == blockchainAddress {
				total -= t.Value
			}

			if t.RecipientBlockchainAddress == blockchainAddress {
				total += t.Value
			}
		}
	}

	return total
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf(
			"%s Chain %d %s\n",
			strings.Repeat("=", 25),
			i,
			strings.Repeat("=", 25),
		)
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}
