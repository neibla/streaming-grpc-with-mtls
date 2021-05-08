# Streaming Bidirectional gRPC Server With mTLS Authentication

## Requirements:
- Golang 1.16

## Install Dependencies
```bash
go get ./...
```

## Setup Development Certificates
```bash
make generate-dev-certs
```

## Start Server
```bash
make start-server
```

## Start Client
```bash
make start-client
```

Configure client certs with client/config.env