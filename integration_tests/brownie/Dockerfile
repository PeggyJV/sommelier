FROM python:3.8-slim as base

# Setup env
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER 1

# Install pipenv and compilation dependencies
RUN pip install pipenv
RUN apt-get update && apt-get install -y --no-install-recommends gcc build-essential python3-dev libssl-dev npm bash

RUN mkdir /code
WORKDIR /code

# Install python dependencies in /.venv
COPY Pipfile .
COPY Pipfile.lock .
RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy

ENV PATH="/code/.venv/bin:$PATH"
ENV WEB3_INFURA_PROJECT_ID d6f22be0f7fd447186086d2495779003

RUN npm install -g ganache-cli
EXPOSE 8545

COPY . /code

#RUN pipenv run brownie accounts import somm0 somm0.json
#RUN pipenv run brownie accounts import somm1 somm1.json
#RUN pipenv run brownie accounts import somm2 somm2.json
#RUN pipenv run brownie accounts import somm3 somm3.json

#RUN
CMD ["sh", "startup.sh"]