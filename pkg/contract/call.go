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

// Call ...
func Call(conf *config.Config, solc, solidityFile, contractAddr, targetContract, methodName string, args []string) {
	var (
		rawContracts map[string]*compiler.Contract
		err          error
	)
	if strings.HasSuffix(solidityFile, ".sol") {
		rawContracts, err = compiler.CompileSolidity(solc, solidityFile)
	} else {
		rawContracts, err = compiler.CompileVyper(solc, solidityFile)
	}

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

	evmABI, err := abi.JSON(strings.NewReader(string(abiBytes)))

	method, ok := evmABI.Methods[methodName]
	if !ok {
		log.Fatalf("method not exists: %s", methodName)
	}

	if len(args) != len(method.Inputs) {
		log.Fatalf("args mismatch, expect #%d, got #%d", len(method.Inputs), len(args))
	}

	var (
		params []interface{}
		param  interface{}
	)
	for i, input := range method.Inputs {

		param, err = encode(input.Type, args[i])
		if err != nil {
			log.Fatalf("encode failed for arg:%s, index:%d", args[i], i)
		}
		params = append(params, param)

	}

	bc := bind.NewBoundContract(common.HexToAddress(contractAddr), evmABI, client, nil, nil)
	var results []interface{}
	err = bc.Call(nil, &results, methodName, params...)
	if err != nil {
		log.Fatalf("bc.Call failed: %v", err)
	}

	for i, result := range results {
		ty := method.Outputs[i].Type
		log.Infof("output %d(type %s) is: %s", i+1, ty.String(), format(ty, result))
	}

}
