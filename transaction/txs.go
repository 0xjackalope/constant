package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	// "crypto/sha256"

	"time"

	"math"

	"github.com/ninjadotorg/cash/common"
	"github.com/ninjadotorg/cash/privacy"
)

// Txs represents a coin-transfer-transaction stored in a block
type Txs struct {
	Version  int8   `json:"Version"`
	Type     string `json:"Type"` // n
	LockTime int64  `json:"LockTime"`
	Fee      uint64 `json:"Fee"` // Fee applies to first js desc

	SigPubKey            []byte `json:"JSPubKey,omitempty"` // 64 bytes
	Signature            []byte `json:"JSSig,omitempty"`    // 64 bytes
	rt                   []byte
	SNInputs             [][]byte // upspent serial number list
	CMOutputs            [][]byte // new commitment list
	PubEncKeys           [][]byte
	EncryptedOutputCoins [][]byte // if privacy,
	PlainOutputCoins     []*privacy.Coin
	Proof                []byte
	// Reward               uint64 `json:"Reward"` // For coinbase tx

	AddressLastByte byte `json:"AddressLastByte"`

	txID       *common.Hash
	sigPrivKey []byte `json:"sigPrivKey,omitempty"` // 64 bytes
}

// SetTxID sets transaction id
func (tx *Txs) SetTxID(txID *common.Hash) {
	tx.txID = txID
}

// GetTxID returns transaction id
func (tx *Txs) GetTxID() *common.Hash {
	return tx.txID
}

