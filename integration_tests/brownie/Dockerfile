FROM python:3.8

# Set up code directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install linux dependencies
RUN apt-get update
RUN apt-get install -y libssl-dev npm bash

RUN npm install -g ganache-cli

COPY Pipfile .
COPY Pipfile.lock .
RUN pipenv sync

WORKDIR /code
COPY . /code

#RUN pipenv run brownie accounts import somm0 somm0.json
#RUN pipenv run brownie accounts import somm1 somm1.json
#RUN pipenv run brownie accounts import somm2 somm2.json
#RUN pipenv run brownie accounts import somm3 somm3.json

#RUN
CMD startup.sh