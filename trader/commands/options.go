package commands

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/tendermint/basecoin-examples/trader/options"
	bcmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
	"github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
	"github.com/urfave/cli"
)

const OptionName = "options"

var (
	OptionAddrFlag = cli.StringFlag{
		Name:  "option",
		Usage: "The address of this option",
	}
	OptionExpireFlag = cli.Uint64Flag{
		Name:  "expire",
		Value: 0,
		Usage: "The block height when the option expires",
	}
	OptionSellToFlag = cli.StringFlag{
		Name:  "sellto",
		Usage: "Who to sell the options to (optional)",
	}
	OptionTradeAmountFlag = cli.IntFlag{
		Name:  "trade",
		Usage: "Amount of coins to trade",
	}
	OptionTradeCoinFlag = cli.StringFlag{
		Name:  "trade-coin",
		Value: "blank",
		Usage: "Specify a coin denomination to trade",
	}
	OptionPriceAmountFlag = cli.IntFlag{
		Name:  "price",
		Usage: "Amount of coins for price",
	}
	OptionPriceCoinFlag = cli.StringFlag{
		Name:  "price-coin",
		Value: "blank",
		Usage: "Specify a coin denomination for price",
	}
)

var (
	OptionsTxCmd = cli.Command{
		Name:  "options",
		Usage: "Create, trade, and exercise currency options",
		Subcommands: []cli.Command{
			OptionsCreateTxCmd,
			OptionsSellTxCmd,
			OptionsBuyTxCmd,
			OptionsExerciseTxCmd,
			OptionsDisolveTxCmd,
			OptionsQueryCmd,
		},
	}

	OptionsCreateTxCmd = cli.Command{
		Name:  "create",
		Usage: "Create a new option by sending money",
		Flags: []cli.Flag{
			OptionExpireFlag,
			OptionTradeAmountFlag,
			OptionTradeCoinFlag,
		},
		Action: func(c *cli.Context) error {
			return cmdOptionCreateTx(c)
		},
	}

	OptionsSellTxCmd = cli.Command{
		Name:  "sell",
		Usage: "Offer to sell this option",
		Flags: []cli.Flag{
			OptionAddrFlag,
			OptionSellToFlag,
			OptionPriceAmountFlag,
			OptionPriceCoinFlag,
		},
		Action: func(c *cli.Context) error {
			return cmdOptionSellTx(c)
		},
	}

	OptionsBuyTxCmd = cli.Command{
		Name:  "buy",
		Usage: "Attempt to buy this option",
		Flags: []cli.Flag{
			OptionAddrFlag,
		},
		Action: func(c *cli.Context) error {
			return cmdOptionBuyTx(c)
		},
	}

	OptionsExerciseTxCmd = cli.Command{
		Name:  "Exercise",
		Usage: "Exercise this option to trade currency at the given rate",
		Flags: []cli.Flag{
			OptionAddrFlag,
		},
		Action: func(c *cli.Context) error {
			return cmdOptionExerciseTx(c)
		},
	}

	OptionsDisolveTxCmd = cli.Command{
		Name:  "disolve",
		Usage: "Attempt to disolve this option (if never sold, or already expired)",
		Flags: []cli.Flag{
			OptionAddrFlag,
		},
		Action: func(c *cli.Context) error {
			return cmdOptionDisolveTx(c)
		},
	}

	OptionsQueryCmd = cli.Command{
		Name:      "query",
		Usage:     "Return the contents of the given escrow",
		ArgsUsage: "<address>",
		Action: func(c *cli.Context) error {
			return cmdOptionQuery(c)
		},
	}

	OptionsPluginFlag = cli.BoolFlag{
		Name:  "options-plugin",
		Usage: "Enable the options plugin",
	}
)

func init() {
	bcmd.RegisterTxPlugin(OptionsTxCmd)
	bcmd.RegisterStartPlugin(OptionsPluginFlag,
		func() types.Plugin { return options.New(OptionName) })
}

