package privacy

import (
	"math/big"

	"github.com/minio/blake2b-simd"
)

// PedersenCommitment represents a commitment that includes 4 generators
type PedersenCommitment interface {
	// Params returns the parameters for the commitment
	Params() *PCParams
	// InitCommitment initialize the parameters
	InitCommitment() *PCParams
	// CommitAll commits
	Commit([CM_CAPACITY][]byte) []byte
	GetHashOfValues([]byte) []byte
	CommitSpecValue([]byte, []byte, byte) []byte
}

// PCParams represents the parameters for the commitment
type PCParams struct {
	G [CM_CAPACITY]EllipticPoint // generators
	// G[0]: public key
	// G[1]: Value
	// G[2]: SerialNumber
	// G[3]: Random
}

const (
	//CM_CAPACITY ...
	CM_CAPACITY = 4
)

//PCParams ...
var Pcm PCParams

// hashGenerator derives new generator from another generator using hash function
func hashGenerator(g EllipticPoint) EllipticPoint {
	// res.X = hash(g.X), res.Y = sqrt(res.X^3 - 3X + B)
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

//GetHashOfValues get blake2b(G0||G1||G2||G3||<inputs>)
func (com PCParams) GetHashOfValues(values [][]byte) []byte {
	hashFunc := blake2b.New256()
	appendStr := append(CompressKey(Pcm.G[0]), CompressKey(Pcm.G[1])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[2])...)
	appendStr = append(appendStr, CompressKey(Pcm.G[3])...)
	for i := 0; i < len(values); i++ {
		appendStr = append(appendStr, values[i]...)
	}
	hashFunc.Write(appendStr)
	hashValue := hashFunc.Sum(nil)
	return hashValue
}

//ComputeYCoord calculates Y coord from X
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

// Params returns parameters of commitment
func (com PCParams) Params() PCParams {
	return com
}

// InitCommitment initializes parameters of Pedersen commitment
func (com *PCParams) InitCommitment() {

	// G0 is the base point of curve P256
	// G1 is the G0's hash
	// G2 is the G1's hash
	// G3 is the G2's hash
	//com.G[0] = EllipticPoint{new(big.Int).SetBytes(107 23 209 242 225 44 66 71 248 188 230 229 99 164 64 242 119 3 125 129 45 235 51 160 244 161 57 69 216 152 194 150)}
	//com.G[1] = EllipticPoint{big.NewInt(23958808978146169065002618177085267768449753078050893909214942889381335385516),
	//						big.NewInt(7683342728189100387735055495645379680280002577849876913406403613392655773246)}
	//com.G[2] = EllipticPoint{big.NewInt(3815079413196168640531957554999832314172554264133353679019997712191027719881),
	//						big.NewInt(77272455227631396963147866972129356047185376243416151568214292723259718584837)}
	//com.G[3] = EllipticPoint{big.NewInt(96857018922751703925448770729453114406629947956893140191061131561860088449600),
	//						big.NewInt(91021963404625493633337876713358362238547713052118385945840294627046689931308)}

	com.G[0] = EllipticPoint{Curve.Params().Gx, Curve.Params().Gy}
	//fmt.Printf("G0.X: %#v\n", com.G[0].X.Bytes())
	//fmt.Printf("G0.Y: %#v\n", com.G[0].Y.Bytes())
	for i := 1; i < CM_CAPACITY; i++ {
		com.G[i] = hashGenerator(com.G[i-1])
		//fmt.Printf("G%v.X: %#v\n", i, com.G[i].X.Bytes())
		//fmt.Printf("G%v.Y: %#v\n", i, com.G[i].Y.Bytes())
	}

	//TODO: hard code parameters
}

// Commit commits a list of CM_CAPACITY value(s)
func (com PCParams) Commit(values [CM_CAPACITY][]byte) []byte {
	var commitment, temp EllipticPoint
	commitment = EllipticPoint{big.NewInt(0), big.NewInt(0)}
	for i := 0; i < CM_CAPACITY; i++ {
		temp.X, temp.Y = Curve.ScalarMult(com.G[i].X, com.G[i].Y, values[i])
		commitment.X, commitment.Y = Curve.Add(commitment.X, commitment.Y, temp.X, temp.Y)
	}

	// convert result from Elliptic to bytes array
	return (CompressKey(commitment))
}

func (com PCParams) CommitSpecValue(value, sRnd []byte, index byte) []byte {
	var commitment, temp EllipticPoint
	commitment = EllipticPoint{big.NewInt(0), big.NewInt(0)}
	temp.X, temp.Y = Curve.ScalarMult(com.G[0].X, com.G[0].Y, sRnd)
	commitment.X, commitment.Y = Curve.Add(commitment.X, commitment.Y, temp.X, temp.Y)
	temp.X, temp.Y = Curve.ScalarMult(com.G[index].X, com.G[index].Y, sRnd)
	commitment.X, commitment.Y = Curve.Add(commitment.X, commitment.Y, temp.X, temp.Y)
	// convert result from Elliptic to bytes array
	return (CompressKey(commitment))
}
