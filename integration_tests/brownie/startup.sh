echo "running brownie script"
ls
pwd
pipenv run brownie run /code/scripts/testnet.py --network mainnet-fork
