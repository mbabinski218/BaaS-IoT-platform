package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost            string
	Port                  string
	BlockchainUrl         string
	BlockchainPrivateKey  string
	MongoDbUri            string
	MongoDbName           string
	MongoDbCollectionName string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:            getEnv("PUBLIC_HOST", ""),
		Port:                  getEnv("PORT", "8080"),
		BlockchainUrl:         getEnv("BLOCKCHAIN_URL", ""),
		BlockchainPrivateKey:  getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
		MongoDbUri:            getEnv("MONGO_URI", ""),
		MongoDbName:           getEnv("MONGO_DB_NAME", ""),
		MongoDbCollectionName: getEnv("MONGO_COLLECTION_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
