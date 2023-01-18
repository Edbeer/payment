protoa:
	@protoc -I ./payment-grpc/proto -I ./auth-grpc/proto --go_out=./auth-grpc/proto --go_opt=paths=source_relative \
	--go-grpc_out=./auth-grpc/proto --go-grpc_opt=paths=source_relative \
	auth-grpc/proto/auth.proto

protop:
	@protoc -I ./payment-grpc/proto --go_out=./payment-grpc/proto --go_opt=paths=source_relative \
	--go-grpc_out=./payment-grpc/proto --go-grpc_opt=paths=source_relative \
	payment-grpc/proto/*.proto

evansp:
	evans --path ./payment-grpc --path . --proto proto/payment.proto

evansa:
	evans -r repl -p 50052

