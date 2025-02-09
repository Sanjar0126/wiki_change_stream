# wiki_change_stream

## Running 
First you need to create discord app and bot for it. <br>
https://discord.com/developers/applications<br>
Need to turn on `MESSAGE CONTENT INTENT` in bot section<br>

Setup .env using .env.sample (for mongodb and discord credentials)<br>
go mod tidy && go run cmd/main.go<br>

Also you can build it from Dockerfile or use binary from releases.<br>
