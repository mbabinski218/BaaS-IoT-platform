package types

const (
	BCNone BlockchainMode = iota
	BCLightCheck
	BCFullCheck
	BCBatchCheck
	TimeLayout               = "2006-01-02T15:04:05.000000"
	ShortTimeLayout          = "2006-01-02"
	BlockchainBatchStartTime = "2025-01-01T00:00:00.000000"
	MongoDuration            = "mongoDuration"
	BlockchainDuration       = "blockchainDuration"
	TotalDuration            = "totalDuration"
	Missed                   = "missed"
)
