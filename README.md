# payment
go-payment

protofiles -> https://github.com/Edbeer/payment-proto

## Launching auth container
```
payment$ cd auth-grpc
payment/auth-grpc$ docker-compose up --build auth
```
## Launching payment container
```
payment$ cd payment-grpc
payment/payment-grpc$ docker-compose up --build payment
```
## Launching api-gateway container
```
payment$ cd api-gateway
payment/api-gateway$ docker-compose up --build api-gateway
```
## Swagger ```http://localhost:8080/swagger/```
