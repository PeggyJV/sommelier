FROM node:14-alpine3.13

RUN apk update
RUN apk add --no-cache git

COPY package.json package.json
COPY yarn.lock yarn.lock
RUN yarn

COPY . .

CMD npx hardhat integration_test_setup