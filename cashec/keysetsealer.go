package cashec

/*type KeySetSealer struct {
	SprivateKey     []byte
	SpublicKey      []byte
	SpendingAddress []byte
	TransmissionKey []byte
	ReceivingKey    []byte
}

func (self *KeySetSealer) GenerateKey(seed []byte) (*KeySetSealer, error) {
	var err error
	self.SpublicKey, self.SprivateKey, err = ed25519.GenerateKey(bytes.NewBuffer(seed))
	if err != nil {
		return self, err
	}
	return self, nil
}

func (self *KeySetSealer) Import(privateKey string) (*KeySetSealer, error) {
	key := ed25519.PrivateKey{}
	base58C := base58.Base58Check{}
	key, _, err := base58C.Decode(privateKey)
	if err != nil {
		return self, err
	}
	self.SpublicKey = key.Public().(ed25519.PaymentAddress)
	self.SprivateKey = key
	return self, nil
}

func (self *KeySetSealer) Verify(data, signature []byte) (bool, error) {
	isValid := false
	isValid = ed25519.Verify(self.SpublicKey, data, signature)
	return isValid, nil
}

func (self *KeySetSealer) Sign(data []byte) ([]byte, error) {
	result := ed25519.Sign(self.SprivateKey, data)
	return result, nil
}

func (self *KeySetSealer) EncodeToString() string {
	val, _ := json.Marshal(self)
	result := base58.Base58Check{}.Encode(val, byte(0x00))
	return result
}

func (self *KeySetSealer) DecodeToKeySet(keystring string) (*KeySetSealer, error) {
	base58C := base58.Base58Check{}
	keyBytes, _, _ := base58C.Decode(keystring)
	json.Unmarshal(keyBytes, self)
	return self, nil
}

func (self *KeySetSealer) GetPaymentAddress() (privacy.PaymentAddress, error) {
	var paymentAddr privacy.PaymentAddress
	paymentAddr.PublicKey = self.SpendingAddress
	paymentAddr.TransmissionKey = self.TransmissionKey
	return paymentAddr, nil
}

func (self *KeySetSealer) GetViewingKey() (privacy.ViewingKey, error) {
	var viewingKey privacy.ViewingKey
	viewingKey.PublicKey = self.SpendingAddress
	viewingKey.ReceivingKey = self.ReceivingKey
	return viewingKey, nil
}*/

/*func ValidateDataB58(pubkey string, sig string, data []byte) error {
	decPubkey, _, err := base58.Base58Check{}.Decode(pubkey)
	if err != nil {
		return errors.New("can't decode public key:" + err.Error())
	}

	validatorKp := KeySetSealer{
		SpublicKey: decPubkey,
	}
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
}*/
