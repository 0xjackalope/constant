package privacy

import (
	"github.com/minio/blake2b-simd"
	"math/big"
	"math/rand"
	"time"
)
type Helper interface {
	InitBasePoint() *BasePoint
}
type proofFactor EllipticPoint
type BasePoint struct {

	G EllipticPoint
	H EllipticPoint
}
type ProofOfProductCommitment struct {
	basePoint BasePoint
	D proofFactor
	D1 proofFactor
	E proofFactor
	f1 big.Int
	z1 big.Int
	f2 big.Int
	z2 big.Int
	z3 big.Int
	cmA 		 	[]byte
	cmB 		 	[]byte
	cmC 		 	[]byte
}
type inputCommitments struct {
	witnessA 	[]byte
	cmA 		 	[]byte
	randA		 	[]byte
	witnessB 	[]byte
	cmB 		 	[]byte
	randB		 	[]byte
	witnessAB []byte
	cmAB 		 	[]byte
	randC		 	[]byte
}
func (basePoint *BasePoint) InitBasePoint() {

	P:= new(EllipticPoint)
	P.X = Curve.Params().Gx
	P.Y = Curve.Params().Gy
	basePoint.G = hashGenerator(*P)
	basePoint.H = hashGenerator(basePoint.G);
}
// Random number modular P
func randIntModP() *big.Int {
	rand.Seed(time.Now().UTC().Unix())
	x:=big.NewInt(rand.Int63())
	x.Mod(x,Curve.Params().P);
	return x
}
func computeCommitment(pointG EllipticPoint, pointH EllipticPoint, val1 *big.Int, val2 *big.Int) proofFactor{
	factor:= new(proofFactor)
	factor.X, factor.Y= Curve.ScalarMult(pointG.X, pointG.Y, val1.Bytes())
	tmp:= new(proofFactor)
	tmp.X, tmp.Y = Curve.ScalarMult(pointH.X, pointH.Y, val2.Bytes())
	factor.X,factor.Y = Curve.Add(factor.X, factor.Y, tmp.X, tmp.Y);
	return *factor;
}
func computeHashString(data [][]byte) []byte{
	str:=make([]byte, 0)
	for i:=0;i<len(data);i++{
		str = append(str,data[i]...)
	}
	hashFunc := blake2b.New256()
	hashFunc.Write(str)
	hashValue := hashFunc.Sum(nil)
	return hashValue
}
func ProveProductCommitment(productCommitment inputCommitments)  ProofOfProductCommitment{
	proof :=  new(ProofOfProductCommitment)
	proof.basePoint.InitBasePoint();
	d := randIntModP();
	e := randIntModP();
	s := randIntModP();
	s1 := randIntModP();
	t := randIntModP();
	//Compute D factor of proof
	D:= computeCommitment(proof.basePoint.G, proof.basePoint.H, d,s);

	//Compute D' factor of proof
	G1 := new(EllipticPoint)
	G1.X, G1.Y = Curve.ScalarMult(proof.basePoint.G.X, proof.basePoint.G.Y, productCommitment.witnessB)
	D1:= computeCommitment(*G1,proof.basePoint.H, d,s1);

	//Compute E factor of proof
	E:=computeCommitment(proof.basePoint.G,proof.basePoint.H, e,t)
	proof.D = D;
	proof.E = E;
 	proof.D1 = D1;
	// x = hash(G||H||D||D1||E)
	data:=[][]byte{
		proof.basePoint.G.X.Bytes(),
		proof.basePoint.G.X.Bytes(),
		proof.basePoint.G.Y.Bytes(),
		proof.basePoint.G.Y.Bytes(),
		D.X.Bytes(),
		D.X.Bytes(),
		D1.X.Bytes(),
		D1.Y.Bytes(),
		E.X.Bytes(),
		E.Y.Bytes(),
	}
	x:=new(big.Int)
	x.SetBytes(computeHashString(data))
	//compute f1
	a:= new(big.Int)
	a.SetBytes(productCommitment.witnessA)
	f1:= a.Mul(a,x)
	f1 = f1.Add(f1,d)
	f1 = f1.Mod(f1,Curve.Params().P);
	proof.f1 = *f1;
	//compute z1
	ra:= new(big.Int)
	ra.SetBytes(productCommitment.randA)
	z1:= ra.Mul(ra,x)
	z1 = z1.Add(z1,s)
	z1 = z1.Mod(z1,Curve.Params().P)
	proof.z1 = *z1;
	//compute f2
	b:= new(big.Int)
	b.SetBytes(productCommitment.witnessB)
	f2:= b.Mul(b,x)
	f2 = f2.Add(f2,e)
	f2 = f2.Mod(f2,Curve.Params().P)
	proof.f2 = *f2;
	//compute z2 = rb*x+t mod p
	rb:= new(big.Int)
	rb.SetBytes(productCommitment.randB)
	z2:= rb.Mul(rb,x)
	z2 = z2.Add(z2,t)
	z2 = z2.Mod(z2,Curve.Params().P)
	proof.z2 = *z2;
	//compute z2 = rb*x+t mod p
	rc:= new(big.Int)
	rc.SetBytes(productCommitment.randC)
	rc = rc.Sub(rc,a.Mul(a,rb))

	z3:= rc.Mul(rc,x)
	z3 = z3.Add(z3,s1)
	z3 = z3.Mod(z3,Curve.Params().P)
	proof.z3 = *z3;
	return *proof;
}

func VerifyProductCommitment (proof ProofOfProductCommitment) bool {
	pts1 := new(EllipticPoint)
	data:=[][]byte{
		proof.basePoint.G.X.Bytes(),
		proof.basePoint.G.X.Bytes(),
		proof.basePoint.G.Y.Bytes(),
		proof.basePoint.G.Y.Bytes(),
		proof.D.X.Bytes(),
		proof.D.X.Bytes(),
		proof.D1.X.Bytes(),
		proof.D1.Y.Bytes(),
		proof.E.X.Bytes(),
		proof.E.Y.Bytes(),
	}
	x:= computeHashString(data)
	A:= new(EllipticPoint)
	A,_ = DecompressKey(proof.cmA);
	pts1.X, pts1.Y =Curve.ScalarMult(A.X, A.Y, x)
	pts1.X, pts1.Y = Curve.Add(pts1.X, pts1.Y, proof.D.X,proof.D.Y);
	com1 := computeCommitment(proof.basePoint.G,proof.basePoint.H, &proof.f1, &proof.z1)

	if (com1.X.Cmp(pts1.X)==1 && com1.Y.Cmp(pts1.Y)==1){
		return true;
		}
}
