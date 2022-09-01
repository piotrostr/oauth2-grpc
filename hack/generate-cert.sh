#!/bin/bash

if [ ! -d "./certs" ]; then
  echo "Cert dir does not exist, creating"
  mkdir -p ./certs
fi

echo "Generating Certificate"
openssl req -x509 \
  -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -sha256 \
  -subj "/CN=localhost" \
  -days 365
