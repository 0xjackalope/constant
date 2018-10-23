package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

// Sign create signature for message with secret key
func Sign(hash []byte, spendingKey []byte) (r, s *big.Int, err error) {

	signingKey := new(ecdsa.PrivateKey)
	signingKey.PublicKey.Curve = Curve
	signingKey.D = new(big.Int).SetBytes(spendingKey)
	signingKey.PublicKey.X, signingKey.PublicKey.Y = Curve.ScalarBaseMult(spendingKey)

	r, s, err = ecdsa.Sign(rand.Reader, signingKey, hash[:])
	if err != nil {
		fmt.Printf("\nSigning Error: %v\n", err)
	}
	return

}

// Verify checks the signature that is signed by secret key corresponding with public key
func Verify(r, s *big.Int, hash []byte, address []byte) bool {

	verKey := new(ecdsa.PublicKey)
	verKey.Curve = Curve

	point := FromByteArrayToPoint(address)
	verKey.X = point.X
	verKey.Y = point.Y

	fmt.Println("verKey.X: %v", verKey.X)
	fmt.Println("verKey.Y: %v", verKey.Y)

	res := ecdsa.Verify(verKey, hash, r, s)
	return res
}
