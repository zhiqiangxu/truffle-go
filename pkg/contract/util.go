package contract

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

func handleRawContracts(rawContracts map[string]*compiler.Contract) map[string]*compiler.Contract {
	contracts := make(map[string]*compiler.Contract)
	for name, contract := range rawContracts {
		nameParts := strings.Split(name, ":")
		contractName := nameParts[len(nameParts)-1]
		if contracts[contractName] != nil {
			log.Fatalf("duplicate contract found: %s", contractName)
		}
		contracts[contractName] = contract
	}
	return contracts
}

func newClient(conf *config.Config) (client *ethclient.Client, err error) {
	client, err = ethclient.Dial(conf.Node)
	if err != nil {
		err = fmt.Errorf("ethclient.Dial failed:%v", err)
		return
	}
	return
}

var (
	defaultGasLimit = 5000000
)

func newTransactOpts(client *ethclient.Client, conf *config.Config) (auth *bind.TransactOpts, err error) {
	privateKey, err := crypto.HexToECDSA(conf.PrivateKey)
	if err != nil {
		err = fmt.Errorf("HexToECDSA failed:%v", err)
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("publicKey cast failed")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		err = fmt.Errorf("client.PendingNonceAt failed:%v", err)
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		err = fmt.Errorf("client.SuggestGasPrice failed:%v", err)
		return
	}

	if conf.ChainID > 0 {
		auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(conf.ChainID)))
		if err != nil {
			err = fmt.Errorf("bind.NewKeyedTransactorWithChainID failed:%v", err)
			return
		}
	} else {
		auth = bind.NewKeyedTransactor(privateKey)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(int64(0))       // in wei
	auth.GasLimit = uint64(defaultGasLimit) // in units
	auth.GasPrice = gasPrice

	return
}

func waitTransactionConfirm(client *ethclient.Client, hash common.Hash) {

	for {
		time.Sleep(time.Second)
		_, ispending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Warnf("TransactionByHash failed:%v, hash:%s", err, hash.String())
			continue
		}
		if ispending == true {
			log.Infof("waitTransactionConfirm pending, tx:%s", hash.String())
			continue
		} else {
			break
		}
	}

}

func encodeToInt(unsigned bool, size int, arg string) (result interface{}, err error) {
	n, ok := big.NewInt(0).SetString(arg, 10)
	if !ok {
		err = fmt.Errorf("invalid arg for int:%s", arg)
		return
	}
	if n.BitLen() > size {
		err = fmt.Errorf("number overflow, got %s, expected size %d", arg, size)
		return
	}
	if unsigned {
		switch size {
		case 8:
			result = uint8(n.Uint64())
			return
		case 16:
			result = uint16(n.Uint64())
			return
		case 32:
			result = uint32(n.Uint64())
			return
		case 64:
			result = uint64(n.Uint64())
			return
		}
	}
	switch size {
	case 8:
		result = int8(n.Int64())
		return
	case 16:
		result = int16(n.Int64())
		return
	case 32:
		result = int32(n.Int64())
		return
	case 64:
		result = int64(n.Int64())
		return
	}
	result = n
	return
}

func encodeToBool(arg string) (result interface{}, err error) {
	switch arg {
	case "true":
		result = true
	case "false":
		result = false
	default:
		err = fmt.Errorf("invalid bool arg:%s", arg)
	}
	return
}

func encode(t abi.Type, arg string) (result interface{}, err error) {
	switch t.T {
	case abi.IntTy:
		return encodeToInt(false, t.Size, arg)
	case abi.UintTy:
		return encodeToInt(true, t.Size, arg)
	case abi.BoolTy:
		return encodeToBool(arg)
	case abi.StringTy:
		return arg, nil
	case abi.SliceTy:

		var jsonArray []interface{}
		err = json.Unmarshal([]byte(arg), &jsonArray)
		if err != nil {
			err = fmt.Errorf("slice arg not valid json:%v", err)
			return
		}

		size := len(jsonArray)
		refSlice := reflect.MakeSlice(t.GetType(), size, size)
		result = refSlice.Interface()
		err = json.Unmarshal([]byte(arg), &result)
		if err != nil {
			err = fmt.Errorf("slice arg Unmarshal failed:%v arg:%s", err, arg)
			return
		}

	case abi.ArrayTy:
		refSlice := reflect.New(t.GetType()).Elem()
		result = refSlice.Interface()
		err = json.Unmarshal([]byte(arg), &result)
		if err != nil {
			err = fmt.Errorf("array arg Unmarshal failed:%v arg:%s", err, arg)
			return
		}
	case abi.TupleTy:
		err = fmt.Errorf("tuple type not supported yet")
	case abi.AddressTy:
		return common.HexToAddress(arg), nil
	case abi.FixedBytesTy:
		return abi.ReadFixedBytes(t, []byte(arg))
	case abi.BytesTy:
		if has0xPrefix(arg) {
			arg = arg[2:]
		}
		result, err = hex.DecodeString(arg)
		return
	case abi.HashTy:
		err = fmt.Errorf("hash type not supported yet")
	case abi.FixedPointTy:
		err = fmt.Errorf("fixedpoint type not supported yet")
	case abi.FunctionTy:
		err = fmt.Errorf("function type not supported yet")
	default:
		err = fmt.Errorf("invalid type:%d", t.T)
	}
	return
}

func format(t abi.Type, value interface{}) string {
	switch t.T {
	case abi.TupleTy:
		panic("tuple not supported")
	case abi.SliceTy:
		panic("slice not supported")
	case abi.ArrayTy:
		panic("array not supported")
	case abi.StringTy:
		return fmt.Sprintf("%v", value)
	case abi.IntTy, abi.UintTy:
		return fmt.Sprintf("%v", value)
	case abi.BoolTy:
		return fmt.Sprintf("%v", value)
	case abi.AddressTy:
		return value.(common.Address).Hex()
	case abi.HashTy:
		return "0x" + value.(common.Hash).Hex()
	case abi.BytesTy:
		return "0x" + hex.EncodeToString(value.([]byte))
	case abi.FixedBytesTy:
		sliceType := reflect.TypeOf([]byte{})
		return "0x" + hex.EncodeToString(reflect.ValueOf(value).Convert(sliceType).Interface().([]byte))
	case abi.FunctionTy:
		rawV := value.([24]byte)
		return "0x" + hex.EncodeToString(rawV[:])
	default:
		panic(fmt.Sprintf("abi: unknown type %v", t.T))
	}
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
