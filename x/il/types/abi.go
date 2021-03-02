// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package types

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TypesABI is the input ABI used to generate the binding from.
const TypesABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_uni_router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"redeemLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"redeemLiquidityETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_b\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"state_tokenContract\",\"type\":\"address\"}],\"name\":\"transferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Types is an auto generated Go binding around an Ethereum contract.
type Types struct {
	TypesCaller     // Read-only binding to the contract
	TypesTransactor // Write-only binding to the contract
	TypesFilterer   // Log filterer for contract events
}

// TypesCaller is an auto generated read-only Go binding around an Ethereum contract.
type TypesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TypesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TypesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TypesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TypesSession struct {
	Contract     *Types            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TypesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TypesCallerSession struct {
	Contract *TypesCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TypesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TypesTransactorSession struct {
	Contract     *TypesTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TypesRaw is an auto generated low-level Go binding around an Ethereum contract.
type TypesRaw struct {
	Contract *Types // Generic contract binding to access the raw methods on
}

// TypesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TypesCallerRaw struct {
	Contract *TypesCaller // Generic read-only contract binding to access the raw methods on
}

// TypesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TypesTransactorRaw struct {
	Contract *TypesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTypes creates a new instance of Types, bound to a specific deployed contract.
func NewTypes(address common.Address, backend bind.ContractBackend) (*Types, error) {
	contract, err := bindTypes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Types{TypesCaller: TypesCaller{contract: contract}, TypesTransactor: TypesTransactor{contract: contract}, TypesFilterer: TypesFilterer{contract: contract}}, nil
}

// NewTypesCaller creates a new read-only instance of Types, bound to a specific deployed contract.
func NewTypesCaller(address common.Address, caller bind.ContractCaller) (*TypesCaller, error) {
	contract, err := bindTypes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TypesCaller{contract: contract}, nil
}

// NewTypesTransactor creates a new write-only instance of Types, bound to a specific deployed contract.
func NewTypesTransactor(address common.Address, transactor bind.ContractTransactor) (*TypesTransactor, error) {
	contract, err := bindTypes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TypesTransactor{contract: contract}, nil
}

// NewTypesFilterer creates a new log filterer instance of Types, bound to a specific deployed contract.
func NewTypesFilterer(address common.Address, filterer bind.ContractFilterer) (*TypesFilterer, error) {
	contract, err := bindTypes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TypesFilterer{contract: contract}, nil
}

// bindTypes binds a generic wrapper to an already deployed contract.
func bindTypes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TypesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Types *TypesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Types.Contract.TypesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Types *TypesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Types.Contract.TypesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Types *TypesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Types.Contract.TypesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Types *TypesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Types.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Types *TypesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Types.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Types *TypesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Types.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Types *TypesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Types.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Types *TypesSession) Owner() (common.Address, error) {
	return _Types.Contract.Owner(&_Types.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Types *TypesCallerSession) Owner() (common.Address, error) {
	return _Types.Contract.Owner(&_Types.CallOpts)
}

// RedeemLiquidity is a paid mutator transaction binding the contract method 0x6f221a7a.
//
// Solidity: function redeemLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns()
func (_Types *TypesTransactor) RedeemLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.contract.Transact(opts, "redeemLiquidity", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RedeemLiquidity is a paid mutator transaction binding the contract method 0x6f221a7a.
//
// Solidity: function redeemLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns()
func (_Types *TypesSession) RedeemLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.Contract.RedeemLiquidity(&_Types.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RedeemLiquidity is a paid mutator transaction binding the contract method 0x6f221a7a.
//
// Solidity: function redeemLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns()
func (_Types *TypesTransactorSession) RedeemLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.Contract.RedeemLiquidity(&_Types.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RedeemLiquidityETH is a paid mutator transaction binding the contract method 0x0f77f6b2.
//
// Solidity: function redeemLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_Types *TypesTransactor) RedeemLiquidityETH(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.contract.Transact(opts, "redeemLiquidityETH", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RedeemLiquidityETH is a paid mutator transaction binding the contract method 0x0f77f6b2.
//
// Solidity: function redeemLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_Types *TypesSession) RedeemLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.Contract.RedeemLiquidityETH(&_Types.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RedeemLiquidityETH is a paid mutator transaction binding the contract method 0x0f77f6b2.
//
// Solidity: function redeemLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_Types *TypesTransactorSession) RedeemLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Types.Contract.RedeemLiquidityETH(&_Types.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Types *TypesTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Types.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Types *TypesSession) RenounceOwnership() (*types.Transaction, error) {
	return _Types.Contract.RenounceOwnership(&_Types.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Types *TypesTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Types.Contract.RenounceOwnership(&_Types.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Types *TypesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Types.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Types *TypesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Types.Contract.TransferOwnership(&_Types.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Types *TypesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Types.Contract.TransferOwnership(&_Types.TransactOpts, newOwner)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xd63dd196.
//
// Solidity: function transferTokens(address _to, uint256 _a, uint256 _b, address state_tokenContract) returns()
func (_Types *TypesTransactor) TransferTokens(opts *bind.TransactOpts, _to common.Address, _a *big.Int, _b *big.Int, state_tokenContract common.Address) (*types.Transaction, error) {
	return _Types.contract.Transact(opts, "transferTokens", _to, _a, _b, state_tokenContract)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xd63dd196.
//
// Solidity: function transferTokens(address _to, uint256 _a, uint256 _b, address state_tokenContract) returns()
func (_Types *TypesSession) TransferTokens(_to common.Address, _a *big.Int, _b *big.Int, state_tokenContract common.Address) (*types.Transaction, error) {
	return _Types.Contract.TransferTokens(&_Types.TransactOpts, _to, _a, _b, state_tokenContract)
}

// TransferTokens is a paid mutator transaction binding the contract method 0xd63dd196.
//
// Solidity: function transferTokens(address _to, uint256 _a, uint256 _b, address state_tokenContract) returns()
func (_Types *TypesTransactorSession) TransferTokens(_to common.Address, _a *big.Int, _b *big.Int, state_tokenContract common.Address) (*types.Transaction, error) {
	return _Types.Contract.TransferTokens(&_Types.TransactOpts, _to, _a, _b, state_tokenContract)
}

// TypesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Types contract.
type TypesOwnershipTransferredIterator struct {
	Event *TypesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TypesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TypesOwnershipTransferred)
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
		it.Event = new(TypesOwnershipTransferred)
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
func (it *TypesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TypesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TypesOwnershipTransferred represents a OwnershipTransferred event raised by the Types contract.
type TypesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Types *TypesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TypesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Types.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TypesOwnershipTransferredIterator{contract: _Types.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Types *TypesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TypesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Types.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TypesOwnershipTransferred)
				if err := _Types.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Types *TypesFilterer) ParseOwnershipTransferred(log types.Log) (*TypesOwnershipTransferred, error) {
	event := new(TypesOwnershipTransferred)
	if err := _Types.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
