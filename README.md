# go-grpc-starter
rpc framework with golang for starter

## steps
0. should install protoc first 
`go get -u github.com/golang/protobuf/protoc-gen-go`

1. defined protobuf file first: ProductInfo.proto

2. run
`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative product/ProductInfo.proto`

3. init mod
`go mod init grpc-demo`
`go mod tidy`

3. create server main.go

4. create client main.go