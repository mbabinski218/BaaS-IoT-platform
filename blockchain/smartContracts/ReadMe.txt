npx solcjs --bin --abi dataHashRegistry.sol -o build/
abigen --bin=build/dataHashRegistry_sol_DataHashRegistry.bin --abi=build/dataHashRegistry_sol_DataHashRegistry.abi --pkg=smartContracts --out=dataHashRegistry.go --type DataHashRegistry

npx solcjs --bin --abi batchRegistry.sol -o build/
abigen --bin=build/batchRegistry_sol_BatchRegistry.bin --abi=build/batchRegistry_sol_BatchRegistry.abi --pkg=smartContracts --out=batchRegistry.go --type BatchRegistry 


sudo kurtosis run github.com/ethpandaops/ethereum-package --args-file /mnt/c/code/BaaS-IoT-platform/blockchain/testnet/network_params.yaml --image-download always