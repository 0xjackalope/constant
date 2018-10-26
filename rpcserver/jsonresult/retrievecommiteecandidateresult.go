package jsonresult

import "github.com/ninjadotorg/cash/blockchain"

type RetrieveCommitteecCandidateResult struct {
	Value     uint64 `json:"Value"`
	Timestamp int64  `json:"Timestamp"`
	ChainID   byte   `json:"ChainID"`
}

func (self *RetrieveCommitteecCandidateResult) Init(obj *blockchain.CommitteeCandidateInfo) {
	self.ChainID = obj.ChainID
	self.Timestamp = obj.Timestamp
	self.ChainID = obj.ChainID
}