func cmdOptionCreateTx(c *cli.Context) error {
	optionAmount := int64(c.Int(OptionTradeAmountFlag.Name))
	optionCoin := c.String(OptionTradeCoinFlag.Name)
	expire := c.Uint64(EscrowExpireFlag.Name)
	parent := c.Parent().Parent()

	tx := options.CreateOptionTx{
		Expiration: expire,
		Trade: types.Coins{{ // yes {{ an array with one element....
			Denom:  optionCoin,
			Amount: optionAmount,
		}},
	}
	data := options.TxBytes(tx)
	return bcmd.AppTx(parent, OptionName, data)
}

func cmdOptionSellTx(c *cli.Context) error {
	addrHex := c.String(OptionAddrFlag.Name)
	buyerHex := c.String(OptionSellToFlag.Name)
	optionAmount := int64(c.Int(OptionPriceAmountFlag.Name))
	optionCoin := c.String(OptionPriceCoinFlag.Name)
	parent := c.Parent().Parent()

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(addrHex))
	if err != nil {
		return errors.New("Recv address is invalid hex: " + err.Error())
	}

	buyer, err := hex.DecodeString(bcmd.StripHex(buyerHex))
	if err != nil { // this is optional, we can ignore it
		buyer = nil
	}

	tx := options.SellOptionTx{
		Addr:      addr,
		NewHolder: buyer,
		Price: types.Coins{{ // yes {{ an array with one element....
			Denom:  optionCoin,
			Amount: optionAmount,
		}},
	}
	data := options.TxBytes(tx)
	return bcmd.AppTx(parent, OptionName, data)
}

func cmdOptionBuyTx(c *cli.Context) error {
	addrHex := c.String(OptionAddrFlag.Name)
	parent := c.Parent().Parent()

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(addrHex))
	if err != nil {
		return errors.New("Recv address is invalid hex: " + err.Error())
	}

	tx := options.BuyOptionTx{
		Addr: addr,
	}
	data := options.TxBytes(tx)
	return bcmd.AppTx(parent, OptionName, data)
}

func cmdOptionExerciseTx(c *cli.Context) error {
	addrHex := c.String(OptionAddrFlag.Name)
	parent := c.Parent().Parent()

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(addrHex))
	if err != nil {
		return errors.New("Recv address is invalid hex: " + err.Error())
	}

	tx := options.ExerciseOptionTx{
		Addr: addr,
	}
	data := options.TxBytes(tx)
	return bcmd.AppTx(parent, OptionName, data)
}

func cmdOptionDisolveTx(c *cli.Context) error {
	addrHex := c.String(OptionAddrFlag.Name)
	parent := c.Parent().Parent()

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(addrHex))
	if err != nil {
		return errors.New("Recv address is invalid hex: " + err.Error())
	}

	tx := options.DisolveOptionTx{
		Addr: addr,
	}
	data := options.TxBytes(tx)
	return bcmd.AppTx(parent, OptionName, data)
}

func cmdOptionQuery(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("account command requires an argument ([address])")
	}
	addrHex := bcmd.StripHex(c.Args()[0])

	// convert destination address to bytes
	addr, err := hex.DecodeString(addrHex)
	if err != nil {
		return errors.New("Recv address is invalid hex: " + err.Error())
	}

	opt, err := getOption(c.String("node"), addr)
	if err != nil {
		return err
	}
	fmt.Println(string(wire.JSONBytes(opt)))
	return nil
}

func getOption(tmAddr string, address []byte) (*options.OptionData, error) {
	prefix := []byte(fmt.Sprintf("%s/", OptionName))
	key := append(prefix, address...)
	response, err := bcmd.Query(tmAddr, key)
	if err != nil {
		return nil, err
	}

	optionBytes := response.Value

	if len(optionBytes) == 0 {
		return nil, fmt.Errorf("Option bytes are empty for address: %X ", address)
	}
	opt, err := options.ParseData(optionBytes)
	if err != nil {
		return nil, fmt.Errorf("Error reading option %X error: %v",
			optionBytes, err.Error())
	}
	return &opt, nil
}