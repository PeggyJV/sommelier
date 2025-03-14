// Code generated by MockGen. DO NOT EDIT.
// Source: x/cork/types/expected_keepers.go

// Package mock_types is a generated GoMock package.
package mock_types

import (
	reflect "reflect"

	math "cosmossdk.io/math"
	bytes "github.com/cometbft/cometbft/libs/bytes"
	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/staking/types"
	common "github.com/ethereum/go-ethereum/common"
	gomock "github.com/golang/mock/gomock"
	types1 "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	types2 "github.com/peggyjv/sommelier/v9/x/pubsub/types"
)

// MockStakingKeeper is a mock of StakingKeeper interface.
type MockStakingKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockStakingKeeperMockRecorder
}

// MockStakingKeeperMockRecorder is the mock recorder for MockStakingKeeper.
type MockStakingKeeperMockRecorder struct {
	mock *MockStakingKeeper
}

// NewMockStakingKeeper creates a new mock instance.
func NewMockStakingKeeper(ctrl *gomock.Controller) *MockStakingKeeper {
	mock := &MockStakingKeeper{ctrl: ctrl}
	mock.recorder = &MockStakingKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStakingKeeper) EXPECT() *MockStakingKeeperMockRecorder {
	return m.recorder
}

// GetBondedValidatorsByPower mocks base method.
func (m *MockStakingKeeper) GetBondedValidatorsByPower(ctx types.Context) []types0.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBondedValidatorsByPower", ctx)
	ret0, _ := ret[0].([]types0.Validator)
	return ret0
}

// GetBondedValidatorsByPower indicates an expected call of GetBondedValidatorsByPower.
func (mr *MockStakingKeeperMockRecorder) GetBondedValidatorsByPower(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).GetBondedValidatorsByPower), ctx)
}

// GetLastTotalPower mocks base method.
func (m *MockStakingKeeper) GetLastTotalPower(ctx types.Context) math.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastTotalPower", ctx)
	ret0, _ := ret[0].(math.Int)
	return ret0
}

// GetLastTotalPower indicates an expected call of GetLastTotalPower.
func (mr *MockStakingKeeperMockRecorder) GetLastTotalPower(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastTotalPower", reflect.TypeOf((*MockStakingKeeper)(nil).GetLastTotalPower), ctx)
}

// GetLastValidatorPower mocks base method.
func (m *MockStakingKeeper) GetLastValidatorPower(ctx types.Context, operator types.ValAddress) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastValidatorPower", ctx, operator)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetLastValidatorPower indicates an expected call of GetLastValidatorPower.
func (mr *MockStakingKeeperMockRecorder) GetLastValidatorPower(ctx, operator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastValidatorPower", reflect.TypeOf((*MockStakingKeeper)(nil).GetLastValidatorPower), ctx, operator)
}

// IterateBondedValidatorsByPower mocks base method.
func (m *MockStakingKeeper) IterateBondedValidatorsByPower(arg0 types.Context, arg1 func(int64, types0.ValidatorI) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateBondedValidatorsByPower", arg0, arg1)
}

// IterateBondedValidatorsByPower indicates an expected call of IterateBondedValidatorsByPower.
func (mr *MockStakingKeeperMockRecorder) IterateBondedValidatorsByPower(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).IterateBondedValidatorsByPower), arg0, arg1)
}

// IterateLastValidators mocks base method.
func (m *MockStakingKeeper) IterateLastValidators(arg0 types.Context, arg1 func(int64, types0.ValidatorI) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateLastValidators", arg0, arg1)
}

// IterateLastValidators indicates an expected call of IterateLastValidators.
func (mr *MockStakingKeeperMockRecorder) IterateLastValidators(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateLastValidators", reflect.TypeOf((*MockStakingKeeper)(nil).IterateLastValidators), arg0, arg1)
}

// IterateValidators mocks base method.
func (m *MockStakingKeeper) IterateValidators(arg0 types.Context, arg1 func(int64, types0.ValidatorI) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateValidators", arg0, arg1)
}

// IterateValidators indicates an expected call of IterateValidators.
func (mr *MockStakingKeeperMockRecorder) IterateValidators(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateValidators", reflect.TypeOf((*MockStakingKeeper)(nil).IterateValidators), arg0, arg1)
}

// Jail mocks base method.
func (m *MockStakingKeeper) Jail(arg0 types.Context, arg1 types.ConsAddress) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Jail", arg0, arg1)
}

// Jail indicates an expected call of Jail.
func (mr *MockStakingKeeperMockRecorder) Jail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Jail", reflect.TypeOf((*MockStakingKeeper)(nil).Jail), arg0, arg1)
}

