FROM golang:1.22-alpine

# Add git for downloading dependencies
RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main cmd/main.go

CMD ["./main"]