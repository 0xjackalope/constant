package blockchain

import (
	"errors"
	"time"

	"github.com/ninjadotorg/cash/common"
	"github.com/ninjadotorg/cash/privacy/client"
	"github.com/ninjadotorg/cash/transaction"
	"github.com/ninjadotorg/cash/privacy"
)

func (blockgen *BlkTmplGenerator) NewBlockTemplate(payToAddress privacy.PaymentAddress, chainID byte) (*BlockTemplate, error) {

	prevBlock := blockgen.chain.BestState[chainID]
	prevBlockHash := blockgen.chain.BestState[chainID].BestBlock.Hash()
	sourceTxns := blockgen.txPool.MiningDescs()

	var txsToAdd []transaction.Transaction
	var txToRemove []transaction.Transaction
	// var actionParamTxs []*transaction.ActionParamTx
	// var feeMap = map[string]uint64{
	// 	fmt.Sprintf(common.AssetTypeCoin):     0,
	// 	fmt.Sprintf(common.AssetTypeBond):     0,
	// 	fmt.Sprintf(common.AssetTypeGovToken): 0,
	// 	fmt.Sprintf(common.AssetTypeDcbToken): 0,
	// }
	totalFee := uint64(0)

	// Get salary per tx
	salaryPerTx := blockgen.rewardAgent.GetSalaryPerTx(chainID)
	// Get basic salary on block
	basicSalary := blockgen.rewardAgent.GetBasicSalary(chainID)

	if len(sourceTxns) < common.MinTxsInBlock {
		// if len of sourceTxns < MinTxsInBlock -> wait for more transactions
		Logger.log.Info("not enough transactions. Wait for more...")
		<-time.Tick(common.MinBlockWaitTime * time.Second)
		sourceTxns = blockgen.txPool.MiningDescs()
		if len(sourceTxns) == 0 {
			<-time.Tick(common.MaxBlockWaitTime * time.Second)
			sourceTxns = blockgen.txPool.MiningDescs()
			if len(sourceTxns) == 0 {
				// return nil, errors.New("No Tx")
				Logger.log.Info("Creating empty block...")
				goto concludeBlock
			}
		}
	}

	for _, txDesc := range sourceTxns {
		tx := txDesc.Tx
		txChainID, _ := common.GetTxSenderChain(tx.GetSenderAddrLastByte())
		if txChainID != chainID {
			continue
		}
		if !tx.ValidateTransaction() {
			txToRemove = append(txToRemove, transaction.Transaction(tx))
			continue
		}
		// txType := tx.GetType()
		// txFee := uint64(0)
		// switch txType {
		// case common.TxActionParamsType:
		// 	// actionParamTxs = append(actionParamTxs, tx.(*transaction.ActionParamTx))
		// 	continue
		// case common.TxCustomTokenType:
		// 	txFee = tx.(*transaction.TxCustomToken).Fee
		// case common.TxVotingType:
		// 	txFee = tx.(*transaction.TxVoting).Fee
		// case common.TxSalaryType:
		// case common.TxNormalType:
		// 	txFee = tx.(*transaction.Tx).Fee
		// }
		totalFee += tx.GetTxFee()
		txsToAdd = append(txsToAdd, tx)
		if len(txsToAdd) == common.MaxTxsInBlock {
			break
		}
	}

	for _, tx := range txToRemove {
		blockgen.txPool.RemoveTx(tx)
	}

	// check len of txs in block
	if len(txsToAdd) == 0 {
		return nil, errors.New("no transaction available for this chain")
	}

concludeBlock:
// Get blocksalary fund from txs
	salaryFundAdd := uint64(0)
	salaryMULTP := uint64(0)
	for _, blockTx := range txsToAdd {
		if blockTx.GetType() == common.TxVotingType {
			tx, ok := blockTx.(*transaction.TxVoting)
			if !ok {
				Logger.log.Error("Transaction not recognized to store in database")
				continue
			}
			salaryFundAdd += tx.GetValue()
		}
		if blockTx.GetTxFee() > 0 {
			salaryMULTP++
		}
	}

	rt := blockgen.chain.BestState[chainID].BestBlock.Header.MerkleRootCommitments.CloneBytes()

	// ------------------------ HOW to GET salary on a block-------------------
	// total salary = tx * (salary per tx) + (basic salary on block)
	// ------------------------------------------------------------------------
	totalSalary := salaryMULTP*salaryPerTx + basicSalary
	// create salary tx to pay constant for block producer
	salaryTx, err := createSalaryTx(totalSalary, &payToAddress, rt, chainID)

	if err != nil {
		return nil, err
	}
	// the 1st tx will be salaryTx
	txsToAdd = append([]transaction.Transaction{salaryTx}, txsToAdd...)

	merkleRoots := Merkle{}.BuildMerkleTreeStore(txsToAdd)
	merkleRoot := merkleRoots[len(merkleRoots)-1]

	block := Block{}
	currentSalaryFund := blockgen.chain.BestState[chainID].BestBlock.Header.SalaryFund
	block.Header = BlockHeader{
		Version:               BlockVersion,
		PrevBlockHash:         *prevBlockHash,
		MerkleRoot:            *merkleRoot,
		MerkleRootCommitments: common.Hash{},
		Timestamp:             time.Now().Unix(),
		BlockCommitteeSigs:    make([]string, common.TotalValidators),
		Committee:             make([]string, common.TotalValidators),
		ChainID:               chainID,
		SalaryFund:            currentSalaryFund - (salaryMULTP * salaryPerTx) + totalFee + salaryFundAdd,
	}
	for _, tx := range txsToAdd {
		if err := block.AddTransaction(tx); err != nil {
			return nil, err
		}
	}

	// Add new commitments to merkle tree and save the root
	newTree := blockgen.chain.BestState[chainID].CmTree.MakeCopy()
	UpdateMerkleTreeForBlock(newTree, &block)
	rt = newTree.GetRoot(common.IncMerkleTreeHeight)
	copy(block.Header.MerkleRootCommitments[:], rt)

	//update the latest AgentDataPoints to block
	// block.AgentDataPoints = agentDataPoints
	// Set height
	block.Height = prevBlock.Height + 1

	blockTemp := &BlockTemplate{
		Block: &block,
	}
	return blockTemp, nil
}

