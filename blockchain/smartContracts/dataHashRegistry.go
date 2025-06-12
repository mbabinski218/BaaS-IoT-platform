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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"id\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"HashStored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"records\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes16\",\"name\":\"iotId\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"id\",\"type\":\"bytes16\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes16\",\"name\":\"iotId\",\"type\":\"bytes16\"}],\"name\":\"storeHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"id\",\"type\":\"bytes16\"},{\"internalType\":\"bytes32\",\"name\":\"providedHash\",\"type\":\"bytes32\"}],\"name\":\"verifyHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506106288061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80632269a1e114610043578063b1932b251461005f578063dd864ab21461008f575b5f5ffd5b61005d600480360381019061005891906103bd565b6100c2565b005b6100796004803603810190610074919061040d565b610297565b6040516100869190610465565b60405180910390f35b6100a960048036038101906100a4919061047e565b6102db565b6040516100b9949392919061051e565b60405180910390f35b5f5f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f20600201541461013d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610134906105bb565b60405180910390fd5b6040518060800160405280838152602001826fffffffffffffffffffffffffffffffff191681526020014281526020013373ffffffffffffffffffffffffffffffffffffffff168152505f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f820151815f01556020820151816001015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c0217905550604082015181600201556060820151816003015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550905050826fffffffffffffffffffffffffffffffff19167fe2d3c54e250a34c193a79d28f9c0ca4dedc53ef3b02a88cdc9d349bb906fe5c88360405161028a91906105d9565b60405180910390a2505050565b5f815f5f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f015414905092915050565b5f602052805f5260405f205f91509050805f015490806001015f9054906101000a900460801b90806002015490806003015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905084565b5f5ffd5b5f7fffffffffffffffffffffffffffffffff0000000000000000000000000000000082169050919050565b61036981610335565b8114610373575f5ffd5b50565b5f8135905061038481610360565b92915050565b5f819050919050565b61039c8161038a565b81146103a6575f5ffd5b50565b5f813590506103b781610393565b92915050565b5f5f5f606084860312156103d4576103d3610331565b5b5f6103e186828701610376565b93505060206103f2868287016103a9565b925050604061040386828701610376565b9150509250925092565b5f5f6040838503121561042357610422610331565b5b5f61043085828601610376565b9250506020610441858286016103a9565b9150509250929050565b5f8115159050919050565b61045f8161044b565b82525050565b5f6020820190506104785f830184610456565b92915050565b5f6020828403121561049357610492610331565b5b5f6104a084828501610376565b91505092915050565b6104b28161038a565b82525050565b6104c181610335565b82525050565b5f819050919050565b6104d9816104c7565b82525050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610508826104df565b9050919050565b610518816104fe565b82525050565b5f6080820190506105315f8301876104a9565b61053e60208301866104b8565b61054b60408301856104d0565b610558606083018461050f565b95945050505050565b5f82825260208201905092915050565b7f4861736820616c72656164792065786973747320666f722074686973204964005f82015250565b5f6105a5601f83610561565b91506105b082610571565b602082019050919050565b5f6020820190508181035f8301526105d281610599565b9050919050565b5f6020820190506105ec5f8301846104a9565b9291505056fea2646970667358221220722367fdcb9a1c714eccbb390ef427d38b7e9d4532e798e8840046bff84ac11a64736f6c634300081e0033",
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

// DataHashRegistryHashStoredIterator is returned from FilterHashStored and is used to iterate over the raw logs and unpacked data for HashStored events raised by the DataHashRegistry contract.
type DataHashRegistryHashStoredIterator struct {
	Event *DataHashRegistryHashStored // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DataHashRegistryHashStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataHashRegistryHashStored)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DataHashRegistryHashStored)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DataHashRegistryHashStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataHashRegistryHashStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataHashRegistryHashStored represents a HashStored event raised by the DataHashRegistry contract.
type DataHashRegistryHashStored struct {
	Id       [16]byte
	DataHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterHashStored is a free log retrieval operation binding the contract event 0xe2d3c54e250a34c193a79d28f9c0ca4dedc53ef3b02a88cdc9d349bb906fe5c8.
//
// Solidity: event HashStored(bytes16 indexed id, bytes32 dataHash)
func (_DataHashRegistry *DataHashRegistryFilterer) FilterHashStored(opts *bind.FilterOpts, id [][16]byte) (*DataHashRegistryHashStoredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _DataHashRegistry.contract.FilterLogs(opts, "HashStored", idRule)
	if err != nil {
		return nil, err
	}
	return &DataHashRegistryHashStoredIterator{contract: _DataHashRegistry.contract, event: "HashStored", logs: logs, sub: sub}, nil
}

// WatchHashStored is a free log subscription operation binding the contract event 0xe2d3c54e250a34c193a79d28f9c0ca4dedc53ef3b02a88cdc9d349bb906fe5c8.
//
// Solidity: event HashStored(bytes16 indexed id, bytes32 dataHash)
func (_DataHashRegistry *DataHashRegistryFilterer) WatchHashStored(opts *bind.WatchOpts, sink chan<- *DataHashRegistryHashStored, id [][16]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _DataHashRegistry.contract.WatchLogs(opts, "HashStored", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataHashRegistryHashStored)
				if err := _DataHashRegistry.contract.UnpackLog(event, "HashStored", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHashStored is a log parse operation binding the contract event 0xe2d3c54e250a34c193a79d28f9c0ca4dedc53ef3b02a88cdc9d349bb906fe5c8.
//
// Solidity: event HashStored(bytes16 indexed id, bytes32 dataHash)
func (_DataHashRegistry *DataHashRegistryFilterer) ParseHashStored(log types.Log) (*DataHashRegistryHashStored, error) {
	event := new(DataHashRegistryHashStored)
	if err := _DataHashRegistry.contract.UnpackLog(event, "HashStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
