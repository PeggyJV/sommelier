FROM node:14-alpine3.13

RUN apk update
RUN apk add --no-cache git

COPY package.json package.json
COPY yarn.lock yarn.lock
RUN yarn install

COPY . .

EXPOSE 8545

CMD npx hardhat integration_test_setup