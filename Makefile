DB_URL =postgresql://postgres:qwerty7867@localhost:5432/simple_bank_db?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)


sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/huzaifac137/simple_bank/db/sqlc Store

db:
	docker run --name simple_bank_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty7867 -e POSTGRES_DB=simple_bank_db -p 5433:5432 -d postgres:16

golang: 
      docker run -p 8080:8080 --name simple_bank_app -e DB_SOURCE="postgresql://postgres:qwerty7867@172.18.0.2:5433/simple_bank_db?sslmode=disable" -d simple_bank:latest
   
proto:
      protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      proto/*.proto

grpc_client:
          grpc-client-cli --proto proto/ localhost:9090   

proto_with_gateway:
                protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	  --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
      proto/*.proto
               

proto_with_gateway_and_swagger:
                        DB_URL =postgresql://postgres:qwerty7867@localhost:5432/simple_bank_db?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)


sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/huzaifac137/simple_bank/db/sqlc Store

db:
	docker run --name simple_bank_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty7867 -e POSTGRES_DB=simple_bank_db -p 5433:5432 -d postgres:16

golang: 
      docker run -p 8080:8080 --name simple_bank_app -e DB_SOURCE="postgresql://postgres:qwerty7867@172.18.0.2:5433/simple_bank_db?sslmode=disable" -d simple_bank:latest
   
proto:
      protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      proto/*.proto

grpc_client:
          grpc-client-cli --proto proto/ localhost:9090   

proto_with_gateway_and_swagger:
                protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	  --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	  --openapiv2_out=swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger
      proto/*.proto
               
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger proto/*.proto
