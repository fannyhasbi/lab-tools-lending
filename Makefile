appname := peminjaman-testing
personalChatID := 284324420

run:
	@go run main.go

ngrok:
	@ngrok http 3000

change-server:
	curl -F "url=$(URL)"  https://api.telegram.org/bot1701903841:AAHBGnkqTsEPggVwNt56oNMVW2ynnWbv2OI/setWebhook

deploy:
	heroku container:push web -a $(appname) && \
	heroku container:release web -a $(appname)