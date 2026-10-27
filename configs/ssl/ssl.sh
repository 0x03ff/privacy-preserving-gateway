#!/bin/bash 


cd "$(dirname "$0")"

echo "--- Generating server ECC P-384 key pair ---"
openssl ecparam -genkey -name secp384r1 -out go_key.pem
openssl req -new -key go_key.pem -out go_csr.pem -subj "/C=HK/ST=HK/L=HK/O=capstone/CN=privary"
openssl x509 -req -days 365 -in go_csr.pem -signkey go_key.pem -out go_cert.pem

echo "--- Generating CA ECC P-384 key pair ---"
openssl ecparam -genkey -name secp384r1 -out ca_key.pem

echo "--- Signing cert with CA ECC P-384 key pair ---"
openssl req -x509 -new -key ca_key.pem -out ca_cert.pem -subj "/C=HK/ST=HK/L=HK/O=capstone/CN=privary-CA" -days 365
openssl x509 -req -in go_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -set_serial 01 -out go_cert.pem -days 365

echo "--- Merging go_cert and ca_cert into fullchain.pem ---"
cat go_cert.pem ca_cert.pem > fullchain.pem