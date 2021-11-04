FROM python:3.8-slim as base

# Setup env
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER 1

#FROM base AS python-deps

# Install pipenv and compilation dependencies
RUN pip install pipenv
RUN apt-get update && apt-get install -y --no-install-recommends gcc build-essential python3-dev libssl-dev npm bash

# Install python dependencies in /.venv
COPY Pipfile .
COPY Pipfile.lock .
RUN PIPENV_VENV_IN_PROJECT=1 pipenv install --deploy


#FROM base AS runtime
# Copy virtual env from python-deps stage
#COPY --from=python-deps /.venv /.venv
ENV PATH="/.venv/bin:$PATH"

RUN npm install -g ganache-cli

COPY . .

#RUN pipenv run brownie accounts import somm0 somm0.json
#RUN pipenv run brownie accounts import somm1 somm1.json
#RUN pipenv run brownie accounts import somm2 somm2.json
#RUN pipenv run brownie accounts import somm3 somm3.json

#RUN
CMD ["sh", "startup.sh"]