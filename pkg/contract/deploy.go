package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

// Deploy ...
func Deploy(conf *config.Config, solc, solidityFile, targetContract string, args []string) {
	rawContracts, err := compiler.CompileSolidity(solc, solidityFile)
	if err != nil {
		log.Fatalf("CompileSolidity failed: %v", err)
	}

	contracts := handleRawContracts(rawContracts)

	contract := contracts[targetContract]
	if contract == nil {
		log.Fatalf(fmt.Sprintf("contract not found:%s", targetContract))
	}

	abiBytes, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		log.Fatalf("json.Marshal(contract.Info.AbiDefinition) failed: %v", err)
	}

	client, err := newClient(conf)
	if err != nil {
		log.Fatalf("newClient failed: %v", err)
	}
	auth, err := newTransactOpts(client, conf)
	if err != nil {
		log.Fatalf("newTransactOpts failed: %v", err)
	}

	evmABI, err := abi.JSON(strings.NewReader(string(abiBytes)))

	if len(args) != len(evmABI.Constructor.Inputs) {
		log.Fatalf("args mismatch, expect #%d, got #%d", len(evmABI.Constructor.Inputs), len(args))
	}

	var (
		params []interface{}
		param  interface{}
	)
	for i, input := range evmABI.Constructor.Inputs {

		param, err = encode(input.Type, args[i])
		if err != nil {
			log.Fatalf("encode failed for arg:%s, index:%d", args[i], i)
		}
		params = append(params, param)

	}

	address, tx, _, err := bind.DeployContract(auth, evmABI, common.FromHex(contract.Code), client, params...)
	if err != nil {
		log.Fatalf("bind.DeployContract failed: %v", err)
	}
	waitTransactionConfirm(client, tx.Hash())

	log.Infof("contract deployed at %s", address.Hex())
	return
}
