package metadata

import (
	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/database"
	"github.com/ninjadotorg/constant/privacy"
)

type SealedLv1DCBVoteProposalMetadata struct {
	SealedLv1VoteProposalMetadata SealedLv1VoteProposalMetadata
	MetadataBase
}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	boardType := DCBBoard
	nextConstitutionIndex := bcr.GetConstitution(boardType).GetConstitutionIndex() + 1
	err := bcr.GetDatabase().AddVoteLv1or2Proposal(boardType.BoardTypeDB(), nextConstitutionIndex, tx.GetMetadata().Hash())
	if err != nil {
		return err
	}
	return nil
}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return sealedLv1DCBVoteProposalMetadata.SealedLv1VoteProposalMetadata.ValidateSanityData(bcr, tx)
}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) ValidateMetadataByItself() bool {
	return sealedLv1DCBVoteProposalMetadata.SealedLv1VoteProposalMetadata.ValidateMetadataByItself()
}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) GetBoardType() BoardType {
	return DCBBoard
}

func NewSealedLv1DCBVoteProposalMetadata(
	sealedVoteProposal []byte,
	lockersPaymentAddress []privacy.PaymentAddress,
	pointerToLv2VoteProposal common.Hash,
	pointerToLv3VoteProposal common.Hash,
) *SealedLv1DCBVoteProposalMetadata {
	return &SealedLv1DCBVoteProposalMetadata{
		SealedLv1VoteProposalMetadata: *NewSealedLv1VoteProposalMetadata(
			sealedVoteProposal,
			lockersPaymentAddress,
			pointerToLv2VoteProposal,
			pointerToLv3VoteProposal,
		),

		MetadataBase: *NewMetadataBase(SealedLv1DCBVoteProposalMeta),
	}
}

func NewSealedLv1VoteProposalMetadataFromRPC(data map[string]interface{}) (Metadata, error) {
	boardType := NewBoardTypeFromString(data["BoardType"].(string))
	sealLv1Data := data["SealLv1Data"].([]byte)
	paymentAddresses := data["PaymentAddresses"].([]privacy.PaymentAddress)
	lv2TxID := data["Lv2TxID"].(common.Hash)
	lv3TxID := data["Lv3TxID"].(common.Hash)
	var meta Metadata
	if boardType == DCBBoard {
		meta = NewSealedLv1DCBVoteProposalMetadata(
			sealLv1Data,
			paymentAddresses,
			lv2TxID,
			lv3TxID,
		)
	} else {
		meta = NewSealedLv1GOVVoteProposalMetadata(
			sealLv1Data,
			paymentAddresses,
			lv2TxID,
			lv3TxID,
		)
	}
	return meta, nil
}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) Hash() *common.Hash {
	record := string(sealedLv1DCBVoteProposalMetadata.SealedLv1VoteProposalMetadata.ToBytes())

	hash := common.DoubleHashH([]byte(record))
	return &hash

}

func (sealedLv1DCBVoteProposalMetadata *SealedLv1DCBVoteProposalMetadata) ValidateTxWithBlockChain(
	tx Transaction,
	bcr BlockchainRetriever,
	shardID byte,
	db database.DatabaseInterface,
) (bool, error) {
	boardType := DCBBoard
	return sealedLv1DCBVoteProposalMetadata.SealedLv1VoteProposalMetadata.ValidateTxWithBlockChain(
		boardType,
		tx,
		bcr,
		shardID,
		db,
	)
}

type SealedLv2DCBVoteProposalMetadata struct {
	SealedLv2VoteProposalMetadata SealedLv2VoteProposalMetadata

	MetadataBase
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	boardType := DCBBoard
	nextConstitutionIndex := bcr.GetConstitution(boardType).GetConstitutionIndex() + 1
	err := bcr.GetDatabase().AddVoteLv1or2Proposal(boardType.BoardTypeDB(), nextConstitutionIndex, tx.GetMetadata().Hash())
	if err != nil {
		return err
	}
	return nil
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return sealedLv2DCBVoteProposalMetadata.SealedLv2VoteProposalMetadata.ValidateSanityData(bcr, tx)
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) ValidateMetadataByItself() bool {
	return sealedLv2DCBVoteProposalMetadata.SealedLv2VoteProposalMetadata.ValidateMetadataByItself()
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) GetBoardType() BoardType {
	return DCBBoard
}

