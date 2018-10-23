package cashec

import (
	"github.com/ninjadotorg/cash/common"
	"github.com/ninjadotorg/cash/privacy"
	"encoding/json"
	"github.com/ninjadotorg/cash/common/base58"
	"errors"
)

type KeySet struct {
	PrivateKey  privacy.SpendingKey
	PublicKey   privacy.PaymentAddress
	ReadonlyKey privacy.ViewingKey
}

/*
GenerateKey - generate key set from seed byte[]
*/
func (self *KeySet) GenerateKey(seed []byte) *KeySet {
	self.PrivateKey = privacy.GenerateSpendingKey(seed)
	self.PublicKey = privacy.GeneratePaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenerateViewingKey(self.PrivateKey)
	return self
}

/*
ImportFromPrivateKeyByte - from private-key byte[], regenerate pub-key and readonly-key
*/
func (self *KeySet) ImportFromPrivateKeyByte(privateKey []byte) {
	copy(self.PrivateKey[:], privateKey)
	self.PublicKey = privacy.GeneratePaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenerateViewingKey(self.PrivateKey)
}

/*
ImportFromPrivateKeyByte - from private-key data, regenerate pub-key and readonly-key
*/
func (self *KeySet) ImportFromPrivateKey(privateKey *privacy.SpendingKey) {
	self.PrivateKey = *privateKey
	self.PublicKey = privacy.GeneratePaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenerateViewingKey(self.PrivateKey)
}

/*
Generate Sealer keyset from privacy key set
*/
func (self *KeySet) CreateSealerKeySet() (*KeySetSealer, error) {
	var sealerKeySet KeySetSealer
	sealerKeySet.GenerateKey(self.PrivateKey[:])
	sealerKeySet.SpendingAddress = self.PublicKey.Address
	sealerKeySet.TransmissionKey = self.PublicKey.TransmissionKey
	sealerKeySet.ReceivingKey = self.ReadonlyKey.ReceivingKey
	return &sealerKeySet, nil
}

func (self *KeySet) Verify(data, signature []byte) (bool, error) {
	isValid := false
	hash := common.HashB(data)
	isValid = privacy.Verify(signature, hash[:], self.PublicKey.Address)
	return isValid, nil
}

func (self *KeySet) Sign(data []byte) ([]byte, error) {
	hash := common.HashB(data)
	signature, err := privacy.Sign(hash[:], self.PrivateKey)
	return signature, err
}

func (self *KeySet) EncodeToString() string {
	val, _ := json.Marshal(self)
	result := base58.Base58Check{}.Encode(val, byte(0x00))
	return result
}

func (self *KeySet) DecodeToKeySet(keystring string) (*KeySet, error) {
	base58C := base58.Base58Check{}
	keyBytes, _, _ := base58C.Decode(keystring)
	json.Unmarshal(keyBytes, self)
	return self, nil
}

func (self *KeySet) GetPaymentAddress() (privacy.PaymentAddress, error) {
	return self.PublicKey, nil
}

func (self *KeySet) GetViewingKey() (privacy.ViewingKey, error) {
	return self.ReadonlyKey, nil
}

func ValidateDataB58(pubkey string, sig string, data []byte) error {
	decPubkey, _, err := base58.Base58Check{}.Decode(pubkey)
	if err != nil {
		return errors.New("can't decode public key:" + err.Error())
	}

	validatorKp := KeySet{}
	validatorKp.PublicKey.Address = decPubkey
	decSig, _, err := base58.Base58Check{}.Decode(sig)
	if err != nil {
		return errors.New("can't decode signature: " + err.Error())
	}

	isValid, err := validatorKp.Verify(data, decSig)
	if err != nil {
		return errors.New("error when verify data: " + err.Error())
	}
	if !isValid {
		return errors.New("Invalid signature")
	}
	return nil
}
