package util

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func FromBigIntToString(v *big.Int) string {
	return v.String()
}

func FromStringToBigInt(v string) (*big.Int, error) {
	n := big.NewInt(0)
	n, ok := n.SetString(v, 10)
	if !ok {
		return nil, fmt.Errorf("convert string to big int err")
	}
	return n, nil
}

func FromStringToUint64(value string) (uint64, error) {
	if value == "" {
		return 0, nil
	}
	number, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func FromUint64ToString(value uint64) string {
	return strconv.Itoa(int(value))
}

func FromStringToInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func FromInt64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FromStringToInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func FromIntToString(value int) string {
	return strconv.Itoa(value)
}

func FromStringToByte20Address(address string) common.Address {
	var byte20Address common.Address
	if address[:2] == "0x" {
		address = address[2:]
	}
	byte20Address = common.HexToAddress(address)

	return byte20Address
}

func FromByte20AddressToCosmosAddr(address common.Address) (sdk.AccAddress, error) {
	var addrStr string
	if address.Hex()[:2] == "0x" {
		addrStr = address.Hex()[2:]
	} else {
		addrStr = address.Hex()
	}

	addr, err := sdk.AccAddressFromHex(addrStr)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func FromStringHexToHash(hashString string) common.Hash {
	return common.HexToHash(hashString)
}

func From0xHexStringToIBignt(hexString string) *big.Int {
	n := new(big.Int)
	n.SetString(hexString[2:], 16)

	return n
}

func FromStringToTypeHexString(value string) string {
	if !strings.Contains(value, "0x") {
		return "0x" + value
	} else {
		return value
	}
}

func ToString(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	if s == "" {
		return defaultValue
	} else {
		return s
	}
}
