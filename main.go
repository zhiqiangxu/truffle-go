package main

import (
	"flag"

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
var solc string

func init() {
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.StringVar(&function, "func", "deploy", "choose function to run: deploy or run")
	flag.StringVar(&sol, "sol", "", "specify solidity file")
	flag.StringVar(&targetContract, "contract", "", "specify target contract")
	flag.StringVar(&contractAddr, "addr", "", "specify contract address")
	flag.StringVar(&targetMethod, "method", "", "specify target method")
	flag.StringVar(&solc, "solc", "", "specify solidity compiler")

	flag.Parse()

	// log.Infof("confFile=%s\nfunction=%s\nsol=%s\ntargetContract=%s\nsolc=%s\n", confFile, function, sol, targetContract, solc)
}

func main() {

	conf, err := config.LoadConfig(confFile)
	if err != nil {
		log.Fatal("LoadConfig fail", err)
	}

	switch function {
	case "deploy":
		contract.Deploy(conf, solc, sol, targetContract, flag.Args())
	case "call":
		contract.Call(conf, solc, sol, contractAddr, targetContract, targetMethod, flag.Args())
	case "rawcall":
		if len(flag.Args()) > 1 {
			log.Fatal("rawcall expects a single combined hex argument")
		}
		var hexarg string
		if len(flag.Args()) == 1 {
			hexarg = flag.Args()[0]
		}
		contract.RawCall(conf, contractAddr, targetMethod, hexarg)
	case "transact":
		contract.Transact(conf, solc, sol, contractAddr, targetContract, targetMethod, flag.Args())
	default:
		log.Fatal("unknown function", function)
	}
}
