// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dataHashRegistry

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DataHashRegistryMetaData contains all meta data concerning the DataHashRegistry contract.
var DataHashRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes16\",\"name\":\"iotId\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"id\",\"type\":\"bytes16\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes16\",\"name\":\"iotId\",\"type\":\"bytes16\"}],\"name\":\"storeHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"id\",\"type\":\"bytes16\"},{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"}],\"name\":\"verifyHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506105c48061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80632269a1e114610043578063b1932b251461005f578063dd864ab21461008f575b5f5ffd5b61005d60048036038101906100589190610372565b6100c2565b005b610079600480360381019061007491906103c2565b61024c565b604051610086919061041a565b60405180910390f35b6100a960048036038101906100a49190610433565b610290565b6040516100b994939291906104d3565b60405180910390f35b5f5f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f20600201541461013d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013490610570565b60405180910390fd5b6040518060800160405280838152602001826fffffffffffffffffffffffffffffffff191681526020014281526020013373ffffffffffffffffffffffffffffffffffffffff168152505f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f820151815f01556020820151816001015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c0217905550604082015181600201556060820151816003015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550905050505050565b5f815f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f015414905092915050565b5f602052805f5260405f205f91509050805f015490806001015f9054906101000a900460801b90806002015490806003015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905084565b5f5ffd5b5f7fffffffffffffffffffffffffffffffff0000000000000000000000000000000082169050919050565b61031e816102ea565b8114610328575f5ffd5b50565b5f8135905061033981610315565b92915050565b5f819050919050565b6103518161033f565b811461035b575f5ffd5b50565b5f8135905061036c81610348565b92915050565b5f5f5f60608486031215610389576103886102e6565b5b5f6103968682870161032b565b93505060206103a78682870161035e565b92505060406103b88682870161032b565b9150509250925092565b5f5f604083850312156103d8576103d76102e6565b5b5f6103e58582860161032b565b92505060206103f68582860161035e565b9150509250929050565b5f8115159050919050565b61041481610400565b82525050565b5f60208201905061042d5f83018461040b565b92915050565b5f60208284031215610448576104476102e6565b5b5f6104558482850161032b565b91505092915050565b6104678161033f565b82525050565b610476816102ea565b82525050565b5f819050919050565b61048e8161047c565b82525050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6104bd82610494565b9050919050565b6104cd816104b3565b82525050565b5f6080820190506104e65f83018761045e565b6104f3602083018661046d565b6105006040830185610485565b61050d60608301846104c4565b95945050505050565b5f82825260208201905092915050565b7f4861736820616c72656164792065786973747320666f722074686973204964005f82015250565b5f61055a601f83610516565b915061056582610526565b602082019050919050565b5f6020820190508181035f8301526105878161054e565b905091905056fea26469706673582212201f0dd6e90a1e423e867c5a3f7d6e5e692cd59c65417b0652c2e4303c58281c1a64736f6c634300081e0033",
}

// DataHashRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use DataHashRegistryMetaData.ABI instead.
var DataHashRegistryABI = DataHashRegistryMetaData.ABI

// DataHashRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DataHashRegistryMetaData.Bin instead.
var DataHashRegistryBin = DataHashRegistryMetaData.Bin

// DeployDataHashRegistry deploys a new Ethereum contract, binding an instance of DataHashRegistry to it.
func DeployDataHashRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DataHashRegistry, error) {
	parsed, err := DataHashRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DataHashRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DataHashRegistry{DataHashRegistryCaller: DataHashRegistryCaller{contract: contract}, DataHashRegistryTransactor: DataHashRegistryTransactor{contract: contract}, DataHashRegistryFilterer: DataHashRegistryFilterer{contract: contract}}, nil
}

// DataHashRegistry is an auto generated Go binding around an Ethereum contract.
type DataHashRegistry struct {
	DataHashRegistryCaller     // Read-only binding to the contract
	DataHashRegistryTransactor // Write-only binding to the contract
	DataHashRegistryFilterer   // Log filterer for contract events
}

// DataHashRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataHashRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataHashRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataHashRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataHashRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataHashRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataHashRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataHashRegistrySession struct {
	Contract     *DataHashRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DataHashRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataHashRegistryCallerSession struct {
	Contract *DataHashRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DataHashRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataHashRegistryTransactorSession struct {
	Contract     *DataHashRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DataHashRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataHashRegistryRaw struct {
	Contract *DataHashRegistry // Generic contract binding to access the raw methods on
}

// DataHashRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataHashRegistryCallerRaw struct {
	Contract *DataHashRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// DataHashRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataHashRegistryTransactorRaw struct {
	Contract *DataHashRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataHashRegistry creates a new instance of DataHashRegistry, bound to a specific deployed contract.
func NewDataHashRegistry(address common.Address, backend bind.ContractBackend) (*DataHashRegistry, error) {
	contract, err := bindDataHashRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataHashRegistry{DataHashRegistryCaller: DataHashRegistryCaller{contract: contract}, DataHashRegistryTransactor: DataHashRegistryTransactor{contract: contract}, DataHashRegistryFilterer: DataHashRegistryFilterer{contract: contract}}, nil
}

// NewDataHashRegistryCaller creates a new read-only instance of DataHashRegistry, bound to a specific deployed contract.
func NewDataHashRegistryCaller(address common.Address, caller bind.ContractCaller) (*DataHashRegistryCaller, error) {
	contract, err := bindDataHashRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataHashRegistryCaller{contract: contract}, nil
}

// NewDataHashRegistryTransactor creates a new write-only instance of DataHashRegistry, bound to a specific deployed contract.
func NewDataHashRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DataHashRegistryTransactor, error) {
	contract, err := bindDataHashRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataHashRegistryTransactor{contract: contract}, nil
}

// NewDataHashRegistryFilterer creates a new log filterer instance of DataHashRegistry, bound to a specific deployed contract.
func NewDataHashRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DataHashRegistryFilterer, error) {
	contract, err := bindDataHashRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataHashRegistryFilterer{contract: contract}, nil
}

