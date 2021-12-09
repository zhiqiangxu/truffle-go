package contract

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

type RawParam struct {
	ContractAddr string
	HexData      string
	Value        int64
	GasLimit     uint64
}

func RawTransact(conf *config.Config, param RawParam) {
	client, err := newClient(conf)
	if err != nil {
		log.Fatalf("newClient failed: %v", err)
	}

	data := common.FromHex(param.HexData)

	toAddr := common.HexToAddress(param.ContractAddr)

	privateKey, err := crypto.HexToECDSA(conf.PrivateKey)
	if err != nil {
		log.Fatalf("HexToECDSA failed: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("publicKey cast failed")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("NonceAt failed: %v", err)
	}

	gasPrice := big.NewInt(0)
	if conf.GasPrice > 0 {
		gasPrice.SetInt64(conf.GasPrice)
	} else {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("client.SuggestGasPrice failed:%v", err)
			return
		}
	}

	value := big.NewInt(param.Value)
	gasLimit := param.GasLimit
	if gasLimit == 0 {
		callMsg := ethereum.CallMsg{
			From: fromAddress, To: &toAddr, Gas: 0, GasPrice: gasPrice,
			Value: value, Data: data,
		}
		gasLimit, err = client.EstimateGas(context.Background(), callMsg)
		if err != nil {
			log.Fatalf("EstimateGas failed: %v", err)
		}
	}

	ntx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, data)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("ChainID failed: %v", err)
	}
	signedTx, err := types.SignTx(ntx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("SignTx failed: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("SendTransaction failed: %v", err)
	}

	waitTransactionConfirm(client, signedTx.Hash())

	receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		log.Fatalf("TransactionReceipt failed: %v", err)
	}
	receiptBytes, err := json.MarshalIndent(receipt, "", "    ")
	if err != nil {
		log.Fatalf("MarshalIndent failed: %v", err)
	}
	fmt.Println("receipt", string(receiptBytes))
}
