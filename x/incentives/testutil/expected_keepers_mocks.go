// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/atro/source/work/sommelier/x/incentives/types/expected_keepers.go

// Package mock_types is a generated GoMock package.
package mock_types

import (
	reflect "reflect"

	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	types1 "github.com/cosmos/cosmos-sdk/x/mint/types"
	types2 "github.com/cosmos/cosmos-sdk/x/staking/types"
	gomock "github.com/golang/mock/gomock"
)

// MockDistributionKeeper is a mock of DistributionKeeper interface.
type MockDistributionKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockDistributionKeeperMockRecorder
}

// MockDistributionKeeperMockRecorder is the mock recorder for MockDistributionKeeper.
type MockDistributionKeeperMockRecorder struct {
	mock *MockDistributionKeeper
}

// NewMockDistributionKeeper creates a new mock instance.
func NewMockDistributionKeeper(ctrl *gomock.Controller) *MockDistributionKeeper {
	mock := &MockDistributionKeeper{ctrl: ctrl}
	mock.recorder = &MockDistributionKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDistributionKeeper) EXPECT() *MockDistributionKeeperMockRecorder {
	return m.recorder
}

// GetFeePool mocks base method.
func (m *MockDistributionKeeper) GetFeePool(ctx types.Context) types0.FeePool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeePool", ctx)
	ret0, _ := ret[0].(types0.FeePool)
	return ret0
}

// GetFeePool indicates an expected call of GetFeePool.
func (mr *MockDistributionKeeperMockRecorder) GetFeePool(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeePool", reflect.TypeOf((*MockDistributionKeeper)(nil).GetFeePool), ctx)
}

// GetValidatorAccumulatedCommission mocks base method.
func (m *MockDistributionKeeper) GetValidatorAccumulatedCommission(ctx types.Context, valAddr types.ValAddress) types0.ValidatorAccumulatedCommission {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorAccumulatedCommission", ctx, valAddr)
	ret0, _ := ret[0].(types0.ValidatorAccumulatedCommission)
	return ret0
}

// GetValidatorAccumulatedCommission indicates an expected call of GetValidatorAccumulatedCommission.
func (mr *MockDistributionKeeperMockRecorder) GetValidatorAccumulatedCommission(ctx, valAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorAccumulatedCommission", reflect.TypeOf((*MockDistributionKeeper)(nil).GetValidatorAccumulatedCommission), ctx, valAddr)
}

// GetValidatorCurrentRewards mocks base method.
func (m *MockDistributionKeeper) GetValidatorCurrentRewards(ctx types.Context, valAddr types.ValAddress) types0.ValidatorCurrentRewards {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorCurrentRewards", ctx, valAddr)
	ret0, _ := ret[0].(types0.ValidatorCurrentRewards)
	return ret0
}

// GetValidatorCurrentRewards indicates an expected call of GetValidatorCurrentRewards.
func (mr *MockDistributionKeeperMockRecorder) GetValidatorCurrentRewards(ctx, valAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorCurrentRewards", reflect.TypeOf((*MockDistributionKeeper)(nil).GetValidatorCurrentRewards), ctx, valAddr)
}

// GetValidatorOutstandingRewards mocks base method.
func (m *MockDistributionKeeper) GetValidatorOutstandingRewards(ctx types.Context, valAddr types.ValAddress) types0.ValidatorOutstandingRewards {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorOutstandingRewards", ctx, valAddr)
	ret0, _ := ret[0].(types0.ValidatorOutstandingRewards)
	return ret0
}

// GetValidatorOutstandingRewards indicates an expected call of GetValidatorOutstandingRewards.
func (mr *MockDistributionKeeperMockRecorder) GetValidatorOutstandingRewards(ctx, valAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorOutstandingRewards", reflect.TypeOf((*MockDistributionKeeper)(nil).GetValidatorOutstandingRewards), ctx, valAddr)
}

