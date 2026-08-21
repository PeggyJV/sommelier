#!/usr/bin/env python3
"""
Turn a mainnet state export into a single-validator rehearsal genesis.

Used by the upgrade rehearsal described in docs/upgrade-rehearsal.md. The point
is to keep ALL mainnet state -- real corks, real params, real accounts, real
validator set -- and change only what is needed for one node whose consensus key
we hold to produce blocks on its own.

What it changes, and why each change is required:

  * Our validator's stake is inflated to `--share` of bonded power so it alone
    exceeds the 2/3 needed to commit a block. Without this the chain cannot make
    progress, because we do not hold the other validators' keys.

  * bonded_tokens_pool and total supply both move by the same delta. x/staking
    and x/bank InitGenesis cross-check module account balances against the sum
    of bonded tokens and panic on a mismatch, so inflating stake alone will not
    boot.

  * delegator_shares and the matching self-delegation move with it, so the
    validator's shares stay consistent with its tokens.

  * last_validator_powers / last_total_power are updated, since they seed the
    initial CometBFT validator set.

  * Governance voting period and min deposit are shrunk so the upgrade proposal
    can be proposed, voted, and passed inside the rehearsal rather than in two
    days.

  * `validators` is emptied: for an exported chain the CometBFT set comes from
    x/staking InitGenesis, and a stale list here only conflicts with it.

Everything else is left byte-for-byte as mainnet had it. In particular the
module state under test is untouched -- that is the whole point.

Usage:

    rehearsal_genesis_surgery.py IN.json OUT.json \\
        --valoper sommvaloper1... \\
        --delegator somm1... \\
        [--chain-id sommelier-rehearsal] [--share 4] [--voting-period 180s]

--share is the multiple of the OTHER bonded validators' stake to give ours;
4 yields 80% of bonded power, comfortably above 2/3 with room for rounding.

Find the module account addresses for a given chain with:

    jq -r '.app_state.auth.accounts[]
           | select(."@type"=="/cosmos.auth.v1beta1.ModuleAccount")
           | "\\(.name) \\(.base_account.address)"' export.json
"""
import argparse
import json
import sys

POWER_REDUCTION = 1_000_000


def find_module_account(genesis: dict, name: str) -> str:
    for acc in genesis["app_state"]["auth"]["accounts"]:
        if acc.get("@type") == "/cosmos.auth.v1beta1.ModuleAccount" and acc.get("name") == name:
            return acc["base_account"]["address"]
    raise SystemExit(f"module account {name!r} not found in export")


def bump_usomm(coins: list, delta: int, what: str) -> None:
    for c in coins:
        if c["denom"] == "usomm":
            c["amount"] = str(int(c["amount"]) + delta)
            return
    raise SystemExit(f"no usomm entry in {what}; cannot keep the books balanced")


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("src")
    ap.add_argument("dst")
    ap.add_argument("--valoper", required=True, help="operator address of the validator we hold the key for")
    ap.add_argument("--delegator", required=True, help="that validator's self-delegation account")
    ap.add_argument("--chain-id", default="sommelier-rehearsal")
    ap.add_argument("--share", type=int, default=4,
                    help="multiple of other bonded stake to assign ours (4 => 80%% of bonded power)")
    ap.add_argument("--voting-period", default="180s")
    ap.add_argument("--min-deposit", default="1000000", help="usomm")
    args = ap.parse_args()

    g = json.load(open(args.src))
    g["chain_id"] = args.chain_id

    st = g["app_state"]["staking"]
    vals = st["validators"]

    try:
        ours = next(v for v in vals if v["operator_address"] == args.valoper)
    except StopIteration:
        raise SystemExit(f"validator {args.valoper} not present in the export")

    old_tokens = int(ours["tokens"])
    others = sum(int(v["tokens"]) for v in vals
                 if v["status"] == "BOND_STATUS_BONDED" and v["operator_address"] != args.valoper)
    if others == 0:
        raise SystemExit("no other bonded validators; nothing to dominate")

    new_tokens = others * args.share
    delta = new_tokens - old_tokens
    if delta <= 0:
        raise SystemExit("validator already dominant; refusing to shrink stake")

    ours["tokens"] = str(new_tokens)
    ours["delegator_shares"] = f"{new_tokens}.000000000000000000"
    ours["jailed"] = False
    ours["status"] = "BOND_STATUS_BONDED"

    # Shares must track tokens, or the validator's delegation is inconsistent.
    try:
        d = next(x for x in st["delegations"]
                 if x["validator_address"] == args.valoper and x["delegator_address"] == args.delegator)
    except StopIteration:
        raise SystemExit(f"no self-delegation from {args.delegator} to {args.valoper}")
    d["shares"] = f"{new_tokens}.000000000000000000"

    # These seed the initial CometBFT validator set.
    new_power = new_tokens // POWER_REDUCTION
    old_power = 0
    for e in st["last_validator_powers"]:
        if e["address"] == args.valoper:
            old_power = int(e["power"])
            e["power"] = str(new_power)
            break
    else:
        st["last_validator_powers"].append({"address": args.valoper, "power": str(new_power)})
    st["last_total_power"] = str(int(st["last_total_power"]) - old_power + new_power)

    # Keep the books balanced or InitGenesis panics.
    bonded_pool = find_module_account(g, "bonded_tokens_pool")
    bank = g["app_state"]["bank"]
    pool = next((b for b in bank["balances"] if b["address"] == bonded_pool), None)
    if pool is None:
        raise SystemExit(f"bonded pool {bonded_pool} has no bank balance entry")
    bump_usomm(pool["coins"], delta, "bonded_tokens_pool")
    bump_usomm(bank["supply"], delta, "bank supply")

    # Let the upgrade proposal complete inside the rehearsal.
    gp = g["app_state"]["gov"]["params"]
    gp["voting_period"] = args.voting_period
    gp["max_deposit_period"] = "600s"
    gp["min_deposit"] = [{"denom": "usomm", "amount": args.min_deposit}]
    if "expedited_voting_period" in gp:
        gp["expedited_voting_period"] = "15s"

    # The comet set comes from staking InitGenesis; a stale list only conflicts.
    g["validators"] = []

    json.dump(g, open(args.dst, "w"))

    pct = 100 * new_tokens // (new_tokens + others)
    print(f"chain-id       : {args.chain_id}", file=sys.stderr)
    print(f"our old tokens : {old_tokens}", file=sys.stderr)
    print(f"other bonded   : {others}", file=sys.stderr)
    print(f"our new tokens : {new_tokens}  ({pct}% of bonded power)", file=sys.stderr)
    print(f"delta applied  : {delta} (to bonded pool and supply)", file=sys.stderr)
    print(f"new power      : {new_power}  last_total_power={st['last_total_power']}", file=sys.stderr)
    if pct < 67:
        print("WARNING: below 2/3; the chain will not commit blocks alone", file=sys.stderr)


if __name__ == "__main__":
    main()