// bindDataHashRegistry binds a generic wrapper to an already deployed contract.
func bindDataHashRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataHashRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataHashRegistry *DataHashRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataHashRegistry.Contract.DataHashRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataHashRegistry *DataHashRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.DataHashRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataHashRegistry *DataHashRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.DataHashRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataHashRegistry *DataHashRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataHashRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataHashRegistry *DataHashRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataHashRegistry *DataHashRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.contract.Transact(opts, method, params...)
}

// Records is a free data retrieval call binding the contract method 0xdd864ab2.
//
// Solidity: function records(bytes16 ) view returns(bytes32 dataHash, bytes16 iotId, uint256 timestamp, address sender)
func (_DataHashRegistry *DataHashRegistryCaller) Records(opts *bind.CallOpts, arg0 [16]byte) (struct {
	DataHash  [32]byte
	IotId     [16]byte
	Timestamp *big.Int
	Sender    common.Address
}, error) {
	var out []interface{}
	err := _DataHashRegistry.contract.Call(opts, &out, "records", arg0)

	outstruct := new(struct {
		DataHash  [32]byte
		IotId     [16]byte
		Timestamp *big.Int
		Sender    common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DataHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.IotId = *abi.ConvertType(out[1], new([16]byte)).(*[16]byte)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Sender = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Records is a free data retrieval call binding the contract method 0xdd864ab2.
//
// Solidity: function records(bytes16 ) view returns(bytes32 dataHash, bytes16 iotId, uint256 timestamp, address sender)
func (_DataHashRegistry *DataHashRegistrySession) Records(arg0 [16]byte) (struct {
	DataHash  [32]byte
	IotId     [16]byte
	Timestamp *big.Int
	Sender    common.Address
}, error) {
	return _DataHashRegistry.Contract.Records(&_DataHashRegistry.CallOpts, arg0)
}

// Records is a free data retrieval call binding the contract method 0xdd864ab2.
//
// Solidity: function records(bytes16 ) view returns(bytes32 dataHash, bytes16 iotId, uint256 timestamp, address sender)
func (_DataHashRegistry *DataHashRegistryCallerSession) Records(arg0 [16]byte) (struct {
	DataHash  [32]byte
	IotId     [16]byte
	Timestamp *big.Int
	Sender    common.Address
}, error) {
	return _DataHashRegistry.Contract.Records(&_DataHashRegistry.CallOpts, arg0)
}

// VerifyHash is a free data retrieval call binding the contract method 0xb1932b25.
//
// Solidity: function verifyHash(bytes16 id, bytes32 providedHash) view returns(bool)
func (_DataHashRegistry *DataHashRegistryCaller) VerifyHash(opts *bind.CallOpts, id [16]byte, providedHash [32]byte) (bool, error) {
	var out []interface{}
	err := _DataHashRegistry.contract.Call(opts, &out, "verifyHash", id, providedHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyHash is a free data retrieval call binding the contract method 0xb1932b25.
//
// Solidity: function verifyHash(bytes16 id, bytes32 providedHash) view returns(bool)
func (_DataHashRegistry *DataHashRegistrySession) VerifyHash(id [16]byte, providedHash [32]byte) (bool, error) {
	return _DataHashRegistry.Contract.VerifyHash(&_DataHashRegistry.CallOpts, id, providedHash)
}

// VerifyHash is a free data retrieval call binding the contract method 0xb1932b25.
//
// Solidity: function verifyHash(bytes16 id, bytes32 providedHash) view returns(bool)
func (_DataHashRegistry *DataHashRegistryCallerSession) VerifyHash(id [16]byte, providedHash [32]byte) (bool, error) {
	return _DataHashRegistry.Contract.VerifyHash(&_DataHashRegistry.CallOpts, id, providedHash)
}

// StoreHash is a paid mutator transaction binding the contract method 0x2269a1e1.
//
// Solidity: function storeHash(bytes16 id, bytes32 dataHash, bytes16 iotId) returns()
func (_DataHashRegistry *DataHashRegistryTransactor) StoreHash(opts *bind.TransactOpts, id [16]byte, dataHash [32]byte, iotId [16]byte) (*types.Transaction, error) {
	return _DataHashRegistry.contract.Transact(opts, "storeHash", id, dataHash, iotId)
}

// StoreHash is a paid mutator transaction binding the contract method 0x2269a1e1.
//
// Solidity: function storeHash(bytes16 id, bytes32 dataHash, bytes16 iotId) returns()
func (_DataHashRegistry *DataHashRegistrySession) StoreHash(id [16]byte, dataHash [32]byte, iotId [16]byte) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.StoreHash(&_DataHashRegistry.TransactOpts, id, dataHash, iotId)
}

// StoreHash is a paid mutator transaction binding the contract method 0x2269a1e1.
//
// Solidity: function storeHash(bytes16 id, bytes32 dataHash, bytes16 iotId) returns()
func (_DataHashRegistry *DataHashRegistryTransactorSession) StoreHash(id [16]byte, dataHash [32]byte, iotId [16]byte) (*types.Transaction, error) {
	return _DataHashRegistry.Contract.StoreHash(&_DataHashRegistry.TransactOpts, id, dataHash, iotId)
}
