package feegrant

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "feegrant"

	// StoreKey is the store key string for supply
	StoreKey = ModuleName

	// RouterKey is the message route for supply
	RouterKey = ModuleName

	// QuerierRoute is the querier route for supply
	QuerierRoute = ModuleName
)

var (
	// FeeAllowanceKeyPrefix is the set of the kvstore for fee allowance data
	FeeAllowanceKeyPrefix = []byte{0x00}

	// FeeAllowanceQueueKeyPrefix is the set of the kvstore for fee allowance keys data
	FeeAllowanceQueueKeyPrefix = []byte{0x01}
)

// FeeAllowanceKey is the canonical key to store a grant from granter to grantee
// We store by grantee first to allow searching by everyone who granted to you
func FeeAllowanceKey(granter sdk.AccAddress, grantee sdk.AccAddress) []byte {
	return append(FeeAllowancePrefixByGrantee(grantee), address.MustLengthPrefix(granter.Bytes())...)
}

// FeeAllowancePrefixByGrantee returns a prefix to scan for all grants to this given address.
func FeeAllowancePrefixByGrantee(grantee sdk.AccAddress) []byte {
	return append(FeeAllowanceKeyPrefix, address.MustLengthPrefix(grantee.Bytes())...)
}

// FeeAllowancePrefixQueue is the canonical key to store grant key.
func FeeAllowancePrefixQueue(exp *time.Time, allowanceKey []byte) []byte {
	allowanceByExpTimeKey := AllowanceByExpTimeKey(exp)
	return append(allowanceByExpTimeKey, allowanceKey[1:]...)
}

func AllowanceByExpTimeKey(exp *time.Time) []byte {
	return append(FeeAllowanceQueueKeyPrefix, sdk.FormatTimeBytes(*exp)...)
}
