package metadata

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ninjadotorg/constant/blockchain/params"
	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/database"
	"github.com/ninjadotorg/constant/privacy"
	"github.com/ninjadotorg/constant/privacy/zeroknowledge"
)

type MetadataBase struct {
	Type int
}

func NewMetadataBase(thisType int) *MetadataBase {
	return &MetadataBase{Type: thisType}
}

func calculateSize(meta Metadata) uint64 {
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return 0
	}
	return uint64(len(metaBytes))
}

func (mb *MetadataBase) CalculateSize() uint64 {
	return 0
}

func (mb *MetadataBase) Validate() error {
	return nil
}

func (mb *MetadataBase) Process() error {
	return nil
}

func (mb *MetadataBase) GetType() int {
	return mb.Type
}

func (mb *MetadataBase) Hash() *common.Hash {
	record := strconv.Itoa(mb.Type)
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (mb *MetadataBase) ValidateBeforeNewBlock(tx Transaction, bcr BlockchainRetriever, shardID byte) bool {
	// TODO: 0xjackalope
	return true
}

func (mb *MetadataBase) CheckTransactionFee(tx Transaction, minFeePerKbTx uint64) bool {
	txFee := tx.GetTxFee()
	fullFee := minFeePerKbTx * tx.GetTxActualSize()
	return !(txFee < fullFee)
}

func (mb *MetadataBase) VerifyMultiSigs(
	tx Transaction,
	db database.DatabaseInterface,
) (bool, error) {
	return true, nil
}

func (mb *MetadataBase) BuildReqActions(tx Transaction, bcr BlockchainRetriever, shardID byte) ([][]string, error) {
	return [][]string{}, nil
}

func (mb *MetadataBase) ProcessWhenInsertBlockShard(tx Transaction, retriever BlockchainRetriever) error {
	return nil
}

// TODO(@0xankylosaurus): move TxDesc to mempool DTO
// This is tx struct which is really saved in tx mempool
type TxDesc struct {
	// Tx is the transaction associated with the entry.
	Tx Transaction

	// Added is the time when the entry was added to the source pool.
	Added time.Time

	// Height is the best block's height when the entry was added to the the source pool.
	Height uint64

	// Fee is the total fee the transaction associated with the entry pays.
	Fee uint64

	// FeePerKB is the fee the transaction pays in coin per 1000 bytes.
	FeePerKB int32
}

// Interface for mempool which is used in metadata
type MempoolRetriever interface {
	GetSerialNumbers() map[common.Hash][][]byte
	GetTxsInMem() map[common.Hash]TxDesc
}

// Interface for blockchain which is used in metadata
type BlockchainRetriever interface {
	GetTxChainHeight(tx Transaction) (uint64, error)
	GetChainHeight(byte) uint64
	GetBeaconHeight() uint64
	GetCustomTokenTxs(*common.Hash) (map[common.Hash]Transaction, error)
	GetDCBParams() params.DCBParams
	GetBoardPubKeys(boardType BoardType) [][]byte
	GetBoardPaymentAddress(boardType BoardType) []privacy.PaymentAddress
	GetGOVParams() params.GOVParams
	GetTransactionByHash(*common.Hash) (byte, *common.Hash, int, Transaction, error)
	GetOracleParams() *params.Oracle
	GetConstitutionStartHeight(boardType BoardType, shardID byte) uint64
	GetConstitutionEndHeight(boardType BoardType, shardID byte) uint64
	GetCurrentBeaconBlockHeight(byte) uint64
	GetBoardEndHeight(boardType BoardType, chainID byte) uint64
	GetAllCommitteeValidatorCandidate() (map[byte][]string, map[byte][]string, []string, []string, []string, []string, []string, []string)
	GetDatabase() database.DatabaseInterface

	// For validating loan metadata
	// GetLoanTxs([]byte) ([][]byte, error)
	GetLoanReq(loanID []byte) (*common.Hash, error)
	GetLoanResps(loanID []byte) ([][]byte, []ValidLoanResponse, error)
	GetNumberOfDCBGovernors() int
	GetNumberOfGOVGovernors() int
	GetLoanPayment([]byte) (uint64, uint64, uint64, error)
	GetLoanRequestMeta(loanID []byte) (*LoanRequest, error)
	GetLoanWithdrawed(loanID []byte) (bool, error)

	// For validating dividend
	GetLatestDividendProposal(bool) (uint64, uint64)
	GetAmountPerAccount(*common.Hash) (uint64, []privacy.PaymentAddress, []uint64, error)
	GetDividendReceiversForID(dividendID uint64, forDCB bool) ([]privacy.PaymentAddress, []uint64, bool, error)

	// For validating crowdsale
	GetCrowdsaleData([]byte) (*params.SaleData, error)

	// For validating reserve
	GetAssetPrice(assetID *common.Hash) uint64

	// For validating cmb
	GetCMB([]byte) (privacy.PaymentAddress, []privacy.PaymentAddress, uint64, *common.Hash, uint8, uint64, error)
	GetBlockHeightByBlockHash(*common.Hash) (uint64, byte, error)
	GetCMBResponse([]byte) ([][]byte, error)
	GetDepositSend([]byte) ([]byte, error)
	GetWithdrawRequest([]byte) ([]byte, uint8, error)
	UpdateDCBBoard(transaction Transaction) error
	UpdateGOVBoard(transaction Transaction) error
	UpdateConstitution(transaction Transaction, boardType BoardType) error
	GetConstitution(boardType BoardType) ConstitutionInterface
	UpdateDCBFund(transaction Transaction)
	GetGovernor(boardType BoardType) GovernorInterface
}

// Interface for all types of metadata in tx
type Metadata interface {
	GetType() int
	Hash() *common.Hash
	CheckTransactionFee(Transaction, uint64) bool
	ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error)
	// isContinue, ok, err
	ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error)
	ValidateMetadataByItself() bool // TODO: need to define the method for metadata
	ValidateBeforeNewBlock(tx Transaction, bcr BlockchainRetriever, shardID byte) bool
	VerifyMultiSigs(Transaction, database.DatabaseInterface) (bool, error)
	BuildReqActions(tx Transaction, bcr BlockchainRetriever, shardID byte) ([][]string, error)
	ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error
	CalculateSize() uint64
}

