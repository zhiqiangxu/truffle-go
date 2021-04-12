# truffle-go, truffle for gophers

## steps

```
# deploy contract
go run main.go -func deploy -solc '/Users/xuzhiqiang/Downloads/solc-macosx-amd64-v0.6.0+commit.26b70077'  -sol '/Users/xuzhiqiang/Desktop/workspace/opensource/solidity_contracts/poly-swap/contracts/core/lock_proxy/selectSwap.sol' -contract Select

# call contract 
go run main.go -func call -solc '/Users/xuzhiqiang/Downloads/solc-macosx-amd64-v0.6.0+commit.26b70077'  -sol '/Users/xuzhiqiang/Desktop/workspace/opensource/solidity_contracts/poly-swap/contracts/core/lock_proxy/selectSwap.sol' -contract Select -addr 0x44F4e537797845a92573c08C9e855352c2CF63B2 -method getSwapFee 0x0687e6392de735B83ed2808797c92051B5dF5618

# transact contract
go run main.go -func transact -solc '/Users/xuzhiqiang/Downloads/solc-macosx-amd64-v0.6.0+commit.26b70077'  -sol '/Users/xuzhiqiang/Desktop/workspace/opensource/solidity_contracts/poly-swap/contracts/core/lock_proxy/selectSwap.sol' -contract Select -addr 0x44F4e537797845a92573c08C9e855352c2CF63B2 -method getSwapFee 0x0687e6392de735B83ed2808797c92051B5dF5618
```