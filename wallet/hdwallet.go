package wallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"github.com/ninjadotorg/cash/cashec"
	"github.com/ninjadotorg/cash/common/base58"
	"github.com/ninjadotorg/cash/common"
)

const (
	PriKeyType      = byte(0x0)
	PubKeyType      = byte(0x1)
	ReadonlyKeyType = byte(0x2)
)

// KeySet represents a bip32 extended Key
type Key struct {
	Depth       byte   // 1 bytes
	ChildNumber []byte // 4 bytes
	ChainCode   []byte // 32 bytes
	KeySet      cashec.KeySet
}

// NewMasterKey creates a new master extended Key from a Seed
func NewMasterKey(seed []byte) (*Key, error) {
	// Generate Key and chaincode
	hmac := hmac.New(sha512.New, seed)
	_, err := hmac.Write(seed)
	if err != nil {
		Logger.log.Error(err)
		return nil, err
	}
	intermediary := hmac.Sum(nil)

	// Split it into our Key and chain code
	keyBytes := intermediary[:32]  // use to create master private/public keypair
	chainCode := intermediary[32:] // be used with public Key (in keypair) for new Child keys

	// Validate Key
	/*err = validatePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}*/

	keySet := (&cashec.KeySet{}).GenerateKey(keyBytes)

	// Create the Key struct
	key := &Key{
		ChainCode:   chainCode,
		KeySet:      *keySet,
		Depth:       0x00,
		ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
	}

	return key, nil
}

// NewChildKey derives a Child Key from a given parent as outlined by bip32
func (key *Key) NewChildKey(childIdx uint32) (*Key, error) {
	intermediary, err := key.getIntermediary(childIdx)
	if err != nil {
		return nil, err
	}

	newSeed := []byte{}
	newSeed = append(newSeed[:], intermediary[:32]...)
	newKeyset := (&cashec.KeySet{}).GenerateKey(newSeed)
	// Create Child KeySet with data common to all both scenarios
	childKey := &Key{
		ChildNumber: uint32Bytes(childIdx),
		ChainCode:   intermediary[32:],
		Depth:       key.Depth + 1,
		KeySet:      *newKeyset,
	}

	return childKey, nil
}

func (key *Key) getIntermediary(childIdx uint32) ([]byte, error) {
	childIndexBytes := uint32Bytes(childIdx)

	var data []byte
	data = append(data, childIndexBytes...)

	hmac := hmac.New(sha512.New, key.ChainCode)
	_, err := hmac.Write(data)
	if err != nil {
		return nil, err
	}
	return hmac.Sum(nil), nil
}

// Serialize a KeySet to a 78 byte byte slice
func (key *Key) Serialize(keyType byte) ([]byte, error) {
	// Write fields to buffer in order
	buffer := new(bytes.Buffer)
	buffer.WriteByte(keyType)
	if keyType == PriKeyType {

		buffer.WriteByte(key.Depth)
		buffer.Write(key.ChildNumber)
		buffer.Write(key.ChainCode)

		// Private keys should be prepended with a single null byte
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(key.KeySet.PrivateKey))) // set length
		keyBytes = append(keyBytes, key.KeySet.PrivateKey[:]...)      // set pri-key
		buffer.Write(keyBytes)
	} else if keyType == PubKeyType {
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(key.KeySet.PaymentAddress.Pk))) // set length Apk
		keyBytes = append(keyBytes, key.KeySet.PaymentAddress.Pk[:]...)      // set Apk

		keyBytes = append(keyBytes, byte(len(key.KeySet.PaymentAddress.Tk))) // set length Pkenc
		keyBytes = append(keyBytes, key.KeySet.PaymentAddress.Tk[:]...)      // set Pkenc
		buffer.Write(keyBytes)
	} else if keyType == ReadonlyKeyType {
		keyBytes := make([]byte, 0)
		keyBytes = append(keyBytes, byte(len(key.KeySet.ReadonlyKey.Pk))) // set length Apk
		keyBytes = append(keyBytes, key.KeySet.ReadonlyKey.Pk[:]...)      // set Apk

		keyBytes = append(keyBytes, byte(len(key.KeySet.ReadonlyKey.Rk))) // set length Skenc
		keyBytes = append(keyBytes, key.KeySet.ReadonlyKey.Rk[:]...)      // set Pkenc
		buffer.Write(keyBytes)
	}

	// Append the standard doublesha256 checksum
	serializedKey, err := addChecksumToBytes(buffer.Bytes())
	if err != nil {
		Logger.log.Error(err)
		return nil, err
	}

	return serializedKey, nil
}

// Base58CheckSerialize encodes the KeySet in the standard Bitcoin base58 encoding
func (key *Key) Base58CheckSerialize(keyType byte) string {
	serializedKey, err := key.Serialize(keyType)
	if err != nil {
		return common.EmptyString
	}

	return base58.Base58Check{}.Encode(serializedKey, byte(0x00))
}

// Deserialize a byte slice into a KeySet
func Deserialize(data []byte) (*Key, error) {
	var key = &Key{}
	keyType := data[0]
	if keyType == PriKeyType {
		key.Depth = data[1]
		key.ChildNumber = data[2:6]
		key.ChainCode = data[6:38]
		keyLength := int(data[38])
		key.KeySet.PrivateKey = make([]byte, keyLength)
		copy(key.KeySet.PrivateKey[:], data[39:39+keyLength])
	} else if keyType == PubKeyType {
		apkKeyLength := int(data[1])
		key.KeySet.PaymentAddress.Pk = make([]byte, apkKeyLength)
		copy(key.KeySet.PaymentAddress.Pk[:], data[2:2+apkKeyLength])
		pkencKeyLength := int(data[apkKeyLength+2])
		key.KeySet.PaymentAddress.Tk = make([]byte, pkencKeyLength)
		copy(key.KeySet.PaymentAddress.Tk[:], data[3+apkKeyLength:3+apkKeyLength+pkencKeyLength])
	} else if keyType == ReadonlyKeyType {
		apkKeyLength := int(data[1])
		key.KeySet.PaymentAddress.Pk = make([]byte, apkKeyLength)
		copy(key.KeySet.ReadonlyKey.Pk[:], data[2:2+apkKeyLength])
		skencKeyLength := int(data[apkKeyLength+2])
		key.KeySet.ReadonlyKey.Rk = make([]byte, skencKeyLength)
		copy(key.KeySet.ReadonlyKey.Rk[:], data[3+apkKeyLength:3+apkKeyLength+skencKeyLength])
	}

	// validate checksum
	cs1 := base58.ChecksumFirst4Bytes(data[0: len(data)-4])
	cs2 := data[len(data)-4:]
	for i := range cs1 {
		if cs1[i] != cs2[i] {
			return nil, NewWalletError(InvalidChecksumErr, nil)
		}
	}
	return key, nil
}

// Base58CheckDeserialize deserializes a KeySet encoded in base58 encoding
func Base58CheckDeserialize(data string) (*Key, error) {
	b, _, err := base58.Base58Check{}.Decode(data)
	if err != nil {
		return nil, err
	}
	return Deserialize(b)
}
