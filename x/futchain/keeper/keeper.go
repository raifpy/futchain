package keeper

import (
	"fmt"
	"net/http"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/raifpy/futchain/x/futchain/keeper/datasource"
	"github.com/raifpy/futchain/x/futchain/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Datasource *datasource.DatasourceFM
	ABI        abi.ABI // base contract abi
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	c DatasourceConfig,
	abi abi.ABI,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		Datasource: &datasource.DatasourceFM{
			Client:  &http.Client{},
			BaseURL: c.ApiURL,
			Headers: func() http.Header {
				var headers = make(http.Header, len(c.Headers))
				for key, value := range c.Headers {
					headers.Set(key, value)
				}
				return headers
			}(),
		},
		ABI: abi,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
