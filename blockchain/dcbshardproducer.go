package blockchain

import (
	"fmt"

	"github.com/ninjadotorg/constant/cashec"
	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/common/base58"
	"github.com/ninjadotorg/constant/metadata"
	"github.com/ninjadotorg/constant/privacy"
	"github.com/ninjadotorg/constant/transaction"
)

func (blockchain *BlockChain) GetAmountPerAccount(tokenID *common.Hash) (uint64, []privacy.PaymentAddress, []uint64, error) {
	tokenHoldersMap, err := blockchain.config.DataBase.GetCustomTokenPaymentAddressesBalanceUnreward(tokenID)
	if err != nil {
		return 0, nil, nil, err
	}

	// Get total token supply
	totalTokenAmount := uint64(0)
	for _, value := range tokenHoldersMap {
		totalTokenAmount += value
	}

	// Get amount per account (only count unrewarded utxo)
	tokenHolders := []privacy.PaymentAddress{}
	amounts := []uint64{}
	for holder, amount := range tokenHoldersMap {
		paymentAddressInBytes, _, _ := base58.Base58Check{}.Decode(holder)
		keySet := cashec.KeySet{}
		keySet.PaymentAddress = privacy.PaymentAddress{}
		keySet.PaymentAddress.SetBytes(paymentAddressInBytes)
		//vouts, err := blockchain.GetUnspentTxCustomTokenVout(keySet, tokenID)
		//if err != nil {
		//	return 0, nil, nil, err
		//}
		//amount := uint64(0)
		//for _, vout := range vouts {
		//	amount += vout.Value
		//}

		if amount > 0 {
			tokenHolders = append(tokenHolders, keySet.PaymentAddress)
			amounts = append(amounts, amount)
		}
	}
	return totalTokenAmount, tokenHolders, amounts, nil
}

func (blockgen *BlkTmplGenerator) buildInstitutionDividendPaymentTxs(forDCB bool, producerPrivateKey *privacy.SpendingKey) ([]*transaction.Tx, error) {
	// Get latest dividend proposal id and amount
	id, cstToPayout := blockgen.chain.BestState.Beacon.GetLatestDividendProposal(forDCB)
	if id == 0 {
		return nil, nil // No Dividend proposal found
	}

	// Check in shard state if DividendSubmit tx has been included in chain
	receivers, amounts, hasValue, err := blockgen.chain.config.DataBase.GetDividendReceiversForID(id, forDCB)
	if err != nil {
		return nil, err
	}
	if !hasValue {
		return nil, nil // Waiting for Dividend submit tx to be included in block
	}

	if len(receivers) == 0 || len(amounts) == 0 {
		return nil, nil // Paid to all receivers for the latest dividend proposal
	}

	// Get dividend info
	tokenID := &common.DCBTokenID
	if !forDCB {
		tokenID = &common.GOVTokenID
	}
	totalTokenOnAllShards, cstToPayout, err := blockgen.chain.BestState.Beacon.GetDividendAggregatedInfo(id, tokenID)
	if err != nil {
		return nil, err
	}

	// Make dividend payments to token holders
	paymentAddresses := []*privacy.PaymentAddress{}
	payoutAmounts := []uint64{}
	for i, receiver := range receivers {
		if i > metadata.MaxDivTxsPerBlock {
			break
		}
		amount := amounts[i]

		receiverCstAmount := amount * cstToPayout / totalTokenOnAllShards
		paymentAddresses = append(paymentAddresses, &receiver)
		payoutAmounts = append(payoutAmounts, receiverCstAmount)
	}

	txs, err := transaction.BuildDividendTxs(
		id,
		tokenID,
		paymentAddresses,
		payoutAmounts,
		producerPrivateKey,
		blockgen.chain.GetDatabase(),
	)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (blockgen *BlkTmplGenerator) buildDividendTxs(producerPrivateKey *privacy.SpendingKey, shardID byte) ([]metadata.Transaction, error) {
	// Build dividend payments for DCB
	forDCB := true
	dcbDividendPaymentTxs, err := blockgen.buildInstitutionDividendPaymentTxs(forDCB, producerPrivateKey)
	if err != nil {
		return nil, err
	}

	// Build dividend payments for GOV
	forDCB = false
	govDividendPaymentTxs, err := blockgen.buildInstitutionDividendPaymentTxs(forDCB, producerPrivateKey)
	if err != nil {
		return nil, err
	}

	txs := []metadata.Transaction{}
	for _, tx := range dcbDividendPaymentTxs {
		txs = append(txs, tx)
	}
	for _, tx := range govDividendPaymentTxs {
		txs = append(txs, tx)
	}

	return txs, nil
}

func buildInstitutionDividendSubmitInst(
	bc *BlockChain,
	forDCB bool,
	shardID byte,
) ([][]string, error) {
	// Get latest dividend proposal id and amount
	id, _ := bc.BestState.Beacon.GetLatestDividendProposal(forDCB)
	if id == 0 {
		// fmt.Printf("[db] no div proposal found: %t\n", forDCB)
		return nil, nil // No Dividend proposal found
	}

	// Check in shard state if DividendSubmit tx has been created
	_, _, hasValue, err := bc.config.DataBase.GetDividendReceiversForID(id, forDCB)
	if err != nil {
		fmt.Printf("[db] buildDivSub err: %v\n", err)
		return nil, err
	}
	if hasValue {
		fmt.Printf("[db] divsub created: %d %t\n", id, forDCB)
		return nil, nil // Already created DividendSubmit tx in previous blocks
	}

	tokenID := &common.DCBTokenID
	if !forDCB {
		tokenID = &common.GOVTokenID
	}
	totalTokenAmount, _, _, err := bc.GetAmountPerAccount(tokenID)
	fmt.Printf("[db] buildDivSubmit: %t, %d, %d, %d\n", forDCB, id, totalTokenAmount, shardID)
	if err != nil {
		return nil, err
	}

	// Create instruction
	return metadata.BuildDividendSubmitInst(tokenID, id, totalTokenAmount, shardID)
}
