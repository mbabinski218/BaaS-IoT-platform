package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost                string
	MongoDbUri                string
	MongoDbName               string
	MongoDbCollectionName     string
	BlockchainEnabled         bool
	BlockchainUrl             string
	BlockchainPrivateKey      string
	BlockchainContractAddress string
	AuditEnabled              bool
	AuditTimeout              int64
	AuditSize                 int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load("../.env")

	return Config{
		PublicHost:                getEnv("PUBLIC_HOST", ""),
		MongoDbUri:                getEnv("MONGO_URI", ""),
		MongoDbName:               getEnv("MONGO_DB_NAME", ""),
		MongoDbCollectionName:     getEnv("MONGO_COLLECTION_NAME", ""),
		BlockchainEnabled:         getEnvAsBool("BLOCKCHAIN_ENABLED", false),
		BlockchainUrl:             getEnv("BLOCKCHAIN_URL", ""),
		BlockchainPrivateKey:      getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
		BlockchainContractAddress: getEnv("BLOCKCHAIN_CONTRACT_ADDRESS", ""),
		AuditEnabled:              getEnvAsBool("AUDIT_ENABLED", false),
		AuditTimeout:              getEnvAsInt("AUDIT_TIMEOUT", 3600000),
		AuditSize:                 getEnvAsInt("AUDIT_SIZE", 1000),
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

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return b
	}

	return fallback
}
