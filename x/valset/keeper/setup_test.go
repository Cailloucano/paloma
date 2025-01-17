package keeper

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/paloma/x/valset/types/mocks"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

type mockedServices struct {
	StakingKeeper *mocks.StakingKeeper
}

func newValsetKeeper(t testing.TB) (*Keeper, mockedServices, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)

	types.RegisterInterfaces(registry)

	paramsSubspace := typesparams.NewSubspace(appCodec,
		types.Amino,
		storeKey,
		memStoreKey,
		"ValsetParams",
	)

	ms := mockedServices{
		StakingKeeper: mocks.NewStakingKeeper(t),
	}
	k := NewKeeper(
		appCodec,
		storeKey,
		memStoreKey,
		paramsSubspace,
		ms.StakingKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, nil)
	ctx = ctx.WithMultiStore(stateStore).WithGasMeter(sdk.NewInfiniteGasMeter())

	ctx = ctx.WithLogger(log.NewTMJSONLogger(os.Stdout))

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ms, ctx
}
