// Code generated by MockGen. DO NOT EDIT.
// Source: x/axelarcork/types/expected_keepers.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	math "cosmossdk.io/math"
	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types2 "github.com/cosmos/cosmos-sdk/x/capability/types"
	types3 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	types4 "github.com/cosmos/cosmos-sdk/x/staking/types"
	types5 "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	types6 "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	types7 "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	gomock "github.com/golang/mock/gomock"
	types8 "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) types0.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(types0.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// GetModuleAccount mocks base method.
func (m *MockAccountKeeper) GetModuleAccount(ctx types.Context, name string) types0.ModuleAccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAccount", ctx, name)
	ret0, _ := ret[0].(types0.ModuleAccountI)
	return ret0
}

// GetModuleAccount indicates an expected call of GetModuleAccount.
func (mr *MockAccountKeeperMockRecorder) GetModuleAccount(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAccount), ctx, name)
}

// GetModuleAddress mocks base method.
func (m *MockAccountKeeper) GetModuleAddress(name string) types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAddress", name)
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetModuleAddress indicates an expected call of GetModuleAddress.
func (mr *MockAccountKeeperMockRecorder) GetModuleAddress(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAddress", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAddress), name)
}

// SetModuleAccount mocks base method.
func (m *MockAccountKeeper) SetModuleAccount(arg0 types.Context, arg1 types0.ModuleAccountI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetModuleAccount", arg0, arg1)
}

// SetModuleAccount indicates an expected call of SetModuleAccount.
func (mr *MockAccountKeeperMockRecorder) SetModuleAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).SetModuleAccount), arg0, arg1)
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

// GetDenomMetaData mocks base method.
func (m *MockBankKeeper) GetDenomMetaData(ctx types.Context, denom string) (types1.Metadata, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDenomMetaData", ctx, denom)
	ret0, _ := ret[0].(types1.Metadata)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDenomMetaData indicates an expected call of GetDenomMetaData.
func (mr *MockBankKeeperMockRecorder) GetDenomMetaData(ctx, denom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDenomMetaData", reflect.TypeOf((*MockBankKeeper)(nil).GetDenomMetaData), ctx, denom)
}

// GetSupply mocks base method.
func (m *MockBankKeeper) GetSupply(ctx types.Context, denom string) types.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupply", ctx, denom)
	ret0, _ := ret[0].(types.Coin)
	return ret0
}