// Hash returns the hash of all fields of the transaction
func (tx *Txs) Hash() *common.Hash {
	record := strconv.Itoa(int(tx.Version))
	record += tx.Type
	record += strconv.FormatInt(tx.LockTime, 10)
	record += strconv.FormatUint(tx.Fee, 10)

	record += string(tx.SigPubKey)
	record += string(tx.rt)

	record += strconv.Itoa(len(tx.SNInputs))
	for _, sn := range tx.SNInputs {
		record += string(sn)
	}

	record += strconv.Itoa(len(tx.CMOutputs))
	for _, cm := range tx.CMOutputs {
		record += string(cm)
	}

	record += strconv.Itoa(len(tx.PubEncKeys))
	for _, pkenc := range tx.PubEncKeys {
		record += string(pkenc)
	}

	record += strconv.Itoa(len(tx.EncryptedOutputCoins))
	for _, encryptedCoin := range tx.EncryptedOutputCoins {
		record += string(encryptedCoin)
	}

	record += strconv.Itoa(len(tx.PlainOutputCoins))
	for _, plainCoin := range tx.PlainOutputCoins {
		plainCoinStr, _ := json.Marshal(plainCoin)
		record += string(plainCoinStr)
	}

	record += string(tx.Proof)

	record += string(tx.AddressLastByte)
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

// ValidateTransaction returns true if transaction is valid:
// - Signature matches the signing public key
// - JSDescriptions are valid (zk-snark proof satisfied)
// Note: This method doesn't check for double spending
func (tx *Txs) ValidateTransaction() bool {
	// TODO:
	return true

	// // Check for tx signature
	// tx.SetTxID(tx.Hash())
	// valid, err := VerifySign(tx)
	// if valid == false {
	// 	if err != nil {
	// 		fmt.Printf("Error verifying signature of tx: %+v", err)
	// 	}
	// 	return false
	// }

	// // Check each js desc
	// for txID, desc := range tx.Descs {
	// 	//if desc.Reward != 0 {
	// 	//	return false // Coinbase tx shouldn't be broadcasted across the network
	// 	//}

	// 	// Apply fee only to the first desc of tx
	// 	fee := uint64(0)
	// 	if txID == 0 {
	// 		fee = tx.Fee
	// 	}

	// 	nf1, nf2 := desc.Nullifiers[0], desc.Nullifiers[1]
	// 	hSig := client.HSigCRH(desc.HSigSeed, nf1, nf2, tx.JSPubKey)
	// 	valid, err := client.Verify(
	// 		desc.Proof,
	// 		desc.Nullifiers,
	// 		desc.Commitments,
	// 		desc.Anchor,
	// 		desc.Vmacs,
	// 		hSig,
	// 		desc.Reward,
	// 		fee,
	// 		tx.AddressLastByte,
	// 	)

	// 	if valid == false {
	// 		if err != nil {
	// 			fmt.Printf("Error validating tx: %+v\n", err)
	// 		}
	// 		return false
	// 	}
	// }

	// return true
}

// GetType returns the type of the transaction
func (tx *Txs) GetType() string {
	return tx.Type
}

// GetTxVirtualSize computes the virtual size of a given transaction in kilobyte
func (tx *Txs) GetTxVirtualSize() uint64 {
	var sizeVersion uint64 = 1  // int8
	var sizeType uint64 = 8     // string
	var sizeLockTime uint64 = 8 // int64
	var sizeFee uint64 = 8      // uint64

	var sizeSigPubKey uint64 = 33 // [33]byte
	var sizeSignature uint64 = 64 // [64]byte
	var sizeRt uint64 = 32        // [32]byte
	var sizeSNInputs = uint64(max(1, len(tx.SNInputs))) * 32
	var sizeCMOutputs = uint64(max(1, len(tx.CMOutputs))) * 32

	estimateTxSizeInByte := sizeVersion + sizeType + sizeLockTime + sizeFee + sizeSigPubKey + sizeSignature + sizeSNInputs + sizeCMOutputs + sizeRt
	return uint64(math.Ceil(float64(estimateTxSizeInByte) / 1024))
}

// func max(x, y int) int {
// 	if x > y {
// 		return x
// 	}
// 	return y
// }

// GetSenderAddrLastByte returns sender address last byte
func (tx *Txs) GetSenderAddrLastByte() byte {
	return tx.AddressLastByte
}

// CreateTxs creates transaction with appropriate proof for a private payment
// rts: mapping from the chainID to the root of the commitment merkle tree at current block
// 		(the latest block of the node creating this tx)
func CreateTxs(
	senderKey *privacy.SpendingKey,
	paymentInfo []*privacy.PaymentInfo,
	rt *common.Hash,
	unspentCoins []*privacy.Coin,
	spentSNs [][]byte,
	fee uint64,
	senderChainID byte,
) (*Txs, error) {
	// var tx *Txs
	tx, err := CreateEmptyTxs()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("List of all commitments before building tx:\n")
	// fmt.Printf("rts: %+v\n", rts)
	// for _, cm := range commitments {
	// 	fmt.Printf("%x\n", cm)
	// }

	var value uint64
	for _, p := range paymentInfo {
		value += p.Amount
		fmt.Printf("[CreateTx] paymentInfo.Value: %+v, paymentInfo.Adreess: %x\n", p.Amount, p.PaymentAddress.Pk)
	}

	// type ChainNote struct {
	// 	note    *client.Note
	// 	chainID byte
	// }

	// Get list of notes to use
	// var inputNotes []*Coin
	// for chainID, chainTxs := range usableTx {
	// 	for _, tx := range chainTxs {
	// 		for _, desc := range tx.Descs {
	// 			for _, note := range desc.Note {
	// 				chainNote := &ChainNote{note: note, chainID: chainID}
	// 				inputNotes = append(inputNotes, chainNote)
	// 				fmt.Printf("[CreateTx] inputNote.Value: %+v\n", note.Value)
	// 			}
	// 		}
	// 	}
	// }

	// Left side value
	// var sumInputValue uint64
	// for _, coin := range unspentCoins {
	// 	sumInputValue += coin.Value
	// }
	// if sumInputValue < value+fee {
	// 	return nil, fmt.Errorf("Input value less than output value")
	// }

	// senderFullKey := cashec.KeySet{}
	// senderFullKey.ImportFromPrivateKeyByte((*senderKey)[:])

	// // Create tx before adding js descs
	// tx, err := CreateEmptyTx()
	// if err != nil {
	// 	return nil, err
	// }
	// tempKeySet := cashec.KeySet{}
	// var temp privacy.SpendingKey
	// copy(temp[:], (*senderKey)[:])
	// tempKeySet.ImportFromPrivateKey(&temp)
	// lastByte := tempKeySet.PaymentAddress.Pk[len(tempKeySet.PaymentAddress.Pk)-1]
	// tx.AddressLastByte = lastByte
	// var latestAnchor map[byte][]byte

	// for len(inputNotes) > 0 || len(paymentInfo) > 0 {
	// 	// Sort input and output notes ascending by value to start building js descs
	// 	sort.Slice(inputNotes, func(i, j int) bool {
	// 		return inputNotes[i].note.Value < inputNotes[j].note.Value
	// 	})
	// 	sort.Slice(paymentInfo, func(i, j int) bool {
	// 		return paymentInfo[i].Amount < paymentInfo[j].Amount
	// 	})

	// 	// Choose inputs to build js desc
	// 	// var inputsToBuildWitness, inputs []*client.JSInput
	// 	inputsToBuildWitness := make(map[byte][]*client.JSInput)
	// 	inputs := make(map[byte][]*client.JSInput)
	// 	inputValue := uint64(0)
	// 	numInputNotes := 0
	// 	for len(inputNotes) > 0 && len(inputs) < NumDescInputs {
	// 		input := &client.JSInput{}
	// 		chainNote := inputNotes[len(inputNotes)-1] // Get note with largest value
	// 		input.InputNote = chainNote.note
	// 		var temp client.SpendingKey
	// 		copy(temp[:], (*senderKey)[:])
	// 		input.Key = &temp
	// 		inputs[chainNote.chainID] = append(inputs[chainNote.chainID], input)
	// 		inputsToBuildWitness[chainNote.chainID] = append(inputsToBuildWitness[chainNote.chainID], input)
	// 		inputValue += input.InputNote.Value

	// 		inputNotes = inputNotes[:len(inputNotes)-1]
	// 		numInputNotes++
	// 		fmt.Printf("Choose input note with value %+v and cm %x\n", input.InputNote.Value, input.InputNote.Cm)
	// 	}

	// 	var feeApply uint64 // Zero fee for js descs other than the first one
	// 	if len(tx.Descs) == 0 {
	// 		// First js desc, applies fee
	// 		feeApply = fee
	// 		tx.Fee = fee
	// 	}
	// 	if len(tx.Descs) == 0 {
	// 		if inputValue < feeApply {
	// 			return nil, fmt.Errorf("Input note values too small to pay fee")
	// 		}
	// 		inputValue -= feeApply
	// 	}

	// 	// Add dummy input note if necessary
	// 	for numInputNotes < NumDescInputs {
	// 		input := &client.JSInput{}
	// 		var temp privacy.SpendingKey
	// 		copy(temp[:], (*senderKey)[:])
	// 		input.InputNote = createDummyNote(&temp)
	// 		var temp1 client.SpendingKey
	// 		copy(temp1[:], (*senderKey)[:])
	// 		input.Key = &temp1
	// 		input.WitnessPath = (&client.MerklePath{}).CreateDummyPath() // No need to build commitment merkle path for dummy note
	// 		dummyNoteChainID := senderChainID                            // Dummy note's chain is the same as sender's
	// 		inputs[dummyNoteChainID] = append(inputs[dummyNoteChainID], input)
	// 		numInputNotes++
	// 		fmt.Printf("Add dummy input note\n")
	// 	}

	// 	// Check if input note's cm is in commitments list
	// 	for chainID, chainInputs := range inputsToBuildWitness {
	// 		for _, input := range chainInputs {
	// 			input.InputNote.Cm = client.GetCommitment(input.InputNote)

	// 			found := false
	// 			for _, c := range commitments[chainID] {
	// 				if bytes.Equal(c, input.InputNote.Cm) {
	// 					found = true
	// 				}
	// 			}
	// 			if found == false {
	// 				return nil, fmt.Errorf("Commitment %x of input note isn't in commitments list of chain %d", input.InputNote.Cm, chainID)
	// 			}
	// 		}
	// 	}

	// 	// Build witness path for the input notes
	// 	newRts, err := client.BuildWitnessPathMultiChain(inputsToBuildWitness, commitments)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// For first js desc, check if provided Rt is the root of the merkle tree contains all commitments
	// 	if latestAnchor == nil {
	// 		for chainID, rt := range newRts {
	// 			if !bytes.Equal(rt, rts[chainID][:]) {
	// 				return nil, fmt.Errorf("Provided anchor doesn't match commitments list: %d %x %x", chainID, rt, rts[chainID][:])
	// 			}
	// 		}
	// 	}
	// 	latestAnchor = newRts
	// 	// Add dummy anchor to for dummy inputs
	// 	if len(latestAnchor[senderChainID]) == 0 {
	// 		latestAnchor[senderChainID] = make([]byte, 32)
	// 	}

	// 	// Choose output notes for the js desc
	// 	outputs := []*client.JSOutput{}
	// 	for len(paymentInfo) > 0 && len(outputs) < NumDescOutputs-1 && inputValue > 0 { // Leave out 1 output note for change
	// 		p := paymentInfo[len(paymentInfo)-1]
	// 		var outNote *client.Note
	// 		var encKey []byte
	// 		if p.Amount <= inputValue { // Enough for one more output note, include it
	// 			outNote = &client.Note{Value: p.Amount, Apk: p.PaymentAddress.Pk}
	// 			encKey = p.PaymentAddress.Tk
	// 			inputValue -= p.Amount
	// 			paymentInfo = paymentInfo[:len(paymentInfo)-1]
	// 			fmt.Printf("Use output value %+v => %x\n", outNote.Value, outNote.Apk)
	// 		} else { // Not enough for this note, send some and save the rest for next js desc
	// 			outNote = &client.Note{Value: inputValue, Apk: p.PaymentAddress.Pk}
	// 			encKey = p.PaymentAddress.Tk
	// 			paymentInfo[len(paymentInfo)-1].Amount = p.Amount - inputValue
	// 			inputValue = 0
	// 			fmt.Printf("Partially send %+v to %x\n", outNote.Value, outNote.Apk)
	// 		}

	// 		var temp client.Tk
	// 		copy(temp[:], encKey[:])
	// 		output := &client.JSOutput{EncKey: temp, OutputNote: outNote}
	// 		outputs = append(outputs, output)
	// 	}

	// 	if inputValue > 0 {
	// 		// Still has some room left, check if one more output note is possible to add
	// 		var p *privacy.PaymentInfo
	// 		if len(paymentInfo) > 0 {
	// 			p = paymentInfo[len(paymentInfo)-1]
	// 		}

	// 		if p != nil && p.Amount == inputValue {
	// 			// Exactly equal, add this output note to js desc
	// 			outNote := &client.Note{Value: p.Amount, Apk: p.PaymentAddress.Pk}
	// 			var temp client.Tk
	// 			copy(temp[:], p.PaymentAddress.Tk[:])
	// 			output := &client.JSOutput{EncKey: temp, OutputNote: outNote}
	// 			outputs = append(outputs, output)
	// 			paymentInfo = paymentInfo[:len(paymentInfo)-1]
	// 			fmt.Printf("Exactly enough, include 1 more output %+v, %x\n", outNote.Value, outNote.Apk)
	// 		} else {
	// 			// Cannot put the output note into this js desc, create a change note instead
	// 			outNote := &client.Note{Value: inputValue, Apk: senderFullKey.PaymentAddress.Pk}
	// 			var temp client.Tk
	// 			copy(temp[:], p.PaymentAddress.Tk[:])
	// 			output := &client.JSOutput{EncKey: temp, OutputNote: outNote}
	// 			outputs = append(outputs, output)
	// 			fmt.Printf("Create change outnote %+v, %x\n", outNote.Value, outNote.Apk)

	// 			// Use the change note to continually send to receivers if needed
	// 			if len(paymentInfo) > 0 {
	// 				// outNote data (R and Rho) will be updated when building zk-proof
	// 				chainNote := &ChainNote{note: outNote, chainID: senderChainID}
	// 				inputNotes = append(inputNotes, chainNote)
	// 				fmt.Printf("Reuse change note later\n")
	// 			}
	// 		}
	// 		inputValue = 0
	// 	}

	// 	// Add dummy output note if necessary
	// 	for len(outputs) < NumDescOutputs {
	// 		outputs = append(outputs, CreateRandomJSOutput())
	// 		fmt.Printf("Create dummy output note\n")
	// 	}

	// 	// TODO: Shuffle output notes randomly (if necessary)

	// 	// Generate proof and sign tx
	// 	var reward uint64 // Zero reward for non-coinbase transaction
	// 	err = tx.BuildNewJSDesc(inputs, outputs, latestAnchor, reward, feeApply, false)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Add new commitments to list to use in next js desc if needed
	// 	for _, output := range outputs {
	// 		fmt.Printf("Add new output cm to list: %x\n", output.OutputNote.Cm)
	// 		commitments[senderChainID] = append(commitments[senderChainID], output.OutputNote.Cm)
	// 	}

	// 	fmt.Printf("Len input and info: %+v %+v\n", len(inputNotes), len(paymentInfo))
	// }

	// // Sign tx
	// tx, err = SignTx(tx)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Printf("jspubkey: %x\n", tx.JSPubKey)
	// fmt.Printf("jssig: %x\n", tx.JSSig)
	return tx, nil
}

// // BuildNewJSDesc creates zk-proof for a js desc and add it to the transaction
// func (tx *Txs) BuildNewJSDesc(
// 	inputMap map[byte][]*client.JSInput,
// 	outputs []*client.JSOutput,
// 	rtMap map[byte][]byte,
// 	reward, fee uint64,
// 	noPrivacy bool,
// ) error {
// 	// Gather inputs from different chains
// 	inputs := []*client.JSInput{}
// 	rts := [][]byte{}
// 	for chainID, inputList := range inputMap {
// 		for _, input := range inputList {
// 			inputs = append(inputs, input)
// 			rt, ok := rtMap[chainID]
// 			if !ok {
// 				return fmt.Errorf("Commitments merkle root not found for chain %d", chainID)
// 			}
// 			rts = append(rts, rt) // Input's corresponding merkle root
// 		}
// 	}
// 	if len(inputs) != NumDescInputs || len(outputs) != NumDescOutputs {
// 		return fmt.Errorf("Cannot build JSDesc with %d inputs and %d outputs", len(inputs), len(outputs))
// 	}

// 	var seed, phi []byte
// 	var outputR [][]byte
// 	proof, hSig, seed, phi, err := client.Prove(inputs, outputs, tx.JSPubKey, rts, reward, fee, seed, phi, outputR, tx.AddressLastByte)
// 	if noPrivacy {
// 		proof = nil
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	var ephemeralPrivKey *client.EphemeralPrivKey // nil ephemeral key, will be randomly created later
// 	err = tx.buildJSDescAndEncrypt(inputs, outputs, proof, rts, reward, hSig, seed, ephemeralPrivKey)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("jsPubKey: %x\n", tx.JSPubKey)
// 	fmt.Printf("jsSig: %x\n", tx.JSSig)
// 	return nil
// }

// func (tx *Txs) buildJSDescAndEncrypt(
// 	inputs []*client.JSInput,
// 	outputs []*client.JSOutput,
// 	proof *zksnark.PHGRProof,
// 	rts [][]byte,
// 	reward uint64,
// 	hSig, seed []byte,
// 	ephemeralPrivKey *client.EphemeralPrivKey,
// ) error {
// 	nullifiers := [][]byte{inputs[0].InputNote.Nf, inputs[1].InputNote.Nf}
// 	commitments := [][]byte{outputs[0].OutputNote.Cm, outputs[1].OutputNote.Cm}
// 	notes := [2]*client.Note{outputs[0].OutputNote, outputs[1].OutputNote}
// 	keys := [2]client.Tk{outputs[0].EncKey, outputs[1].EncKey}

// 	ephemeralPubKey := new(client.EphemeralPubKey)
// 	if ephemeralPrivKey == nil {
// 		ephemeralPrivKey = new(client.EphemeralPrivKey)
// 		*ephemeralPubKey, *ephemeralPrivKey = client.GenEphemeralKey()
// 	} else { // Genesis block only
// 		ephemeralPrivKey.GenPubKey()
// 		*ephemeralPubKey = ephemeralPrivKey.GenPubKey()
// 	}
// 	fmt.Printf("hSig: %x\n", hSig)
// 	fmt.Printf("ephemeralPrivKey: %x\n", *ephemeralPrivKey)
// 	fmt.Printf("ephemeralPubKey: %x\n", *ephemeralPubKey)
// 	fmt.Printf("tranmissionKey[0]: %x\n", keys[0])
// 	fmt.Printf("tranmissionKey[1]: %x\n", keys[1])
// 	fmt.Printf("notes[0].Value: %+v\n", notes[0].Value)
// 	fmt.Printf("notes[0].Rho: %x\n", notes[0].Rho)
// 	fmt.Printf("notes[0].R: %x\n", notes[0].R)
// 	fmt.Printf("notes[0].Memo: %+v\n", notes[0].Memo)
// 	fmt.Printf("notes[1].Value: %+v\n", notes[1].Value)
// 	fmt.Printf("notes[1].Rho: %x\n", notes[1].Rho)
// 	fmt.Printf("notes[1].R: %x\n", notes[1].R)
// 	fmt.Printf("notes[1].Memo: %+v\n", notes[1].Memo)
// 	var noteciphers [][]byte
// 	if proof != nil {
// 		noteciphers = client.EncryptNote(notes, keys, *ephemeralPrivKey, *ephemeralPubKey, hSig)
// 	}

// 	//Calculate vmacs to prove this transaction is signed by this user
// 	vmacs := make([][]byte, 2)
// 	for i := range inputs {
// 		ask := make([]byte, 32)
// 		copy(ask[:], inputs[i].Key[:])
// 		vmacs[i] = client.PRF_pk(uint64(i), ask, hSig)
// 	}

// 	desc := &JoinSplitDesc{
// 		Anchor:          rts,
// 		Nullifiers:      nullifiers,
// 		Commitments:     commitments,
// 		Proof:           proof,
// 		EncryptedData:   noteciphers,
// 		EphemeralPubKey: ephemeralPubKey[:],
// 		HSigSeed:        seed,
// 		Type:            common.AssetTypeCoin,
// 		Reward:          reward,
// 		Vmacs:           vmacs,
// 	}
// 	tx.Descs = append(tx.Descs, desc)
// 	if desc.Proof == nil { // no privacy
// 		desc.Note = []*client.Note{outputs[0].OutputNote, outputs[1].OutputNote}
// 	}

// 	fmt.Println("desc:")
// 	fmt.Printf("Anchor: %x\n", desc.Anchor)
// 	fmt.Printf("Nullifiers: %x\n", desc.Nullifiers)
// 	fmt.Printf("Commitments: %x\n", desc.Commitments)
// 	fmt.Printf("Proof: %x\n", desc.Proof)
// 	fmt.Printf("EncryptedData: %x\n", desc.EncryptedData)
// 	fmt.Printf("EphemeralPubKey: %x\n", desc.EphemeralPubKey)
// 	fmt.Printf("HSigSeed: %x\n", desc.HSigSeed)
// 	fmt.Printf("Type: %+v\n", desc.Type)
// 	fmt.Printf("Reward: %+v\n", desc.Reward)
// 	fmt.Printf("Vmacs: %x %x\n", desc.Vmacs[0], desc.Vmacs[1])
// 	return nil
// }

// // CreateRandomJSInput creates a dummy input with 0 value note that belongs to a random address
// func CreateRandomJSInput(spendingKey *client.SpendingKey) *client.JSInput {
// 	if spendingKey == nil {
// 		randomKey := client.RandSpendingKey()
// 		spendingKey = &randomKey
// 	}
// 	var temp privacy.SpendingKey
// 	copy(temp, spendingKey[:])
// 	input := new(client.JSInput)
// 	input.InputNote = createDummyNote(&temp)
// 	input.Key = spendingKey
// 	input.WitnessPath = (&client.MerklePath{}).CreateDummyPath()
// 	return input
// }

// // CreateRandomJSOutput creates a dummy output with 0 value note that is sended to a random address
// func CreateRandomJSOutput() *client.JSOutput {
// 	randomKey := client.RandSpendingKey()
// 	output := new(client.JSOutput)
// 	var spendingKey privacy.SpendingKey
// 	copy(spendingKey, randomKey[:])
// 	output.OutputNote = createDummyNote(&spendingKey)
// 	paymentAddr := client.GenPaymentAddress(randomKey)
// 	output.EncKey = paymentAddr.Pkenc
// 	return output
// }

// func createDummyNote(spendingKey *privacy.SpendingKey) *client.Note {
// 	addr := privacy.GeneratePublicKey(*spendingKey)
// 	var rho, r [32]byte
// 	copy(rho[:], client.RandBits(32*8))
// 	copy(r[:], client.RandBits(32*8))
// 	var temp client.SpendingKey
// 	copy(temp[:], (*spendingKey)[:])
// 	note := &client.Note{
// 		Value: 0,
// 		Apk:   addr,
// 		Rho:   rho[:],
// 		R:     r[:],
// 		Nf:    client.GetNullifier(temp, rho),
// 	}
// 	return note
// }

// Sign signs transaction using ECDSA
func Sign(tx *Txs) (*Txs, error) {
	//Check input transaction
	if tx.Signature != nil {
		return nil, errors.New("Input transaction must be an unsigned one")
	}

	// Hash transaction
	tx.SetTxID(tx.Hash())
	hash := tx.GetTxID()
	data := make([]byte, common.HashSize)
	copy(data, hash[:])

	// Sign
	signature, err := privacy.Sign(data[:], tx.sigPrivKey)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature
	return tx, nil
}

// Verify verify the transaction's signature
func Verify(tx *Txs) (bool, error) {
	//Check input transaction
	if tx.Signature == nil || tx.SigPubKey == nil {
		return false, errors.New("input transaction must be an signed one")
	}

	// Hash origin transaction
	hash := tx.GetTxID()
	data := make([]byte, common.HashSize)
	copy(data, hash[:])

	valid := privacy.Verify(tx.Signature, data[:], tx.SigPubKey)
	return valid, nil
}

// // GenerateProofForGenesisTx creates zk-proof and build the transaction (without signing) for genesis block
// func GenerateProofForGenesisTx(
// 	inputs []*client.JSInput,
// 	outputs []*client.JSOutput,
// 	rts [][]byte,
// 	reward uint64,
// 	seed, phi []byte,
// 	outputR [][]byte,
// 	ephemeralPrivKey client.EphemeralPrivKey,
// ) (*Tx, error) {
// 	// Generate JoinSplit key pair to act as a dummy key (since we don't sign genesis tx)
// 	privateSignKey := [32]byte{1}
// 	keySet := &cashec.KeySet{}
// 	keySet.ImportFromPrivateKeyByte(privateSignKey[:])
// 	sigPubKey := keySet.PaymentAddress.Pk[:]

// 	// Get last byte of genesis sender's address
// 	tempKeySet := cashec.KeySet{}
// 	var temp privacy.SpendingKey
// 	copy(temp[:], inputs[0].Key[:])
// 	tempKeySet.ImportFromPrivateKey(&temp)
// 	addressLastByte := tempKeySet.PaymentAddress.Pk[len(tempKeySet.PaymentAddress.Pk)-1]

// 	tx, err := CreateEmptyTx()
// 	if err != nil {
// 		return nil, err
// 	}
// 	tx.JSPubKey = sigPubKey
// 	tx.AddressLastByte = addressLastByte
// 	fmt.Printf("JSPubKey: %x\n", tx.JSPubKey)

// 	var fee uint64 // Zero fee for genesis tx
// 	proof, hSig, seed, phi, err := client.Prove(
// 		inputs,
// 		outputs,
// 		tx.JSPubKey,
// 		rts,
// 		reward,
// 		fee,
// 		seed,
// 		phi,
// 		outputR,
// 		addressLastByte,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = tx.buildJSDescAndEncrypt(inputs, outputs, proof, rts, reward, hSig, seed, &ephemeralPrivKey)
// 	return tx, err
// }

// func PubKeyToByteArray(pubKey *client.PaymentAddress) []byte {
// 	var pub []byte
// 	pubX := pubKey.X.Bytes()
// 	pubY := pubKey.Y.Bytes()
// 	pub = append(pub, pubX...)
// 	pub = append(pub, pubY...)
// 	return pub
// }

// func JSSigToByteArray(jsSig *client.EcdsaSignature) []byte {
// 	var jssig []byte
// 	r := jsSig.R.Bytes()
// 	s := jsSig.S.Bytes()
// 	jssig = append(jssig, r...)
// 	jssig = append(jssig, s...)
// 	return jssig
// }

// func SortArrayTxs(data []Tx, sortType int, sortAsc bool) {
// 	if len(data) == 0 {
// 		return
// 	}
// 	switch sortType {
// 	case NoSort:
// 		{
// 			// do nothing
// 		}
// 	case SortByAmount:
// 		{
// 			sort.SliceStable(data, func(i, j int) bool {
// 				desc1 := data[i].Descs
// 				amount1 := uint64(0)
// 				for _, desc := range desc1 {
// 					for _, note := range desc.GetNote() {
// 						amount1 += note.Value
// 					}
// 				}
// 				desc2 := data[j].Descs
// 				amount2 := uint64(0)
// 				for _, desc := range desc2 {
// 					for _, note := range desc.GetNote() {
// 						amount2 += note.Value
// 					}
// 				}
// 				if !sortAsc {
// 					return amount1 > amount2
// 				} else {
// 					return amount1 <= amount2
// 				}
// 			})
// 		}
// 	default:
// 		{
// 			// do nothing
// 		}
// 	}
// }

// EstimateTxSize returns the estimated size of the tx in kilobyte
func EstimateTxsSize(usableTx []*Tx, payments []*privacy.PaymentInfo) uint64 {
	var sizeVersion uint64 = 1  // int8
	var sizeType uint64 = 8     // string
	var sizeLockTime uint64 = 8 // int64
	var sizeFee uint64 = 8      // uint64
	var sizeDescs = uint64(max(1, (len(usableTx)+len(payments)-3))) * EstimateJSDescSize()
	var sizejSPubKey uint64 = 64 // [64]byte
	var sizejSSig uint64 = 64    // [64]byte
	estimateTxSizeInByte := sizeVersion + sizeType + sizeLockTime + sizeFee + sizeDescs + sizejSPubKey + sizejSSig
	return uint64(math.Ceil(float64(estimateTxSizeInByte) / 1024))
}

// CreateEmptyTx returns a new Tx initialized with default data
func CreateEmptyTxs() (*Txs, error) {
	//Generate key pair
	sigPrivKey, sigPubKey := privacy.GenerateKey()

	tx := &Txs{
		Version:  TxVersion,
		Type:     common.TxNormalType,
		LockTime: time.Now().Unix(),
		Fee:      0,

		SigPubKey: sigPubKey,
		Signature: nil,
		rt:        nil,
		SNInputs:  nil,
		CMOutputs: nil,

		AddressLastByte: 0,

		txID:       nil,
		sigPrivKey: sigPrivKey,
	}
	return tx, nil
}
