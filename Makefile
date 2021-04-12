appname := peminjaman-testing
personalChatID := 284324420
port := 3000

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