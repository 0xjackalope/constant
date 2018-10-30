package privacy

import (
	"fmt"
	"math/big"
)

//inversePoint return inverse point of ECC Point input
func inversePoint(eccPoint EllipticPoint) (*EllipticPoint, error) {
	//Check that input is ECC point
	if !Curve.IsOnCurve(eccPoint.X, eccPoint.Y) {
		return nil, fmt.Errorf("Input is not ECC Point")
	}
	//Create result point
	resPoint := new(EllipticPoint)
	resPoint.X = big.NewInt(0)
	resPoint.Y = big.NewInt(0)

	//inverse point of A(x,y) in ECC is A'(x, P - y) with P is order of Curve
	resPoint.X.SetBytes(eccPoint.X.Bytes())
	resPoint.Y.SetBytes(eccPoint.Y.Bytes())
	resPoint.Y.Sub(Curve.Params().P, resPoint.Y)

	return resPoint, nil
}
