package privacy

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/minio/blake2b-simd"
)

// Commitment represents a commitment that includes 4 generators
type Commitment interface {
	// Params returns the parameters for the commitment
	Params() *CommitmentParams
	// InitCommitment initialize the parameters
	InitCommitment() *CommitmentParams
	// Commit commits
	Commit(prdnum, address, value, serialNumber []byte) []byte
}

// CommitmentParams represents the parameters for the commitment
type CommitmentParams struct {
	G [4]EllipticPoint // generators
	// G[0]: public key
	// G[1]: Value
	// G[2]: SerialNumber
	// G[3]: Random
}



// Cm is a global variable, it is initialized only one time
// var Cm Commitment
// Cm = Cm.InitCommitment()

// hashGenerator create new generator from a generator using hash function
func hashGenerator(g EllipticPoint) EllipticPoint {
	// TODO: res.X = hash(g.X), res.Y = sqrt(res.X^3 - 3X + B)
	// done
	var res = new(EllipticPoint)
	res.X = big.NewInt(0)
	res.Y = big.NewInt(0)
	res.X.SetBytes(g.X.Bytes())
	for {
		hashMachine := blake2b.New256()
		hashMachine.Write(res.X.Bytes())
		res.X.SetBytes(hashMachine.Sum(nil))
		res.Y = ComputeYCoord(res.X)
		if (res.Y != nil) && (Curve.IsOnCurve(res.X, res.Y)) {
			break
		}
	}
	//check Point of degree 2
	pointToChecked := new(EllipticPoint)
	pointToChecked.X, pointToChecked.Y = Curve.Double(res.X, res.Y)

	if pointToChecked.X == nil || pointToChecked.Y == nil {
		//fmt.Errorf("Point at infinity")
		return *new(EllipticPoint)
	}
	return *res
}

//ComputeYCoord calculate Y coord from X
func ComputeYCoord(x *big.Int) *big.Int {
	Q := Curve.Params().P
	temp := new(big.Int)
	xTemp := new(big.Int)

	// Y = +-sqrt(x^3 - 3*x + B)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	x3.Add(x3, Curve.Params().B)
	x3.Sub(x3, xTemp.Mul(x, new(big.Int).SetInt64(3)))
	x3.Mod(x3, Curve.Params().P)

	//check P = 3 mod 4?
	if temp.Mod(Q, new(big.Int).SetInt64(4)).Cmp(new(big.Int).SetInt64(3)) == 0 {
		//		fmt.Println("Ok!!!")
	}

	// Now calculate sqrt mod p of x^3 - 3*x + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, PAdd1Div4(Q), Q)
	// Check that y is a square root of x^3  - 3*x + B.
	y2 := new(big.Int).Mul(y, y)
	y2.Mod(y2, Curve.Params().P)
	//fmt.Printf("y2: %X\n", y2)
	//fmt.Printf("x3: %X\n", x3)
	if y2.Cmp(x3) != 0 {
		return nil
	}
	return y
}

// Params return the parameters of commitment
func (com *CommitmentParams) Params() *CommitmentParams {
	return com
}

// InitCommitment initial
func (com *CommitmentParams) InitCommitment() {
	fmt.Println(Curve.Params().Gx, "___", Curve.Params().Gy)
	com.G[0] = EllipticPoint{Curve.Params().Gx, Curve.Params().Gy}
	for i := 1; i < 4; i++{
		 com.G[i] = hashGenerator(com.G[i-1])
	}
}

// Commit commits a preoud random number and any attributes of coin
func (com *CommitmentParams) Commit(values map[string][]byte) ([]byte, error){

	if len(values) > 4{
		return nil, errors.New("len of values to commit must less than or equal 4")
	}

	var point, commitment EllipticPoint
	//fmt.Printf("commitment.X: %+v\n", commitment.X)
	//fmt.Printf("commitment.Y: %+v\n", commitment.Y)

	i := 0
	for value := range values {
		switch value {
		case "pk":
			point.X, point.Y = Curve.ScalarMult(com.G[0].X, com.G[0].Y, values["pk"])
		case "v":
			point.X, point.Y = Curve.ScalarMult(com.G[1].X, com.G[1].Y, values[""])
		case "sn":
			point.X, point.Y = Curve.ScalarMult(com.G[2].X, com.G[2].Y, values["sn"])
		case "r":
			point.X, point.Y = Curve.ScalarMult(com.G[3].X, com.G[3].Y, values["r"])
		}
		if i == 0 {
			commitment.X = point.X
			commitment.Y = point.Y
		} else{
			commitment.X, commitment.Y = Curve.Add(commitment.X, commitment.Y, point.X, point.Y)
		}
		i++
	}

	// convert result from Elliptic to bytes array
	res := CompressKey(commitment)
	return res, nil
}


