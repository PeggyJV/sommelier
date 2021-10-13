# Creating a new contract

## Install cosmwasm rust tooling
```bash
cargo install cargo-generate --features vendored-openssl
cargo install cargo-run-script
```

## Create new contract from template
```bash
cd contracts/
cargo generate --git https://github.com/CosmWasm/cw-template.git --name PROJECT_NAME
```