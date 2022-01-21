package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Vulcanize/ipld-eth-db-validator/pkg/validator"
)

// stateValidatorCmd represents the stateValidator command
var stateValidatorCmd = &cobra.Command{
	Use:   "stateValidator",
	Short: "Validate ethereum state",
	Long:  `Usage ./ipld-eth-db-validator stateValidator --config={path to toml config file}`,

	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		stateValidator()
	},
}

func stateValidator() {
	cfg, err := validator.NewConfig()
	if err != nil {
		logWithCommand.Fatal(err)
	}

	height := viper.GetUint64("validate.block-height")
	if height < 1 {
		logWithCommand.Fatalf("block height cannot be less the 1")
	}

	trail := viper.GetUint64("validate.trail")
	toblock := viper.GetUint64("validate.to-block")
	var chaincfg *params.ChainConfig = nil
	if viper.GetBool("chain.set") {
		chaincfg = setChainConfig()
	}
	srvc := validator.NewService(cfg.DB, height, trail, toblock, chaincfg)

	_, err = srvc.Start(context.Background())
	if err != nil {
		logWithCommand.Fatal(err)
	}

	logWithCommand.Println("state validation complete")
}

func init() {
	rootCmd.AddCommand(stateValidatorCmd)

	stateValidatorCmd.PersistentFlags().String("block-height", "1", "block height to initiate state validation")
	stateValidatorCmd.PersistentFlags().String("trail", "0", "trail of block height to validate")
	stateValidatorCmd.PersistentFlags().String("to-block", "0", "validate till block number")

	_ = viper.BindPFlag("validate.block-height", stateValidatorCmd.PersistentFlags().Lookup("block-height"))
	_ = viper.BindPFlag("validate.trail", stateValidatorCmd.PersistentFlags().Lookup("trail"))
	_ = viper.BindPFlag("validate.to-block", stateValidatorCmd.PersistentFlags().Lookup("to-block"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Printf("Using config file: %s", viper.ConfigFileUsed())
		} else {
			log.Fatal(fmt.Sprintf("Couldn't read config file: %s", err.Error()))
		}
	} else {
		log.Warn("No config file passed with --config flag")
	}
}

func setChainConfig() *params.ChainConfig {
	chaincfg := &params.ChainConfig{
		ChainID:             big.NewInt(viper.GetInt64("chain.chainid")),
		HomesteadBlock:      big.NewInt(viper.GetInt64("chain.homestead-block")),
		DAOForkBlock:        big.NewInt(viper.GetInt64("chain.fork-block")),
		DAOForkSupport:      viper.GetBool("chain.fork-support"),
		EIP150Block:         big.NewInt(viper.GetInt64("chain.eip150-block")),
		EIP150Hash:          common.HexToHash(viper.GetString("chain.eip150-hash")),
		EIP155Block:         big.NewInt(viper.GetInt64("chain.eip155-block")),
		EIP158Block:         big.NewInt(viper.GetInt64("chain.eip158-block")),
		ByzantiumBlock:      big.NewInt(viper.GetInt64("chain.byzantium-block")),
		ConstantinopleBlock: big.NewInt(viper.GetInt64("chain.constantinople-block")),
		PetersburgBlock:     big.NewInt(viper.GetInt64("chain.petersburg-block")),
		IstanbulBlock:       big.NewInt(viper.GetInt64("chain.istanbul-block")),
		MuirGlacierBlock:    big.NewInt(viper.GetInt64("chain.muirGlacier-block")),
		BerlinBlock:         big.NewInt(viper.GetInt64("chain.berlin-block")),
		LondonBlock:         big.NewInt(viper.GetInt64("chain.london-block")),
		ArrowGlacierBlock:   big.NewInt(viper.GetInt64("chain.arrowGlacier-block")),
		MergeForkBlock:      big.NewInt(viper.GetInt64("chain.mergeFork-block")),
	}

	if viper.GetBool("chain.ethash") {
		chaincfg.Ethash = new(params.EthashConfig)
	} else {
		chaincfg.Clique = &params.CliqueConfig{
			Period: viper.GetUint64("chain.clique-period"),
			Epoch:  viper.GetUint64("chain.clique-epoch"),
		}
	}
	return chaincfg
}
