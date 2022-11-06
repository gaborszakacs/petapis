# PetAPIs

Protobuf/gRPC playground based on https://github.com/bufbuild/buf-tour

## Features

- Generate Go gRPC client and server stubs based on Proto definitions.
- Use Envoy to get a JSON API for ~free via [transcoding](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter)
- Generate OpenAPI description and serve interactive documentation.

## Usage

```
# Lint Protobuf
buf lint

# Detect breaking changes
buf breaking --against ./.git#branch=master

# Generate Go stubs and OpenAPI documentation
buf generate

# Generate descriptor for Envoy to consume
buf build --as-file-descriptor-set --output pet/v1/pet.pb

# Run gRPC server
go run ./server/main.go

# Run envoy
./start-envoy.sh

# Make gRPC request
go run ./client/main.go

# Make JSON request
curl localhost:8082/v1/pets/Bolyhos

# Browse OpenAPI doc (and make JSON requests)
http://localhost:8082/docs/

```

## Connect Demo

```
# Start BE
go run ./connect-server/main.go


# Call it with curl
curl \
    --header "Content-Type: application/json" \
    --data '{"pet_id": "pet2"}' \
    http://localhost:8080/pet.v1.PetStoreService/GetPet


# Call it with grpcurl
grpcurl \
    --plaintext \
    --d '{"pet_id": "pet2"}' \
    localhost:8080 \
    pet.v1.PetStoreService/GetPet


# Start FE
cd ./connect-frontend && yarn dev

# Make changes (e.g add age to put)
buf push
go get go.buf.build/bufbuild/connect-go/gaborszakacs-bitrise/pet@latest
npm install @buf/bufbuild_connect-web_gaborszakacs-bitrise_pet

```


