#!/usr/bin/env bash

docker run -it --rm --name envoy -p=8082:8082 \
  -v "$(pwd)/pet/v1/pet.pb:/data/pet.pb:ro" \
  -v "$(pwd)/envoy-config.yml:/etc/envoy/envoy.yaml:ro" \
  envoyproxy/envoy:v1.21-latest -l debug -c /etc/envoy/envoy.yaml
