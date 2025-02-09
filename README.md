# wiki_change_stream

## Running 
First you need to create discord app and bot for it. <br>
https://discord.com/developers/applications<br>
Need to turn on `MESSAGE CONTENT INTENT` `PRESENCE INTENT` `SERVER MEMBERS INTENT` in bot section<br>

Setup .env using .env.sample (for mongodb and discord credentials)<br>
go mod tidy && go run cmd/main.go<br>

Also you can build it from Dockerfile or use binary from releases.<br>

## Using
After adding bot to server, send commands in bot's dm.
Commands:
- !ping for testing connection
- !setLang [language_code]: Sets a default language for the user/server session. !setLang en (e.g., ru, fr, es, etc.).
- !recent <offset> <limit>: Retrieves the most recent changes for the current language. Default offset=0, limit=10
- !stats [yyyy-mm-dd]: Displays how many changes occurred on that date for the chosen language.

##Workflow
Used programming language is Go.<br>
For consuming eventsource standard go http client is used. Consuming function is run in goroutine and sends data to channel. Processing function receives data from channel and pushes to db.<br>
For storing wiki and discord user data mongodb is used.<db>
For discord integration `https://github.com/bwmarrin/discordgo` is used.<br>
