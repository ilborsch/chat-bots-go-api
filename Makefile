run:
	go run ./cmd/chat-bots-api/main.go --config=./config/development.yaml

migrate:
	go run ./cmd/migrator/main.go --config=./config/development.yaml --migrations-path=migrations

mysql:
	sudo /usr/local/mysql/support-files/mysql.server start

swagger:
	swag init -g cmd/chat-bots-api/main.go -o cmd/chat-bots-api/docs

docker-build:
	docker build --no-cache -t backend-api:latest -f Dockerfile.backend .
	docker build --no-cache -t migrator:latest -f Dockerfile.migrator .
	docker-compose build --no-cache