type BlkTmplGenerator struct {
	txPool      TxPool
	chain       *BlockChain
	rewardAgent RewardAgent
	// chainParams *blockchain.Params
	// policy      *Policy
}

type BlockTemplate struct {
	Block *Block
}

// txPool represents a source of transactions to consider for inclusion in
// new blocks.
//
// The interface contract requires that all of these methods are safe for
// concurrent access with respect to the source.
type TxPool interface {
	// LastUpdated returns the last time a transaction was added to or
	// removed from the source pool.
	LastUpdated() time.Time

	// MiningDescs returns a slice of mining descriptors for all the
	// transactions in the source pool.
	MiningDescs() []*transaction.TxDesc

	// HaveTransaction returns whether or not the passed transaction hash
	// exists in the source pool.
	HaveTransaction(hash *common.Hash) bool

	// RemoveTx remove tx from tx resource
	RemoveTx(tx transaction.Transaction) error
}

type RewardAgent interface {
	GetBasicSalary(chainID byte) uint64
	GetSalaryPerTx(chainID byte) uint64
}

func (self BlkTmplGenerator) Init(txPool TxPool, chain *BlockChain, rewardAgent RewardAgent) (*BlkTmplGenerator, error) {
	return &BlkTmplGenerator{
		txPool:      txPool,
		chain:       chain,
		rewardAgent: rewardAgent,
	}, nil
}

// createSalaryTx
// Blockchain use this tx to pay a reward(salary) to miner of chain
// #1 - salary:
// #2 - receiverAddr:
// #3 - rt
// #4 - chainID
func createSalaryTx(
	salary uint64,
	receiverAddr *privacy.PaymentAddress,
	rt []byte,
	chainID byte,
) (*transaction.Tx, error) {
	// Create Proof for the joinsplit op
	inputs := make([]*client.JSInput, 2)
	inputs[0] = transaction.CreateRandomJSInput(nil)
	inputs[1] = transaction.CreateRandomJSInput(inputs[0].Key)
	dummyAddress := client.GenPaymentAddress(*inputs[0].Key)

	// Create new notes: first one is salary UTXO, second one has 0 value
	var temp []byte
	copy(temp, receiverAddr.Pk[:])
	outNote := &client.Note{Value: salary, Apk: temp}
	placeHolderOutputNote := &client.Note{Value: 0, Apk: temp}

	var temp2 client.TransmissionKey
	copy(temp2[:], receiverAddr.Tk[:])
	outputs := []*client.JSOutput{&client.JSOutput{}, &client.JSOutput{}}
	outputs[0].EncKey = temp2
	outputs[0].OutputNote = outNote
	outputs[1].EncKey = temp2
	outputs[1].OutputNote = placeHolderOutputNote

	// Generate proof and sign tx
	tx, err := transaction.CreateEmptyTx(common.TxSalaryType)
	if err != nil {
		return nil, err
	}
	tx.AddressLastByte = dummyAddress.Apk[len(dummyAddress.Apk)-1]
	rtMap := map[byte][]byte{chainID: rt}
	inputMap := map[byte][]*client.JSInput{chainID: inputs}

	// NOTE: always pay salary with constant coin
	assetTypeToPaySalary := common.AssetTypeCoin
	err = tx.BuildNewJSDesc(inputMap, outputs, rtMap, salary, 0, assetTypeToPaySalary, true)
	if err != nil {
		return nil, err
	}
	tx, err = transaction.SignTx(tx)
	if err != nil {
		return nil, err
	}
	return tx, err
}
