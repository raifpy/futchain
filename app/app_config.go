package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	"cosmossdk.io/depinject/appconfig"
	_ "github.com/cosmos/cosmos-sdk/x/bank" // import for side-effects
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	_ "github.com/cosmos/cosmos-sdk/x/distribution" // import for side-effects
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	_ "github.com/cosmos/cosmos-sdk/x/staking" // import for side-effects
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/raifpy/futchain/x/futchain/module"
	futchainmoduletypes "github.com/raifpy/futchain/x/futchain/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	_ "github.com/cosmos/cosmos-sdk/x/auth"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

var (
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		// this line is used by starport scaffolding # stargate/app/maccPerms
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		// We allow the following module accounts to receive funds:
		// govtypes.ModuleName
	}

	// application configuration (used by depinject)
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: Name,
					// NOTE: upgrade module is required to be prioritized
					PreBlockers: []string{
						authtypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/preBlockers
					},
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					BeginBlockers: []string{
						distrtypes.ModuleName,
						stakingtypes.ModuleName,
						// chain modules
						futchainmoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/beginBlockers
					},
					EndBlockers: []string{
						stakingtypes.ModuleName,
						// chain modules
						futchainmoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/endBlockers
					},
					// The following is mostly only needed when ModuleName != StoreKey name.
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					// NOTE: The genutils module must occur after staking so that pools are
					// properly initialized with tokens from genesis accounts.
					// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
					InitGenesis: []string{
						consensustypes.ModuleName,
						authtypes.ModuleName,
						banktypes.ModuleName,
						distrtypes.ModuleName,
						stakingtypes.ModuleName,
						genutiltypes.ModuleName,
						// chain modules
						futchainmoduletypes.ModuleName,
						// this line is used by starport scaffolding # stargate/app/initGenesis
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:                AccountAddressPrefix,
					ModuleAccountPermissions:    moduleAccPerms,
					EnableUnorderedTransactions: false,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name: stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{
					// NOTE: specifying a prefix is only necessary when using bech32 addresses
					// If not specfied, the auth Bech32Prefix appended with "valoper" and "valcons" is used by default
					Bech32PrefixValidator: AccountAddressPrefix + "valoper",
					Bech32PrefixConsensus: AccountAddressPrefix + "valcons",
				}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name: futchainmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&futchainmoduletypes.Module{
					ApiUrl: "https://www.fotmob.com", //TODO: implement default values from config
					Headers: map[string]string{
						"accept":             "*/*",
						"accept-language":    "en-US,en;q=0.7",
						"cache-control":      "no-cache",
						"dnt":                "1",
						"pragma":             "no-cache",
						"priority":           "u=1, i",
						"referer":            "https://www.fotmob.com/",
						"sec-ch-ua":          `"Chromium";v="140", "Not=A?Brand";v="24", "Brave";v="140"`,
						"sec-ch-ua-mobile":   "?0",
						"sec-ch-ua-platform": `"macOS"`,
						"sec-fetch-dest":     "empty",
						"sec-fetch-mode":     "cors",
						"sec-fetch-site":     "same-origin",
						"sec-gpc":            "1",
						"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36",
						"x-mas":              "eyJib2R5Ijp7InVybCI6Ii9hcGkvZGF0YS9tYXRjaGVzP2RhdGU9MjAyNTA5MDYmdGltZXpvbmU9RXVyb3BlJTJGSXN0YW5idWwmY2NvZGUzPUdCUiIsImNvZGUiOjE3NTcxNjIzMjEzNTEsImZvbyI6InByb2R1Y3Rpb246ZTUyYTNmYzE5Y2Y0YmY0NTY3ZTBlMzA3N2Q1OWQzNjVhNGEyYjNkNiJ9LCJzaWduYXR1cmUiOiIwOUVDNUVBNTNDQkNDNEFCODU3M0ZFMDQxN0REMTUwRCJ9",
					},
				}),
			},
			// this line is used by starport scaffolding # stargate/app/moduleConfig
		},
	})
)
