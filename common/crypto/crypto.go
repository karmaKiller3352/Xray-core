// Package crypto provides common crypto libraries for Xray.
package crypto // import "github.com/karmaKiller3352/Xray-core/common/crypto"

import (
	"crypto/rand"
	"math/big"
)

func RandBetween(from int64, to int64) int64 {
	if from == to {
		return from
	}
	bigInt, _ := rand.Int(rand.Reader, big.NewInt(to-from))
	return from + bigInt.Int64()
}
