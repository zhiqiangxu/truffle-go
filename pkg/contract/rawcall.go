package contract

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

func RawCall(conf *config.Config, contractAddr, hexdata string) {

	client, err := newClient(conf)
	if err != nil {
		log.Fatalf("newClient failed: %v", err)
	}

	data, err := hex.DecodeString(hexdata)
	if err != nil {
		log.Fatalf("DecodeString hexarg failed: %v", err)
	}

	toAddr := common.HexToAddress(contractAddr)
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{To: &toAddr, Data: data}, nil)
	if err != nil {
		log.Fatalf("CallContract failed: %v", err)
	}

	fmt.Println(hex.EncodeToString(result))
}
