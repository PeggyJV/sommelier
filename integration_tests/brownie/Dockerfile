FROM python:3.8-slim as base

# Setup env
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER 1

# Install pipenv and compilation dependencies
RUN pip install pipenv
RUN apt-get update && apt-get install -y --no-install-recommends gcc build-essential python3-dev libssl-dev npm bash

RUN mkdir /testnet
WORKDIR /testnet

# Install python dependencies in /.venv
COPY Pipfile /testnet
COPY Pipfile.lock /testnet
RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy

#ENV PATH="/code/.venv/bin:$PATH"
ENV WEB3_INFURA_PROJECT_ID d6f22be0f7fd447186086d2495779003

RUN npm install -g ganache-cli
EXPOSE 8545

COPY . /testnet

#RUN
CMD ["sh", "startup.sh"]