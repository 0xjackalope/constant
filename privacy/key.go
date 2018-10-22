package privacy

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// Using curve Secp256k1 with package github.com/ethereum/go-ethereum/crypto/secp256k1
var curve = secp256k1.S256()

// fmt.Printf("N: %v\n", curve.N)
// fmt.Printf("P: %v\n", curve.P)
// fmt.Printf("B: %v\n", curve.B)
// fmt.Printf("Gx: %v\n", curve.Gx)
// fmt.Printf("Gy: %v\n", curve.Gy)
// fmt.Printf("BitSize: %v\n", curve.BitSize)

// secret x: SpendingKey: []byte, 32 bytes
type SpendingKey []byte

// EllipticPoint represents an point of ellipctic secp256k1
type EllipticPoint struct {
	X, Y *big.Int
}

// ViewingKey represents an key that be used to view transactions
type ViewingKey struct {
	Address      []byte // 33 bytes, use to receive coin
	ReceivingKey []byte // 32 bytes, use to decrypt data
}

// PaymentAddress represents an payment address of receiver
type PaymentAddress struct {
	Address         []byte // 33 bytes, use to receive coin
	TransmissionKey []byte // 33 bytes, use to encrypt data
}
type PaymentInfo struct {
	PaymentAddress PaymentAddress
	Amount         uint64
}

// RandBits generates random bits and return as bytes; zero out redundant bits
func RandBits(n int) []byte {
	m := 1 + (n-1)/8
	b := make([]byte, m)
	rand.Read(b)

	if n%8 > 0 {
		b[m-1] &= ((1 << uint(n%8)) - 1)
	}
	return b
}

// GenSpendingKey generates a random SpendingKey
// SpendingKey: 32 bytes
func GenSpendingKey() []byte {
	spendingKey := RandBits(256)
	return spendingKey
}

// GenAddress computes an address corresponding with spendingKey
// Address : 33 bytes
func GenAddress(spendingKey []byte) []byte {
	var p EllipticPoint
	p.X, p.Y = curve.ScalarBaseMult(spendingKey)
	address := secp256k1.CompressPubkey(p.X, p.Y)
	return address
}

// GenReceivingKey computes a receiving key corresponding with spendingKey
// ReceivingKey : 32 bytes
func GenReceivingKey(spendingKey []byte) []byte {
	hash := sha256.Sum256(spendingKey)
	receivingKey := make([]byte, 32)
	copy(receivingKey, hash[:])
	return receivingKey
}

// GenTransmissionKey computes a transmission key corresponding with receivingKey
// TransmissionKey : 33 bytes
func GenTransmissionKey(receivingKey []byte) []byte {
	var p, generator EllipticPoint
	random := RandBits(256)
	//create new generator from base generator
	generator.X, generator.Y = curve.ScalarBaseMult(random)

	p.X, p.Y = curve.ScalarMult(generator.X, generator.Y, receivingKey)
	transmissionKey := secp256k1.CompressPubkey(p.X, p.Y)
	return transmissionKey
}

// GenViewingKey generates a viewingKey corressponding with spendingKey
func GenViewingKey(spendingKey []byte) ViewingKey {
	var viewingKey ViewingKey
	viewingKey.Address = GenAddress(spendingKey)
	viewingKey.ReceivingKey = GenReceivingKey(spendingKey)
	return viewingKey
}

// GenPaymentAddress generates a payment address corressponding with spendingKey
func GenPaymentAddress(spendingKey []byte) PaymentAddress {
	var paymentAddress PaymentAddress
	paymentAddress.Address = GenAddress(spendingKey)
	paymentAddress.TransmissionKey = GenTransmissionKey(GenReceivingKey(spendingKey))
	return paymentAddress
}
