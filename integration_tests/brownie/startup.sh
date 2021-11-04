echo "starting ganache"
ganache-cli --fork https://mainnet.infura.io/v3/d6f22be0f7fd447186086d2495779003 -p 8545 &

echo "running brownie script"
pipenv run brownie run ./scripts/testnet.py --network=development

echo "brownie complete"