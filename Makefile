protoa:
	@protoc -I ./payment-grpc/proto -I ./auth-grpc/proto --go_out=./auth-grpc/proto --go_opt=paths=source_relative \
	--go-grpc_out=./auth-grpc/proto --go-grpc_opt=paths=source_relative \
	auth-grpc/proto/auth.proto

protop:
	@protoc -I proto --go_out=proto --go_opt=paths=source_relative \
	--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
	proto/*.proto