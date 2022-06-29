package walletrequests

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	Value                      *string `json:"value"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPublicKey == nil ||
		tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.Value == nil {
		return false
	}

	return true
}
