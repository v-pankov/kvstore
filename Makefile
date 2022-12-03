generate:
	protoc --proto_path=api/grpc schema.proto --go_out=. --go-grpc_out=.