# wiki_change_stream

## Running 
First you need to create discord app and bot for it.
https://discord.com/developers/applications
Need to turn on `MESSAGE CONTENT INTENT` in bot section

Setup .env using .env.sample (for mongodb and discord credentials)
go mod tidy && go run cmd/main.go

Also you can build it from Dockerfile or use binary from releases.