func NewSealedLv2DCBVoteProposalMetadata(
	sealedVoteProposal []byte,
	lockerPaymentAddress []privacy.PaymentAddress,
	pointerToLv3VoteProposal common.Hash,
) *SealedLv2DCBVoteProposalMetadata {
	return &SealedLv2DCBVoteProposalMetadata{
		SealedLv2VoteProposalMetadata: *NewSealedLv2VoteProposalMetadata(
			sealedVoteProposal,
			lockerPaymentAddress,
			pointerToLv3VoteProposal,
		),

		MetadataBase: *NewMetadataBase(SealedLv2DCBVoteProposalMeta),
	}
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) Hash() *common.Hash {
	record := sealedLv2DCBVoteProposalMetadata.SealedLv2VoteProposalMetadata.ToBytes()
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (sealedLv2DCBVoteProposalMetadata *SealedLv2DCBVoteProposalMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	boardType := sealedLv2DCBVoteProposalMetadata.GetBoardType()
	return sealedLv2DCBVoteProposalMetadata.SealedLv2VoteProposalMetadata.ValidateTxWithBlockChain(
		boardType,
		tx,
		bcr,
		shardID,
		db,
	)

}

type SealedLv3DCBVoteProposalMetadata struct {
	SealedLv3VoteProposalMetadata SealedLv3VoteProposalMetadata

	MetadataBase
}

func (sealedLv3DCBVoteProposalMetadata *SealedLv3DCBVoteProposalMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	boardType := DCBBoard
	nextConstitutionIndex := bcr.GetConstitution(boardType).GetConstitutionIndex() + 1
	err := bcr.GetDatabase().AddVoteLv3Proposal(boardType.BoardTypeDB(), nextConstitutionIndex, tx.GetMetadata().Hash())
	if err != nil {
		return err
	}
	return nil
}

func (sealedLv3DCBVoteProposalMetadata *SealedLv3DCBVoteProposalMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return sealedLv3DCBVoteProposalMetadata.SealedLv3VoteProposalMetadata.ValidateTxWithBlockChain(tx, bcr, b, db)
}

func (sealedLv3DCBVoteProposalMetadata *SealedLv3DCBVoteProposalMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return sealedLv3DCBVoteProposalMetadata.SealedLv3VoteProposalMetadata.ValidateSanityData(bcr, tx)
}

func (sealedLv3DCBVoteProposalMetadata *SealedLv3DCBVoteProposalMetadata) ValidateMetadataByItself() bool {
	return sealedLv3DCBVoteProposalMetadata.SealedLv3VoteProposalMetadata.ValidateMetadataByItself()
}

func (sealedLv3DCBVoteProposalMetadata *SealedLv3DCBVoteProposalMetadata) GetBoardType() BoardType {
	return DCBBoard
}

func NewSealedLv3DCBVoteProposalMetadata(
	sealedVoteProposal []byte,
	lockerPaymentAddress []privacy.PaymentAddress,
) *SealedLv3DCBVoteProposalMetadata {
	return &SealedLv3DCBVoteProposalMetadata{
		SealedLv3VoteProposalMetadata: *NewSealedLv3VoteProposalMetadata(
			sealedVoteProposal, lockerPaymentAddress,
		),
		MetadataBase: *NewMetadataBase(SealedLv3DCBVoteProposalMeta),
	}
}

type NormalDCBVoteProposalFromSealerMetadata struct {
	NormalVoteProposalFromSealerMetadata NormalVoteProposalFromSealerMetadata

	MetadataBase
}

func (normalDCBVoteProposalFromSealerMetadata *NormalDCBVoteProposalFromSealerMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	boardType := DCBBoard
	nextConstitutionIndex := bcr.GetConstitution(boardType).GetConstitutionIndex() + 1
	err := bcr.GetDatabase().AddVoteNormalProposalFromOwner(
		boardType.BoardTypeDB(),
		nextConstitutionIndex,
		&normalDCBVoteProposalFromSealerMetadata.NormalVoteProposalFromSealerMetadata.PointerToLv3VoteProposal,
		normalDCBVoteProposalFromSealerMetadata.NormalVoteProposalFromSealerMetadata.VoteProposal.ToBytes(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (normalDCBVoteProposalFromSealerMetadata *NormalDCBVoteProposalFromSealerMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return normalDCBVoteProposalFromSealerMetadata.NormalVoteProposalFromSealerMetadata.ValidateSanityData(bcr, tx)
}

func (normalDCBVoteProposalFromSealerMetadata *NormalDCBVoteProposalFromSealerMetadata) ValidateMetadataByItself() bool {
	return normalDCBVoteProposalFromSealerMetadata.NormalVoteProposalFromSealerMetadata.ValidateMetadataByItself()
}

func NewNormalDCBVoteProposalFromSealerMetadata(
	voteProposal VoteProposalData,
	lockerPaymentAddress []privacy.PaymentAddress,
	pointerToLv1VoteProposal common.Hash,
	pointerToLv3VoteProposal common.Hash,
) *NormalDCBVoteProposalFromSealerMetadata {
	return &NormalDCBVoteProposalFromSealerMetadata{
		NormalVoteProposalFromSealerMetadata: *NewNormalVoteProposalFromSealerMetadata(
			voteProposal,
			lockerPaymentAddress,
			pointerToLv1VoteProposal,
			pointerToLv3VoteProposal,
		),
		MetadataBase: *NewMetadataBase(NormalDCBVoteProposalFromSealerMeta),
	}
}

func (normalDCBVoteProposalFromSealerMetadata *NormalDCBVoteProposalFromSealerMetadata) GetBoardType() BoardType {
	return DCBBoard
}

func (normalDCBVoteProposalFromSealerMetadata *NormalDCBVoteProposalFromSealerMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	//Validate these pubKeys are in board
	boardType := normalDCBVoteProposalFromSealerMetadata.GetBoardType()
	return normalDCBVoteProposalFromSealerMetadata.NormalVoteProposalFromSealerMetadata.ValidateTxWithBlockChain(
		boardType,
		tx,
		bcr,
		shardID,
		db,
	)
}

type NormalDCBVoteProposalFromOwnerMetadata struct {
	NormalVoteProposalFromOwnerMetadata NormalVoteProposalFromOwnerMetadata
	MetadataBase
}

func (normalDCBVoteProposalFromOwnerMetadata *NormalDCBVoteProposalFromOwnerMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	boardType := DCBBoard
	nextConstitutionIndex := bcr.GetConstitution(boardType).GetConstitutionIndex() + 1
	err := bcr.GetDatabase().AddVoteNormalProposalFromOwner(
		boardType.BoardTypeDB(),
		nextConstitutionIndex,
		&normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.PointerToLv3VoteProposal,
		normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.VoteProposal.ToBytes(),
	)
	if err != nil {
		return err
	}
	return nil

}

func (normalDCBVoteProposalFromOwnerMetadata *NormalDCBVoteProposalFromOwnerMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.ValidateSanityData(bcr, tx)
}

func (normalDCBVoteProposalFromOwnerMetadata *NormalDCBVoteProposalFromOwnerMetadata) ValidateMetadataByItself() bool {
	return normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.ValidateMetadataByItself()
}

func NewNormalDCBVoteProposalFromOwnerMetadata(
	voteProposal VoteProposalData,
	lockerPaymentAddress []privacy.PaymentAddress,
	pointerToLv3VoteProposal common.Hash,
) *NormalDCBVoteProposalFromOwnerMetadata {
	return &NormalDCBVoteProposalFromOwnerMetadata{
		NormalVoteProposalFromOwnerMetadata: *NewNormalVoteProposalFromOwnerMetadata(
			voteProposal,
			lockerPaymentAddress,
			pointerToLv3VoteProposal,
		),
		MetadataBase: *NewMetadataBase(NormalDCBVoteProposalFromOwnerMeta),
	}
}

func (normalDCBVoteProposalFromOwnerMetadata *NormalDCBVoteProposalFromOwnerMetadata) Hash() *common.Hash {
	record := normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.ToBytes()

	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (normalDCBVoteProposalFromOwnerMetadata *NormalDCBVoteProposalFromOwnerMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	//Validate these pubKeys are in board
	boardType := DCBBoard
	return normalDCBVoteProposalFromOwnerMetadata.NormalVoteProposalFromOwnerMetadata.ValidateTxWithBlockChain(
		boardType,
		tx,
		bcr,
		shardID,
		db,
	)
}

type PunishDCBDecryptMetadata struct {
	PunishDecryptMetadata PunishDecryptMetadata
	MetadataBase
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	// todo @0xjackalope
	return nil
}

func NewPunishDCBDecryptMetadata(paymentAddress privacy.PaymentAddress) *PunishDCBDecryptMetadata {
	return &PunishDCBDecryptMetadata{
		PunishDecryptMetadata: PunishDecryptMetadata{
			PaymentAddress: paymentAddress,
		},
		MetadataBase: *NewMetadataBase(PunishDCBDecryptMeta),
	}
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) Hash() *common.Hash {
	record := string(punishDCBDecryptMetadata.PunishDecryptMetadata.ToBytes())
	record += punishDCBDecryptMetadata.MetadataBase.Hash().String()

	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) ValidateTxWithBlockChain(Transaction, BlockchainRetriever, byte, database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) ValidateSanityData(BlockchainRetriever, Transaction) (bool, bool, error) {
	return true, true, nil
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) ValidateMetadataByItself() bool {
	return true
}

func (punishDCBDecryptMetadata *PunishDCBDecryptMetadata) CalculateSize() uint64 {
	return calculateSize(punishDCBDecryptMetadata)
}
