protogen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/acknowledge.proto

start-server:
	cd server && go run .  

start-client:
	cd client && go run .  
