package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Environment string

	MongoDBHost     string
	MongoDBPassword string
	MongoDBDatabase string
	MongoDBUser     string
	MongoDBPort     int

	DiscordAppID     string
	DiscordPublicKey string
	DiscordBotToken  string
	DiscordChannelId string
}

func Load() Config {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Println("failed to retrieve absolute path for .env")
		panic(err)
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		log.Print("No .env file found")
	}

	config := Config{}
	config.Environment = cast.ToString(env("ENVIRONMENT", "develop"))

	config.MongoDBHost = cast.ToString(env("MONGO_DB_HOST", "localhost"))
	config.MongoDBHost = cast.ToString(env("MONGO_DB_HOST", "localhost"))
	config.MongoDBPort = cast.ToInt(env("MONGO_DB_PORT", "27017"))
	config.MongoDBDatabase = cast.ToString(env("MONGO_DB_DATABASE", "wiki"))
	config.MongoDBUser = cast.ToString(env("MONGO_DB_USER", "mongo"))
	config.MongoDBPassword = cast.ToString(env("MONGO_DB_PASSWORD", "mongo"))

	config.DiscordAppID = cast.ToString(env("DISCORD_APP_ID", "YOUR_APP_ID"))
	config.DiscordPublicKey = cast.ToString(env("DISCORD_PUBLIC_KEY", "YOUR_PUBLIC_KEY"))
	config.DiscordBotToken = cast.ToString(env("DISCORD_BOT_TOKEN", "YOUR_BOT_TOKEN"))
	config.DiscordChannelId = cast.ToString(env("DISCORD_CHANNEL_ID", "DISCORD_CHANNEL_ID"))

	return config
}

func env(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
