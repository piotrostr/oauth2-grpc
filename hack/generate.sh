#!/bin/bash

module_name=github.com/piotrostr/oauth2-grpc

goa gen $module_name/api/design

# TODO(piotrostr): wait to see if the below is idempotent or would it overwrite the existing
# files
# goa example $module_name/api/design

