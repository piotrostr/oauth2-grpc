#!/bin/bash

if [ ! -d "./certs" ]; then
  echo "Cert dir does not exist, creating"
  mkdir -p ./certs
fi

echo "Generating Certificate"
openssl req -x509 \
  -newkey rsa:4096 \
  -keyout certs/key.pem \
  -out certs/cert.pem \
  -sha256 \
  -subj "/CN=localhost" \
  -days 365
