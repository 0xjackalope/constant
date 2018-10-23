package cashec

import (
	"github.com/ninjadotorg/cash/common"
	"github.com/ninjadotorg/cash/privacy"
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
