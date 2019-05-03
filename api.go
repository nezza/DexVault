package main

import (
	"github.com/binance-chain/go-sdk/types/msg"
)

type BasicMessage struct {
	Wallet string
}

type SignedMessage struct {
	BasicMessage
	BroadcastHost    string
	BroadcastNetwork int
	ChainId          string
	AccountNumber    int64
	Sequence         int64
}

type CreateOrder struct {
	SignedMessage
	BaseAssetSymbol  string
	QuoteAssetSymbol string
	Op               int8
	Price            int64
	Quantity         int64
}

type CancelOrder struct {
	SignedMessage
	BaseAssetSymbol  string
	QuoteAssetSymbol string
	RefId            string
}

type TokenBurn struct {
	SignedMessage
	Symbol string
	Amount int64
}

type DepositProposal struct {
	SignedMessage
	ProposalID int64
	Amount     int64
}

type FreezeToken struct {
	SignedMessage
	Symbol string
	Amount int64
}

type IssueToken struct {
	SignedMessage
	Name   string
	Symbol string
	// OriginalSymbol string
	Supply   int64
	Mintable bool
	// Owner string
}

type ListPair struct {
	SignedMessage
	ProposalID       int64
	BaseAssetSymbol  string
	QuoteAssetSymbol string
	InitPrice        int64
}

type MintToken struct {
	SignedMessage
	Symbol string
	Amount int64
}

type SendToken struct {
	SignedMessage
	Transfers []msg.Transfer
}

type SubmitProposal struct {
	SignedMessage
	Title          string
	Description    string
	ProposalType   byte
	InitialDeposit int64
	VotingPeriod   int64
}

type UnfreezeToken struct {
	SignedMessage
	Symbol string
	Amount int64
}

type VoteProposal struct {
	SignedMessage
	ProposalID int64
	Option     byte
}
