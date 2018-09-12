package domainModel

import "time"

type Response struct {
	Payload []byte
	Message string
	Status  int32
}

type TransactionProposal struct {
	TxId    string
	Header  []byte
	Payload []byte
}

type ProposalResponse struct {
	Endroser        string
	Signature       []byte
	Version         int32
	Status          int32
	ChaincodeStatus int32
	TimeStamp       time.Time
	Response        Response
}

type ChainCodeResponse struct {
	TxId               string
	Payload            []byte
	TxValidationCode   int32
	ChaincodeStatus    int32
	ProposalResponses  []ProposalResponse
	TxProposalResponse TransactionProposal
}

type BlockTransaction struct {
	Type           string `json:"type"`
	TxId           string `json:"txId"`
	ValidationCode string `json:"validationCode"`
}

type BlockInfo struct {
	BlockNumber      uint64             `json:"blockNumber"`
	ChannelId        string             `json:"channelId"`
	Source           string             `json:"source"`
	NoOfTransactions int                `json:"noOfTransactions"`
	Transactions     []BlockTransaction `json:"transactions"`
}
