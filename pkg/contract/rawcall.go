package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

func RawCall(conf *config.Config, param RawParam) {

	client, err := newClient(conf)
	if err != nil {
		log.Fatalf("newClient failed: %v", err)
	}

	data := common.FromHex(param.HexData)

	toAddr := common.HexToAddress(param.ContractAddr)
	value := big.NewInt(param.Value)
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{To: &toAddr, Gas: param.GasLimit, Value: value, Data: data}, nil)
	if err != nil {
		log.Fatalf("CallContract failed: %v", err)
	}

	fmt.Println(hex.EncodeToString(result))
}
