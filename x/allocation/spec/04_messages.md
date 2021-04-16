<!--
order: 4
-->

# Messages

## MsgDelegateFeedConsent

Validators may also elect to delegate voting rights to another key as to not require the validator operator key to be kept online. To do so, validators must submit a `MsgDelegateFeedConsent`, delegating their oracle voting rights to a `Delegate` that can sign `MsgExchangeRatePrevote` and `MsgExchangeRateVote` on behalf of the validator.

The `Validator` field contains the operator address of the validator (prefixed `cosmos1...`). The `Delegate` field is the account address (prefixed `cosmos1...`) of the delegate account that will be submitting exchange rate related votes and prevotes on behalf of the `Operator`.

```proto
// MsgDelegateFeedConsent - sdk.Msg for delegating oracle voting rights from a validator
// to another address, must be signed by an active validator
message MsgDelegateFeedConsent {
    string delegate  = 1;
    string validator = 2;
}
```

## MsgOracleDataPrevote

`hashes` is an array of hex byte arrays of the SHA256 hash of a string of the format `{salt}:{data_cannonical_json}:{signer}`. This is a commitment to the actual `MsgOracleDataVote` a validator will submit in the subsequent `VotePeriod`. You can use the `oracletypes.DataHash(salt string, jsn string, signer sdk.AccAddress)` function to help encode it. Note that since in the subsequent `MsgAggregateExchangeRateVote`, the salts for each data type will have to be revealed. Salt used must be regenerated for each prevote submission.

```proto
// MsgOracleDataPrevote - sdk.Msg for prevoting on an array of oracle data types.
// The purpose of the prevote is to hide vote for data with hashes formatted as hex string: 
// SHA256("{salt}:{data_cannonical_json}:{voter}")
message MsgOracleDataPrevote {
    repeated bytes  hashes = 1;
    string          signer = 2;
}
```

## MsgOracleDataVote

The `MsgOracleDataVote` contains the an array of the different `OracleData` for a validators vote. The `salt` parameter must match the salt used to create the prevote, otherwise the voter cannot be rewarded. Note: The indexes of the salt and oracle_data must be the same as the indexes of the hashes that were submitted in the `MsgOraclePrevote`.

```proto
// MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
message MsgOracleDataVote {
    repeated string              salt       = 1;
    repeated google.protobuf.Any oracle_data = 2 [
        (cosmos_proto.accepts_interface) = "OracleData"
    ];
    string                       signer     = 3;
}
```