// SetFeePool mocks base method.
func (m *MockDistributionKeeper) SetFeePool(ctx types.Context, feePool types0.FeePool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFeePool", ctx, feePool)
}

// SetFeePool indicates an expected call of SetFeePool.
func (mr *MockDistributionKeeperMockRecorder) SetFeePool(ctx, feePool interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFeePool", reflect.TypeOf((*MockDistributionKeeper)(nil).SetFeePool), ctx, feePool)
}

// SetValidatorAccumulatedCommission mocks base method.
func (m *MockDistributionKeeper) SetValidatorAccumulatedCommission(ctx types.Context, valAddr types.ValAddress, commission types0.ValidatorAccumulatedCommission) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValidatorAccumulatedCommission", ctx, valAddr, commission)
}

// SetValidatorAccumulatedCommission indicates an expected call of SetValidatorAccumulatedCommission.
func (mr *MockDistributionKeeperMockRecorder) SetValidatorAccumulatedCommission(ctx, valAddr, commission interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValidatorAccumulatedCommission", reflect.TypeOf((*MockDistributionKeeper)(nil).SetValidatorAccumulatedCommission), ctx, valAddr, commission)
}

// SetValidatorCurrentRewards mocks base method.
func (m *MockDistributionKeeper) SetValidatorCurrentRewards(ctx types.Context, valAddr types.ValAddress, rewards types0.ValidatorCurrentRewards) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValidatorCurrentRewards", ctx, valAddr, rewards)
}

// SetValidatorCurrentRewards indicates an expected call of SetValidatorCurrentRewards.
func (mr *MockDistributionKeeperMockRecorder) SetValidatorCurrentRewards(ctx, valAddr, rewards interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValidatorCurrentRewards", reflect.TypeOf((*MockDistributionKeeper)(nil).SetValidatorCurrentRewards), ctx, valAddr, rewards)
}

// SetValidatorOutstandingRewards mocks base method.
func (m *MockDistributionKeeper) SetValidatorOutstandingRewards(ctx types.Context, valAddr types.ValAddress, rewards types0.ValidatorOutstandingRewards) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValidatorOutstandingRewards", ctx, valAddr, rewards)
}

// SetValidatorOutstandingRewards indicates an expected call of SetValidatorOutstandingRewards.
func (mr *MockDistributionKeeperMockRecorder) SetValidatorOutstandingRewards(ctx, valAddr, rewards interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValidatorOutstandingRewards", reflect.TypeOf((*MockDistributionKeeper)(nil).SetValidatorOutstandingRewards), ctx, valAddr, rewards)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// BlockedAddr mocks base method.
func (m *MockBankKeeper) BlockedAddr(addr types.AccAddress) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockedAddr", addr)
	ret0, _ := ret[0].(bool)
	return ret0
}

// BlockedAddr indicates an expected call of BlockedAddr.
func (mr *MockBankKeeperMockRecorder) BlockedAddr(addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockedAddr", reflect.TypeOf((*MockBankKeeper)(nil).BlockedAddr), addr)
}

// GetAllBalances mocks base method.
func (m *MockBankKeeper) GetAllBalances(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalances", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// GetAllBalances indicates an expected call of GetAllBalances.
func (mr *MockBankKeeperMockRecorder) GetAllBalances(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAllBalances), ctx, addr)
}

// GetBalance mocks base method.
func (m *MockBankKeeper) GetBalance(ctx types.Context, addr types.AccAddress, denom string) types.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, addr, denom)
	ret0, _ := ret[0].(types.Coin)
	return ret0
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockBankKeeperMockRecorder) GetBalance(ctx, addr, denom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBankKeeper)(nil).GetBalance), ctx, addr, denom)
}

