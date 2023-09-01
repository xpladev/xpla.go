package util

import "math/big"

func MulUint64(val1 uint64, val2 uint64) uint64 {
	return val1 * val2
}

func MulBigInt(val1 *big.Int, val2 *big.Int) *big.Int {
	result := big.NewInt(0)
	return result.Mul(val1, val2)
}
