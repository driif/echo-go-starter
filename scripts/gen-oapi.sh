#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

oapi-codegen -generate types -o "$output_dir/$service/openapi_types.gen.go" -package "$package" "api/paths/$service.yml"
oapi-codegen -generate server,spec -o "$output_dir/$service/openapi_api.gen.go" -package "$package" "api/paths/$service.yml"
