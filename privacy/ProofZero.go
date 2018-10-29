package privacy

import (
	"crypto/rand"
	"math/big"
)

/*Protocol for opening a commitment to 0
Prove:
	commitmentValue is commitment value of Zero, that is statement needed to prove
	commitmentValue is calculated by Comm_ck(Value,PRDNumber)
	commitmentRnd is PRDNumber, which is used to calculate commitmentValue
	s <- Zp; P is Curve's order
	B <- Comm_ck(0,s);  Comm_ck is PedersenCommit function using public params - Curve.Params() (G0,G1...)
						but is just commit special value (in this case, special value is 0),
						which is stick with G[index] (in this case, index is the index stick with commitmentValue)
						B is a.k.a commitmentZeroS
	x <- Hash(G0||G1||G2||G3||commitmentvalue) x is pseudorandom number, which could be computed easily by Verifier
	z <- rx + s; z in Zp, r is commitmentRnd
	return commitmentZeroS, z

Verify:
	commitmentValue is commitment value of Zero, that is statement needed to prove
	commitmentValue is calculated by Comm_ck(Value,PRDNumber), a.k.a A
	commitmentZeroS, z are output of Prove function, commitmentZeroS is a.k.a B
	x <- Hash(G0||G1||G2||G3||commitmentvalue)
	boolValue <- (Comm_ck(0,z) == A.x + B); in this case, A and B needed to convert to EllipticPoint
	return boolValue
)
*/

//ProveIsZero generate a proof prove that the commitment is zero
func ProveIsZero(commitmentValue, commitmentRnd []byte, index byte) ([]byte, *big.Int) {
	//var x big.Int
	//s is a random number in Zp, with p is P, which is order of Curve
	sRnd, err := rand.Int(rand.Reader, Curve.Params().P)
	if err != nil {
		panic(err)
	}
	sRnd.Bytes()
	zeroInt := big.NewInt(0)
	commitmentZeroS := Pcm.CommitSpecValue(zeroInt.Bytes(), sRnd.Bytes(), index)
	xRnd := big.NewInt(0)
	xRnd.SetBytes(Pcm.GetHashOfValues([][]byte{commitmentValue}))
	xRnd.Mod(xRnd, Curve.Params().P)
	z := big.NewInt(0)
	z.SetBytes(commitmentRnd)
	z.Mul(z, xRnd)
	z.Mod(z, Curve.Params().P)
	z.Add(z, sRnd)
	z.Mod(z, Curve.Params().P)
	return commitmentZeroS, z
}

//VerifyIsZero verify that under commitment is zero
func VerifyIsZero(commitmentValue, commitmentZeroS []byte, index byte, z *big.Int) bool {
	xRnd := big.NewInt(0)
	xRnd.SetBytes(Pcm.GetHashOfValues([][]byte{commitmentValue}))
	xRnd.Mod(xRnd, Curve.Params().P)
	commitmentValuePoint, err := DecompressCommitment(commitmentValue)
	if err != nil {
		return false
	}
	if (!Curve.IsOnCurve(commitmentValuePoint.X, commitmentValuePoint.Y)) || (z.Cmp(Curve.Params().P) > -1) {
		return false
	}

	commitmentZeroSPoint, err := DecompressCommitment(commitmentZeroS)
	if err != nil {
		return false
	}
	if (!Curve.IsOnCurve(commitmentZeroSPoint.X, commitmentZeroSPoint.Y)) || (z.Cmp(Curve.Params().P) > -1) {
		return false
	}

	zeroInt := big.NewInt(0)
	commitmentZeroZ := Pcm.CommitSpecValue(zeroInt.Bytes(), z.Bytes(), index)
	verifyPoint := new(EllipticPoint)
	verifyPoint.X.SetBytes(commitmentValuePoint.X.Bytes())
	verifyPoint.Y.SetBytes(commitmentValuePoint.Y.Bytes())
	verifyPoint.X, verifyPoint.Y = Curve.ScalarMult(verifyPoint.X, verifyPoint.Y, xRnd.Bytes())

	verifyPoint.X, verifyPoint.Y = Curve.Add(verifyPoint.X, verifyPoint.Y, commitmentZeroSPoint.X, commitmentZeroSPoint.Y)

	commitmentZeroZPoint, err := DecompressCommitment(commitmentZeroZ)
	if err != nil {
		return false
	}
	if (!Curve.IsOnCurve(commitmentZeroSPoint.X, commitmentZeroSPoint.Y)) || (z.Cmp(Curve.Params().P) > -1) {
		return false
	}

	if commitmentZeroZPoint.X.CmpAbs(verifyPoint.X) != 0 {
		return false
	}
	if commitmentZeroZPoint.Y.CmpAbs(verifyPoint.Y) != 0 {
		return false
	}

	return true
}
