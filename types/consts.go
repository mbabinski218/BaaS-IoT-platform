package types

const (
	BCNone BlockchainMode = iota
	BCLightCheck
	BCFullCheck
	BCBatchCheck
	TimeLayout               = "2006-01-02T15:04:05.000000"
	BlockchainBatchStartTime = "2025-01-01T00:00:00.000000"
)
