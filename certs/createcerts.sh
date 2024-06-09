#!/bin/bash

set -e

echo "Generating CA key and certificate:"
openssl req -x509 -sha256 -noenc -days 7 -newkey rsa:2048 \
  -keyout ca.key -out ca.pem \
  -subj "/O=Nodda Certificate Authority/"

echo "Generating CSR"
openssl req -out cert.csr -new -newkey rsa:2048 -noenc -keyout private.key \
  -subj "/O=quic-go-samples/"

echo "Sign certificate:"
openssl x509 -req -sha256 -days 7 -in cert.csr  -out cert.pem \
  -CA ca.pem -CAkey ca.key -CAcreateserial \
  -extfile <(printf "subjectAltName=DNS:localhost")

# debug output the certificate
openssl x509 -noout -text -in cert.pem

# we don't need the CA key, the serial number and the CSR any more
# rm ca.key cert.csr ca.srl
