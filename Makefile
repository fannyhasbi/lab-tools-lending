appname := peminjaman-testing
personalChatID := 284324420
port := 3000

include .env
export $(shell sed 's/=.*//' .env)

run:
	@go run main.go

test:
	go test -v ./...

ngrok:
	@ngrok http ${port}

change-server:
	curl -F "url=$(URL)"  https://api.telegram.org/bot1701903841:AAHBGnkqTsEPggVwNt56oNMVW2ynnWbv2OI/setWebhook

deploy: test
	heroku container:push web -a $(appname) && \
	heroku container:release web -a $(appname)

migrate-up:
	migrate -path database/migration \
		-database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" \
		-verbose up

migrate-down:
	migrate -path database/migration \
		-database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" \
		-verbose down
