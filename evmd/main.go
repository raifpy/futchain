package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	examplechain "github.com/raifpy/futchainevm"
	"github.com/raifpy/futchainevm/evmd/cmd"
	chainconfig "github.com/raifpy/futchainevm/evmd/config"
)

func main() {
	setupSDKConfig()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "futchaind", examplechain.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func setupSDKConfig() {
	config := sdk.GetConfig()
	chainconfig.SetBech32Prefixes(config)
	config.Seal()
}
