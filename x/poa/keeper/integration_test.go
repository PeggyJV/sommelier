package keeper_test

// Integration tests for x/poa require simapp scaffolding with at least one
// genesis-bonded validator (otherwise x/staking's InitGenesis panics with
// "validator set is empty after InitGenesis"). That setup — key generation,
// gentx, delegation, and a tailored genesis — is substantial and is tracked
// as follow-up work.
//
// What is covered today:
//   - x/poa/keeper unit tests: ComputeMultiplier, authority storage,
//     snapshot store + pruning, WrappedStakingKeeper rescaling and Slash
//     normalisation, EndBlocker boost / no-op / halt-on-empty / disabled.
//   - app/app_test.go::TestSommelierAppExport: exercises NewSommelierApp +
//     InitChain through the full module wiring including the PoA keeper,
//     proving the surgery in app/app.go assembles and initialises without
//     panic.
//
// Deferred integration scenarios (follow-up):
//   - Full InitChain with bonded authority + community validators; multi-
//     block run asserting authority share >= floor every block.
//   - Gravity-bridge SignerSetTx weight check (authority signers >= 67%).
//   - Downtime / double-sign slash on an authority validator producing
//     `slash = fraction * actual_bonded` (not boosted).
//   - Governance MsgUpdateAuthoritySet end-to-end through the gov module.