// PowerReduction mocks base method.
func (m *MockStakingKeeper) PowerReduction(ctx types.Context) math.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PowerReduction", ctx)
	ret0, _ := ret[0].(math.Int)
	return ret0
}

// PowerReduction indicates an expected call of PowerReduction.
func (mr *MockStakingKeeperMockRecorder) PowerReduction(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PowerReduction", reflect.TypeOf((*MockStakingKeeper)(nil).PowerReduction), ctx)
}

// Slash mocks base method.
func (m *MockStakingKeeper) Slash(arg0 types.Context, arg1 types.ConsAddress, arg2, arg3 int64, arg4 types.Dec) math.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Slash", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(math.Int)
	return ret0
}

// Slash indicates an expected call of Slash.
func (mr *MockStakingKeeperMockRecorder) Slash(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Slash", reflect.TypeOf((*MockStakingKeeper)(nil).Slash), arg0, arg1, arg2, arg3, arg4)
}

// Validator mocks base method.
func (m *MockStakingKeeper) Validator(arg0 types.Context, arg1 types.ValAddress) types0.ValidatorI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator", arg0, arg1)
	ret0, _ := ret[0].(types0.ValidatorI)
	return ret0
}

// Validator indicates an expected call of Validator.
func (mr *MockStakingKeeperMockRecorder) Validator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockStakingKeeper)(nil).Validator), arg0, arg1)
}

// ValidatorByConsAddr mocks base method.
func (m *MockStakingKeeper) ValidatorByConsAddr(arg0 types.Context, arg1 types.ConsAddress) types0.ValidatorI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorByConsAddr", arg0, arg1)
	ret0, _ := ret[0].(types0.ValidatorI)
	return ret0
}

// ValidatorByConsAddr indicates an expected call of ValidatorByConsAddr.
func (mr *MockStakingKeeperMockRecorder) ValidatorByConsAddr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorByConsAddr", reflect.TypeOf((*MockStakingKeeper)(nil).ValidatorByConsAddr), arg0, arg1)
}

// MockGravityKeeper is a mock of GravityKeeper interface.
type MockGravityKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockGravityKeeperMockRecorder
}

// MockGravityKeeperMockRecorder is the mock recorder for MockGravityKeeper.
type MockGravityKeeperMockRecorder struct {
	mock *MockGravityKeeper
}

// NewMockGravityKeeper creates a new mock instance.
func NewMockGravityKeeper(ctrl *gomock.Controller) *MockGravityKeeper {
	mock := &MockGravityKeeper{ctrl: ctrl}
	mock.recorder = &MockGravityKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGravityKeeper) EXPECT() *MockGravityKeeperMockRecorder {
	return m.recorder
}

// CreateContractCallTx mocks base method.
func (m *MockGravityKeeper) CreateContractCallTx(ctx types.Context, invalidationNonce uint64, invalidationScope bytes.HexBytes, address common.Address, payload []byte, tokens, fees []types1.ERC20Token) *types1.ContractCallTx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContractCallTx", ctx, invalidationNonce, invalidationScope, address, payload, tokens, fees)
	ret0, _ := ret[0].(*types1.ContractCallTx)
	return ret0
}

// CreateContractCallTx indicates an expected call of CreateContractCallTx.
func (mr *MockGravityKeeperMockRecorder) CreateContractCallTx(ctx, invalidationNonce, invalidationScope, address, payload, tokens, fees interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContractCallTx", reflect.TypeOf((*MockGravityKeeper)(nil).CreateContractCallTx), ctx, invalidationNonce, invalidationScope, address, payload, tokens, fees)
}

// GetEthereumOrchestratorAddress mocks base method.
func (m *MockGravityKeeper) GetEthereumOrchestratorAddress(ctx types.Context, ethAddr common.Address) types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEthereumOrchestratorAddress", ctx, ethAddr)
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetEthereumOrchestratorAddress indicates an expected call of GetEthereumOrchestratorAddress.
func (mr *MockGravityKeeperMockRecorder) GetEthereumOrchestratorAddress(ctx, ethAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEthereumOrchestratorAddress", reflect.TypeOf((*MockGravityKeeper)(nil).GetEthereumOrchestratorAddress), ctx, ethAddr)
}

// GetOrchestratorValidatorAddress mocks base method.
func (m *MockGravityKeeper) GetOrchestratorValidatorAddress(ctx types.Context, orchAddr types.AccAddress) types.ValAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrchestratorValidatorAddress", ctx, orchAddr)
	ret0, _ := ret[0].(types.ValAddress)
	return ret0
}

