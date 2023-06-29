package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"

	txfuzz "github.com/MariusVanDerWijden/tx-fuzz"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
)

var sendCommand = &cli.Command{
	Name:   "send-blob",
	Usage:  "Sends a single blob-carrying transaction",
	Action: runSendBlob,
	Flags: []cli.Flag{
		skFlag,
		rpcFlag,
	},
}

func main() {
	commands := []*cli.Command{sendCommand}
	app := &cli.App{
		Commands: commands,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	maxDataPerTx = 1 << 17 // 128Kb
)

// PUSH0, DATAHASH, PUSH0, DATAHASH, SSTORE
var TxData4844 = []byte{0x5f, 0x49, 0x5f, 0x49, 0x55}

func runSendBlob(cc *cli.Context) error {
	cl, sk := getRealBackend()
	backend := ethclient.NewClient(cl)
	sender := crypto.PubkeyToAddress(sk.PublicKey)
	nonce, err := backend.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return err
	}
	chainid, err := backend.ChainID(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Nonce: %v\n", nonce)
	gp, _ := backend.SuggestGasPrice(context.Background())
	tip, _ := backend.SuggestGasTipCap(context.Background())
	blob, _ := randomBlobData()
	nonce = nonce - 2
	tx := txfuzz.New4844Tx(nonce, nil, 500000, chainid, tip.Mul(tip, common.Big1), gp.Mul(gp, common.Big1), common.Big1, TxData4844, big.NewInt(1000000), blob, make(types.AccessList, 0))
	signedTx, _ := types.SignTx(&tx.Transaction, types.NewCancunSigner(chainid), sk)
	tx.Transaction = *signedTx
	if err := backend.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}
	rlpData, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	cl.CallContext(context.Background(), nil, "eth_sendRawTransaction", hexutil.Encode(rlpData))
	return nil
}

func getRealBackend() (*rpc.Client, *ecdsa.PrivateKey) {
	sk := crypto.ToECDSAUnsafe(common.FromHex(skFlag.Value))
	cl, err := rpc.Dial(rpcFlag.Value)
	if err != nil {
		panic(err)
	}
	return cl, sk
}

func randomBlobData() ([]byte, error) {
	size := rand.Intn(maxDataPerTx)
	data := make([]byte, size)
	n, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, fmt.Errorf("could not create random blob data with size %d: %v", size, err)
	}
	return data, nil
}
