package cashec

import (
	"github.com/ninjadotorg/cash/common"
	"github.com/ninjadotorg/cash/privacy"
)

type KeySet struct {
	// SealerKeyPair KeyPair
	PrivateKey  privacy.SpendingKey
	PublicKey   privacy.PaymentAddress
	ReadonlyKey privacy.ViewingKey
}

/*
GenerateKey - generate key set from seed byte[]
*/
func (self *KeySet) GenerateKey(seed []byte) *KeySet {
	hash := common.HashB(seed)
	hash[len(hash)-1] &= 0x0F // Private key only has 252 bits
	copy(self.PrivateKey[:], hash)
	self.PublicKey = privacy.GenPaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenViewingKey(self.PrivateKey)
	return self
}

/*
ImportFromPrivateKeyByte - from private-key byte[], regenerate pub-key and readonly-key
*/
func (self *KeySet) ImportFromPrivateKeyByte(privateKey []byte) {
	copy(self.PrivateKey[:], privateKey)
	self.PublicKey = privacy.GenPaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenViewingKey(self.PrivateKey)
}

/*
ImportFromPrivateKeyByte - from private-key data, regenerate pub-key and readonly-key
*/
func (self *KeySet) ImportFromPrivateKey(privateKey *privacy.SpendingKey) {
	self.PrivateKey = *privateKey
	self.PublicKey = privacy.GenPaymentAddress(self.PrivateKey)
	self.ReadonlyKey = privacy.GenViewingKey(self.PrivateKey)
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