// GetOrchestratorValidatorAddress indicates an expected call of GetOrchestratorValidatorAddress.
func (mr *MockGravityKeeperMockRecorder) GetOrchestratorValidatorAddress(ctx, orchAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrchestratorValidatorAddress", reflect.TypeOf((*MockGravityKeeper)(nil).GetOrchestratorValidatorAddress), ctx, orchAddr)
}

// GetValidatorEthereumAddress mocks base method.
func (m *MockGravityKeeper) GetValidatorEthereumAddress(ctx types.Context, valAddr types.ValAddress) common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorEthereumAddress", ctx, valAddr)
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetValidatorEthereumAddress indicates an expected call of GetValidatorEthereumAddress.
func (mr *MockGravityKeeperMockRecorder) GetValidatorEthereumAddress(ctx, valAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorEthereumAddress", reflect.TypeOf((*MockGravityKeeper)(nil).GetValidatorEthereumAddress), ctx, valAddr)
}

// SetOrchestratorValidatorAddress mocks base method.
func (m *MockGravityKeeper) SetOrchestratorValidatorAddress(ctx types.Context, val types.ValAddress, orchAddr types.AccAddress) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOrchestratorValidatorAddress", ctx, val, orchAddr)
}

// SetOrchestratorValidatorAddress indicates an expected call of SetOrchestratorValidatorAddress.
func (mr *MockGravityKeeperMockRecorder) SetOrchestratorValidatorAddress(ctx, val, orchAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOrchestratorValidatorAddress", reflect.TypeOf((*MockGravityKeeper)(nil).SetOrchestratorValidatorAddress), ctx, val, orchAddr)
}

// SetOutgoingTx mocks base method.
func (m *MockGravityKeeper) SetOutgoingTx(ctx types.Context, outgoing types1.OutgoingTx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOutgoingTx", ctx, outgoing)
}

// SetOutgoingTx indicates an expected call of SetOutgoingTx.
func (mr *MockGravityKeeperMockRecorder) SetOutgoingTx(ctx, outgoing interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOutgoingTx", reflect.TypeOf((*MockGravityKeeper)(nil).SetOutgoingTx), ctx, outgoing)
}

// MockPubsubKeeper is a mock of PubsubKeeper interface.
type MockPubsubKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockPubsubKeeperMockRecorder
}

// MockPubsubKeeperMockRecorder is the mock recorder for MockPubsubKeeper.
type MockPubsubKeeperMockRecorder struct {
	mock *MockPubsubKeeper
}

// NewMockPubsubKeeper creates a new mock instance.
func NewMockPubsubKeeper(ctrl *gomock.Controller) *MockPubsubKeeper {
	mock := &MockPubsubKeeper{ctrl: ctrl}
	mock.recorder = &MockPubsubKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubsubKeeper) EXPECT() *MockPubsubKeeperMockRecorder {
	return m.recorder
}

// DeleteDefaultSubscription mocks base method.
func (m *MockPubsubKeeper) DeleteDefaultSubscription(ctx types.Context, subscriptionID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteDefaultSubscription", ctx, subscriptionID)
}

// DeleteDefaultSubscription indicates an expected call of DeleteDefaultSubscription.
func (mr *MockPubsubKeeperMockRecorder) DeleteDefaultSubscription(ctx, subscriptionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDefaultSubscription", reflect.TypeOf((*MockPubsubKeeper)(nil).DeleteDefaultSubscription), ctx, subscriptionID)
}

// GetPublisher mocks base method.
func (m *MockPubsubKeeper) GetPublisher(ctx types.Context, publisherDomain string) (types2.Publisher, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublisher", ctx, publisherDomain)
	ret0, _ := ret[0].(types2.Publisher)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPublisher indicates an expected call of GetPublisher.
func (mr *MockPubsubKeeperMockRecorder) GetPublisher(ctx, publisherDomain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublisher", reflect.TypeOf((*MockPubsubKeeper)(nil).GetPublisher), ctx, publisherDomain)
}

// SetDefaultSubscription mocks base method.
func (m *MockPubsubKeeper) SetDefaultSubscription(ctx types.Context, defaultSubscription types2.DefaultSubscription) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefaultSubscription", ctx, defaultSubscription)
}

// SetDefaultSubscription indicates an expected call of SetDefaultSubscription.
func (mr *MockPubsubKeeperMockRecorder) SetDefaultSubscription(ctx, defaultSubscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultSubscription", reflect.TypeOf((*MockPubsubKeeper)(nil).SetDefaultSubscription), ctx, defaultSubscription)
}
