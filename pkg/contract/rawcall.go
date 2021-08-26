package contract

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

func RawCall(conf *config.Config, contractAddr, methodSig string, hexarg string) {

	client, err := newClient(conf)
	if err != nil {
		log.Fatalf("newClient failed: %v", err)
	}

	arg, err := hex.DecodeString(hexarg)
	if err != nil {
		log.Fatalf("DecodeString hexarg failed: %v", err)
	}
	mid := crypto.Keccak256([]byte(methodSig))[0:4]

	toAddr := common.HexToAddress(contractAddr)
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{To: &toAddr, Data: append(mid, arg...)}, nil)
	if err != nil {
		log.Fatalf("CallContract failed: %v", err)
	}

	fmt.Println(hex.EncodeToString(result))
}