// GetSupply indicates an expected call of GetSupply.
func (mr *MockBankKeeperMockRecorder) GetSupply(ctx, denom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupply", reflect.TypeOf((*MockBankKeeper)(nil).GetSupply), ctx, denom)
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

// MockICS4Wrapper is a mock of ICS4Wrapper interface.
type MockICS4Wrapper struct {
	ctrl     *gomock.Controller
	recorder *MockICS4WrapperMockRecorder
}

// MockICS4WrapperMockRecorder is the mock recorder for MockICS4Wrapper.
type MockICS4WrapperMockRecorder struct {
	mock *MockICS4Wrapper
}

// NewMockICS4Wrapper creates a new mock instance.
func NewMockICS4Wrapper(ctrl *gomock.Controller) *MockICS4Wrapper {
	mock := &MockICS4Wrapper{ctrl: ctrl}
	mock.recorder = &MockICS4WrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICS4Wrapper) EXPECT() *MockICS4WrapperMockRecorder {
	return m.recorder
}

// GetAppVersion mocks base method.
func (m *MockICS4Wrapper) GetAppVersion(ctx types.Context, portID, channelID string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppVersion", ctx, portID, channelID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetAppVersion indicates an expected call of GetAppVersion.
func (mr *MockICS4WrapperMockRecorder) GetAppVersion(ctx, portID, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppVersion", reflect.TypeOf((*MockICS4Wrapper)(nil).GetAppVersion), ctx, portID, channelID)
}

// SendPacket mocks base method.
func (m *MockICS4Wrapper) SendPacket(ctx types.Context, chanCap *types2.Capability, sourcePort, sourceChannel string, timeoutHeight types6.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPacket", ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendPacket indicates an expected call of SendPacket.
func (mr *MockICS4WrapperMockRecorder) SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPacket", reflect.TypeOf((*MockICS4Wrapper)(nil).SendPacket), ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// WriteAcknowledgement mocks base method.
func (m *MockICS4Wrapper) WriteAcknowledgement(ctx types.Context, chanCap *types2.Capability, packet exported.PacketI, acknowledgement exported.Acknowledgement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAcknowledgement", ctx, chanCap, packet, acknowledgement)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteAcknowledgement indicates an expected call of WriteAcknowledgement.
func (mr *MockICS4WrapperMockRecorder) WriteAcknowledgement(ctx, chanCap, packet, acknowledgement interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAcknowledgement", reflect.TypeOf((*MockICS4Wrapper)(nil).WriteAcknowledgement), ctx, chanCap, packet, acknowledgement)
}

// MockChannelKeeper is a mock of ChannelKeeper interface.
type MockChannelKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockChannelKeeperMockRecorder
}

// MockChannelKeeperMockRecorder is the mock recorder for MockChannelKeeper.
type MockChannelKeeperMockRecorder struct {
	mock *MockChannelKeeper
}

// NewMockChannelKeeper creates a new mock instance.
func NewMockChannelKeeper(ctrl *gomock.Controller) *MockChannelKeeper {
	mock := &MockChannelKeeper{ctrl: ctrl}
	mock.recorder = &MockChannelKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChannelKeeper) EXPECT() *MockChannelKeeperMockRecorder {
	return m.recorder
}

// GetChannel mocks base method.
func (m *MockChannelKeeper) GetChannel(ctx types.Context, portID, channelID string) (types7.Channel, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannel", ctx, portID, channelID)
	ret0, _ := ret[0].(types7.Channel)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetChannel indicates an expected call of GetChannel.
func (mr *MockChannelKeeperMockRecorder) GetChannel(ctx, portID, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannel", reflect.TypeOf((*MockChannelKeeper)(nil).GetChannel), ctx, portID, channelID)
}

// GetChannelClientState mocks base method.
func (m *MockChannelKeeper) GetChannelClientState(ctx types.Context, portID, channelID string) (string, exported.ClientState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannelClientState", ctx, portID, channelID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(exported.ClientState)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetChannelClientState indicates an expected call of GetChannelClientState.
func (mr *MockChannelKeeperMockRecorder) GetChannelClientState(ctx, portID, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelClientState", reflect.TypeOf((*MockChannelKeeper)(nil).GetChannelClientState), ctx, portID, channelID)
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

// GetBondedValidatorsByPower mocks base method.
func (m *MockStakingKeeper) GetBondedValidatorsByPower(ctx types.Context) []types4.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBondedValidatorsByPower", ctx)
	ret0, _ := ret[0].([]types4.Validator)
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
func (m *MockStakingKeeper) IterateBondedValidatorsByPower(arg0 types.Context, arg1 func(int64, types4.ValidatorI) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateBondedValidatorsByPower", arg0, arg1)
}

// IterateBondedValidatorsByPower indicates an expected call of IterateBondedValidatorsByPower.
func (mr *MockStakingKeeperMockRecorder) IterateBondedValidatorsByPower(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).IterateBondedValidatorsByPower), arg0, arg1)
}

// IterateLastValidators mocks base method.
func (m *MockStakingKeeper) IterateLastValidators(arg0 types.Context, arg1 func(int64, types4.ValidatorI) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateLastValidators", arg0, arg1)
}

// IterateLastValidators indicates an expected call of IterateLastValidators.
func (mr *MockStakingKeeperMockRecorder) IterateLastValidators(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateLastValidators", reflect.TypeOf((*MockStakingKeeper)(nil).IterateLastValidators), arg0, arg1)
}

// IterateValidators mocks base method.
func (m *MockStakingKeeper) IterateValidators(arg0 types.Context, arg1 func(int64, types4.ValidatorI) bool) {
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
func (m *MockStakingKeeper) Validator(arg0 types.Context, arg1 types.ValAddress) types4.ValidatorI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator", arg0, arg1)
	ret0, _ := ret[0].(types4.ValidatorI)
	return ret0
}

// Validator indicates an expected call of Validator.
func (mr *MockStakingKeeperMockRecorder) Validator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockStakingKeeper)(nil).Validator), arg0, arg1)
}

