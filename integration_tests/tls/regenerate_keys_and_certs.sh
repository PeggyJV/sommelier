#!/bin/sh

set -ex

# These commands are correct for generating certificates using openssl 1.1.1l
# Check your version by running `openssl version` and if it's LibreSSL, install
# the proper version of openssl using a tool like homebrew

# Set up folders if they don't exist yet

mkdir -p server
mkdir -p client

# Create the private keys

openssl ecparam -name secp384r1 -noout -out server/test_server_ca_key_non-pkcs8.pem -genkey
openssl ecparam -name secp384r1 -noout -out server/test_server_key_non-pkcs8.pem -genkey
openssl ecparam -name secp384r1 -noout -out client/test_client_ca_key_non-pkcs8.pem -genkey
openssl ecparam -name secp384r1 -noout -out client/test_client_key_non-pkcs8.pem -genkey

# Create PKCS8 versions of the private keys, which the Rust libraries expect

openssl pkcs8 -in server/test_server_ca_key_non-pkcs8.pem -out server/test_server_ca_key_pkcs8.pem -topk8 -nocrypt
openssl pkcs8 -in server/test_server_key_non-pkcs8.pem -out server/test_server_key_pkcs8.pem -topk8 -nocrypt
openssl pkcs8 -in client/test_client_ca_key_non-pkcs8.pem -out client/test_client_ca_key_pkcs8.pem -topk8 -nocrypt
openssl pkcs8 -in client/test_client_key_non-pkcs8.pem -out client/test_client_key_pkcs8.pem -topk8 -nocrypt

# Create CAs

set +x
echo
echo "You're going to be asked to fill in fields for the server and client CAs."
echo "These values don't really matter, just go with the defaults."
read -p "Press enter to continue."
echo
echo "====================="
echo "Server CA Certificate"
echo "====================="
echo
set -x

openssl req -x509 -new -key server/test_server_ca_key_pkcs8.pem -out server/test_server_ca.crt -sha384 -days 730

set +x
echo
echo "====================="
echo "Client CA Certificate"
echo "====================="
echo
set -x

openssl req -x509 -new -key client/test_client_ca_key_pkcs8.pem -out client/test_client_ca.crt -sha384 -days 730

# Create CSRs

set +x
echo
echo "You're going to be asked to fill in fields first for the server and client certificates."
echo "The only important setting is to make sure the Common Name is set to localhost"
read -p "Press enter to continue."
echo
echo "=================="
echo "Server Certificate"
echo "=================="
echo
set -x

openssl req -new -sha384 -key server/test_server_key_pkcs8.pem -out server/test_server.csr

set +x
echo
echo "=================="
echo "Client Certificate"
echo "=================="
echo
set -x

openssl req -new -sha384 -key client/test_client_key_pkcs8.pem -out client/test_client.csr

# Sign the server and client certificates with the respective CAs
# the v3.ext file makes sure we include the subjectAltName=localhost extension

openssl x509 -req -in server/test_server.csr -CA server/test_server_ca.crt -CAkey server/test_server_ca_key_pkcs8.pem -CAcreateserial -out server/test_server.crt -sha384 -extfile v3.ext -days 730
openssl x509 -req -in client/test_client.csr -CA client/test_client_ca.crt -CAkey client/test_client_ca_key_pkcs8.pem -CAcreateserial -out client/test_client.crt -sha384 -extfile v3.ext -days 730

# Clean up the certificate requests

rm server/test_server.csr || true
rm client/test_client.csr || true
