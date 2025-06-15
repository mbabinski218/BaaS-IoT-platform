// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smartContracts

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

// BatchRegistryMetaData contains all meta data concerning the BatchRegistry contract.
var BatchRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"roots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"storeRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"providedMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"verifyRoot\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061035e8061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c806396d4dd7a14610043578063b0cd48021461005f578063c2b40ae41461008f575b5f5ffd5b61005d600480360381019061005891906101c6565b6100bf565b005b610079600480360381019061007491906101c6565b61012c565b604051610086919061021e565b60405180910390f35b6100a960048036038101906100a49190610237565b610148565b6040516100b69190610271565b60405180910390f35b5f5f1b5f5f8481526020019081526020015f205414610113576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161010a9061030a565b60405180910390fd5b805f5f8481526020019081526020015f20819055505050565b5f815f5f8581526020019081526020015f205414905092915050565b5f602052805f5260405f205f915090505481565b5f5ffd5b5f819050919050565b61017281610160565b811461017c575f5ffd5b50565b5f8135905061018d81610169565b92915050565b5f819050919050565b6101a581610193565b81146101af575f5ffd5b50565b5f813590506101c08161019c565b92915050565b5f5f604083850312156101dc576101db61015c565b5b5f6101e98582860161017f565b92505060206101fa858286016101b2565b9150509250929050565b5f8115159050919050565b61021881610204565b82525050565b5f6020820190506102315f83018461020f565b92915050565b5f6020828403121561024c5761024b61015c565b5b5f6102598482850161017f565b91505092915050565b61026b81610193565b82525050565b5f6020820190506102845f830184610262565b92915050565b5f82825260208201905092915050565b7f416c72656164792073746f72656420666f7220746869732062617463682074695f8201527f6d65000000000000000000000000000000000000000000000000000000000000602082015250565b5f6102f460228361028a565b91506102ff8261029a565b604082019050919050565b5f6020820190508181035f830152610321816102e8565b905091905056fea26469706673582212207f1f2b72c89c869b8a3db3ed06725b5ec950751876d9b51cd3705f90a865025d64736f6c634300081e0033",
}

// BatchRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchRegistryMetaData.ABI instead.
var BatchRegistryABI = BatchRegistryMetaData.ABI

// BatchRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchRegistryMetaData.Bin instead.
var BatchRegistryBin = BatchRegistryMetaData.Bin

// DeployBatchRegistry deploys a new Ethereum contract, binding an instance of BatchRegistry to it.
func DeployBatchRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BatchRegistry, error) {
	parsed, err := BatchRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchRegistry{BatchRegistryCaller: BatchRegistryCaller{contract: contract}, BatchRegistryTransactor: BatchRegistryTransactor{contract: contract}, BatchRegistryFilterer: BatchRegistryFilterer{contract: contract}}, nil
}

// BatchRegistry is an auto generated Go binding around an Ethereum contract.
type BatchRegistry struct {
	BatchRegistryCaller     // Read-only binding to the contract
	BatchRegistryTransactor // Write-only binding to the contract
	BatchRegistryFilterer   // Log filterer for contract events
}

// BatchRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchRegistrySession struct {
	Contract     *BatchRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BatchRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchRegistryCallerSession struct {
	Contract *BatchRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BatchRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchRegistryTransactorSession struct {
	Contract     *BatchRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BatchRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchRegistryRaw struct {
	Contract *BatchRegistry // Generic contract binding to access the raw methods on
}

// BatchRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchRegistryCallerRaw struct {
	Contract *BatchRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// BatchRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchRegistryTransactorRaw struct {
	Contract *BatchRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchRegistry creates a new instance of BatchRegistry, bound to a specific deployed contract.
func NewBatchRegistry(address common.Address, backend bind.ContractBackend) (*BatchRegistry, error) {
	contract, err := bindBatchRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchRegistry{BatchRegistryCaller: BatchRegistryCaller{contract: contract}, BatchRegistryTransactor: BatchRegistryTransactor{contract: contract}, BatchRegistryFilterer: BatchRegistryFilterer{contract: contract}}, nil
}

// NewBatchRegistryCaller creates a new read-only instance of BatchRegistry, bound to a specific deployed contract.
func NewBatchRegistryCaller(address common.Address, caller bind.ContractCaller) (*BatchRegistryCaller, error) {
	contract, err := bindBatchRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchRegistryCaller{contract: contract}, nil
}

// NewBatchRegistryTransactor creates a new write-only instance of BatchRegistry, bound to a specific deployed contract.
func NewBatchRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchRegistryTransactor, error) {
	contract, err := bindBatchRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchRegistryTransactor{contract: contract}, nil
}

// NewBatchRegistryFilterer creates a new log filterer instance of BatchRegistry, bound to a specific deployed contract.
func NewBatchRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchRegistryFilterer, error) {
	contract, err := bindBatchRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchRegistryFilterer{contract: contract}, nil
}