// ValidatorByConsAddr mocks base method.
func (m *MockStakingKeeper) ValidatorByConsAddr(arg0 types.Context, arg1 types.ConsAddress) types4.ValidatorI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorByConsAddr", arg0, arg1)
	ret0, _ := ret[0].(types4.ValidatorI)
	return ret0
}

// ValidatorByConsAddr indicates an expected call of ValidatorByConsAddr.
func (mr *MockStakingKeeperMockRecorder) ValidatorByConsAddr(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorByConsAddr", reflect.TypeOf((*MockStakingKeeper)(nil).ValidatorByConsAddr), arg0, arg1)
}

// MockTransferKeeper is a mock of TransferKeeper interface.
type MockTransferKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockTransferKeeperMockRecorder
}

// MockTransferKeeperMockRecorder is the mock recorder for MockTransferKeeper.
type MockTransferKeeperMockRecorder struct {
	mock *MockTransferKeeper
}

// NewMockTransferKeeper creates a new mock instance.
func NewMockTransferKeeper(ctrl *gomock.Controller) *MockTransferKeeper {
	mock := &MockTransferKeeper{ctrl: ctrl}
	mock.recorder = &MockTransferKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferKeeper) EXPECT() *MockTransferKeeperMockRecorder {
	return m.recorder
}

// Transfer mocks base method.
func (m *MockTransferKeeper) Transfer(goCtx context.Context, msg *types5.MsgTransfer) (*types5.MsgTransferResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", goCtx, msg)
	ret0, _ := ret[0].(*types5.MsgTransferResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transfer indicates an expected call of Transfer.
func (mr *MockTransferKeeperMockRecorder) Transfer(goCtx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockTransferKeeper)(nil).Transfer), goCtx, msg)
}

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
func (m *MockDistributionKeeper) GetFeePool(ctx types.Context) types3.FeePool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeePool", ctx)
	ret0, _ := ret[0].(types3.FeePool)
	return ret0
}

// GetFeePool indicates an expected call of GetFeePool.
func (mr *MockDistributionKeeperMockRecorder) GetFeePool(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeePool", reflect.TypeOf((*MockDistributionKeeper)(nil).GetFeePool), ctx)
}

// SetFeePool mocks base method.
func (m *MockDistributionKeeper) SetFeePool(ctx types.Context, feePool types3.FeePool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFeePool", ctx, feePool)
}

// SetFeePool indicates an expected call of SetFeePool.
func (mr *MockDistributionKeeperMockRecorder) SetFeePool(ctx, feePool interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFeePool", reflect.TypeOf((*MockDistributionKeeper)(nil).SetFeePool), ctx, feePool)
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
func (m *MockPubsubKeeper) GetPublisher(ctx types.Context, publisherDomain string) (types8.Publisher, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublisher", ctx, publisherDomain)
	ret0, _ := ret[0].(types8.Publisher)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPublisher indicates an expected call of GetPublisher.
func (mr *MockPubsubKeeperMockRecorder) GetPublisher(ctx, publisherDomain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublisher", reflect.TypeOf((*MockPubsubKeeper)(nil).GetPublisher), ctx, publisherDomain)
}

// SetDefaultSubscription mocks base method.
func (m *MockPubsubKeeper) SetDefaultSubscription(ctx types.Context, defaultSubscription types8.DefaultSubscription) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefaultSubscription", ctx, defaultSubscription)
}

// SetDefaultSubscription indicates an expected call of SetDefaultSubscription.
func (mr *MockPubsubKeeperMockRecorder) SetDefaultSubscription(ctx, defaultSubscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultSubscription", reflect.TypeOf((*MockPubsubKeeper)(nil).SetDefaultSubscription), ctx, defaultSubscription)
}
