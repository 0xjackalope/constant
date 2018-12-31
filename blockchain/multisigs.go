package blockchain

import (
	"errors"

	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/metadata"
	"github.com/ninjadotorg/constant/transaction"
)

func (blockGen *BlkTmplGenerator) registerMultiSigsAddresses(
	txs []metadata.Transaction,
) error {
	if len(txs) == 0 {
		return nil
	}
	msRegs := map[string]*metadata.MultiSigsRegistration{}
	sortedTxs := transaction.SortTxsByLockTime(txs, false)
	for _, tx := range sortedTxs {
		meta := tx.GetMetadata()
		if meta == nil {
			continue
		}
		multiSigsReg, ok := meta.(*metadata.MultiSigsRegistration)
		if !ok {
			return errors.New("Could not parse MultiSigsRegistration metadata")
		}
		msRegs[string(multiSigsReg.PaymentAddress.Pk)] = multiSigsReg
	}
	// store msRegs to db
	// TODO: should use batch-write to ensure data consistency
	for _, msReg := range msRegs {
		err := blockGen.chain.config.DataBase.StoreMultiSigsRegistration(msReg.PaymentAddress.Pk, common.ToBytes(*msReg))
		if err != nil {
			return err
		}
	}
	return nil
}