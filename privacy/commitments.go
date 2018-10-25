package privacy

import (
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
	G0, G1, G2, H EllipticPoint // generator
}

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
	return *res
}

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
	// TODO: how to generate generators independently
	fmt.Println(Curve.Params().Gx, "___", Curve.Params().Gy)
	com.G0 = EllipticPoint{Curve.Params().Gx, Curve.Params().Gy}
	com.G1 = hashGenerator(com.G0)
	com.G2 = hashGenerator(com.G1)
	com.H = hashGenerator(com.G2)
}

// Commit commits a preoud random number and 3 attributes of coin
func (com *CommitmentParams) Commit(prdnum, address, value, serialNumber []byte) []byte {
	var res []byte
	// TODO: using Pedersen commitment
	//var commitRPoint EllipticPoint
	/*fmt.Print(Curve.IsOnCurve(com.G0.X, com.G0.Y))
	fmt.Print("__")
	fmt.Print(com.G0.X)
	fmt.Print("__")
	fmt.Println(com.G0.Y)
	fmt.Print(Curve.IsOnCurve(com.G1.X, com.G1.Y))
	fmt.Print("__")
	fmt.Print(com.G1.X)
	fmt.Print("__")
	fmt.Println(com.G1.Y)
	fmt.Print(Curve.IsOnCurve(com.G2.X, com.G2.Y))
	fmt.Print("__")
	fmt.Print(com.G2.X)
	fmt.Print("__")
	fmt.Println(com.G2.Y)
	fmt.Print(Curve.IsOnCurve(com.H.X, com.H.Y))
	fmt.Print("__")
	fmt.Print(com.H.X)
	fmt.Print("__")
	fmt.Println(com.H.Y) //*/

	commx, commy := Curve.ScalarMult(com.H.X, com.H.Y, prdnum)
	commxtemp, commytemp := Curve.ScalarMult(com.G0.X, com.G0.Y, address)
	commx, commy = Curve.Add(commx, commy, commxtemp, commytemp)
	commxtemp, commytemp = Curve.ScalarMult(com.G1.X, com.G1.Y, value)
	commx, commy = Curve.Add(commx, commy, commxtemp, commytemp)
	commxtemp, commytemp = Curve.ScalarMult(com.G2.X, com.G2.Y, serialNumber)
	commx, commy = Curve.Add(commx, commy, commxtemp, commytemp)

	// TODO: convert result from Elliptic to bytes array
	var resPoint EllipticPoint
	resPoint.X = commx
	resPoint.Y = commy
	res = CompressKey(resPoint)
	return res
}

/*

func HashGenerator(g EllipticPoint) EllipticPoint {
	// TODO: res.X = hash(g.X), res.Y = sqrt(res.X^3 - 3X + B)
	res := new (EllipticPoint)
	res.X = big.NewInt(rand.Int63());
	res.Y = big.NewInt(rand.Int63());
	i := 0;
	for !Curve.IsOnCurve(res.X, res.Y) {
		if i==0 {
			*res.X = *g.X;
			*res.Y = *g.Y;
		}
		var h = sha256.New()
		h.Write(res.X.Bytes())
		res.X = new(big.Int).SetBytes(h.Sum(nil));
		res.X.Mod(res.X, Curve.Params().P);
		temp := ComputeYCoord(res.X);
		if temp != nil {
			res.Y = temp;
		}
		//fmt.Println(res.X);
		//fmt.Println("Loop:",i);
		//fmt.Println("X = ",res.X);
		//fmt.Println("Y = ",res.Y);
		i++;
	}

	return *res
}

*/
