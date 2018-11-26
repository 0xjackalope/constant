package zkp

import (
	"math/big"

	"github.com/ninjadotorg/constant/privacy-protocol"
)

const (
	CMListProve = 256
)

// Prove creates big proof
func Prove(inputCoins []*privacy.InputCoin, outputCoins []*privacy.OutputCoin) {
	// Commit each component of coins being spent
	cmSK := make([]*privacy.EllipticPoint, len(inputCoins))
	cmValue := make([]*privacy.EllipticPoint, len(inputCoins))
	cmSND := make([]*privacy.EllipticPoint, len(inputCoins))

	randSK := make([]*big.Int, len(inputCoins))
	randValue := make([]*big.Int, len(inputCoins))
	randSND := make([]*big.Int, len(inputCoins))
	for i, inputCoin := range inputCoins {
		randSK[i] = privacy.RandInt()
		randValue[i] = privacy.RandInt()
		randSND[i] = privacy.RandInt()

		cmSK[i] = privacy.PedCom.CommitAtIndex(inputCoin.SpendingKey, randSK[i], privacy.SK)
		cmValue[i] = privacy.PedCom.CommitAtIndex(inputCoin.CoinInfo.Value, randValue[i], privacy.VALUE)
		cmSND[i] = privacy.PedCom.CommitAtIndex(inputCoin.CoinInfo.SNDerivator, randSND[i], privacy.SND)
	}

	// Summing all commitments of each input coin into one commitment and proving the knowledge of its openings
	cmSum := make([]*privacy.EllipticPoint, len(inputCoins))
	for i := 0; i < len(inputCoins); i++ {
		cmSum[i] = cmSK[i]
		cmSum[i].X, cmSum[i].Y = privacy.Curve.Add(cmSum[i].X, cmSum[i].Y, cmValue[i].X, cmValue[i].Y)
		cmSum[i].X, cmSum[i].Y = privacy.Curve.Add(cmSum[i].X, cmSum[i].Y, cmSND[i].X, cmSND[i].Y)
	}

	// Call protocol proving knowledge of each sum commitment's openings

	// Proving one-out-of-N commitments is a commitment to the coins being spent
	cmSumInverse := make([]*privacy.EllipticPoint, len(inputCoins))
	cmLists := make([][]*privacy.EllipticPoint, len(inputCoins))
	//witnessOneOutOfN := make([]*PKOne, len(inputCoins))
	for i := 0; i < len(inputCoins); i++ {
		// get sum commitment inverse
		cmSumInverse[i], _ = cmSum[i].Inverse()

		// Prepare list of commitments for each commitmentSum that includes 2^8 commiments
		// Get all commitments in inputCoin[i]'s BlockHeight and other block (if needed)
		cmLists[i] = make([]*privacy.EllipticPoint, CMListProve)
		cmLists[i] = GetCMList(inputCoins[i].CoinInfo.CoinCommitment, inputCoins[i].BlockHeight)
		for j := 0; j < CMListProve; j++ {
			cmLists[i][j].X, cmLists[i][j].Y = privacy.Curve.Add(cmLists[i][j].X, cmLists[i][j].Y, cmSumInverse[i].X, cmSumInverse[i].Y)
		}

		// Prepare witness for protocol one-out-of-N
		//witnessOneOutOfN[i].Set()

	}

	// Proving that serial number is derived from the committed derivator

	// Proving that output values do not exceed v_max

	cmValueIn := new(privacy.EllipticPoint)
	cmValueIn.X, cmValueIn.Y = big.NewInt(0), big.NewInt(0)
	cmValueRndIn := big.NewInt(0)
	cmValueOut := new(privacy.EllipticPoint)
	cmValueOut.X, cmValueOut.Y = big.NewInt(0), big.NewInt(0)
	for i := 0; i < len(inputCoins); i++ {
		cmValueIn.X, cmValueIn.Y = privacy.Curve.Add(cmValueIn.X, cmValueIn.Y, cmValue[i].X, cmValue[i].Y)
		cmValueRndIn = cmValueRndIn.Add(cmValueRndIn, randValue)
		cmValueRndIn = cmValueRndIn.Mod(cmValueRndIn, privacy.Curve.Params().N)
	}

	//cmEqualValue.X, cmEqualValue.Y = big.NewInt(0), big.NewInt(0)
	cmEqualValue, _ := cmValueIn.Inverse()
	cmEqualValue.X, cmEqualValue.Y = privacy.Curve.Add(cmEqualValue.X, cmEqualValue.Y, cmValueOut.X, cmValueOut.Y)

	witnessEqualValue := new(PKComZeroWitness)
	witnessEqualValue.Set(cmEqualValue, &privacy.VALUE, cmValueRndIn)
	proofEqualValue, _ := witnessEqualValue.Prove()
	// Proving that sum of inputs equals sum of outputs
	// @Hy
	//prove ( cmvaluein cmvalueout) (commit + s...)
}

// GetCMList returns list CMListProve (2^8) commitments that includes cm in blockHeight
func GetCMList(cm *privacy.EllipticPoint, blockHeight *big.Int) []*privacy.EllipticPoint {
	return nil
}