// LockedCoins mocks base method.
func (m *MockBankKeeper) LockedCoins(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockedCoins", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// LockedCoins indicates an expected call of LockedCoins.
func (mr *MockBankKeeperMockRecorder) LockedCoins(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockedCoins", reflect.TypeOf((*MockBankKeeper)(nil).LockedCoins), ctx, addr)
}

// SendCoinsFromAccountToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx types.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromAccountToModule indicates an expected call of SendCoinsFromAccountToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromAccountToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromAccountToModule), ctx, senderAddr, recipientModule, amt)
}

// SendCoinsFromModuleToAccount mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToAccount", ctx, senderModule, recipientAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToAccount indicates an expected call of SendCoinsFromModuleToAccount.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToAccount", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToAccount), ctx, senderModule, recipientAddr, amt)
}

// SendCoinsFromModuleToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToModule(ctx types.Context, senderModule, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToModule", ctx, senderModule, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToModule indicates an expected call of SendCoinsFromModuleToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToModule), ctx, senderModule, recipientModule, amt)
}

// SpendableCoins mocks base method.
func (m *MockBankKeeper) SpendableCoins(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpendableCoins", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// SpendableCoins indicates an expected call of SpendableCoins.
func (mr *MockBankKeeperMockRecorder) SpendableCoins(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpendableCoins", reflect.TypeOf((*MockBankKeeper)(nil).SpendableCoins), ctx, addr)
}

// MockMintKeeper is a mock of MintKeeper interface.
type MockMintKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockMintKeeperMockRecorder
}

// MockMintKeeperMockRecorder is the mock recorder for MockMintKeeper.
type MockMintKeeperMockRecorder struct {
	mock *MockMintKeeper
}

// NewMockMintKeeper creates a new mock instance.
func NewMockMintKeeper(ctrl *gomock.Controller) *MockMintKeeper {
	mock := &MockMintKeeper{ctrl: ctrl}
	mock.recorder = &MockMintKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMintKeeper) EXPECT() *MockMintKeeperMockRecorder {
	return m.recorder
}

// BondedRatio mocks base method.
func (m *MockMintKeeper) BondedRatio(ctx types.Context) types.Dec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BondedRatio", ctx)
	ret0, _ := ret[0].(types.Dec)
	return ret0
}

// BondedRatio indicates an expected call of BondedRatio.
func (mr *MockMintKeeperMockRecorder) BondedRatio(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BondedRatio", reflect.TypeOf((*MockMintKeeper)(nil).BondedRatio), ctx)
}

// GetParams mocks base method.
func (m *MockMintKeeper) GetParams(ctx types.Context) types1.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams", ctx)
	ret0, _ := ret[0].(types1.Params)
	return ret0
}

// GetParams indicates an expected call of GetParams.
func (mr *MockMintKeeperMockRecorder) GetParams(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParams", reflect.TypeOf((*MockMintKeeper)(nil).GetParams), ctx)
}

// StakingTokenSupply mocks base method.
func (m *MockMintKeeper) StakingTokenSupply(ctx types.Context) math.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StakingTokenSupply", ctx)
	ret0, _ := ret[0].(math.Int)
	return ret0
}

// StakingTokenSupply indicates an expected call of StakingTokenSupply.
func (mr *MockMintKeeperMockRecorder) StakingTokenSupply(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StakingTokenSupply", reflect.TypeOf((*MockMintKeeper)(nil).StakingTokenSupply), ctx)
}

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

// ValidatorByConsAddr mocks base method.
func (m *MockStakingKeeper) ValidatorByConsAddr(ctx types.Context, consAddr types.ConsAddress) types2.ValidatorI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorByConsAddr", ctx, consAddr)
	ret0, _ := ret[0].(types2.ValidatorI)
	return ret0
}

// ValidatorByConsAddr indicates an expected call of ValidatorByConsAddr.
func (mr *MockStakingKeeperMockRecorder) ValidatorByConsAddr(ctx, consAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorByConsAddr", reflect.TypeOf((*MockStakingKeeper)(nil).ValidatorByConsAddr), ctx, consAddr)
}
