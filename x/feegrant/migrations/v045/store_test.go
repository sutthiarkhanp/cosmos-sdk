package v045_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	v045 "github.com/cosmos/cosmos-sdk/x/feegrant/migrations/v045"
	"github.com/stretchr/testify/require"
)

func TestMigration(t *testing.T) {
	encCfg := simapp.MakeTestEncodingConfig()
	cdc := encCfg.Codec
	feegrantKey := sdk.NewKVStoreKey(v045.ModuleName)
	ctx := testutil.DefaultContext(feegrantKey, sdk.NewTransientStoreKey("transient_test"))
	granter1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	grantee1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	granter2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	grantee2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	spendLimit := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000)))
	now := ctx.BlockTime()
	oneDay := ctx.BlockTime().AddDate(0, 0, 1)

	grants := []struct {
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		spendLimit sdk.Coins
		expiration *time.Time
	}{
		{
			granter:    granter1,
			grantee:    grantee1,
			spendLimit: spendLimit,
			expiration: &oneDay,
		},
		{
			granter:    granter2,
			grantee:    grantee2,
			spendLimit: spendLimit,
			expiration: &now,
		},
		{
			granter:    granter1,
			grantee:    grantee2,
			spendLimit: spendLimit,
		},
		{
			granter:    granter2,
			grantee:    grantee1,
			expiration: &now,
		},
	}

	store := ctx.KVStore(feegrantKey)
	for _, grant := range grants {
		newGrant, err := feegrant.NewGrant(grant.granter, grant.grantee, &feegrant.BasicAllowance{
			SpendLimit: grant.spendLimit,
			Expiration: grant.expiration,
		})
		require.NoError(t, err)

		bz, err := cdc.Marshal(&newGrant)
		require.NoError(t, err)

		store.Set(feegrant.FeeAllowanceKey(grant.granter, grant.grantee), bz)
	}

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(1 * time.Hour))
	require.NoError(t, v045.MigrateStore(ctx, feegrantKey, cdc))
	store = ctx.KVStore(feegrantKey)

	require.NotNil(t, store.Get(feegrant.FeeAllowanceKey(granter1, grantee1)))
	require.Nil(t, store.Get(feegrant.FeeAllowanceKey(granter2, grantee2)))
	require.NotNil(t, store.Get(feegrant.FeeAllowanceKey(granter1, grantee2)))
	require.Nil(t, store.Get(feegrant.FeeAllowanceKey(granter2, grantee1)))
}
