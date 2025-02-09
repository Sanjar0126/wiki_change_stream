# wiki_change_stream

## Running 
### Discord app setup
First you need to create discord app and bot for it. <br>
https://discord.com/developers/applications<br>
Don't forget to turn on `MESSAGE CONTENT INTENT`, `PRESENCE INTENT`, `SERVER MEMBERS INTENT` in bot section<br>

### Setting up environmental variables
Setup .env using .env.sample (for mongodb and discord credentials)<br>
You can export env variables with:<br>
`set -a && source .env && set +a`

### Running from source
Clone repository:<br>
`https://github.com/Sanjar0126/wiki_change_stream.git`<br>
Get dependencies using:<br>
```
go mod tidy
```
And run cmd/main.go<br>
```
go run cmd/main.go
```

Building binary file from source
```
go build -o builds/main.o cmd/main.go
./builds/main.o
```

### Pre-built binary file
You can download pre-built binary files from:<br>
https://github.com/Sanjar0126/wiki_change_stream/releases<br>
Before running, put .env file in the same directory as binary file's or export env variables.

### Running from Docker
Build Docker image.
```
docker build -t wiki-streaming .
```
Run the container
```
docker run --env-file .env wiki-streaming
```

## Usage
After adding bot to server, send commands in bot's dm.
Commands:
- !ping for testing connection
- !setLang [language_code]: Sets a default language for the user/server session. !setLang en (e.g., ru, fr, es, etc.).
- !recent <offset> <limit>: Retrieves the most recent changes for the current language. Default offset=0, limit=10
- !stats [yyyy-mm-dd]: Displays how many changes occurred on that date for the chosen language.

## Workflow
Used programming language is Go.<br>
For consuming eventsource standard go http client is used. Consuming function is run in goroutine and sends data to channel. Processing function receives data from channel and pushes to db. I used to different goroutines for parallel and independent services and if the database slows down, the event consumer isn’t directly affected. If connection is disconnected or buffer is malfunctioned, goroutine restarts and starts to consume events with latest saved timestamp<br>
For graceful shutdown of goroutines and avoid race conditions, I used standard `sync` package<br>
In order to avoid duplications, I make `meta.id` field unique in database.<br>
For storing wiki and discord user data mongodb is used.<db>
For discord integration https://github.com/bwmarrin/discordgo library is used.<br>

## Design decisions and Trade-offs
This project is built in monolith architecture for simplicity and time constraints. In order to add new event source, code modifications are required. 

I would consider to split app into microservices and use message brokers or event streaming/processing tools. It would help to handle high-traffic services (like your Wikipedia consumer) separately without affecting others, push updates to one service without redeploying the entire app, failure in one service (e.g. Discord bot crash) won’t bring down the entire app, other parts of the system continue to function even if one microservice fails and also they work seamlessly with Docker & Kubernetes for automated deployments.

## Scaling for higher loads
I would use Stream Processing Framework<br>
github.com/ThreeDotsLabs/watermill<br>
Because it is easy to use, universal (messaging, stream processing, CQRS). It can be also used for integrating stream processing system something like Apache Kafka or Redis streams.<br>

Also I would consider to use a time-series database like InfluxDB/TimescaleDB or Cassandra db, because their write throughput is much higher. Also time-series database uses better compression for saving data storage. 