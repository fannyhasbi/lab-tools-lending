appname := peminjaman-testing
port := 3000

include .env
export $(shell sed 's/=.*//' .env)

# migration var
ifdef N
DOWN := $(N)
else
DOWN := 1
endif

run:
	@go run main.go

test:
	go test -race -v ./...

container:
	docker-compose build && docker-compose up

ngrok:
	@ngrok http ${port}

change-server:
	curl -F "url=$(URL)"  https://api.telegram.org/${BOT_TOKEN}/setWebhook

deploy: test
	heroku container:push web -a $(appname) && \
	heroku container:release web -a $(appname)

migrate-up:
	@migrate -path database/migration \
		-database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" \
		-verbose up

migrate-down:
	@migrate -path database/migration \
		-database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" \
		-verbose down ${DOWN}

migrate-force:
	@migrate -path database/migration \
		-database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" \
		-verbose force $(VERSION)