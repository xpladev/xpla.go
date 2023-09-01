package util

import (
	"math/big"
	"strconv"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xpladev/xpla.go/types"
)

func GetAddrByPrivKey(priv cryptotypes.PrivKey) (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromHex(priv.PubKey().Address().String())
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GasLimitAdjustment(gasUsed uint64, gasAdjustment string) (string, error) {
	gasAdj, err := strconv.ParseFloat(gasAdjustment, 64)
	if err != nil {
		return "", err
	}
	return FromIntToString(int(gasAdj * float64(gasUsed))), nil
}

func GrpcUrlParsing(normalUrl string) string {
	if strings.Contains(normalUrl, "http://") || strings.Contains(normalUrl, "https://") {
		parsedUrl := strings.Split(normalUrl, "://")
		return parsedUrl[1]
	} else {
		return normalUrl
	}
}

func DenomAdd(amount string) string {
	if strings.Contains(amount, types.XplaDenom) {
		return amount
	} else {
		return amount + types.XplaDenom
	}
}

func DenomRemove(amount string) string {
	if strings.Contains(amount, types.XplaDenom) {
		returnAmount := strings.Split(amount, types.XplaDenom)
		return returnAmount[0]
	} else {
		return amount
	}
}

func ConvertEvmChainId(chainId string) (*big.Int, error) {
	conv1 := strings.Split(chainId, "_")
	conv2 := strings.Split(conv1[1], "-")
	returnChainId, err := FromStringToBigInt(conv2[0])
	if err != nil {
		return nil, err
	}
	return returnChainId, nil
}

func Bech32toValidatorAddress(validators []string) ([]sdk.ValAddress, error) {
	vals := make([]sdk.ValAddress, len(validators))
	for i, validator := range validators {
		addr, err := sdk.ValAddressFromBech32(validator)
		if err != nil {
			return nil, err
		}
		vals[i] = addr
	}
	return vals, nil
}

func MakeQueryLcdUrl(metadata string) string {
	return "/" + strings.Replace(metadata, "query.proto", "", -1)
}

func MakeQueryLabels(labels ...string) string {
	return strings.Join(labels, "/")
}