// Interface for all type of transaction
type Transaction interface {
	Hash() *common.Hash
	ValidateTransaction(bool, database.DatabaseInterface, byte, *common.Hash) bool
	GetMetadataType() int
	GetType() string
	GetLockTime() int64
	GetTxActualSize() uint64
	GetSenderAddrLastByte() byte
	GetTxFee() uint64
	ListNullifiers() [][]byte
	CheckTxVersion(int8) bool
	CheckTransactionFee(minFeePerKbTx uint64) bool
	IsSalaryTx() bool
	ValidateTxWithCurrentMempool(MempoolRetriever) error
	ValidateTxWithBlockChain(BlockchainRetriever, byte, database.DatabaseInterface) error
	ValidateSanityData(BlockchainRetriever) (bool, error)
	ValidateTxByItself(bool, database.DatabaseInterface, BlockchainRetriever, byte) bool
	ValidateType() bool
	GetMetadata() Metadata
	SetMetadata(Metadata)
	GetInfo() []byte
	ValidateConstDoubleSpendWithBlockchain(BlockchainRetriever, byte, database.DatabaseInterface) error

	GetSigPubKey() []byte
	IsPrivacy() bool
	IsCoinsBurning() bool
	CalculateTxValue() uint64
	GetProof() *zkp.PaymentProof

	// Get receivers' data for tx
	GetReceivers() ([][]byte, []uint64)
	GetUniqueReceiver() (bool, []byte, uint64)

	// Get receivers' data for custom token tx (nil for normal tx)
	GetTokenReceivers() ([][]byte, []uint64)
	GetTokenUniqueReceiver() (bool, []byte, uint64)
	GetAmountOfVote() (uint64, error)
	GetVoterPaymentAddress() (*privacy.PaymentAddress, error)

	GetMetadataFromVinsTx(BlockchainRetriever) (Metadata, error)
}
