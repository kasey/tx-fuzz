package main

import "github.com/urfave/cli/v2"

var (
	skFlag = &cli.StringFlag{
		Name:    "secret-key",
		Aliases: []string{"private-key", "sk"},
		Usage:   "Secret/Private key",
		Value:   "0x2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622",
	}
	rpcFlag = &cli.StringFlag{
		Name:    "rpc-url",
		Aliases: []string{"rpc"},
		Usage:   "RPC provider",
		Value:   "http://127.0.0.1:8545",
	}
	gasLimitFlag = &cli.StringFlag{
		Name:  "gas-limit",
		Usage: "Gas Limit tx field",
		Value: "1000000",
	}
	gasPriceFlag = &cli.StringFlag{
		Name:  "gas-price",
		Usage: "Gas Price tx field",
		Value: "100000000000",
	}
	blobFileFlag = &cli.StringFlag{
		Name:  "blob-file",
		Usage: "Path to file to use as raw blob bytes.",
	}
)
