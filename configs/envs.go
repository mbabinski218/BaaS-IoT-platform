package configs

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
)

type Config struct {
	PublicHost                      string
	MongoDbUri                      string
	MongoDbName                     string
	MongoDbCollectionName           string
	MongoContextTimeout             int64
	BlockchainMode                  types.BlockchainMode
	BlockchainUrl                   string
	BlockchainPrivateKey            string
	BlockchainContractAddress       string
	BlockchainContextTimeout        int64
	BlockchainGasLimit              int64
	BlockchainGasTipCap             int64
	BlockchainGasFeeCap             int64
	BlockchainBatchInterval         int64
	BlockchainBatchContractAddress  string
	BlockchainSecondsPerBlock       int64
	BlockchainCheckpoints           []uint64
	BlockchainCheckpointCallRepeats int64
	BlockchainValidators            string
	BlockchainServerIP              string
	BlockchainServerPort            string
	BlockchainAsyncMode             bool
	BlockchainCustomBatchStartTime  string
	IotSimulatorCommand             string
	IotSimulatorPath                string
	IotSimulatorParams              string
	AuditEnabled                    bool
	AuditTimeout                    int64
	AuditSize                       int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load("../.env")

	return Config{
		PublicHost:                      getEnv("PUBLIC_HOST", ""),
		MongoDbUri:                      getEnv("MONGO_URI", ""),
		MongoDbName:                     getEnv("MONGO_DB_NAME", ""),
		MongoDbCollectionName:           getEnv("MONGO_COLLECTION_NAME", ""),
		MongoContextTimeout:             getEnvAsInt("MONGO_CONTEXT_TIMEOUT", 5), // Default is 5 seconds
		BlockchainMode:                  getEnvAsBCMode("BLOCKCHAIN_MODE", types.BCNone),
		BlockchainUrl:                   getEnv("BLOCKCHAIN_URL", ""),
		BlockchainPrivateKey:            getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
		BlockchainContractAddress:       getEnv("BLOCKCHAIN_CONTRACT_ADDRESS", ""),
		BlockchainContextTimeout:        getEnvAsInt("BLOCKCHAIN_CONTEXT_TIMEOUT", 30), // Default is 30 seconds
		BlockchainGasLimit:              getEnvAsInt("BLOCKCHAIN_GAS_LIMIT", 0),        // Default 0 means it will be set by the network
		BlockchainGasTipCap:             getEnvAsInt("BLOCKCHAIN_GAS_TIP_CAP", 0),      // Default 0 means it will be set by the network
		BlockchainGasFeeCap:             getEnvAsInt("BLOCKCHAIN_GAS_FEE_CAP", 0),      // Default 0 means it will be set by the network
		BlockchainBatchInterval:         getEnvAsInt("BLOCKCHAIN_BATCH_INTERVAL", 15),  // Default is 15 seconds
		BlockchainBatchContractAddress:  getEnv("BLOCKCHAIN_BATCH_CONTRACT_ADDRESS", ""),
		BlockchainSecondsPerBlock:       getEnvAsInt("BLOCKCHAIN_SECONDS_PER_BLOCK", 15), // Default is 15 seconds
		BlockchainCheckpoints:           getEnvAsUintArray("BLOCKCHAIN_CHECKPOINTS", []uint64{}),
		BlockchainCheckpointCallRepeats: getEnvAsInt("BLOCKCHAIN_CHECKPOINT_CALL_REPEATS", 10), // Default is 10 repeats
		BlockchainValidators:            getEnv("BLOCKCHAIN_VALIDATORS", ""),
		BlockchainServerIP:              getEnv("BLOCKCHAIN_SERVER_IP", "127.0.0.1"),
		BlockchainServerPort:            getEnv("BLOCKCHAIN_SERVER_PORT", "22"),
		BlockchainAsyncMode:             getEnvAsBool("BLOCKCHAIN_ASYNC_MODE", true), // Default is true
		BlockchainCustomBatchStartTime:  getEnv("BLOCKCHAIN_CUSTOM_BATCH_START_TIME", ""),
		IotSimulatorCommand:             getEnv("IOT_SIMULATOR_COMMAND", ""),
		IotSimulatorPath:                getEnv("IOT_SIMULATOR_PATH", ""),
		IotSimulatorParams:              getEnv("IOT_SIMULATOR_PARAMS", ""),
		AuditEnabled:                    getEnvAsBool("AUDIT_ENABLED", false),
		AuditTimeout:                    getEnvAsInt("AUDIT_TIMEOUT", 3600000), // Default is 1 hour in milliseconds
		AuditSize:                       getEnvAsInt("AUDIT_SIZE", 1000),
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

func getEnvAsBCMode(key string, fallback types.BlockchainMode) types.BlockchainMode {
	if value, ok := os.LookupEnv(key); ok {
		switch value {
		case "None":
			return types.BCNone
		case "Light":
			return types.BCLightCheck
		case "Full":
			return types.BCFullCheck
		case "Batch":
			return types.BCBatchCheck
		default:
			return fallback
		}
	}

	return fallback
}

func getEnvAsUintArray(key string, fallback []uint64) []uint64 {
	if value, ok := os.LookupEnv(key); ok {
		var result []uint64
		for _, v := range splitAndTrim(value, ",") {
			if i, err := strconv.ParseUint(v, 10, 64); err == nil {
				result = append(result, i)
			}
		}
		return result
	}

	return fallback
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