// bindBatchRegistry binds a generic wrapper to an already deployed contract.
func bindBatchRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchRegistry *BatchRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchRegistry.Contract.BatchRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchRegistry *BatchRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchRegistry.Contract.BatchRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchRegistry *BatchRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchRegistry.Contract.BatchRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchRegistry *BatchRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchRegistry *BatchRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchRegistry *BatchRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchRegistry.Contract.contract.Transact(opts, method, params...)
}

// Roots is a free data retrieval call binding the contract method 0xc2b40ae4.
//
// Solidity: function roots(uint256 ) view returns(bytes32)
func (_BatchRegistry *BatchRegistryCaller) Roots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BatchRegistry.contract.Call(opts, &out, "roots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Roots is a free data retrieval call binding the contract method 0xc2b40ae4.
//
// Solidity: function roots(uint256 ) view returns(bytes32)
func (_BatchRegistry *BatchRegistrySession) Roots(arg0 *big.Int) ([32]byte, error) {
	return _BatchRegistry.Contract.Roots(&_BatchRegistry.CallOpts, arg0)
}

// Roots is a free data retrieval call binding the contract method 0xc2b40ae4.
//
// Solidity: function roots(uint256 ) view returns(bytes32)
func (_BatchRegistry *BatchRegistryCallerSession) Roots(arg0 *big.Int) ([32]byte, error) {
	return _BatchRegistry.Contract.Roots(&_BatchRegistry.CallOpts, arg0)
}

// VerifyRoot is a free data retrieval call binding the contract method 0xb0cd4802.
//
// Solidity: function verifyRoot(uint256 batchTime, bytes32 providedMerkleRoot) view returns(bool)
func (_BatchRegistry *BatchRegistryCaller) VerifyRoot(opts *bind.CallOpts, batchTime *big.Int, providedMerkleRoot [32]byte) (bool, error) {
	var out []interface{}
	err := _BatchRegistry.contract.Call(opts, &out, "verifyRoot", batchTime, providedMerkleRoot)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyRoot is a free data retrieval call binding the contract method 0xb0cd4802.
//
// Solidity: function verifyRoot(uint256 batchTime, bytes32 providedMerkleRoot) view returns(bool)
func (_BatchRegistry *BatchRegistrySession) VerifyRoot(batchTime *big.Int, providedMerkleRoot [32]byte) (bool, error) {
	return _BatchRegistry.Contract.VerifyRoot(&_BatchRegistry.CallOpts, batchTime, providedMerkleRoot)
}

// VerifyRoot is a free data retrieval call binding the contract method 0xb0cd4802.
//
// Solidity: function verifyRoot(uint256 batchTime, bytes32 providedMerkleRoot) view returns(bool)
func (_BatchRegistry *BatchRegistryCallerSession) VerifyRoot(batchTime *big.Int, providedMerkleRoot [32]byte) (bool, error) {
	return _BatchRegistry.Contract.VerifyRoot(&_BatchRegistry.CallOpts, batchTime, providedMerkleRoot)
}

// StoreRoot is a paid mutator transaction binding the contract method 0x96d4dd7a.
//
// Solidity: function storeRoot(uint256 batchTime, bytes32 merkleRoot) returns()
func (_BatchRegistry *BatchRegistryTransactor) StoreRoot(opts *bind.TransactOpts, batchTime *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _BatchRegistry.contract.Transact(opts, "storeRoot", batchTime, merkleRoot)
}

// StoreRoot is a paid mutator transaction binding the contract method 0x96d4dd7a.
//
// Solidity: function storeRoot(uint256 batchTime, bytes32 merkleRoot) returns()
func (_BatchRegistry *BatchRegistrySession) StoreRoot(batchTime *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _BatchRegistry.Contract.StoreRoot(&_BatchRegistry.TransactOpts, batchTime, merkleRoot)
}

// StoreRoot is a paid mutator transaction binding the contract method 0x96d4dd7a.
//
// Solidity: function storeRoot(uint256 batchTime, bytes32 merkleRoot) returns()
func (_BatchRegistry *BatchRegistryTransactorSession) StoreRoot(batchTime *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _BatchRegistry.Contract.StoreRoot(&_BatchRegistry.TransactOpts, batchTime, merkleRoot)
}
