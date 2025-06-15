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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"roots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"storeRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"leaf\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"verifyProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"providedMerkleRoot\",\"type\":\"bytes32\"}],\"name\":\"verifyRoot\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506105aa8061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c806396d4dd7a1461004e578063b0cd48021461006a578063c2b40ae41461009a578063c48555c4146100ca575b5f5ffd5b610068600480360381019061006391906102c8565b6100fa565b005b610084600480360381019061007f91906102c8565b610167565b6040516100919190610320565b60405180910390f35b6100b460048036038101906100af9190610339565b610183565b6040516100c19190610373565b60405180910390f35b6100e460048036038101906100df91906103ed565b610197565b6040516100f19190610320565b60405180910390f35b5f5f1b5f5f8481526020019081526020015f20541461014e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610145906104de565b60405180910390fd5b805f5f8481526020019081526020015f20819055505050565b5f815f5f8581526020019081526020015f205414905092915050565b5f602052805f5260405f205f915090505481565b5f5f8490505f5f90505b8484905081101561023b575f8585838181106101c0576101bf6104fc565b5b905060200201359050808310156102015782816040516020016101e4929190610549565b60405160208183030381529060405280519060200120925061022d565b8083604051602001610214929190610549565b6040516020818303038152906040528051906020012092505b5080806001019150506101a1565b505f5f8781526020019081526020015f20548114915050949350505050565b5f5ffd5b5f5ffd5b5f819050919050565b61027481610262565b811461027e575f5ffd5b50565b5f8135905061028f8161026b565b92915050565b5f819050919050565b6102a781610295565b81146102b1575f5ffd5b50565b5f813590506102c28161029e565b92915050565b5f5f604083850312156102de576102dd61025a565b5b5f6102eb85828601610281565b92505060206102fc858286016102b4565b9150509250929050565b5f8115159050919050565b61031a81610306565b82525050565b5f6020820190506103335f830184610311565b92915050565b5f6020828403121561034e5761034d61025a565b5b5f61035b84828501610281565b91505092915050565b61036d81610295565b82525050565b5f6020820190506103865f830184610364565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126103ad576103ac61038c565b5b8235905067ffffffffffffffff8111156103ca576103c9610390565b5b6020830191508360208202830111156103e6576103e5610394565b5b9250929050565b5f5f5f5f606085870312156104055761040461025a565b5b5f61041287828801610281565b9450506020610423878288016102b4565b935050604085013567ffffffffffffffff8111156104445761044361025e565b5b61045087828801610398565b925092505092959194509250565b5f82825260208201905092915050565b7f416c72656164792073746f72656420666f7220746869732062617463682074695f8201527f6d65000000000000000000000000000000000000000000000000000000000000602082015250565b5f6104c860228361045e565b91506104d38261046e565b604082019050919050565b5f6020820190508181035f8301526104f5816104bc565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f819050919050565b61054361053e82610295565b610529565b82525050565b5f6105548285610532565b6020820191506105648284610532565b602082019150819050939250505056fea2646970667358221220129abe3071ad51802ab9e75e49a8335254cc717b845151856afb169f3bbd11d664736f6c634300081e0033",
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

// VerifyProof is a free data retrieval call binding the contract method 0xc48555c4.
//
// Solidity: function verifyProof(uint256 batchTime, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_BatchRegistry *BatchRegistryCaller) VerifyProof(opts *bind.CallOpts, batchTime *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _BatchRegistry.contract.Call(opts, &out, "verifyProof", batchTime, leaf, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0xc48555c4.
//
// Solidity: function verifyProof(uint256 batchTime, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_BatchRegistry *BatchRegistrySession) VerifyProof(batchTime *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _BatchRegistry.Contract.VerifyProof(&_BatchRegistry.CallOpts, batchTime, leaf, proof)
}

// VerifyProof is a free data retrieval call binding the contract method 0xc48555c4.
//
// Solidity: function verifyProof(uint256 batchTime, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_BatchRegistry *BatchRegistryCallerSession) VerifyProof(batchTime *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _BatchRegistry.Contract.VerifyProof(&_BatchRegistry.CallOpts, batchTime, leaf, proof)
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
