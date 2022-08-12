gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/image_storage.proto

build:
	docker-compose build app

run:
	docker-compose up app