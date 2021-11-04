FROM python:3.8

RUN mkdir -p /go/src/github.com/peggyjv/sommelier
WORKDIR /code

RUN apt-get update && apt-get install -y libssl-dev
RUN apt-get update && apt-get install -y npm
RUN npm install -g ganache-cli

COPY requirements.txt .
RUN pip install -r requirements.txt
RUN pip install eth-brownie
#Hangs
#RUN brownie run testnet --network mainnet-fork
