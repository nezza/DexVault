package main

import (
	"github.com/binance-chain/go-sdk/keys"
	// "github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"time"
)

func signMessage(sm SignedMessage, memo string, m msg.Msg, keyManager keys.KeyManager) ([]byte, error) {
	err := m.ValidateBasic()
	if err != nil {
		return nil, err
	}

	signMsg := tx.StdSignMsg{
		ChainID:       sm.ChainId,
		AccountNumber: sm.AccountNumber,
		Sequence:      sm.Sequence,
		Memo:          memo, // Only transfer supports memo
		Msgs:          []msg.Msg{m},
		Source:        types.GoSdkSource,
	}

	hexTx, err := keyManager.Sign(signMsg)
	if err != nil {
		return nil, err
	}

	return hexTx, nil
}

func createSignedCreateOrderMessage(keyManager keys.KeyManager, co *CreateOrder) ([]byte, error) {
	fromAddr := keyManager.GetAddr()

	newOrderMessage := msg.NewCreateOrderMsg(
		fromAddr,
		string(fromAddr)+"-"+string(co.Sequence),
		// co.OrderId, // TODO: Difference between sequence & orderid?!
		co.Op,
		co.CombinedSymbol(),
		co.Price,
		co.Quantity,
	)

	hexTx, err := signMessage(co.SignedMessage, "", newOrderMessage, keyManager)
	return hexTx, err
}

// cancelOrderResult, err := client.CancelOrder(tradeSymbol, nativeSymbol, createOrderResult.OrderId, createOrderResult.OrderId, true)

func createSignedCancelOrderMsg(keyManager keys.KeyManager, co *CancelOrder) ([]byte, error) {
	fromAddr := keyManager.GetAddr()

	cancelOrderMessage := msg.NewCancelOrderMsg(
		fromAddr,
		co.CombinedSymbol(),
		co.RefId,
	)
	hexTx, err := signMessage(co.SignedMessage, "", cancelOrderMessage, keyManager)
	return hexTx, err
}

func createSignedTokenBurnMsg(keyManager keys.KeyManager, tb *TokenBurn) ([]byte, error) {
	burnMsg := msg.NewTokenBurnMsg(
		keyManager.GetAddr(),
		tb.Symbol,
		tb.Amount)
	hexTx, err := signMessage(tb.SignedMessage, "", burnMsg, keyManager)
	return hexTx, err
}
func createSignedDepositMsg(keyManager keys.KeyManager, dp *DepositProposal) ([]byte, error) {
	coins := types.Coins{types.Coin{Denom: types.NativeSymbol, Amount: dp.Amount}}
	depositMsg := msg.NewDepositMsg(keyManager.GetAddr(), dp.ProposalID, coins)
	hexTx, err := signMessage(dp.SignedMessage, "", depositMsg, keyManager)
	return hexTx, err
}

func createSignedFreezeTokenMsg(keyManager keys.KeyManager, ft *FreezeToken) ([]byte, error) {
	freezeMsg := msg.NewFreezeMsg(
		keyManager.GetAddr(),
		ft.Symbol,
		ft.Amount,
	)
	hexTx, err := signMessage(ft.SignedMessage, "", freezeMsg, keyManager)
	return hexTx, err
}

func createSignedIssueTokenMsg(keyManager keys.KeyManager, it *IssueToken) ([]byte, error) {
	issueMsg := msg.NewTokenIssueMsg(
		keyManager.GetAddr(),
		it.Name,
		it.Symbol,
		it.Supply,
		it.Mintable)

	hexTx, err := signMessage(it.SignedMessage, "", issueMsg, keyManager)
	return hexTx, err
}

func createSignedListPairMsg(keyManager keys.KeyManager, lp *ListPair) ([]byte, error) {
	msg := msg.NewDexListMsg(keyManager.GetAddr(), lp.ProposalID, lp.BaseAssetSymbol, lp.QuoteAssetSymbol, lp.InitPrice)
	hexTx, err := signMessage(lp.SignedMessage, "", msg, keyManager)
	return hexTx, err
}

func createSignedMintTokenMsg(keyManager keys.KeyManager, mt *MintToken) ([]byte, error) {
	msg := msg.NewMintMsg(
		keyManager.GetAddr(),
		mt.Symbol,
		mt.Amount)
	hexTx, err := signMessage(mt.SignedMessage, "", msg, keyManager)
	return hexTx, err
}

func createSignedSendTokenMsg(keyManager keys.KeyManager, st *SendToken) ([]byte, error) {
	fromCoins := types.Coins{}
	for _, t := range st.Transfers {
		fromCoins = fromCoins.Plus(t.Coins)
	}
	sendMsg := msg.CreateSendMsg(
		keyManager.GetAddr(),
		fromCoins,
		st.Transfers)
	hexTx, err := signMessage(st.SignedMessage, "", sendMsg, keyManager)
	return hexTx, err
}

func createSignedSubmitProposalMsg(keyManager keys.KeyManager, sp *SubmitProposal) ([]byte, error) {
	coins := types.Coins{types.Coin{Denom: types.NativeSymbol, Amount: sp.InitialDeposit}}
	proposalMsg := msg.NewMsgSubmitProposal(sp.Title, sp.Description, msg.ProposalKind(sp.ProposalType), keyManager.GetAddr(), coins, time.Duration(sp.VotingPeriod))
	hexTx, err := signMessage(sp.SignedMessage, "", proposalMsg, keyManager)
	return hexTx, err
}

func createUnfreezeTokenMsg(keyManager keys.KeyManager, ut *UnfreezeToken) ([]byte, error) {
	unfreezeMsg := msg.NewUnfreezeMsg(
		keyManager.GetAddr(),
		ut.Symbol,
		ut.Amount)
	hexTx, err := signMessage(ut.SignedMessage, "", unfreezeMsg, keyManager)
	return hexTx, err
}

func createSignedVoteProposalMsg(keyManager keys.KeyManager, vp *VoteProposal) ([]byte, error) {
	voteMsg := msg.NewMsgVote(keyManager.GetAddr(), vp.ProposalID, msg.VoteOption(vp.Option))
	hexTx, err := signMessage(vp.SignedMessage, "", voteMsg, keyManager)
	return hexTx, err
}
