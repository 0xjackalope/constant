package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

// Sign create signature for message with secret key
func Sign(hash []byte, spendingKey []byte) (signature []byte, err error) {

	signingKey := new(ecdsa.PrivateKey)
	signingKey.PublicKey.Curve = Curve
	signingKey.D = new(big.Int).SetBytes(spendingKey)
	signingKey.PublicKey.X, signingKey.PublicKey.Y = Curve.ScalarBaseMult(spendingKey)

	r, s, err := ecdsa.Sign(rand.Reader, signingKey, hash[:])
	if err != nil {
		return nil, err
	}
	signature = SigToByteArray(r, s)
	return

}

// Verify checks the signature that is signed by secret key corresponding with public key
func Verify(signature []byte, hash []byte, address []byte) bool {
	r, s := FromByteArrayToSig(signature)

	verKey := new(ecdsa.PublicKey)
	verKey.Curve = Curve

	point := FromByteArrayToPoint(address)
	verKey.X = point.X
	verKey.Y = point.Y

	fmt.Printf("verKey.X: %v", verKey.X)
	fmt.Printf("verKey.Y: %v", verKey.Y)

	res := ecdsa.Verify(verKey, hash, r, s)
	return res
}

// SigToByteArray converts signature to byte array
func SigToByteArray(r, s *big.Int) (sig []byte) {
	sig = append(sig, r.Bytes()...)
	sig = append(sig, s.Bytes()...)
	return
}

// FromByteArrayToSig converts a byte array to signature
func FromByteArrayToSig(sig []byte) (r, s *big.Int) {
	r = new(big.Int).SetBytes(sig[0:32])
	s = new(big.Int).SetBytes(sig[32:64])
	return
}
