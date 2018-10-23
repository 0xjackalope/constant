package privacy

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ninjadotorg/cash/common"
)

// Curve P256
var Curve = elliptic.P256()

// fmt.Printf("N: %v\n", curve.N)
// fmt.Printf("P: %v\n", curve.P)
// fmt.Printf("B: %v\n", curve.B)
// fmt.Printf("Gx: %v\n", curve.Gx)
// fmt.Printf("Gy: %v\n", curve.Gy)
// fmt.Printf("BitSize: %v\n", curve.BitSize)

// SpendingKey 32 bytes
type SpendingKey []byte

// EllipticPoint represents an point of ellipctic secp256k1
type EllipticPoint struct {
	X, Y *big.Int
}

// ViewingKey represents an key that be used to view transactions
type ViewingKey struct {
	Address      []byte // 64 bytes, use to receive coin
	ReceivingKey []byte // 32 bytes, use to decrypt pointByte
}

// PaymentAddress represents an payment address of receiver
type PaymentAddress struct {
	Address         []byte // 64 bytes, use to receive coin
	TransmissionKey []byte // 64 bytes, use to encrypt pointByte
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

// GenerateSpendingKey generates a random SpendingKey
// SpendingKey: 32 bytes
func GenerateSpendingKey(seed []byte) []byte {
	temp := new(big.Int)
	spendingKey := make([]byte, 32)
	spendingKey = common.HashB(seed)
	for temp.SetBytes(spendingKey).Cmp(Curve.Params().N) == 1 {
		spendingKey = common.HashB(spendingKey)
	}

	return spendingKey
}

// GenerateAddress computes an address corresponding with spendingKey
// Address : 64 bytes
func GenerateAddress(spendingKey []byte) []byte {
	var p EllipticPoint
	p.X, p.Y = Curve.ScalarBaseMult(spendingKey)
	fmt.Printf("p.X: %v", p.X)
	fmt.Printf("p.Y: %v", p.Y)
	address := FromPointToByteArray(p)
	return address
}

// GenerateReceivingKey computes a receiving key corresponding with spendingKey
// ReceivingKey : 32 bytes
func GenerateReceivingKey(spendingKey []byte) []byte {
	hash := sha256.Sum256(spendingKey)
	receivingKey := make([]byte, 32)
	copy(receivingKey, hash[:])
	return receivingKey
}

// GenerateTransmissionKey computes a transmission key corresponding with receivingKey
// TransmissionKey : 64 bytes
func GenerateTransmissionKey(receivingKey []byte) []byte {
	var p, generator EllipticPoint
	random := RandBits(256)
	//create new generator from base generator
	generator.X, generator.Y = Curve.Params().ScalarBaseMult(random)

	p.X, p.Y = Curve.Params().ScalarMult(generator.X, generator.Y, receivingKey)
	transmissionKey := FromPointToByteArray(p)
	return transmissionKey
}

// GenerateViewingKey generates a viewingKey corressponding with spendingKey
func GenerateViewingKey(spendingKey []byte) ViewingKey {
	var viewingKey ViewingKey
	viewingKey.Address = GenerateAddress(spendingKey)
	viewingKey.ReceivingKey = GenerateReceivingKey(spendingKey)
	return viewingKey
}

// GeneratePaymentAddress generates a payment address corressponding with spendingKey
func GeneratePaymentAddress(spendingKey []byte) PaymentAddress {
	var paymentAddress PaymentAddress
	paymentAddress.Address = GenerateAddress(spendingKey)
	paymentAddress.TransmissionKey = GenerateTransmissionKey(GenerateReceivingKey(spendingKey))
	return paymentAddress
}

// FromPointToByteArray converts an elliptic point to byte array
func FromPointToByteArray(p EllipticPoint) []byte {
	var pointByte []byte
	x := p.X.Bytes()
	y := p.Y.Bytes()
	pointByte = append(pointByte, x...)
	pointByte = append(pointByte, y...)
	return pointByte
}

// FromByteArrayToPoint converts a byte array to elliptic point
func FromByteArrayToPoint(pointByte []byte) EllipticPoint {
	point := new(EllipticPoint)
	point.X = new(big.Int).SetBytes(pointByte[0:32])
	point.Y = new(big.Int).SetBytes(pointByte[32:64])
	return *point
}
