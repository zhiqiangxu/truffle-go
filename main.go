package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/zhiqinagxu/truffle-go/config"
	"github.com/zhiqinagxu/truffle-go/pkg/contract"
	"github.com/zhiqinagxu/truffle-go/pkg/log"
)

var confFile string
var function string
var sol string
var targetContract string
var contractAddr string
var targetMethod string
var rawFile string
var solc string
var pkFile string

func init() {
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.StringVar(&function, "func", "deploy", "choose function to run: deploy or run")
	flag.StringVar(&sol, "sol", "", "specify solidity file")
	flag.StringVar(&targetContract, "contract", "", "specify target contract")
	flag.StringVar(&contractAddr, "addr", "", "specify contract address")
	flag.StringVar(&targetMethod, "method", "", "specify target method")
	flag.StringVar(&rawFile, "raw", "", "specify raw data file")
	flag.StringVar(&pkFile, "pk", "", "specify pk file")
	flag.StringVar(&solc, "solc", "", "specify solidity compiler")

	flag.Parse()

	// log.Infof("confFile=%s\nfunction=%s\nsol=%s\ntargetContract=%s\nsolc=%s\n", confFile, function, sol, targetContract, solc)
}

func main() {

	conf, err := config.LoadConfig(confFile)
	if err != nil {
		log.Fatal("LoadConfig fail", err)
	}
	if pkFile != "" {
		pkBytes, err := ioutil.ReadFile(pkFile)
		if err != nil {
			log.Fatal("ReadFile fail", err)
		}
		conf.PrivateKey = string(pkBytes)
	}

	switch function {
	case "deploy":
		contract.Deploy(conf, solc, sol, targetContract, flag.Args())
	case "call":
		contract.Call(conf, solc, sol, contractAddr, targetContract, targetMethod, flag.Args())
	case "rawcall", "rawtransact":
		if len(rawFile) == 0 {
			log.Fatal("raw data file not specified")
		}
		isToml := strings.HasSuffix(rawFile, "toml")
		rawDataBytes, err := ioutil.ReadFile(rawFile)
		if err != nil {
			log.Fatal("ReadFile fail", err)
		}

		var param contract.RawParam
		if isToml {
			_, err = toml.Decode(string(rawDataBytes), &param)
		} else {
			err = json.Unmarshal(rawDataBytes, &param)
		}

		if err != nil {
			log.Fatal("LoadConfig fail", err)
		}

		switch function {
		case "rawcall":
			contract.RawCall(conf, param)
		case "rawtransact":
			contract.RawTransact(conf, param)
		}

	case "transact":
		contract.Transact(conf, solc, sol, contractAddr, targetContract, targetMethod, flag.Args())
	default:
		log.Fatal("unknown function", function)
	}
}
