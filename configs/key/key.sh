#!/bin/bash

cd "$(dirname "$0")"

# Generate KEY PAIR 1: For JWT SIGNING (ECDSA)
echo "Generating ECDSA key pair for JWT signing (secp384r1)..."
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:secp384r1 -out ECDSA_384_private.pem 
openssl pkey -in ECDSA_384_private.pem -pubout -out ECDSA_384_public.pem 


# Generate AES key (for symmetric encryption)
echo "Generating secret AES key..."
openssl rand -hex 32 > secret.key 

# Verify both key pairs (PKCS#8 format is standard for JWT)
echo "Verifying ECDSA signing key pair:"
openssl ec -in ECDSA_384_private.pem -text -noout
openssl pkey -in ECDSA_384_public.pem -pubin -text -noout